package communication

import (
	. "../globalChans"
	. "../globalStructs"
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

var broadcastAddr string = "255.255.255.255:20059" //"129.241.187.255:20059"
var commonPort string = "20059"
var broadcastPort string = "30059"
var listOfElevatorsInNetwork []string
var localAddr string
var connectionList map[string]*net.UDPConn

const (
	arrayX             = 7
	arrayY             = 7
	arrayZ             = 7
	ACK_DEADLINE       = 4
	BROADCAST_DEADLINE = 3
)

type ackTimer struct {
	Message  MessageWithHeader
	DeadLine time.Time
	IpList   []string
}

type addrAndTimer struct {
	DeadLine    time.Time
	NetWorkAddr string
}

type MessageWithHeader struct {
	Tag        string
	Ack        bool
	SenderAddr string
	Data       []byte
}

//Test udpSize by sending this on the network.
type HugeStruct struct {
	HugeNumber int
	HugeBool   bool
	HugeArray  [arrayX][arrayY][arrayZ]int
	HugeName   string
}

func CommNeedBetterName() {
	newShoutFromElevatorChan := make(chan string)
	newElevatorChan := make(chan string)
	newConnectionChan := make(chan string)
	elevatorGoneChan := make(chan string)
	endConnectionChan := make(chan string)
	sendNetworkMessageChan := make(chan []byte)
	messageFromNetworkChan := make(chan []byte)

	ackdByAllChan := make(chan MessageWithHeader)
	resendMessageChan := make(chan MessageWithHeader)
	newAckFromNetworkChan := make(chan MessageWithHeader)
	newAckStartChan := make(chan MessageWithHeader)
	elevatorListChangedChan := make(chan bool)

	localAddr = GetLocalIP()
	if localAddr == "" {
		fmt.Println("Problem using getLocalIP, expecting office-pc")
		localAddr = "129.241.154.78"
	}

	go broadcastPrecense(broadcastPort)
	go listenForBroadcast(broadcastPort, newShoutFromElevatorChan)

	go readUpdateElevatorOverview(newShoutFromElevatorChan, newElevatorChan, elevatorGoneChan)
	go SendMessagesToAllElevators(sendNetworkMessageChan, newConnectionChan, endConnectionChan)
	go readMessagesFromNetwork(localAddr, commonPort, messageFromNetworkChan)
	go decodeMessagesFromNetwork(messageFromNetworkChan, newAckFromNetworkChan, sendNetworkMessageChan)
	go encodeMessagesToNetwork(sendNetworkMessageChan, newAckStartChan, resendMessageChan)
	go setDeadlinesForAcks(resendMessageChan, ackdByAllChan, newAckStartChan, newAckFromNetworkChan, elevatorListChangedChan)

	go func() {
		for {
			select {
			case elevGone := <-elevatorGoneChan:
				fmt.Println("ip list before going through the list. Case: elevgone", listOfElevatorsInNetwork)
				pos := -1
				for i, v := range listOfElevatorsInNetwork {
					if v == elevGone {
						pos = i
						break
					}
				}
				if pos == -1 {
					fmt.Println("main go func: elevator not found, pos == -1")
				} else {
					fmt.Println("position in list", pos)
					fmt.Println("list before change: ", listOfElevatorsInNetwork)
					listOfElevatorsInNetwork = append(listOfElevatorsInNetwork[:pos], listOfElevatorsInNetwork[pos+1:]...)
					fmt.Println("list after change: ", listOfElevatorsInNetwork)

				}
				fmt.Println("Elevator gone, address: ", elevGone)
				endConnectionChan <- elevGone
				elevatorListChangedChan <- true
				FromNetworkElevGoneChan <- elevGone
				if len(listOfElevatorsInNetwork) == 0 {
					FromNetworkNetworkDownChan <- true
				}

			case newElev := <-newElevatorChan:
				fmt.Println("New elevator, address: ", newElev)
				fmt.Println("list before change: ", listOfElevatorsInNetwork)
				listOfElevatorsInNetwork = append(listOfElevatorsInNetwork, newElev)
				fmt.Println("after appending ip to list : ", listOfElevatorsInNetwork)
				newConnectionChan <- newElev
				elevatorListChangedChan <- true
				FromNetworkNewElevChan <- newElev
				if len(listOfElevatorsInNetwork) == 1 {
					FromNetworkNetworkUpChan <- true
				}
			}
		}
	}()
	for {
		select {
		case <-ackdByAllChan:
			fmt.Println("ackd by all!")
		case <-FromNetworkNetworkUpChan:
			fmt.Println("Network is up!")

		}
	}
}

func readUpdateElevatorOverview(newShoutFromElevatorChan chan string, newElevatorChan chan string, elevatorGoneChan chan string) {
	elevatorTimerList := []addrAndTimer{}
	for {
		select {
		case newShout := <-newShoutFromElevatorChan:
			elevatorInList := false
			index := 0
			//Index might be invalid if list is altered somewhere else during the for loop.
			for i, v := range elevatorTimerList {
				if v.NetWorkAddr == newShout {
					elevatorInList = true
					index = i
					break
				}
			}
			if elevatorInList {
				elevatorTimerList[index].DeadLine = time.Now().Add(time.Second * BROADCAST_DEADLINE)
			} else {
				elevatorTimerList = append(elevatorTimerList, addrAndTimer{DeadLine: time.Now().Add(time.Second * BROADCAST_DEADLINE), NetWorkAddr: newShout})
				newElevatorChan <- newShout
			}
		default:
			for i, v := range elevatorTimerList {
				if time.Now().After(v.DeadLine) {
					elevatorGoneChan <- v.NetWorkAddr
					elevatorTimerList = append(elevatorTimerList[:i], elevatorTimerList[i+1:]...) //From slicetricks. Remove element.
				}
			}
		}
		time.Sleep(time.Millisecond * 200)
	}
}

func connectToElevator(remoteIp string, remotePort string) *net.UDPConn {
	fullAddr := remoteIp + ":" + remotePort
	remoteUDPAddr, _ := net.ResolveUDPAddr("udp4", fullAddr)
	connection, _ := net.DialUDP("udp4", nil, remoteUDPAddr)

	return connection

}

func GetLocalIP() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Whoops..")
		panic(err)
	}
	for _, i := range ifaces {
		addrs, _ := i.Addrs()

		ipLow := net.ParseIP("129.241.187.000")
		ipHigh := net.ParseIP("129.241.187.255")
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if bytes.Compare(ip, ipLow) >= 0 && bytes.Compare(ip, ipHigh) <= 0 {
				//fmt.Println("Printing type: ", reflect.TypeOf(ip))
				return ip.String()
			}
		}
	}
	return ""
}

func broadcastPrecense(broadcastPort string) {
	broadcastAddr := "255.255.255.255" + ":" + broadcastPort
	broadcastUDPAddr, _ := net.ResolveUDPAddr("udp", broadcastAddr)
	connection, err := net.DialUDP("udp", nil, broadcastUDPAddr)
	if err != nil {
		fmt.Println("You messed up in spam presence")
	}
	defer connection.Close()
	for {
		_, err := connection.Write([]byte("I am here"))
		if err != nil {
			fmt.Println("error in bcastpres: ", err)
		}
		time.Sleep(time.Millisecond * 500)
	}
}

func listenForBroadcast(broadcastPort string, newShoutFromElevatorChan chan string) {
	buffer := make([]byte, 2048)
	listenBroadcastAddress := "0.0.0.0" + ":" + broadcastPort
	broadcastUDPAddr, _ := net.ResolveUDPAddr("udp4", listenBroadcastAddress)
	connection, _ := net.ListenUDP("udp4", broadcastUDPAddr)
	defer connection.Close()
	for {
		_, senderAddr, err := connection.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading msg in listenForElevators, discard message")
		} else {
			newShoutFromElevatorChan <- senderAddr.IP.String()
		}
	}

}

/* I think you only need to listen on one port */
/* test to iterate through connections if this doesnt work*/

func readMessagesFromNetwork(localAddr string, commonPort string, messageFromNetworkChan chan []byte) {
	buffer := make([]byte, 2048)
	fullAddr := localAddr + ":" + commonPort
	listenConnAddress, _ := net.ResolveUDPAddr("udp4", fullAddr)
	listenConn, err := net.ListenUDP("udp4", listenConnAddress)
	if err != nil {
		fmt.Println("Error setting up listen-connection")
	}
	for {
		packetLength, _, err := listenConn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading message, do nothing")
		} else {
			//Send messages to decodeFunc here.
			byteSliceToDecoder := make([]byte, packetLength)
			copy(byteSliceToDecoder, buffer[:packetLength])
			messageFromNetworkChan <- byteSliceToDecoder
		}
	}
}

//
func SendMessagesToAllElevators(sendNetworkMessageChan chan []byte, newConnectionChan chan string, endConnectionChan chan string) {
	connectionList := make(map[string]*net.UDPConn)
	for _, deferConn := range connectionList {
		defer deferConn.Close()
	}
	for {
		select {
		case newElevator := <-newConnectionChan:
			connectionList[newElevator] = connectToElevator(newElevator, commonPort)
			defer connectionList[newElevator].Close() //This might be enough
		case deadElevator := <-endConnectionChan:
			connectionList[deadElevator].Close()
			delete(connectionList, deadElevator)
		case newMessage := <-sendNetworkMessageChan:
			for _, conn := range connectionList {
				_, err := conn.Write(newMessage)
				if err != nil {
					fmt.Println("Error in sending func, maybe some error handling later")
					fmt.Println("Err: ", err)
				}
			}
		}
	}
}
AssignedTo string
	SentFrom   string
	Order      Order

//Not completely tested yet
func encodeMessagesToNetwork(sendToNetworkChan chan []byte, sendToAckTimerChan chan MessageWithHeader, resendMessageChan chan MessageWithHeader) {
	go func() {
		for {
			expiredMessage := <-resendMessageChan
			expiredMessage.Ack = false
			expiredMessage.SenderAddr = localAddr
			encodedPacket, err := json.Marshal(expiredMessage)
			if err != nil {
				fmt.Println("error when encoding: ", err)
			}
			fmt.Println("Ack expired, resend.")
			
			/*if expiredMessage.Tag == "ordTo"{
				var orderAss OrderAssigned
				_ := json.Unmarshal(message.Data, &orderAss)
				fmt.Println("Content of expired orderAss message.: ")
				fmt.Println("assigned to: ",orderAss.AssignedTo)
				fmt.Println("sent from", orderAss.SentFrom)
				fmt.Println("order", orderAss.Order)

			}*/
			sendToNetworkChan <- encodedPacket
		}

	}()
	for {
		var tag string = ""
		var encodedData []byte = nil
		var err error
		fmt.Println("encodemessages: start of select")
		select {
		case newOrder := <-ExternalButtonPressedChan:
			tag = "newOr"
			encodedData, err = json.Marshal(newOrder)
			fmt.Println("encodeMessagesToNetwork: order encoded1")
			if err != nil {
				fmt.Println("Error when encoding: ", err)
			}
		case orderTo := <-ToNetworkOrderAssignedToChan:
			tag = "ordTo"
			encodedData, err = json.Marshal(orderTo)
			fmt.Println("encodeMessagesToNetwork: order encoded2")
			if err != nil {
				fmt.Println("Error when encoding: ", err)
			}
		case orderServed := <-ToNetworkOrderServedChan:
			tag = "ordSe"
			encodedData, err = json.Marshal(orderServed)
			fmt.Println("encodeMessagesToNetwork: order encoded3")
			if err != nil {
				fmt.Println("Error when encoding: ", err)
			}
		case elevState := <-ToNetworkNewElevStateChan:
			tag = "elSta"
			encodedData, err = json.Marshal(elevState)
			fmt.Println("encodeMessagesToNetwork: order encoded4")
			if err != nil {
				fmt.Println("Error when encoding: ", err)
			}
		case externalArray := <-ToNetworkExternalArrayChan:
			tag = "extAr"
			encodedData, err = json.Marshal(externalArray)
			fmt.Println("encodeMessagesToNetwork: order encoded5")
			if err != nil {
				fmt.Println("Error when encoding: ", err)
			}
		}
		//alt felles gjøres her
		var newPacket MessageWithHeader = MessageWithHeader{Data: encodedData, Tag: tag, Ack: false, SenderAddr: localAddr}
		encodedPacket, err := json.Marshal(newPacket)
		if err != nil {
			fmt.Println("error when encoding: ", err)
		}
		//Send packet to network. Send copy to local center that keeps track of Ack's
		//Better names needed all over the place
		sendToNetworkChan <- encodedPacket
		fmt.Println("packet sent to network    LENGTH OF PACKET:  ", len(encodedPacket))
		sendToAckTimerChan <- newPacket
		fmt.Println("End of encodedfunc")
	}
}

func sendAck(message MessageWithHeader, sendToNetworkChan chan []byte) {
	ackMessage := MessageWithHeader{}
	ackMessage.Ack = true
	ackMessage.Tag = message.Tag
	ackMessage.SenderAddr = localAddr
	ackMessage.Data = make([]byte, len(message.Data))
	copy(ackMessage.Data, message.Data)
	encodedPacket, err := json.Marshal(ackMessage)
	if err != nil {
		fmt.Println("error encoding resend message")
	}
	fmt.Println("Ack sent")
	sendToNetworkChan <- encodedPacket
}

func decodeMessagesFromNetwork(messageFromNetworkChan chan []byte, newAckFromNetworkChan chan MessageWithHeader, sendToNetworkChan chan []byte) {
	//If no ack. Respond immediately somewhere with ack.
	//Switch on tag
	for {
		packetFromNetwork := <-messageFromNetworkChan
		var message MessageWithHeader = MessageWithHeader{}
		err := json.Unmarshal(packetFromNetwork, &message)
		if err != nil {
			fmt.Println("Something went wrong when unmarshalling packet, do nothing (for now)")
			fmt.Println(err)
		} else if message.Ack == true {
			newAckFromNetworkChan <- message
			//Things will happen here
		} else {
			//Maybe have this ack-echo somewhere else. For now its here
			sendAck(message, sendToNetworkChan)
			switch message.Tag {
			case "newOr":
				var newOrder Order
				err := json.Unmarshal(message.Data, &newOrder)
				if err != nil {
					fmt.Println(err)
				} else {
					FromNetworkNewOrderChan <- newOrder
				}
			case "ordSe":
				var orderServed Order
				err := json.Unmarshal(message.Data, &orderServed)
				if err != nil {
					fmt.Println(err)
				} else {
					FromNetworkOrderServedChan <- orderServed
				}
			case "ordTo":
				var orderAss OrderAssigned
				err := json.Unmarshal(message.Data, &orderAss)
				if err != nil {
					fmt.Println(err)
				} else {
					FromNetworkOrderAssignedToChan <- orderAss
				}
			case "elSta":
				var newState ElevatorState
				err := json.Unmarshal(message.Data, &newState)
				if err != nil {
					fmt.Println(err)
				} else {
					FromNetworkNewElevStateChan <- newState
				}
			case "extAr":
				var newExternalArray [NUM_FLOORS][NUM_BUTTONS-1]int
				err := json.Unmarshal(message.Data, &newExternalArray)
				if err != nil {
					fmt.Println(err)
				} else {
					FromNetworkExternalArrayChan <- newExternalArray
				}
			}
		}
		fmt.Println("Message ackd and sent to local")
	}
}

/*
Ting skal helst kun sendes en gang over nett. Det er greit å discarde meldinger, men

*/
func setDeadlinesForAcks(resendMessageChan chan MessageWithHeader, ackdByAllChan chan MessageWithHeader, newMessageSentChan chan MessageWithHeader, newAckChan chan MessageWithHeader, elevatorListChangedChan chan bool) {
	var unAckdMessages []ackTimer
	/*

		NEED HERE:
		-A structure that contains deadlines for all the elevators for each message
		-To tell someone when the deadline is triggered and package needs to be resent
		-To tell someone when the outgoing message is received by everyone.
		-Need to have a notion of who has sent these acks. Need to include IP in message.
	*/
	for {
		select {
		case <-elevatorListChangedChan:
			//Resend everything
			for i, v := range unAckdMessages {
				unAckdMessages[i].DeadLine = time.Now().Add(time.Second * ACK_DEADLINE)
				localIPlist := GetElevList()
				unAckdMessages[i].IpList = localIPlist
				resendMessageChan <- v.Message //With nil entry in senderaddress, Maybe not though.
			}

			//fmt.Println("listOfElevatorsInNetwork when elevator change: ", listOfElevatorsInNetwork)
		case newAck := <-newAckChan:
			senderOfAck := newAck.SenderAddr

			//fmt.Println("sender of ack: ", senderOfAck)
			//fmt.Println("Ready to acknowledge")
			//maybe add some testing when element isnt in map, or do nothing
			//Find same message, then find sender and remove from the list of senders still waiting to be acked
			//If there are no more acks, assume everyone has received it. Confirm as known in network
		Loop:
			for i, v := range unAckdMessages {
				if bytes.Equal(v.Message.Data, newAck.Data) {
					for j, addr := range unAckdMessages[i].IpList {
						if addr == senderOfAck {
							//fmt.Println("Address found:", listOfElevatorsInNetwork)
							unAckdMessages[i].IpList = append(unAckdMessages[i].IpList[:j], unAckdMessages[i].IpList[j+1:]...)
							//fmt.Println("After change: ", listOfElevatorsInNetwork)
							if len(unAckdMessages[i].IpList) == 0 {
								ackdByAllChan <- unAckdMessages[i].Message
								unAckdMessages = append(unAckdMessages[:i], unAckdMessages[i+1:]...)
								break Loop
							}
						}
					}
				}
			}
		case newMessage := <-newMessageSentChan:
			notInUnackdMessages := true
			for _, v := range unAckdMessages {
				if bytes.Equal(v.Message.Data, newMessage.Data) {
					notInUnackdMessages = false
				}

			}
			if notInUnackdMessages {
				localIPlist := make([]string, len(listOfElevatorsInNetwork))
				copy(localIPlist, listOfElevatorsInNetwork)
				newAckTimer := ackTimer{Message: newMessage, DeadLine: time.Now().Add(time.Second * ACK_DEADLINE), IpList: localIPlist}
				unAckdMessages = append(unAckdMessages, newAckTimer)
			}
		default:
			for i, v := range unAckdMessages {
				if time.Now().After(v.DeadLine) {
					resendMessageChan <- v.Message
					unAckdMessages[i].DeadLine = time.Now().Add(time.Second * ACK_DEADLINE)
					localIPlist := make([]string, len(listOfElevatorsInNetwork))
					copy(localIPlist, listOfElevatorsInNetwork)
					unAckdMessages[i].IpList = localIPlist
				}
			}
		}
		time.Sleep(time.Millisecond * 100)
	}
}

func printHugeStruct(bigAssMessage HugeStruct) {
	fmt.Println("HugeBool: ", bigAssMessage.HugeBool)
	fmt.Println("HugeName aka localAddr: ", bigAssMessage.HugeName)
	fmt.Println("HugeNumber: ", bigAssMessage.HugeNumber)
	for i := 0; i < arrayX; i++ {
		for j := 0; j < arrayY; j++ {
			for k := 0; k < arrayZ; k++ {
				fmt.Println(bigAssMessage.HugeArray[i][j][k])
			}
		}
	}
}

func makeHugeStruct(localAddr string) HugeStruct {
	hugeS := HugeStruct{}
	hugeS.HugeBool = true
	hugeS.HugeNumber = 322322
	hugeS.HugeName = localAddr
	for i := 0; i < arrayX; i++ {
		for j := 0; j < arrayY; j++ {
			for k := 0; k < arrayZ; k++ {
				hugeS.HugeArray[i][j][k] = i*100 + j*10 + k
			}
		}
	}
	return hugeS
}

func GetElevList() []string {
	copyElevList := make([]string, len(listOfElevatorsInNetwork))
	copy(copyElevList, listOfElevatorsInNetwork)
	return copyElevList
}
