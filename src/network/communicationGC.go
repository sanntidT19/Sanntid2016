package network

import (
	. "chansnstructs"
	. "encoding/json"
	"fmt"
	. "net"
	//"time"
)

var InCommChans InternalCommunicationChannels

type InternalCommunicationChannels struct {
	newExternalList              chan []Order
	slaveToStateExMasterChanshan chan int //send input to statemachine
}

func Select_send_master() {

	for {
		select {
		//Master
		case externalOrderList := <-ExMasterChans.ToCommOrderListChan:
			Send_order(externalOrderList)
		case order := <-ExMasterChans.ToCommOrderExecutedConfirmedChan:
			Send_order_executed_confirmation(order)
		case message := <-ExMasterChans.ToCommImMasterChan:
			Send_im_master(message)
		}
	}
}
func Select_send_slave() {
	for {
		select {
		//Slave
		case slave := <-ExSlaveChans.ToCommSlaveChan:
			Send_slave(slave)
		case order := <-ExSlaveChans.ToCommOrderListReceivedChan:
			Send_order_received(order)
		case order := <-ExSlaveChans.ToCommOrderExecutedChan:
			Send_order_executed(order)
		case order := <-ExSlaveChans.ToCommOrderExecutedReConfirmedChan:
			Send_order_executed_reconfirmed(order)
		case order := <-ExSlaveChans.ToCommExternalButtonPushedChan:
			Send_ex_button_push(order)
		case ipOrder := <-ExSlaveChans.ToCommImSlaveChan:
			ip := ipOrder.Ip
			Send_im_slave(ip)

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

func internal_comm_chans_init() {
	InCommChans.newExternalList = make(chan []Order)
	InCommChans.slaveToStateExMasterChanshan = make(chan int) //send input to statemachine
	//network
}

//Master
func Send_order(externalOrderList []Order) { //send exectuionOrderList
	byteOrder, _ := Marshal(externalOrderList)
	prefix, _ := Marshal("ord")
	byteOrder = append(prefix, byteOrder...)
	ExNetChans.ToNetwork <- append(prefix, byteOrder...)
}

//To master
func Send_order_received(order Order) {
	byteMessage, _ := Marshal(order)
	prefix, _ := Marshal("ore")
	ExNetChans.ToNetwork <- append(prefix, byteMessage...)
}

//To master
func Send_order_executed(order Order) {
	byteMessage, _ := Marshal(order)
	prefix, _ := Marshal("oex")
	ExNetChans.ToNetwork <- append(prefix, byteMessage...)
}

//To slave
func Send_order_executed_confirmation(order Order) {
	byteMessage, _ := Marshal(order)
	prefix, _ := Marshal("eco")
	ExNetChans.ToNetwork <- append(prefix, byteMessage...)
}

//
func Send_order_executed_reconfirmed(order Order) {
	byteMessage, _ := Marshal(order)
	prefix, _ := Marshal("oce")
	ExNetChans.ToNetwork <- append(prefix, byteMessage...)
}

//To master
func Send_slave(s Slave) {
	byteSlave, _ := Marshal(s)
	prefix, _ := Marshal("sla")
	ExNetChans.ToNetwork <- append(prefix, byteSlave...)
}

//To master
func Send_ex_button_push(order Order) {
	byteMessage, _ := Marshal(order)
	prefix, _ := Marshal("ebp")
	ExNetChans.ToNetwork <- append(prefix, byteMessage...)
}

func Send_im_master(message string) { //send I am master
	//ExNetChans.ConnChan <- c
	fmt.Println("send master")
	byteMessage, _ := Marshal(message)
	prefix, _ := Marshal("iam")
	ExNetChans.ToNetwork <- append(prefix, byteMessage...)
	fmt.Println("end of send im master")

}
func Send_im_slave(ip *UDPAddr) {
	ipOrder := IpOrderMessage{ip, Order{}}
	byteIpOrder, _ := Marshal(ipOrder)
	prefix, _ := Marshal("ias")
	ExNetChans.ToNetwork <- append(prefix, byteIpOrder...)
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
		fmt.Println("iam trigger")
		noPrefix := message[5:]
		stringMessage := string(noPrefix)
		ExCommChans.ToSlaveImMasterChan <- stringMessage

	case string(message[1:4]) == "ias":
		noPrefix := message[5:]
		ipOrder := IpOrderMessage{}
		_ = Unmarshal(noPrefix, &ipOrder)
		ExCommChans.ToMasterImSlaveChan <- ipOrder
	default:
		fmt.Println("ingen caser utløst")
	}

}
