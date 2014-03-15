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
	slaveToCommOrderRecevedChan := make(chan string)  //notity master that slave has receved order
	slaveToCommOrderExecutedChan := make(chan string) //notify master that slave has executed order
	slaveToCommConfirmedExecutuinChan := make(chan string)

	slaveToStateMChan := make(chan int) //send input to statemachine

	//Master
	masterToCommOrderChan := make(chan [][]int)              //sends orders from slave to comm
	masterToCommConfirmChan := make(chan bool)               //confirms that master har receved that slave has confirmed/receved order
	masterToCommImMasterChan := make(chan string)            // sends i am master
	masterToCommRecevedConfirmationChan := make(chan string) // master confirms that slave has receved order

	//communication channels
	commToSlaveOrderChan := make(chan [][]int)             //receves orders from comm
	commToSlaveMastersBackConfirmChan := make(chan string) //master confirms that order is receved
	commToSlaveImMasterChan := make(chan string)           //im master from master
	commToSlaveRecevedConfirmationChan := make(chan string)

	commToMasterFloorRChan := make(chan bool)               //floor reached from slave to master
	commToMasterSlaveStructChan := make(chan Slave)         //sends slave struct
	commToMasterOrderExecuredChan := make(chan string)      //order executed sucessfully
	commToMasterOrderRecevedChan := make(chan string)       //Slave confirmes that order is recived
	commToMasterConfirmedExecutionChan := make(chan string) //slave confirmes order executed

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

func Send_order_receved(slaveToCommOrderRecivedChan chan string, commToNetwork chan []byte) {
	byteMessage, err := Marshal(<-slaveToCommOrderRecivedChan)
	prefix, err := Marshal("mre")
	commToNetwork <- append(prefix, byteMessage)
}

func Send_order_executed(slaveToCommOrderExecuredChan chan string, commToNetwork chan []byte) {
	byteMessage, err := Marshal(<-slaveToCommOrderExecuredChan)
	prefix, err := Marshal("exe")
	commToNetwork <- append(prefix, byteMessage)
}

func Send_receved_confirmation(masterToCommRecevedConfirmationChan chan string, commToNetwork chan []byte) {
	byteMessage, err := Marshal(<-masterToCommRecevedConfirmationChan)
	prefix, err := Marshal("sre")
	commToNetwork <- append(prefix, byteMessage)

}
func Send_slave(slaveToCommSlaveStructChan chan Slave, commToNetwork chan []byte) {
	byteSlave, err := Marshal(<-slaveToCommSlaveStructChan)
	prefix, err := Marshal("sch")
	commToNetwork <- append(prefix, byteSlave)
}

func Decrypt_message(message []byte) {
	switch {

	case HasPrefix(message, "ord"): //externalorderlist
		str := string(message)
		str = TrimPrefix(str, "ord")
		externalOrderList := Unmarshal([]byte(message), [][]int)
		commToSlaveOrderChan <- externalOrderList

	case HasPrefix(message, "iam"): //I am master
		str := string(message)
		str = TrimPrefix(str, "iam")
		commToSlaveImMasterChan <- str

	case HasPrefix(message, "mre"): //confirm recived order from master
		str := string(message)
		str = TrimPrefix(str, "mre")
		commToSlaveOrderRecevedChan <- str

	case HasPrefix(message, "exe"): //Order performed from slave to master
		str := string(message)
		str = TrimPrefix(str, "exe")
		commToMasterOrderExecuredChan <- str

	case HasPrefix(message, "sre"): //confirmes recived order from slave to master
		str := string(message)
		str = TrimPrefix(str, "sre")
		commToSlaveRecevedConfirmationChan <- str

	case HasPrefix(message, "sch"): // receves a slave struct
		str := string(message)
		str = TrimPrefix(str, "sch")
		messageNoPrefix := []byte(str)
		slave := UnMarshal(messageNoPrefix, &Slave)
		ommToMasterSlaveStructChan <- slave
	}

}
