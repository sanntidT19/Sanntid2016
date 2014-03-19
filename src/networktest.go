package main

import (
	//. "./Network"
	"fmt"
	//. "time"
	"os"
	"net"
)
/* LookupHost
 */

func main() {
 if len(os.Args) != 3 {
 fmt.Fprintf(os.Stderr,"Usage: %s network-type service\n",os.Args[0])
 os.Exit(1)
 }
 networkType := os.Args[1]
 service := os.Args[2]
 port, err := net.LookupPort(networkType, service)
 if err != nil {
 fmt.Println("Error: ", err.Error())
 os.Exit(2)
 }
 fmt.Println("Service port ", port)
 os.Exit(0)
}


/*
type Slave struct {
	nr           int
	internalList []bool
	externalList [][]int
	currentFloor int //get from driver/IO
	direction    int // get from driver/IO

}

func main() {

	//Channels_init()
	//slaveToCommSlaveChan := make(chan Slave)                   //"sla"
	slaveToCommOrderReceivedChan := make(chan []int)           //"ore"
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
	commToMasterOrderReceivedChan := make(chan []int)           //"ore"
	commToMasterOrderExecutedChan := make(chan []int)           //"oex"
	commToMasterOrderConfirmedReceivedChan := make(chan []int)  //"ocr"
	commToMasterOrderConfirmedExecutionChan := make(chan []int) //"oce"

	commToSlaveOrderListChan := make(chan [][]int)          //"exo"
	commToSlaveImMasterChan := make(chan string)            //"iam"
	commToSlaveReceivedConfirmationChan := make(chan []int) //"rco"
	commToSlaveExecutedConfirmationChan := make(chan []int) //"eco"

	//newExternalList := make(chan [][]int)
	//slaveToStateMChan := make(chan int) //send input to statemachine
	//network
	commToNetwork := make(chan []byte)
	networkToComm := make(chan []byte)

	go Select_send(commToNetwork, slaveToCommOrderReceivedChan, slaveToCommOrderExecutedChan, slaveToCommOrderConfirmedReceivedChan, slaveToCommOrderConfirmedExecutuinChan, masterToCommOrderListChan, masterToCommImMasterChan, masterToCommReceivedConfirmationChan, masterToCommExecutedConfirmationChan)

	masterToCommImMasterChan <- "hu og hei"

	go Select_receive(networkToComm /*commToMasterSlaveChan chan Slave,, commToMasterOrderReceivedChan, commToMasterOrderExecutedChan, commToMasterOrderConfirmedReceivedChan, commToMasterOrderConfirmedExecutionChan, commToSlaveOrderListChan, commToSlaveImMasterChan, commToSlaveReceivedConfirmationChan, commToSlaveExecutedConfirmationChan)
	fmt.Println("commToNetwork")

	sending := <-commToNetwork
	c := Network_init()
	go Send(c, sending)

	go Receive(networkToComm)

	go Decrypt_message(<-networkToComm, commToMasterOrderReceivedChan, commToMasterOrderExecutedChan, commToMasterOrderConfirmedReceivedChan, commToMasterOrderConfirmedExecutionChan, commToSlaveOrderListChan, commToSlaveImMasterChan, commToSlaveReceivedConfirmationChan, commToSlaveExecutedConfirmationChan)
	//fmt.Println(string(<-networkToComm))
	fmt.Println("finito")
	fmt.Println("channel commToSlave", <-commToSlaveImMasterChan)
	Sleep(10 * Second)

}
*/