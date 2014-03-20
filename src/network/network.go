package network

import (
	"fmt"
	"math/rand"
	. "net"
	"os"
	"time"
)

const (
	MAXWAIT = time.Second
	PORT    = ":20019"
)

var ExNetChans NetworkExternalChannels

var ExSlaveChans ExternalSlaveChannels
var ExMasterChans ExternalMasterChannels
var InSlaveChans InternalSlaveChannels
var InMasterChans InternalMasterChannels

var ExStateMChans ExternalStateMachineChannels

type Master struct {
	m          []Slave
	OrderQueue [][]int
}

type ExternalSlaveChannels struct {
	ToCommSlaveChan                   chan Slave  //"sla"
	ToCommOrderReceivedChan           chan []int  //"ore"
	ToCommOrderExecutedChan           chan []int  //"oex"
	ToCommOrderListReceivedChan       chan []int  //"ocr"
	ToCommOrderConfirmedExecutionChan chan []int  //"oce"
	ToCommExternalButtonPushedChan    chan []int  //"ebp"
	ToCommImMasterChan                chan string //"iam"

}
type ExternalMasterChannels struct {
	ToCommOrderListChan            chan [][]int //"exo"
	ToCommReceivedConfirmationChan chan []int   //"rco"
	ToCommExecutedConfirmationChan chan []int   //"eco"

}
type InternalSlaveChannels struct {
	OrderConfirmedExecutedChan chan []int
	InteruptChan               chan os.Signal
}
type InternalMasterChannels struct {
	OptimizationInitChan    chan Master
	OptimizationTriggerChan chan ipOrderMessage
	OptimizationReturnChan  chan [][]int
	OrderReceivedMangerChan chan ipOrderMessage
}

func Slave_external_chans_init() {
	ExSlaveChans.ToCommSlaveChan = make(chan Slave)                   //"sla"
	ExSlaveChans.ToCommOrderReceivedChan = make(chan []int)           //"ore"
	ExSlaveChans.ToCommOrderExecutedChan = make(chan []int)           //"oex"
	ExSlaveChans.ToCommOrderListReceivedChan = make(chan []int)       //"ocr"
	ExSlaveChans.ToCommOrderConfirmedExecutionChan = make(chan []int) //"oce"
	ExSlaveChans.ToCommExternalButtonPushedChan = make(chan []int)    //"ebp"

}
func Master_external_chans_init() {
	ExMasterChans.ToCommOrderListChan = make(chan [][]int)          //"exo"
	ExMasterChans.ToCommReceivedConfirmationChan = make(chan []int) //"rco"
	ExMasterChans.ToCommExecutedConfirmationChan = make(chan []int) //"eco"

}
func Slave_internal_chans_init() {
	InSlaveChans.OrderConfirmedExecutedChan = make(chan []int)
	InSlaveChans.InteruptChan = make(chan os.Signal, 1) //must be buffered see package declaration
}
func Master_internal_chans_init() {
	InMasterChans.OptimizationInitChan = make(chan Master)
	InMasterChans.OptimizationTriggerChan = make(chan ipOrderMessage)
	InMasterChans.OptimizationReturnChan = make(chan [][]int)
	InMasterChans.OrderReceivedMangerChan = make(chan ipOrderMessage)
}

type ExternalStateMachineChannels struct {
	ExternalButtonPressed chan []int
	OrderServed           chan []int
	CurrentState          chan []int
	GetSlaveStruct        chan Slave
	DirectionUpdate       chan int
}

func External_state_machine_channels_init() {
	ExStateMChans.ExternalButtonPressed = make(chan []int)
	ExStateMChans.OrderServed = make(chan []int)
	ExStateMChans.CurrentState = make(chan []int)
	ExStateMChans.GetSlaveStruct = make(chan Slave)
	ExStateMChans.DirectionUpdate = make(chan int)
}

type Button struct {
	floor      int
	buttonType int
	turnOn     bool
}

type State struct {
	direction    int
	currentFloor int
}

/*
func main() {
	nr := backup()

	primary(nr)

}
*/
func (m Master) Master_communication() {

	for {
		select {
		//triggers new optimization when new order received
		case order := <-ExCommChans.ToMasterExternalButtonPushedChan:

			InMasterChans.OptimizationTriggerChan <- order

			//Need to have same queueing system as order executed if different orders are coming in
			//Same things need to be done, but we must also calculate some optimizationø
			//We compute optimization again if the queue is not empty
			//receives new optimized orderList
		case orderList := <-InMasterChans.OptimizationReturnChan:
			//send to slaves master
			ExMasterChans.ToCommOrderListChan <- orderList

		case order := <-ExCommChans.ToMasterOrderExecutedChan: //to spesific IP
			if m.notInQueue(order) {
				m.OrderQueue = appendElement(m.OrderQueue, order)
				//externalPushChannel <- externalPushQueue //Where its sending must sending the first element in the list if its not empty. otherwise just update
			}

			//ExMasterChans.ToCommExecutedConfirmationChan <- order

			//Master gets a message that order is executed.
			//Save the order in a temp variable
			//calls ordere_executed_manager
			//order_exe sends to channel when its done

			//if any other incoming orders are coming while order_exe is running, queue them in a list
			//reset the temp var if order_exe is done and queue is empty
			// if not empty, extract first in queue and set to temp var

		//Respond on orderList received
		case orderList := <-ExCommChans.ToMasterOrderListReceivedChan: //with spesific IP
			InMasterChans.OrderReceivedMangerChan <- orderList
			//Done

		case tempIpSlave := <-ExCommChans.ToMasterSlaveChan:
			slave := tempIpSlave.s
			nr := slave.nr
			m[nr] = slave
		}
	}
}
func (s Slave) Slave_communication() {

	for {
		select {

		//These two needs must trigger a send_state that doesnt end until master has confirmed receiving it.
		case slave := <-ExStateChans.DirectionUpdate:
			ExCommChans.ToMasterSlaveChan <- slave
		case slave := <-ExStateChans.CurrentFloorUpdate:
			ExCommChans.ToMasterSlaveChan <- slave

		//checks new button pressed, send to master until confirmation
		case order := <-ExStateChans.ToSlaveExButtonPushedChan:
			ExSlaveChans.ToCommExternalButtonPushedChan <- order

		//receives order list
		case orderList := <-ExCommChan.ToSlaveOrderListChan:
			s.externalList = orderList
			//###### set lights on external floor
			ExSlaveChans.ToCommOrderListReceivedChan <- order
			ExStateChans.NewOrderListChan <- orderList //set lights

		//Triggers if order is executed, sends confirmation further up
		case order := <-ExStateChans.ToSlaveExecutedOrderChan:
			ExSlaveChans.ToCommOrderExecutedChan <- order
		}
	}
}

type NetworkExternalChannels struct {
	ToNetwork  chan []byte
	ToComm     chan []byte
	ConnChan   chan Conn
	ToCommAddr chan *UDPAddr
}

func network_external_chan_init() {
	ExNetChans.ToNetwork = make(chan []byte)
	ExNetChans.ToComm = make(chan []byte)
	ExNetChans.ConnChan = make(chan Conn)
	ExNetChans.ToCommAddr = make(chan *UDPAddr)
}

func Network_init() (Conn, Conn) {
	fmt.Println("gi")
	addr, err := ResolveUDPAddr("udp", "129.241.187.255"+PORT) //leser bare fra porten generellt
	c1, err := DialUDP("udp", nil, addr)

	if err != nil {
		fmt.Println(err.Error())
	}

	defer c1.Close()
	addr2, _ := ResolveUDPAddr("udp", PORT)
	c2, err := ListenUDP("udp", addr2)

	//c2.SetReadDeadline(time.Now().Add(300 * time.Millisecond)) //returns error if deadline is reached
	//_, _, err = c2.ReadFromUDP(buf)

	return c1, c2

}

func Send(to_writing []byte, c Conn) {
	_, err := c.Write(to_writing)

	if err != nil {
		fmt.Println(err.Error())
	} else {
		//break
	}

}

func Receive() { //will error trigger if just read fails? or will it only go on deadline?
	buf := make([]byte, 1024)
	addr, _ := ResolveUDPAddr("udp", PORT)
	c, err := ListenUDP("udp", addr)

	defer c.Close()
	//this will also check if the master is still there.
	c.SetReadDeadline(time.Now().Add(300 * time.Millisecond)) //returns error if deadline is reached
	n, sendingAddr, err := c.ReadFromUDP(buf)                 //n contanis numbers of used bytes, fills buf with content on the connection

	fmt.Println(sendingAddr)

	if err == nil { //if error is nil, read from buffer
		ExNetChans.ToComm <- buf[0:n]
		ExNetChans.ToCommAddr <- addr

		//ExSlaveChans.ToSlaveImMasterChan <- true
	} else {
		//ExSlaveChans.ToSlaveImMasterChan <- false
	}

}

///CORRECT FORMAT??
func (m Master) notInQueue(order) bool {
	for i := 0; i < len(m.OrderQueue); i++ {
		//	if m.OrderQueue[i][] == order {
		//	return true
		//}
	}
	return false

}
func appendElement(slice [][]int, order ipOrderMessage) [][]int {
	for _, item := range order[1] {
		slice = Extend(slice, order[1])
	}
	return slice
}

/*
func Choose_master() {
	go Slave_elevator()
}

func Slave_elevator() {

	buf := make([]byte, 1024)
	addr, _ := ResolveUDPAddr("udp", PORT)
	c, _ := ListenUDP("udp", addr)

	defer c.Close()

	for {
		//rand := Random_int(600, 1000)
		c.SetReadDeadline(time.Now().Add(300 * Millisecond))
		_, _, err := c.ReadFromUDP(buf) //n contanis numbers of used bytes

		if err == nil { // of readdeadline dont kicks in
			//decrypt buf
			//if decryptet buf equals iam
			//keep on serching
			Decrypt_message(buf)
			<-ExComChans.ToSlaveImMasterChan

		} else { // if readdeadline kicks in
			//first one here becomes master(?)

			//this will just be called in case of there is no master
			go Master_elevator()
			break
		}

	}
}

func Master_elevator() {

	for {
		MC.ToCommImMasterChan <- true
		time.Sleep(50 * time.Millisecond)
	}
}
*/
func Random_init(min int, max int) int { //gives a random int for waiting
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}
