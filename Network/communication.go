package main

import (
	. "encoding/json"
	"fmt"
	. "strings"
)

func main() {
	fmt.Println("hei")
}
func Channels_init() {
	//Slave
	slaveToCommFloorRChan := make(chan bool)          //send floor reached to master
	slaveToCommSlaveStructChan := make(chan Slave)    // send slave struct to master
	slaveToCommOrderReceivedChan := make(chan string) //notity master that slave has received order
	slaveToCommOrderExecutedChan := make(chan string) //notify master that slave has executed order
	slaveToCommConfirmedExecutuinChan := make(chan string)

	slaveToStateMChan := make(chan int) //send input to statemachine

	//Master
	masterToCommOrderChan := make(chan [][]int)               //sends orders from slave to comm
	masterToCommConfirmChan := make(chan bool)                //confirms that master har received that slave has confirmed/Received order
	masterToCommImMasterChan := make(chan string)             // sends i am master
	masterToCommReceivedConfirmationChan := make(chan []int) // master confirms that slave has received order
	masterToCommExecutedConfirmationChan := make(chan []int)
	//communication channels
	commToSlaveOrderChan := make(chan [][]int)             //receves orders from comm
	commToSlaveMastersBackConfirmChan := make(chan string) //master confirms that order is received
	commToSlaveImMasterChan := make(chan string)           //im master from master
	commToSlaveReceivedConfirmationChan := make(chan string)
	commToSlaveExecutedConfirmationChan := make(chan string)

	commToMasterFloorRChan := make(chan bool)               //floor reached from slave to master
	commToMasterSlaveStructChan := make(chan Slave)         //sends slave struct
	commToMasterOrderExecuredChan := make(chan []int)       //order executed sucessfully - sends array of ints: dir og floor
	commToMasterOrderReceivedChan := make(chan []int)      //Slave confirmes that order is recived
	commToMasterConfirmedExecutionChan := make(chan string) //slave confirmes order executed

	newExternalList := make(chan [][]int)

	//network
	commToNetwork := make(chan []byte)
	networkToComm := make(chan []byte)
}

//Master
func Send_order(masterToCommOrderChan chan [][]int, commToNetwork chan []byte) {
	byteOrder, err := Marshal(<-masterToCommOrderChan)
	prefix, err := Marshal("ord")
	commToNetwork <- append(prefix, byteOrder)
}

func Send_im_master(masterToCommImMasterChan chan string, commToNetwork chan []byte) {
	byteMessage, err := Marshal(<-masterToCommImMasterChan)
	prefix, err := Marshal("iam")
	commToNetwork <- append(prefix, byteMessage)
}

func Send_order_received(slaveToCommOrderRecivedChan chan []int, commToNetwork chan []byte) {
	byteMessage, err := Marshal(<-slaveToCommOrderRecivedChan)
	prefix, err := Marshal("mre")
	commToNetwork <- append(prefix, byteMessage)
}

func Send_order_executed(slaveToCommOrderExecutedChan chan []int, commToNetwork chan []byte) {
	byteMessage, err := Marshal(<-slaveToCommOrderExecuredChan)
	prefix, err := Marshal("exe")
	commToNetwork <- append(prefix, byteMessage)
}

func Send_received_confirmation(masterToCommReceivedConfirmationChan chan []int, commToNetwork chan []byte) {
	byteMessage, err := Marshal(<-masterToCommReceivedConfirmationChan)
	prefix, err := Marshal("sre")
	commToNetwork <- append(prefix, byteMessage)

}
func Send_executed_confirmation(masterToCommExecutedConfirmationChan chan []int, commToNetwork chan []byte) {
	byteMessage, err := Marshal(<-masterToCommExecutedConfirmationChan)
	prefix, err := Marshal("see")
	commToNetwork <- append(prefix, byteMessage)
}
func Send_slave(slaveToCommSlaveStructChan chan Slave, commToNetwork chan []byte) {
	byteSlave, err := Marshal(<-slaveToCommSlaveStructChan)
	prefix, err := Marshal("sch")
	commToNetwork <- append(prefix, byteSlave)
}

func Decrypt_message(message []byte) {
	switch {

	case string(message[1:4]) == "ord": //externalorderlist
		str := string(message)
		str = TrimPrefix(str, "ord")
		externalOrderList := Unmarshal([]byte(message), [][]int)
		commToSlaveOrderChan <- externalOrderList

	case string(message[1:4]) == "iam": //I am master
		noPrefix = message[5:]
		commToSlaveImMasterChan <- noPrefix

	case string(message[1:4]) == "mre": //confirm recived order from master
		noPrefix = message[5:]
		arr := Unmarshal(noPrefix, []int)
		commToSlaveOrderReceivedChan <- noPrefix

	case string(message[1:4]) == "exe": //Order performed from slave to master
		noPrefix = message[5:]
		arr := Unmarshal(noPrefix, []int)
		commToMasterOrderExecuredChan <- noPrefix

	case string(message[1:4]) ==  "sre": //confirmes recived order from slave to master
		noPrefix = message[5:]
		arr := Unmarshal(noPrefix, []int)
		commToSlaveReceivedConfirmationChan <- noPrefix

	case string(message[1:4]) ==  "see": //confirmes recived order from slave to master
		noPrefix = message[5:]
		arr := Unmarshal(noPrefix, []int)
		commToSlaveExecutedConfirmationChan <- noPrefix

	case string(message[1:4] == "sch"): // receves a slave struct
		noPrefix = message[5:]
		slave := UnMarshal(noPrefix, &Slave)
		commToMasterSlaveStructChan <- slave
	}
}