package network

import (
	. "chansnstructs"
	. "encoding/json"
	"fmt"
	. "net"
	"time"
)

var InCommChans InternalCommunicationChannels

type InternalCommunicationChannels struct {
	newExternalList              chan []Order
	slaveToStateExMasterChanshan chan int //send input to statemachine
}

func internal_comm_chans_init() {
	InCommChans.newExternalList = make(chan []Order)
	InCommChans.slaveToStateExMasterChanshan = make(chan int) //send input to statemachine
	//network
}

//Master
func Send_order(externalOrderList []Order, c Conn) { //send exectuionOrderList
	byteOrder, _ := Marshal(externalOrderList)
	prefix, _ := Marshal("ord")
	byteOrder = append(prefix, byteOrder...)
	for {
		Push_to_network(byteOrder, c)
		select {
		case <-ExCommChans.ToMasterOrderListReceivedChan:
			return
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}

//To master
func Send_order_received(order Order, c Conn) {
	byteMessage, _ := Marshal(order)
	prefix, _ := Marshal("ore")

	ExNetChans.ToNetwork <- append(prefix, byteMessage...)
}

//To master
func Send_order_executed(order Order, c Conn) {
	byteMessage, _ := Marshal(order)
	prefix, _ := Marshal("oex")
	ExNetChans.ConnChan <- c
	ExNetChans.ToNetwork <- append(prefix, byteMessage...)
}

//To slave
func Send_order_executed_confirmation(order Order, c Conn) {
	byteMessage, _ := Marshal(order)
	prefix, _ := Marshal("eco")
	ExNetChans.ConnChan <- c
	ExNetChans.ToNetwork <- append(prefix, byteMessage...)
}

//
func Send_order_executed_reconfirmed(order Order, c Conn) {
	byteMessage, _ := Marshal(order)
	prefix, _ := Marshal("oce")
	ExNetChans.ConnChan <- c
	ExNetChans.ToNetwork <- append(prefix, byteMessage...)
}

//To master
func Send_slave(s Slave, c Conn) {
	byteSlave, _ := Marshal(s)
	prefix, _ := Marshal("sla")
	ExNetChans.ConnChan <- c
	ExNetChans.ToNetwork <- append(prefix, byteSlave...)
}

//To master
func Send_ex_button_push(order Order, c Conn) {
	byteMessage, _ := Marshal(order)
	prefix, _ := Marshal("ebp")
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
		case order := <-ExMasterChans.ToCommOrderExecutedConfirmedChan:

			Send_order_executed_confirmation(order, c)
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
		case order := <-ExSlaveChans.ToCommOrderListReceivedChan:
			Send_order_received(order, c)
		case order := <-ExSlaveChans.ToCommOrderExecutedChan:
			Send_order_executed(order, c)
		case order := <-ExSlaveChans.ToCommOrderExecutedReConfirmedChan:
			Send_order_executed_reconfirmed(order, c)
		case order := <-ExSlaveChans.ToCommExternalButtonPushedChan:
			Send_ex_button_push(order, c)
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
	case string(message[1:4]) == "ord":
		noPrefix := message[5:]
		orderList := make([]Order, 2)
		_ = Unmarshal(noPrefix, &orderList)
		ExCommChans.ToSlaveOrderListChan <- orderList

	case string(message[1:4]) == "ore":
		noPrefix := message[5:]
		order := Order{}
		_ = Unmarshal(noPrefix, &order)
		ExCommChans.ToMasterOrderListReceivedChan <- IpOrderMessage{addr, order}

	case string(message[1:4]) == "oex":
		noPrefix := message[5:]
		order := Order{}
		_ = Unmarshal(noPrefix, &order)
		ExCommChans.ToMasterOrderExecutedChan <- IpOrderMessage{addr, order}

	case string(message[1:4]) == "eco":
		noPrefix := message[5:]
		order := Order{}
		_ = Unmarshal(noPrefix, &order)
		ExCommChans.ToSlaveOrderExecutedConfirmedChan <- IpOrderMessage{addr, order}

	case string(message[1:4]) == "sla":
		noPrefix := message[5:]
		var s Slave
		_ = Unmarshal(noPrefix, &s)
		ExCommChans.ToMasterSlaveChan <- IpSlave{addr, s}

	case string(message[1:4]) == "ebp":
		noPrefix := message[5:]
		order := Order{}
		_ = Unmarshal(noPrefix, &order)
		ExCommChans.ToMasterExternalButtonPushedChan <- IpOrderMessage{addr, order}

	case string(message[1:4]) == "oce":
		noPrefix := message[5:]
		order := Order{}
		_ = Unmarshal(noPrefix, &order)
		ExCommChans.ToMasterOrderExecutedReConfirmedChan <- IpOrderMessage{addr, order}

	case string(message[1:4]) == "iam":
		noPrefix := message[5:]
		stringMessage := string(noPrefix)
		ExCommChans.ToSlaveImMasterChan <- stringMessage

	}
}
