package network

import (
	. "encoding/json"
	"fmt"
)

const (
	FLOORS = 4
)

func main() {
	fmt.Println("hei")
}

func Channels_init() {

	slaveToCommSlaveChan := make(chan Slave)                   //"sla"
	slaveToCommOrd_eceivedChan := make(chan []int)             //"ore"
	slaveToCommOrderExecutedChan := make(chan []int)           //"oex"
	slaveToCommOrderConfirmedReceivedChan := make(chan []int)  //"ocr"
	slaveToCommOrderConfirmedExecutuinChan := make(chan []int) //"oce"

	//Master
	masterToCommOrderListChan := make(chan [][]int)          //"exo"
	masterToCommImMasterChan := make(chan string)            //"iam"
	masterToCommReceivedConfirmationChan := make(chan []int) //"rco"
	masterToCommExecutedConfirmationChan := make(chan []int) //"eco"

	//communication channels
	//commToMasterSlaveChan := make(chan Slave)                   //"sla"
	commToMasterOrd_eceivedChan := make(chan []int)             //"ore"
	commToMasterOrderExecutedChan := make(chan []int)           //"oex"
	commToMasterOrderConfirmedReceivedChan := make(chan []int)  //"ocr"
	commToMasterOrderConfirmedExecutionChan := make(chan []int) //"oce"

	commToSlaveOrderListChan := make(chan [][]int)          //"exo"
	commToSlaveImMasterChan := make(chan string)            //"iam"
	commToSlaveReceivedConfirmationChan := make(chan []int) //"rco"
	commToSlaveExecutedConfirmationChan := make(chan []int) //"eco"

	newExternalList := make(chan [][]int)
	slaveToStateMChan := make(chan int) //send input to statemachine
	//network
	commToNetwork := make(chan []byte)
	networkToComm := make(chan []byte)
}

//Master
func Send_order(externalOrderList [][]int, commToNetwork chan []byte) { //send exectuionOrderList
	byteOrder, _ := Marshal(externalOrderList)
	prefix, _ := Marshal("exo")
	commToNetwork <- append(prefix, byteOrder...)
}

func Send_im_master(message string, commToNetwork chan []byte) { //send I am master
	byteMessage, _ := Marshal(message)
	prefix, _ := Marshal("iam")
	fmt.Println("to network", string(byteMessage))
	commToNetwork <- append(prefix, byteMessage...)

}
func Send_received_confirmation(order []int, commToNetwork chan []byte) {
	byteMessage, _ := Marshal(order)
	prefix, _ := Marshal("rco")
	commToNetwork <- append(prefix, byteMessage...)
}

func Send_executed_confirmation(order []int, commToNetwork chan []byte) {
	byteMessage, _ := Marshal(order)
	prefix, _ := Marshal("eco")
	commToNetwork <- append(prefix, byteMessage...)
}

/*
//Slave
func Send_slave(s Slave, commToNetwork chan []byte) {
	byteSlave, _ := Marshal(s)
	prefix, _ := Marshal("sla")
	commToNetwork <- append(prefix, byteSlave...)
}
*/
func Send_order_received(order []int, commToNetwork chan []byte) {
	byteMessage, _ := Marshal(order)
	prefix, _ := Marshal("ore")
	commToNetwork <- append(prefix, byteMessage...)
}

func Send_order_executed(order []int, commToNetwork chan []byte) {
	byteMessage, _ := Marshal(order)
	prefix, _ := Marshal("oex")
	commToNetwork <- append(prefix, byteMessage...)
}

func Send_order_confirmed_received(order []int, commToNetwork chan []byte) {
	byteMessage, _ := Marshal(order)
	prefix, _ := Marshal("ocr")
	commToNetwork <- append(prefix, byteMessage...)
}

func Send_order_confirmed_executed(order []int, commToNetwork chan []byte) {
	byteMessage, _ := Marshal(order)
	prefix, _ := Marshal("oce")
	commToNetwork <- append(prefix, byteMessage...)
}

func Decrypt_message(message []byte, commToMasterOrderReceivedChan chan []int, commToMasterOrderExecutedChan chan []int, commToMasterOrderConfirmedReceivedChan chan []int, commToMasterOrderConfirmedExecutionChan chan []int, commToSlaveOrderListChan chan [][]int, commToSlaveImMasterChan chan string, commToSlaveReceivedConfirmationChan chan []int, commToSlaveExecutedConfirmationChan chan []int) {

	switch {
	//Master
	/*case string(message[1:4] == "sla"):
	noPrefix := message[5:]
	var s Slave
	_ := Unmarshal(noPrefix, &s)
	commToMasterSlaveChan <- s
	*/
	case string(message[1:4]) == "ore":
		noPrefix := message[5:]
		order := make([]int, 2)
		_ = Unmarshal(noPrefix, &order)
		commToMasterOrderReceivedChan <- order

	case string(message[1:4]) == "oex":
		noPrefix := message[5:]
		order := make([]int, 2)
		_ = Unmarshal(noPrefix, &order)
		commToMasterOrderExecutedChan <- order

	case string(message[1:4]) == "ocr":
		noPrefix := message[5:]
		order := make([]int, 2)
		_ = Unmarshal(noPrefix, &order)
		commToMasterOrderConfirmedReceivedChan <- order

	case string(message[1:4]) == "oce":
		noPrefix := message[5:]
		order := make([]int, 2)
		_ = Unmarshal(noPrefix, &order)
		commToMasterOrderConfirmedExecutionChan <- order

	//Slave
	case string(message[1:4]) == "exo":
		noPrefix := message[5:]
		externalOrderList := make([][]int, FLOORS)
		_ = Unmarshal(noPrefix, &externalOrderList)
		commToSlaveOrderListChan <- externalOrderList

	case string(message[1:4]) == "iam":
		fmt.Println("iam trigger")
		noPrefix := message[5:]
		stringMessage := string(noPrefix)
		fmt.Println(stringMessage)
		commToSlaveImMasterChan <- stringMessage
		fmt.Println("channel output")

	case string(message[1:4]) == "rco":
		noPrefix := message[5:]
		order := make([]int, 2)
		_ = Unmarshal(noPrefix, &order)
		commToSlaveReceivedConfirmationChan <- order

	case string(message[1:4]) == "eco":
		noPrefix := message[5:]
		order := make([]int, 2)
		_ = Unmarshal(noPrefix, &order)
		commToSlaveExecutedConfirmationChan <- order
	}
}
func Select_send(commToNetwork chan []byte, slaveToCommOrderReceivedChan chan []int, slaveToCommOrderExecutedChan chan []int, slaveToCommOrderConfirmedReceivedChan chan []int, slaveToCommOrderConfirmedExecutuinChan chan []int, masterToCommOrderListChan chan [][]int, masterToCommImMasterChan chan string, masterToCommReceivedConfirmationChan chan []int, masterToCommExecutedConfirmationChan chan []int) {

	for {
		select {
		//Master
		case externalOrderList := <-masterToCommOrderListChan:
			Send_order(externalOrderList, commToNetwork)
		case message := <-masterToCommImMasterChan:
			fmt.Println("meessage", message, "select send")
			Send_im_master(message, commToNetwork)
		case order := <-masterToCommReceivedConfirmationChan:
			Send_received_confirmation(order, commToNetwork)
		case order := <-masterToCommExecutedConfirmationChan:
			Send_executed_confirmation(order, commToNetwork)
		//Slave
		/*
			case slave := <-slaveToCommSlaveChan:
				Send_slave(slave, commToNetwork)
		*/
		case order := <-slaveToCommOrderReceivedChan:
			Send_order_received(order, commToNetwork)
		case order := <-slaveToCommOrderConfirmedReceivedChan:
			Send_order_executed(order, commToNetwork)
		case order := <-slaveToCommOrderConfirmedReceivedChan:
			Send_order_confirmed_received(order, commToNetwork)
		case order := <-slaveToCommOrderConfirmedExecutuinChan:
			Send_order_confirmed_executed(order, commToNetwork)
		}
	}
}
func Select_receive(networkToComm chan []byte /*commToMasterSlaveChan chan Slave,*/, commToMasterOrderReceivedChan chan []int, commToMasterOrderExecutedChan chan []int, commToMasterOrderConfirmedReceivedChan chan []int, commToMasterOrderConfirmedExecutionChan chan []int, commToSlaveOrderListChan chan [][]int, commToSlaveImMasterChan chan string, commToSlaveReceivedConfirmationChan chan []int, commToSlaveExecutedConfirmationChan chan []int) {
	var barr []byte
	fmt.Println("Select_receive")
	for {
		barr = <-networkToComm
		Decrypt_message(barr, commToMasterOrderReceivedChan, commToMasterOrderExecutedChan, commToMasterOrderConfirmedReceivedChan, commToMasterOrderConfirmedExecutionChan, commToSlaveOrderListChan, commToSlaveImMasterChan, commToSlaveReceivedConfirmationChan, commToSlaveExecutedConfirmationChan)
	}
}
