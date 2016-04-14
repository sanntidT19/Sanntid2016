package messages

import (
	. "../globalChans"
	. "../globalStructs"
	"../network"
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type MessageWithHeader struct {
	Tag        string
	Ack        bool
	SenderAddr string
	Data       []byte
}

type ackTimer struct {
	Message  MessageWithHeader
	DeadLine time.Time
	IpList   []string
}

const ACK_DEADLINE = 2

var commonPort string = "20059"
var localAddr string

func MessagesTopAndWaitForNetworkChanges() {
	newEncodedMessageChan := make(chan []byte)
	newConnectionChan := make(chan string)
	endConnectionChan := make(chan string)
	messageFromNetworkChan := make(chan []byte)
	newAckReceivedChan := make(chan MessageWithHeader)
	resendLostMessageChan := make(chan MessageWithHeader)
	newMessageDeadlineChan := make(chan MessageWithHeader)
	resendUnAckdMessagesChan := make(chan bool)
	networkDownResetAckListChan := make(chan bool)

	localAddr = network.FindLocalIP()

	go sendMessagesOverNetwork(newEncodedMessageChan, newConnectionChan, endConnectionChan)
	go receiveMessagesFromNetwork(localAddr, commonPort, messageFromNetworkChan)
	go encodeMessages(newEncodedMessageChan, newMessageDeadlineChan, resendLostMessageChan)
	go decodeMessages(messageFromNetworkChan, newAckReceivedChan, newEncodedMessageChan) //FÅ ACKEN UT HERIFRA
	go setDeadlinesForAcks(resendLostMessageChan, newMessageDeadlineChan, newAckReceivedChan, resendUnAckdMessagesChan, networkDownResetAckListChan)

	for {
		select {
		case elevAddr := <-ToMessagesDeadElevChan:
			resendUnAckdMessagesChan <- true
			endConnectionChan <- elevAddr
		case elevAddr := <-ToMessagesNewElevChan:
			resendUnAckdMessagesChan <- true
			newConnectionChan <- elevAddr
		case <-ToMessagesNetworkDownChan:
			networkDownResetAckListChan <- true
		}

	}
}

func sendMessagesOverNetwork(sendNetworkMessageChan chan []byte, newConnectionChan chan string, endConnectionChan chan string) {
	connectionList := make(map[string]*net.UDPConn)
	for {
		select {
		case newElevator := <-newConnectionChan:
			connectionList[newElevator] = network.ConnectToElevator(newElevator, commonPort)
			defer connectionList[newElevator].Close() //This might be enough
		case deadElevator := <-endConnectionChan:
			connectionList[deadElevator].Close()
			delete(connectionList, deadElevator)
		case newMessage := <-sendNetworkMessageChan:
			for _, conn := range connectionList {
				_, err := conn.Write(newMessage)
				if err != nil {
					fmt.Println("Error sending message to network: ", err)
				}
			}
		}
	}
}

func receiveMessagesFromNetwork(localAddr string, commonPort string, messageFromNetworkChan chan []byte) {
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

func encodeMessages(sendToNetworkChan chan []byte, sendToAckTimerChan chan MessageWithHeader, resendMessageChan chan MessageWithHeader) {
	go func() {
		for {
			expiredMessage := <-resendMessageChan
			expiredMessage.Ack = false
			copyOfData := make([]byte, len(expiredMessage.Data))
			copy(copyOfData, expiredMessage.Data)
			expiredMessage.Data = copyOfData
			expiredMessage.SenderAddr = localAddr
			encodedMessage, err := json.Marshal(expiredMessage)
			if err != nil {
				fmt.Println("error when encoding: ", err)
			}
			fmt.Println("Ack expired, resend: ", expiredMessage.Tag)
			sendToNetworkChan <- encodedMessage
		}

	}()
	for {
		var tag string = ""
		var encodedData []byte = nil
		var err error
		select {
		case newOrder := <-ExternalButtonPressedChan:
			tag = "newOr"
			encodedData, err = json.Marshal(newOrder)
			//fmt.Println("encodeMessagesToNetwork: order encoded1")
			if err != nil {
				fmt.Println("Error when encoding: ", err)
			}
		case orderTo := <-ToNetworkOrderAssignedToChan:
			tag = "ordTo"
			encodedData, err = json.Marshal(orderTo)
			//fmt.Println("encodeMessagesToNetwork: order encoded2")
			if err != nil {
				fmt.Println("Error when encoding: ", err)
			}
		case orderServed := <-ToNetworkOrderServedChan:
			tag = "ordSe"
			encodedData, err = json.Marshal(orderServed)
			//fmt.Println("encodeMessagesToNetwork: order encoded3")
			if err != nil {
				fmt.Println("Error when encoding: ", err)
			}
		case elevState := <-ToNetworkNewElevStateChan:
			tag = "elSta"
			encodedData, err = json.Marshal(elevState)
			//fmt.Println("encodeMessagesToNetwork: order encoded4")
			if err != nil {
				fmt.Println("Error when encoding: ", err)
			}
		case externalArray := <-ToNetworkExternalArrayChan:
			tag = "extAr"
			encodedData, err = json.Marshal(externalArray)
			//fmt.Println("encodeMessagesToNetwork: order encoded5")
			if err != nil {
				fmt.Println("Error when encoding: ", err)
			}
		}
		//alt felles gjøres her
		var newPacket MessageWithHeader = MessageWithHeader{Data: encodedData, Tag: tag, Ack: false, SenderAddr: localAddr}
		encodedMessage, err := json.Marshal(newPacket)
		if err != nil {
			fmt.Println("error when encoding: ", err)
		}
		//Send packet to network. Send copy to local center that keeps track of Ack's
		//Better names needed all over the place
		sendToNetworkChan <- encodedMessage
		sendToAckTimerChan <- newPacket
	}
}

//Takes a message and sends out an ack
func sendAck(message MessageWithHeader, sendToNetworkChan chan []byte) {
	ackMessage := MessageWithHeader{}
	ackMessage.Ack = true
	ackMessage.Tag = message.Tag
	ackMessage.SenderAddr = localAddr
	ackMessage.Data = make([]byte, len(message.Data))
	copy(ackMessage.Data, message.Data)
	encodedMessage, err := json.Marshal(ackMessage)
	if err != nil {
		fmt.Println("error encoding resend message")
	}
	//fmt.Println("Ack sent")
	sendToNetworkChan <- encodedMessage
}

func decodeMessages(messageFromNetworkChan chan []byte, newAckFromNetworkChan chan MessageWithHeader, sendToNetworkChan chan []byte) {
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
			//fmt.Println("Tag before check: ", message.Tag)
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
				var newExternalArray [NUM_FLOORS][NUM_BUTTONS - 1]int
				err := json.Unmarshal(message.Data, &newExternalArray)
				if err != nil {
					fmt.Println(err)
				} else {
					FromNetworkExternalArrayChan <- newExternalArray
				}
			}
		}
		//fmt.Println("Message ackd and sent to local: ", message.Tag)
	}
}

/*
there is no need for ackdbyallchan


*/
func setDeadlinesForAcks(resendMessageChan chan MessageWithHeader, newMessageSentChan chan MessageWithHeader, newAckChan chan MessageWithHeader, elevatorListChangedChan chan bool, networkDownResetAckListChan chan bool) {
	var unAckdMessages []ackTimer
	var networkIsDown bool
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
			networkIsDown = false
			for i, v := range unAckdMessages {
				unAckdMessages[i].DeadLine = time.Now().Add(time.Second * ACK_DEADLINE)
				localIPlist := network.ElevsSeen()
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
				localIPlist := network.ElevsSeen()
				newAckTimer := ackTimer{Message: newMessage, DeadLine: time.Now().Add(time.Second * ACK_DEADLINE), IpList: localIPlist}
				unAckdMessages = append(unAckdMessages, newAckTimer)
			}
		case <-networkDownResetAckListChan:
			networkIsDown = true
		default:
			if networkIsDown {
				unAckdMessages = nil
			}
			for i, v := range unAckdMessages {
				if time.Now().After(v.DeadLine) {
					resendMessageChan <- v.Message
					unAckdMessages[i].DeadLine = time.Now().Add(time.Second * ACK_DEADLINE)
					localIPlist := network.ElevsSeen()
					fmt.Println("localiplist: ", localIPlist)
					unAckdMessages[i].IpList = localIPlist
					fmt.Println("on element number ", i)
				}
			}
		}
		time.Sleep(time.Millisecond * 100)
	}
}

//Kanskje dele opp i en til som også er top loop av et slag??? En som i sin for select tar imot når endringer i nettverk skjer. Må si ifra til alle sine subrutiner.
