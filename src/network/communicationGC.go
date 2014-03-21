package network

import (
	. "chansnstructs"
	. "encoding/json"
	. "net"
)

var InCommChans InternalCommunicationChannels

type InternalCommunicationChannels struct {
	newExternalList              chan [][]int
	slaveToStateExMasterChanshan chan int //send input to statemachine
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
	byteOrder = append(prefix, byteOrder...)
	for {
		Send(byteOrder, c)
		select {
		case <-ExCommChans.ToMasterOrderListReceivedChan:
			return
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
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

	byteMessage = append(prefix, byteMessage...)

	Send(byteMessage, c)
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

func Send_im_master(c Conn) { //send I am master
	byteMessage, _ := Marshal("im master")
	prefix, _ := Marshal("iam")
	ExNetChans.ConnChan <- c
	ExNetChans.ToNetwork <- append(prefix, byteMessage...)

}

func Select_send_master(c Conn) {

	for {
		select {
		//Master

		case externalOrderList := <-ExMasterChans.ToCommOrderListChan:
			Send_order(externalOrderList, c)
		case order := <-ExMasterChans.ToCommReceivedConfirmationChan:
			Send_received_confirmation(order, c)
		case order := <-ExMasterChans.ToCommExecutedConfirmationChan:
			Send_executed_confirmation(order, c)
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
		case order := <-ExSlaveChans.ToCommOrderListReceivedChan:
			Send_order_executed(order, c)
		case order := <-ExSlaveChans.ToCommOrderConfirmedExecutionChan:
			Send_order_confirmed_executed(order, c)
		}
	}
}
func Select_receive() {
	fmt.Println("Select_receive")
	for {
		barr := <-ExNetChans.ToComm
		addr := <-ExNetChans.ToCommAddr
		Decrypt_message(barr, addr)
	}
}

func Decrypt_message(message []byte, addr *UDPAddr) {

	switch {
	//Master
	case string(message[1:4]) == "sla":
		noPrefix := message[5:]
		var s Slave
		_ = Unmarshal(noPrefix, &s)
		ExCommChans.ToMasterSlaveChan <- ipSlave{addr, s}

	case string(message[1:4]) == "ore":
		noPrefix := message[5:]
		order := make([]int, 2)
		_ = Unmarshal(noPrefix, &order)
		ExCommChans.ToMasterOrderListReceivedChan <- ipOrderMessage{addr, order}

	case string(message[1:4]) == "oex":
		noPrefix := message[5:]
		order := make([]int, 2)
		_ = Unmarshal(noPrefix, &order)
		ExCommChans.ToMasterOrderExecutedChan <- ipOrderMessage{addr, order}

	case string(message[1:4]) == "ocr":
		noPrefix := message[5:]
		order := make([]int, 2)
		_ = Unmarshal(noPrefix, &order)
		ExCommChans.ToMasterOrderListReceivedChan <- ipOrderMessage{addr, order}

	case string(message[1:4]) == "oce":
		noPrefix := message[5:]
		order := make([]int, 2)
		_ = Unmarshal(noPrefix, &order)
		ExCommChans.ToMasterOrderConfirmedExecutionChan <- ipOrderMessage{addr, order}

	//Slave
	case string(message[1:4]) == "ebp":
		noPrefix := message[5:]
		order := make([]int, 2)
		_ = Unmarshal(noPrefix, &order)
		ExCommChans.ToMasterExternalButtonPushedChan <- ipOrderMessage{addr, order}

	case string(message[1:4]) == "exo":
		noPrefix := message[5:]
		externalOrderList := make([][]int, N_FLOORS)
		_ = Unmarshal(noPrefix, &externalOrderList)
		ExCommChans.ToSlaveOrderListChan <- externalOrderList

	case string(message[1:4]) == "eco":
		noPrefix := message[5:]
		order := make([]int, 2)
		_ = Unmarshal(noPrefix, &order)
		ExCommChans.ToSlaveExecutedConfirmationChan <- ipOrderMessage{addr, order}
		/*
			case string(message[1:4]) == "iam":
				noPrefix := message[5:]
				stringMessage := string(noPrefix)
				ExCommChans.ToSlaveImMasterChan <- stringMessage
		*/
	}
}
