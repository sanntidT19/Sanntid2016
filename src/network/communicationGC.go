package network

import (
	. "encoding/json"
	"fmt"
	//. "time"
	."toplayer"
)

const (
	FLOORS = 4
)

var ExCommChans ExternalCommunicationChannels
var InCommChans InternalCommunicationChannels

type ExternalCommunicationChannels struct {

	//communication channels
	ToMasterSlaveChan                   chan Slave //"sla"
	ToMasterOrderReceivedChan           chan []int //"ore"
	ToMasterOrderExecutedChan           chan []int //"oex"
	ToMasterOrderConfirmedReceivedChan  chan []int //"ocr"
	ToMasterOrderConfirmedExecutionChan chan []int //"oce"
	ToMasterExternalButtonPushed 		chan []int //"ebp"

	ToSlaveOrderListChan            chan [][]int //"exo"
	//ToSlaveImMasterChan             chan string  //"iam"
	ToSlaveReceivedConfirmationChan chan []int   //"rco"
	ToSlaveExecutedConfirmationChan chan []int   //"eco"

}
type InternalCommunicationChannels struct {
	newExternalList   chan [][]int
	slaveToStateExMasterChanshan chan int //send input to statemachine
}

func external_comm_channels_init() {
	ExCommChans.ToMasterSlaveChan = make(chan Slave)                   //"sla"
	ExCommChans.ToMasterOrderReceivedChan = make(chan []int)           //"ore"
	ExCommChans.ToMasterOrderExecutedChan = make(chan []int)           //"oex"
	ExCommChans.ToMasterOrderConfirmedReceivedChan = make(chan []int)  //"ocr"
	ExCommChans.ToMasterOrderConfirmedExecutionChan = make(chan []int) //"oce"
	ExCommChans.ToMasterExternalButtonPushedChan = make(chan []int)			//"ebp"

	ExCommChans.ToSlaveOrderListChan = make(chan [][]int)          //"exo"
	//ExCommChans.ToSlaveImMasterChan = make(chan string)            //"iam"
	ExCommChans.ToSlaveReceivedConfirmationChan = make(chan []int) //"rco"
	ExCommChans.ToSlaveExecutedConfirmationChan = make(chan []int) //"eco"



}
func internal_comm_chans_init() {
	InCommChans.newExternalList = make(chan [][]int)
	InCommChans.slaveToStateExMasterChanshan = make(chan int) //send input to statemachine
	//network

}

//Master
func Send_order(externalOrderList [][]int, c Conn) { //send exectuionOrderList
	byteOrder, _ := Marshal(externalOrderList)
	prefix, _ := Marshal("exo")
	ExNetChans.ConnChan <- c
	ExNetChans.ToNetwork <- append(prefix, byteOrder...)
}

func Send_im_master(c Conn) { //send I am master
	byteMessage, _ := Marshal("im master")
	prefix, _ := Marshal("iam")
	ExNetChans.ConnChan <- c
	ExNetChans.ToNetwork <- append(prefix, byteMessage...)

}
func Send_received_confirmation(order []int, c Conn) {
	byteMessage, _ := Marshal(order)
	prefix, _ := Marshal("rco")
	ExNetChans.ConnChan <- c
	ExNetChans.ToNetwork <- append(prefix, byteMessage...)
}

func Send_executed_confirmation(order []int, c Conn) {
	byteMessage, _ := Marshal(order)
	prefix, _ := Marshal("eco")
	ExNetChans.ConnChan <- c
	ExNetChans.ToNetwork <- append(prefix, byteMessage...)
}

//Slave
func Send_slave(s Slave, c Conn) {
	byteSlave, _ := Marshal(s)
	prefix, _ := Marshal("sla")
	ExNetChans.ConnChan <- c
	ExNetChans.ToNetwork <- append(prefix, byteSlave...)
}
func Send_ex_button_push(order []int, c Conn) {
	byteMessage, _ := Marshal(order)
	prefix, _ := Marshal("ebp")
	ExNetChans.ConnChan <- c
	ExNetChans.ToNetwork <- append(prefix, byteMessage...)
}

func Send_order_received(order []int, c Conn) {
	byteMessage, _ := Marshal(order)
	prefix, _ := Marshal("ore")
	ExNetChans.ConnChan <- c
	ExNetChans.ToNetwork <- append(prefix, byteMessage...)
}

func Send_order_executed(order []int, c Conn) {
	byteMessage, _ := Marshal(order)
	prefix, _ := Marshal("oex")
	ExNetChans.ConnChan <- c
	ExNetChans.ToNetwork <- append(prefix, byteMessage...)
}

func Send_order_confirmed_received(order []int, c Conn) {
	byteMessage, _ := Marshal(order)
	prefix, _ := Marshal("ocr")
	ExNetChans.ConnChan <- c
	ExNetChans.ToNetwork <- append(prefix, byteMessage...)
}

func Send_order_confirmed_executed(order []int, c Conn) {
	byteMessage, _ := Marshal(order)
	prefix, _ := Marshal("oce")
	ExNetChans.ConnChan <- c
	ExNetChans.ToNetwork <- append(prefix, byteMessage...)
}

func Decrypt_message(message []byte) {

	switch {
	//Master
	case string(message[1:4]) == "sla":
		noPrefix := message[5:]
		var s Slave
		_ = Unmarshal(noPrefix, &s)
		ExCommChans.ToMasterSlaveChan <- s

	case string(message[1:4]) == "ore":
		noPrefix := message[5:]
		order := make([]int, 2)
		_ = Unmarshal(noPrefix, &order)
		ExCommChans.ToMasterOrderReceivedChan <- order

	case string(message[1:4]) == "oex":
		noPrefix := message[5:]
		order := make([]int, 2)
		_ = Unmarshal(noPrefix, &order)
		ExCommChans.ToMasterOrderExecutedChan <- order

	case string(message[1:4]) == "ocr":
		noPrefix := message[5:]
		order := make([]int, 2)
		_ = Unmarshal(noPrefix, &order)
		ExCommChans.ToMasterOrderConfirmedReceivedChan <- order

	case string(message[1:4]) == "oce":
		noPrefix := message[5:]
		order := make([]int, 2)
		_ = Unmarshal(noPrefix, &order)
		ExCommChans.ToMasterOrderConfirmedExecutionChan <- order

	//Slave
	case string(message[1:4]) == "ebp":
		noPrefix := message[5:]
		order := make([]int,2)
		_ = Unmarshal(noPrefix, &order)
		ExCommChans.ToMasterExternalButtonPushedChan

	case string(message[1:4]) == "exo":
		noPrefix := message[5:]
		externalOrderList := make([][]int, FLOORS)
		_ = Unmarshal(noPrefix, &externalOrderList)
		ExCommChans.ToSlaveOrderListChan <- externalOrderList

	case string(message[1:4]) == "iam":
		fmt.Println("iam trigger")
		noPrefix := message[5:]
		stringMessage := string(noPrefix)
		fmt.Println(stringMessage)
		ExCommChans.ToSlaveImMasterChan <- stringMessage
		fmt.Println("channel output")

	case string(message[1:4]) == "rco":
		noPrefix := message[5:]
		order := make([]int, 2)
		_ = Unmarshal(noPrefix, &order)
		ExCommChans.ToSlaveReceivedConfirmationChan <- order

	case string(message[1:4]) == "eco":
		noPrefix := message[5:]
		order := make([]int, 2)
		_ = Unmarshal(noPrefix, &order)
		ExCommChans.ToSlaveExecutedConfirmationChan <- order
	}
}
func Select_send_master(c Conn) {

	for {
		select {
		//Master

		case externalOrderList := <-ExMasterChans.ToCommOrderListChan:
			Send_order(externalOrderList, c)
		case order := <-ExMasterChans.ToCommReceivedConfirmationChan:
			Send_received_confirmation(order,c)
		case order := <-ExMasterChans.ToCommExecutedConfirmationChan:
			Send_executed_confirmation(order,c)
		default:
			Send_im_master(c)
		}
	}
}
func Select_send_slave(c Conn) {
	for {
		select {
		//Slave
		case slave := <-ExSlaveChans.ToCommSlaveChan:
			Send_slave(slave, c)
		case order := <-ExSlaveChans.ToCommExternalButtonPushedChan:
			Send_ex_button_push(order, c)
		case order := <-ExSlaveChans.ToCommOrderReceivedChan:
			Send_order_received(order, c)
		case order := <-ExSlaveChans.ToCommOrderConfirmedReceivedChan:
			Send_order_executed(order, c)
		case order := <-ExSlaveChans.ToCommOrderConfirmedReceivedChan:
			Send_order_confirmed_received(order, c)
		case order := <-ExSlaveChans.ToCommOrderConfirmedExecutuinChan:
			Send_order_confirmed_executed(order, c)
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
