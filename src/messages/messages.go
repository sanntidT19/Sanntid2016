package messages

import (
	. "../globalStructs"
	"../network"
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

/*
Functions dealing with the informational messages being sent between the elevators
*/

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

const ACK_DEADLINE = 500

var commonPort string = "20059"
var localAddr string

func InitAndRun(fromDecode MessageChans,toEncode MessageChans, newConnectionChan chan string, endConnectionChan chan string, resendUnAckdMessagesChan chan bool) {
	newEncodedMessageChan := make(chan []byte)
	messageFromNetworkChan := make(chan []byte)
	newAckReceivedChan := make(chan MessageWithHeader)
	resendLostMessageChan := make(chan MessageWithHeader)
	newMessageDeadlineChan := make(chan MessageWithHeader)
	localAddr = network.FindLocalIP()

	go sendMessagesOverNetwork(newEncodedMessageChan, newConnectionChan, endConnectionChan)
	go receiveMessagesFromNetwork(localAddr, commonPort, messageFromNetworkChan)
	go encodeMessages(newEncodedMessageChan, newMessageDeadlineChan, resendLostMessageChan, toEncode)
	go decodeMessages(messageFromNetworkChan, newAckReceivedChan, newEncodedMessageChan, fromDecode) 
	go setDeadlinesForAcks(resendLostMessageChan, newMessageDeadlineChan, newAckReceivedChan, resendUnAckdMessagesChan)
}

func sendMessagesOverNetwork(sendNetworkMessageChan chan []byte, newConnectionChan chan string, endConnectionChan chan string) {
	connectionList := make(map[string]*net.UDPConn)
	for {
		select {
		case newElevator := <-newConnectionChan:
			connectionList[newElevator] = network.ConnectToElevator(newElevator, commonPort)
			defer connectionList[newElevator].Close() 
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
			byteSliceToDecoder := make([]byte, packetLength)
			copy(byteSliceToDecoder, buffer[:packetLength])
			messageFromNetworkChan <- byteSliceToDecoder
		}
	}
}



func encodeMessages(sendToNetworkChan chan []byte, sendToAckTimerChan chan MessageWithHeader, resendMessageChan chan MessageWithHeader, encode MessageChans) {
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
			sendToNetworkChan <- encodedMessage
		}
	}()
	for {
		var tag string = ""
		var encodedData []byte = nil
		var err error
		select {
		case newOrder := <-encode.NewOrderChan:
			tag = "newOr"
			encodedData, err = json.Marshal(newOrder)
			if err != nil {
				fmt.Println("Error when encoding: ", err)
			}
		case orderTo := <-encode.OrderAssChan:
			tag = "ordTo"
			encodedData, err = json.Marshal(orderTo)
			if err != nil {
				fmt.Println("Error when encoding: ", err)
			}
		case orderServed := <-encode.OrderServedChan:
			tag = "ordSe"
			encodedData, err = json.Marshal(orderServed)
			if err != nil {
				fmt.Println("Error when encoding: ", err)
			}
		case elevState := <-encode.ElevStateChan:
			tag = "elSta"
			encodedData, err = json.Marshal(elevState)
			if err != nil {
				fmt.Println("Error when encoding: ", err)
			}
		case externalArray := <-encode.ExternalArrayChan:
			tag = "extAr"
			encodedData, err = json.Marshal(externalArray)
			if err != nil {
				fmt.Println("Error when encoding: ", err)
			}
		}
		var newPacket MessageWithHeader = MessageWithHeader{Data: encodedData, Tag: tag, Ack: false, SenderAddr: localAddr}
		encodedMessage, err := json.Marshal(newPacket)
		if err != nil {
			fmt.Println("error when encoding: ", err)
		}
		sendToNetworkChan <- encodedMessage
		sendToAckTimerChan <- newPacket
	}
}


func sendAck(message MessageWithHeader, sendToNetworkChan chan []byte) {
	ackMessage := MessageWithHeader{}
	ackMessage.Ack = true
	ackMessage.Tag = message.Tag
	ackMessage.SenderAddr = localAddr
	ackMessage.Data = make([]byte, len(message.Data))
	copy(ackMessage.Data, message.Data)
	encodedMessage, err := json.Marshal(ackMessage)
	if err != nil {
		fmt.Println("Error encoding ack")
	}
	sendToNetworkChan <- encodedMessage
}

func decodeMessages(messageFromNetworkChan chan []byte, newAckFromNetworkChan chan MessageWithHeader, sendToNetworkChan chan []byte, decode MessageChans) {

	for {
		packetFromNetwork := <-messageFromNetworkChan
		var message MessageWithHeader = MessageWithHeader{}
		err := json.Unmarshal(packetFromNetwork, &message)
		if err != nil {
			fmt.Println(err)
		} else if message.Ack == true {
			newAckFromNetworkChan <- message
		} else {
			switch message.Tag {
			case "newOr":
				var newOrder Order
				err := json.Unmarshal(message.Data, &newOrder)
				if err != nil {
					fmt.Println(err)
				} else {
					sendAck(message, sendToNetworkChan)
					decode.NewOrderChan <- newOrder
				}
			case "ordSe":
				var orderServed Order
				err := json.Unmarshal(message.Data, &orderServed)
				if err != nil {
					fmt.Println(err)
				} else {
					sendAck(message, sendToNetworkChan)
					decode.OrderServedChan <- orderServed
				}
			case "ordTo":
				var orderAss OrderAssigned
				err := json.Unmarshal(message.Data, &orderAss)
				if err != nil {
					fmt.Println(err)
				} else {
					sendAck(message, sendToNetworkChan)
					decode.OrderAssChan <- orderAss
				}
			case "elSta":
				var newState ElevatorState
				err := json.Unmarshal(message.Data, &newState)
				if err != nil {
					fmt.Println(err)
				} else {
					sendAck(message, sendToNetworkChan)
					decode.ElevStateChan <- newState
				}
			case "extAr":
				var newExternalArray [NUM_FLOORS][NUM_BUTTONS - 1]int
				err := json.Unmarshal(message.Data, &newExternalArray)
				if err != nil {
					fmt.Println(err)
				} else {
					sendAck(message, sendToNetworkChan)
					decode.ExternalArrayChan <- newExternalArray
				}
			}
		}
	}
}

//
func setDeadlinesForAcks(resendMessageChan chan MessageWithHeader, newMessageSentChan chan MessageWithHeader, newAckChan chan MessageWithHeader, resendUnAckdMessagesChan chan bool) {
	var unAckdMessages []ackTimer
	var resendUnackd bool
	for {
		select {
		case resendUnackd = <-resendUnAckdMessagesChan:
			if resendUnackd{
				for i, v := range unAckdMessages {
					unAckdMessages[i].DeadLine = time.Now().Add(time.Millisecond * ACK_DEADLINE)
					localIPlist := network.ElevsSeen()
					unAckdMessages[i].IpList = localIPlist
					resendMessageChan <- v.Message 
				}
			}
		case newAck := <-newAckChan:
			senderOfAck := newAck.SenderAddr
		Loop:
			for i, v := range unAckdMessages {
				if bytes.Equal(v.Message.Data, newAck.Data) {
					for j, addr := range unAckdMessages[i].IpList {
						if addr == senderOfAck {
							unAckdMessages[i].IpList = append(unAckdMessages[i].IpList[:j], unAckdMessages[i].IpList[j+1:]...)
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
				newAckTimer := ackTimer{Message: newMessage, DeadLine: time.Now().Add(time.Millisecond * ACK_DEADLINE), IpList: localIPlist}
				unAckdMessages = append(unAckdMessages, newAckTimer)
			}
		default:
			if !resendUnackd {
				unAckdMessages = nil
			}
			for i, v := range unAckdMessages {
				if time.Now().After(v.DeadLine) {
					resendMessageChan <- v.Message
					unAckdMessages[i].DeadLine = time.Now().Add(time.Millisecond * ACK_DEADLINE)
					localIPlist := network.ElevsSeen()
					unAckdMessages[i].IpList = localIPlist
				}
			}
		}
		time.Sleep(time.Millisecond * 100)
	}
}

