package main //communication for når denne skal kjøres, utelukkende

import (
	"bytes"
	"fmt"
	"net"
	"reflect"
	"time"
	"encoding/json"
)

var broadcastAddr string = "255.255.255.255:20059" //"129.241.187.255:20059"
var commonPort string = "20059"
var broadcastPort string = "30059"
var listOfElevatorsInNetwork []string;
var localAddr string

var connectionList map[string] *net.UDPConn

const(
	arrayX = 7
	arrayY = 7
	arrayZ = 7
)


type addrAndTimer struct{
	DeadLine time.Time
	NetWorkAddr string
}

type MessageWithHeader struct{
	Tag string
	SenderAddr string
	Ack bool
	Data []byte 	
}

//Test udpSize by sending this on the network. 
type HugeStruct struct{
	HugeNumber int
	HugeBool bool
	HugeArray [arrayX][arrayY][arrayZ]int
	HugeName string
}

func main() {
	newShoutFromElevatorChan := make(chan string)
	newElevatorChan := make(chan string)
	newConnectionChan := make(chan string)
	elevatorGoneChan := make(chan string)
	endConnectionChan := make(chan string)
	sendNetworkMessageChan := make(chan []byte)
	messageFromNetworkChan := make(chan []byte)

	localAddr := getLocalIP()

	if localAddr == "" {
		fmt.Println("Problem using getLocalIP, expecting office-pc")
		localAddr = "129.241.154.78"
	}
	fmt.Println(localAddr)

	go broadcastPrecense(broadcastPort)
	go listenForBroadcast(broadcastPort, newShoutFromElevatorChan)
	go readUpdateElevatorOverview(newShoutFromElevatorChan, newElevatorChan, elevatorGoneChan)
	go SendMessagesToAllElevators(sendNetworkMessageChan, newConnectionChan, endConnectionChan)
	go readMessagesFromNetwork(localAddr, commonPort, messageFromNetworkChan)
	go decodeMessagesFromNetwork(messageFromNetworkChan)
	go func(){
		for{
			select{
			case elevGone := <- elevatorGoneChan:
				fmt.Println("Elevator gone, address: ", elevGone)
				endConnectionChan <- elevGone
			case newElev := <- newElevatorChan:
				fmt.Println("New elevator, address: ",newElev)
				newConnectionChan <- newElev
			}
		}
	}()
	time.Sleep(time.Second*1)
	//Temp encoding
	
	hugeS := makeHugeStruct(localAddr)
	//var solo int = 322
	localData, err := json.Marshal(hugeS) 
	if err != nil{
		fmt.Println("Something went wrong when marshaling")
	}
	messageToBeSent := MessageWithHeader{
		Tag: "hugeS",
		Ack: false,
		Data: localData,
	}
	encodedMessage, err := json.Marshal(messageToBeSent)
	if err != nil{
		fmt.Println("Something went wrong in second marshal")
	}
	sendNetworkMessageChan <- encodedMessage
	/*
	go func(){
		for{
			sendNetworkMessageChan <- []byte("I am yolomaster")
			time.Sleep(time.Millisecond*500)
		}
	}()*/

	//defer connection.Close()
	fmt.Println("End of main")
	time.Sleep(time.Second*100)
}




func readUpdateElevatorOverview(newShoutFromElevatorChan chan string, newElevatorChan chan string, elevatorGoneChan chan string){
	elevatorTimerList := []addrAndTimer{}
	for{
		select{
		case newShout := <- newShoutFromElevatorChan:
			elevatorInList := false
				index := 0;
			//Index might be invalid if list is altered somewhere else during the for loop.
			for i, v := range elevatorTimerList{
				if v.NetWorkAddr == newShout{
					elevatorInList = true
					index = i;
					break;
				}
			}
			if elevatorInList{
				elevatorTimerList[index].DeadLine = time.Now().Add(time.Second * 3) 
			}else{
				elevatorTimerList = append(elevatorTimerList, addrAndTimer{DeadLine: time.Now().Add(time.Second *3), NetWorkAddr: newShout})
				newElevatorChan <- newShout
			}
		default:
			for i, v := range elevatorTimerList{
				if time.Now().After(v.DeadLine){
					elevatorGoneChan <- v.NetWorkAddr
					elevatorTimerList = append(elevatorTimerList[:i], elevatorTimerList[i+1:]...) //From slicetricks. Remove element.
				}
			}
		}
		time.Sleep(time.Millisecond*200)
	}
}

func connectToElevator(remoteIp string, remotePort string) *net.UDPConn {
	fullAddr := remoteIp + ":" + remotePort
	remoteUDPAddr, _ := net.ResolveUDPAddr("udp4", fullAddr)
	connection, _ := net.DialUDP("udp4", nil, remoteUDPAddr)

	return connection

}

func getLocalIP() string {
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
				fmt.Println("Printing type: ", reflect.TypeOf(ip))
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
			fmt.Println("BroadcastPrecense: msg not sent")
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

func readMessagesFromNetwork(localAddr string, commonPort string, messageFromNetworkChan  chan []byte) {
	buffer := make([]byte, 4096)
	fullAddr := localAddr + ":" + commonPort
	listenConnAddress, _ := net.ResolveUDPAddr("udp4",fullAddr)
	listenConn, err := net.ListenUDP("udp4",listenConnAddress)
	if err != nil{
		fmt.Println("Error setting up listen-connection")
	}
	for {
		packetLength, senderAddr, err := listenConn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading message, do nothing")
		}else{
			//Send messages to decodeFunc here.
			messageFromNetworkChan<-buffer[:packetLength]
			fmt.Println("Reading message from: ", senderAddr.IP.String())
			fmt.Println("Size of packet: ", packetLength)
		}
	}
}

//
func SendMessagesToAllElevators(sendNetworkMessageChan chan []byte, newConnectionChan chan string, endConnectionChan chan string){
	connectionList := make(map[string] *net.UDPConn)
	for _,deferConn := range connectionList{
		defer deferConn.Close()
	}
	for{
		select{
		case newElevator := <-newConnectionChan:
			connectionList[newElevator] = connectToElevator(newElevator, commonPort)
			defer connectionList[newElevator].Close() //This might be enough
		case deadElevator := <-endConnectionChan:
			connectionList[deadElevator].Close() 
		case newMessage := <-sendNetworkMessageChan:
				for _,conn := range connectionList{
					_, err := conn.Write(newMessage)
					if err != nil {
						fmt.Println("Error in sending func, maybe some error handling later")
				}
			}
		}
	}
}

func sumAllIncomingData(intsIncomingChan chan int) {
	var sum int = 0
	for {
		newInt := <-intsIncomingChan
		sum += newInt
	}
	fmt.Println("Current sum received ", sum)
}

func countBoolsIncoming(boolsIncomingChan chan bool) {
	var falseCounter int = 0
	var trueCounter int = 0
	for {
		newBool := <-boolsIncomingChan
		if newBool {
			trueCounter++
		} else {
			falseCounter++
		}
		fmt.Println("Number of true received:", trueCounter)
		fmt.Println("Number of false received: ", falseCounter)
	}
}

//Not completely tested yet
func encodeMessagesToNetwork(someChannelchan chan int, sendToNetworkChan chan []byte, sendToNetworkCenterChan chan MessageWithHeader){
	for{
		select{
			//The only difference between the cases is the tag that is coming along
		case newVariable := <-someChannelchan:
			encodedData, err := json.Marshal(newVariable)
			if err != nil{
				fmt.Println("error when encoding: ", err)
			}
			var newPacket MessageWithHeader = MessageWithHeader{Data : encodedData, Tag: "someTag", Ack: false, SenderAddr: localAddr}
			encodedPacket, err := json.Marshal(newPacket)
			if err != nil{
				fmt.Println("error when encoding: ", err)
			}
			//Send packet to network. Send copy to local center that keeps track of Ack's
			//Better names needed all over the place
			sendToNetworkChan <- encodedPacket
			sendToNetworkCenterChan <- newPacket
		default:
			fmt.Println("Fuck the haters")
		}
	}
}
func decodeMessagesFromNetwork(messageFromNetworkChan chan []byte ){
	//If no ack. Respond immediately somewhere with ack.
	//Switch on tag
	for{
		packetFromNetwork := <-messageFromNetworkChan
		var message MessageWithHeader
		err:= json.Unmarshal(packetFromNetwork, &message)
		if err != nil {
			fmt.Println("Something went from when unmarshalling packet, do nothing (for now)")
			fmt.Println(err)
		}else if message.Ack == true{
			fmt.Println("Send to goroutine handling acks")
			//Things will happen here
		}else{
			//Should send an ack immediately
			switch message.Tag{
				case "hugeS": //for now, will add many more in the future
					var bigAssMessage HugeStruct
					err := json.Unmarshal(message.Data, &bigAssMessage)
					if err != nil{
						fmt.Println("Something went wrong when unmarshalling in the switch. Do nothing (for now)")
						fmt.Println(err)
					}else{
						printHugeStruct(bigAssMessage)
					}
				case "int":
					var localInt int 
					err := json.Unmarshal(message.Data, &localInt)
					if err != nil{
						fmt.Println("Something went wrong when unmarshalling in the switch. Do nothing (for now)")
					}else{
						fmt.Println("its an int!: ", localInt)
					}

			}
		}
	}
}

func setDeadlinesForAcks(deadLineTriggeredChan chan MessageWithHeader, messageReceivedByAllOtherChan chan MessageWithHeader, newMessageSentChan, newAckChan chan MessageWithHeader){
	/*
	NEED HERE:
	-A structure that contains deadlines for all the elevators for each message
	-To tell someone when the deadline is triggered and package needs to be resent
	-To tell someone when the outgoing message is received by everyone.
	-Set light and everything.
	-Need to have a notion of who has sent these acks. Need to include IP in message. 
	*/
}

func printHugeStruct(bigAssMessage HugeStruct){
	fmt.Println("HugeBool: ", bigAssMessage.HugeBool)
	fmt.Println("HugeName aka localAddr: ", bigAssMessage.HugeName)
	fmt.Println("HugeNumber: ", bigAssMessage.HugeNumber)
	for i := 0; i < arrayX; i++{
		for j := 0; j < arrayY; j++{
			for k := 0; k < arrayZ; k++{
			fmt.Println(bigAssMessage.HugeArray[i][j][k])
			}
		}	
	}
}

func makeHugeStruct(localAddr string) HugeStruct{
	hugeS := HugeStruct{}
	hugeS.HugeBool = true
	hugeS.HugeNumber = 322322
	hugeS.HugeName = localAddr
	for i := 0; i < arrayX; i++{
		for j := 0; j < arrayY; j++{
			for k := 0; k < arrayZ; k++{
			hugeS.HugeArray[i][j][k] = i*100 + j*10 + k
			}
		}	
	}
	return hugeS
}