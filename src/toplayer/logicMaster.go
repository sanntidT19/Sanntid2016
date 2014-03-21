package toplayer

import "chansnstructs"

var InMasterChans InternalMasterChannels
var InSlaveChans InternalSlaveChannels

type InternalMasterChannels struct {
	OptimizationInitChan    chan Master
	OptimizationTriggerChan chan ipOrderMessage
	OptimizationReturnChan  chan [][]int
	OrderReceivedMangerChan chan ipOrderMessage
}
type InternalSlaveChannels struct {
	OrderConfirmedExecutedChan chan []int
	InteruptChan               chan os.Signal
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

func Slave_init() {
	buf := make([]byte, 1024)
	s := Slave{}
	connSend, connReceive := Network_init()
	err := connReceive.SetReadDeadline(time.Now().Add(10 * time.Millisecond)) //Random_init(10, 100)

	//	_, _, err = connReceive.ReadFromUDP(buf) //n contanis numbers of used bytes
	if err != nil { //run master if connection fails
		go Master_init(s)
	} else {
		s = Slave{}
		go Select_send_slave(connSend)
		go Select_receive()

		//communicate with statemachine
		//go Recive_externalList()

		// ##1 Send order received <- update external list in struct

		// ##2 send floor reached <- get order from statemachine
		// <- receive confirmation on order executed and set lights

	}
}

func Master_init(s Slave) {
	connSend, _ := Network_init()
	ExNetChans.ConnChan <- connSend

	//intial sending contianing: ipadress, initialization,
	//just listen to this
	m := Master{}
	go Select_send_master(connSend)
	go Select_receive()
	//if call for new optimalization
	//recive it and  iterate slaves laves
	go Interuption_killer()
	// ##1 wait for order received <- stop send orders

	// ##2 wait for floor reacehed <- delete order from masters externalList
	// <- send to all order executed <- all slaves must update the externalpanel

}

//retrice orders if network fails in some way
func Init_Orders() {

}

//handles signals like ctrl-c
func Interuption_killer() {
	signal.Notify(InSlaveChans.InteruptChan, os.Interrupt)
	signal.Notify(InSlaveChans.InteruptChan, syscall.SIGTERM)
	<-InSlaveChans.InteruptChan
	//what should be done when ctrl-c is pressed?????
	//<- goes here.
	//SystemInit()
	fmt.Println("Got ctrl-c signal")
	os.Exit(0)
}

func Error_manager() { // handle error here?

}
func (m Master) Master_communication() {

	for {
		select {
		//triggers new optimization when new order received
		case order := <-ExCommChans.ToMasterExternalButtonPushed:

			InMasterChans.OptimizationTriggerChan <- order

			//Need to have same queueing system as order executed if different orders are coming in
			//Same things need to be done, but we must also calculate some optimizationø
			//We compute optimization again if the queue is not empty
			//receives new optimized orderList
		case orderList := <-InMasterChans.OptimizationReturnChan:
			//send to slaves master
			ExMasterChans.ToCommOrderListChan <- orderList

		case order := <-ExCommChans.ToMasterOrderExecutedChan: //to spesific IP
			if notInQueue(order) {
				externalPushQueue = appendElement(externalPushQueue, order)
				externalPushChannel <- externalPushQueue //Where its sending must sending the first element in the list if its not empty. otherwise just update
			}

			ExMasterChans.ToCommExecutedConfirmationChan <- order
			//Master gets a message that order is executed.
			//Save the order in a temp variable
			//calls ordere_executed_manager
			//order_exe sends to channel when its done

			//if any other incoming orders are coming while order_exe is running, queue them in a list
			//reset the temp var if order_exe is done and queue is empty
			// if not empty, extract first in queue and set to temp var

		//Respond on orderList received
		case order := <-ExCommChans.ToMasterOrderListReceivedChan: //with spesific IP
			InMasterChans.OrderReceivedMangerChan <- order
			//Done

		case slave <- ExCommChans.ToMasterSlaveChan:
			m[slave.nr] = slave
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
func Order_received_manager() {
	for {
		order := <-InMasterChans.OrderReceivedManagerChan
		/*
			need to send a confirmation message to the ip that sent this every time we get the order.
			After some time without fun incoming on the channel we can assume that the confirmation has been received.
		*/

	}
}

func (s Slave) Send_slave_to_state() { //send next floor to statemachine
	if s.externalList[s.currentFloor][1] == 1 || s.internalList[s.currentFloor] == 1 {
		ExStateChan.slaveStateChan <- s.currentFloor

	} else if s.direction == 1 { //heading upwards -> can take higher orders
		for i := s.currentFloor; i < N_FLOORS; i++ {
			if s.externalList[i][0] == 1 || s.internalList[i] == 1 { // any orders on higher floors
				ExStateChan.slaveStateChan <- i
				break
			}
		}
	} else if s.direction == -1 { // heading downwards -> can take lower orders
		for i := 0; i < s.currentFloor; i++ { //any orders on lower floors
			if s.externalList[i][0] == 1 || s.internalList[i] == 1 {
				ExSateChan.slaveStateChan <- i
				break
			}
		}
	} else {
		Sleep(10 * Millisecond)
	}
}
