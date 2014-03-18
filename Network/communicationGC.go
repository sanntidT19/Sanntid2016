package network

import (
	. "encoding/json"
	"fmt"
	//. "time"
)

const (
	FLOORS = 4
)

var ExComChans ExternalCommunicationChannels
var InComChans InternalCommunicationChannels

type ExternalCommunicationChannels struct {

	//communication channels
	ToMasterSlaveChan                   chan Slave //"sla"
	ToMasterOrderReceivedChan           chan []int //"ore"
	ToMasterOrderExecutedChan           chan []int //"oex"
	ToMasterOrderConfirmedReceivedChan  chan []int //"ocr"
	ToMasterOrderConfirmedExecutionChan chan []int //"oce"

	ToSlaveOrderListChan            chan [][]int //"exo"
	ToSlaveImMasterChan             chan string  //"iam"
	ToSlaveReceivedConfirmationChan chan []int   //"rco"
	ToSlaveExecutedConfirmationChan chan []int   //"eco"

}
type InternalCommunicationChannels struct {
	newExternalList   chan [][]int
	slaveToStateMChan chan int //send input to statemachine
}

func external_comm_channels_init() {
	ExComChans.ToMasterSlaveChan = make(chan Slave)                   //"sla"
	ExComChans.ToMasterOrderReceivedChan = make(chan []int)           //"ore"
	ExComChans.ToMasterOrderExecutedChan = make(chan []int)           //"oex"
	ExComChans.ToMasterOrderConfirmedReceivedChan = make(chan []int)  //"ocr"
	ExComChans.ToMasterOrderConfirmedExecutionChan = make(chan []int) //"oce"

	ExComChans.ToSlaveOrderListChan = make(chan [][]int)          //"exo"
	ExComChans.ToSlaveImMasterChan = make(chan string)            //"iam"
	ExComChans.ToSlaveReceivedConfirmationChan = make(chan []int) //"rco"
	ExComChans.ToSlaveExecutedConfirmationChan = make(chan []int) //"eco"

}
func internal_comm_chans_init() {
	InComChans.newExternalList = make(chan [][]int)
	InComChans.slaveToStateMChan = make(chan int) //send input to statemachine
	//network

}

//Master
func Send_order(externalOrderList [][]int) { //send exectuionOrderList
	byteOrder, _ := Marshal(externalOrderList)
	prefix, _ := Marshal("exo")
	ExNetChans.ToNetwork <- append(prefix, byteOrder...)
}

func Send_im_master(message string) { //send I am master
	byteMessage, _ := Marshal(message)
	prefix, _ := Marshal("iam")
	fmt.Println("to network", string(byteMessage))
	ExNetChans.ToNetwork <- append(prefix, byteMessage...)

}
func Send_received_confirmation(order []int) {
	byteMessage, _ := Marshal(order)
	prefix, _ := Marshal("rco")
	ExNetChans.ToNetwork <- append(prefix, byteMessage...)
}

func Send_executed_confirmation(order []int) {
	byteMessage, _ := Marshal(order)
	prefix, _ := Marshal("eco")
	ExNetChans.ToNetwork <- append(prefix, byteMessage...)
}

//Slave
func Send_slave(s Slave, commToNetwork chan []byte) {
	byteSlave, _ := Marshal(s)
	prefix, _ := Marshal("sla")
	ExNetChans.ToNetwork <- append(prefix, byteSlave...)
}

func Send_order_received(order []int) {
	byteMessage, _ := Marshal(order)
	prefix, _ := Marshal("ore")
	ExNetChans.ToNetwork <- append(prefix, byteMessage...)
}

func Send_order_executed(order []int) {
	byteMessage, _ := Marshal(order)
	prefix, _ := Marshal("oex")
	ExNetChans.ToNetwork <- append(prefix, byteMessage...)
}

func Send_order_confirmed_received(order []int) {
	byteMessage, _ := Marshal(order)
	prefix, _ := Marshal("ocr")
	ExNetChans.ToNetwork <- append(prefix, byteMessage...)
}

func Send_order_confirmed_executed(order []int) {
	byteMessage, _ := Marshal(order)
	prefix, _ := Marshal("oce")
	ExNetChans.ToNetwork <- append(prefix, byteMessage...)
}

func Decrypt_message(message []byte) {

	switch {
	//Master
	case string(message[1:4]) == "sla":
		noPrefix := message[5:]
		var s Slave
		_ = Unmarshal(noPrefix, &s)
		ExComChans.ToMasterSlaveChan <- s

	case string(message[1:4]) == "ore":
		noPrefix := message[5:]
		order := make([]int, 2)
		_ = Unmarshal(noPrefix, &order)
		ExComChans.ToMasterOrderReceivedChan <- order

	case string(message[1:4]) == "oex":
		noPrefix := message[5:]
		order := make([]int, 2)
		_ = Unmarshal(noPrefix, &order)
		ExComChans.ToMasterOrderExecutedChan <- order

	case string(message[1:4]) == "ocr":
		noPrefix := message[5:]
		order := make([]int, 2)
		_ = Unmarshal(noPrefix, &order)
		ExComChans.ToMasterOrderConfirmedReceivedChan <- order

	case string(message[1:4]) == "oce":
		noPrefix := message[5:]
		order := make([]int, 2)
		_ = Unmarshal(noPrefix, &order)
		ExComChans.ToMasterOrderConfirmedExecutionChan <- order

	//Slave
	case string(message[1:4]) == "exo":
		noPrefix := message[5:]
		externalOrderList := make([][]int, FLOORS)
		_ = Unmarshal(noPrefix, &externalOrderList)
		ExComChans.ToSlaveOrderListChan <- externalOrderList

	case string(message[1:4]) == "iam":
		fmt.Println("iam trigger")
		noPrefix := message[5:]
		stringMessage := string(noPrefix)
		fmt.Println(stringMessage)
		ExComChans.ToSlaveImMasterChan <- stringMessage
		fmt.Println("channel output")

	case string(message[1:4]) == "rco":
		noPrefix := message[5:]
		order := make([]int, 2)
		_ = Unmarshal(noPrefix, &order)
		ExComChans.ToSlaveReceivedConfirmationChan <- order

	case string(message[1:4]) == "eco":
		noPrefix := message[5:]
		order := make([]int, 2)
		_ = Unmarshal(noPrefix, &order)
		ExComChans.ToSlaveExecutedConfirmationChan <- order
	}
}
func Select_send() {

	for {
		select {
		//Master
		case externalOrderList := <-MC.ToCommOrderListChan:
			Send_order(externalOrderList)
		case message := <-MC.ToCommImMasterChan:
			fmt.Println("meessage", message, "select send")
			Send_im_master(message)
		case order := <-MC.ToCommReceivedConfirmationChan:
			Send_received_confirmation(order)
		case order := <-MC.ToCommExecutedConfirmationChan:
			Send_executed_confirmation(order)
		//Slave
		case slave := <-SC.ToCommSlaveChan:
			Send_slave(slave, commToNetwork)
		case order := <-SC.ToCommOrderReceivedChan:
			Send_order_received(order)
		case order := <-SC.ToCommOrderConfirmedReceivedChan:
			Send_order_executed(order)
		case order := <-SC.ToCommOrderConfirmedReceivedChan:
			Send_order_confirmed_received(order)
		case order := <-SC.ToCommOrderConfirmedExecutuinChan:
			Send_order_confirmed_executed(order)
		}
	}
}
func Select_receive() {
	var barr []byte
	fmt.Println("Select_receive")
	for {
		barr = <-ExNetChans.ToComm
		Decrypt_message(barr)
	}
}
