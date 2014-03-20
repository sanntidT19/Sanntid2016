package toplayer

//export GOPATH =$HOME/Elevator-progress-/
import(
	."time"
	."os"
)
const (
	N_FLOORS = 4
)


//Global chans
var ExSlaveChans ExternalSlaveChannels
var ExMasterChans ExternalMasterChannels
var InSlaveChans InternalSlaveChannels
var InMasterChans InternalMasterChannels

type Slave struct {
	ip IP
	internalList []int
	externalList [][]int
	currentFloor int 
	direction    int 
}

type Master struct {
	m []Slave

}
type ExternalSlaveChannels struct {
	ToCommSlaveChan                   chan Slave //"sla"
	ToCommOrderReceivedChan           chan []int //"ore"
	ToCommOrderExecutedChan           chan []int //"oex"
	ToCommOrderListReceivedChan  chan []int //"ocr"
	ToCommOrderConfirmedExecutionChan chan []int //"oce"
	ToCommExternalButtonPushedChan chan []int //"ebp"
}
type ExternalMasterChannels struct {
	ToCommOrderListChan            chan [][]int //"exo"
	ToCommReceivedConfirmationChan chan []int   //"rco"
	ToCommExecutedConfirmationChan chan []int   //"eco"

}
type InternalSlaveChannels struct {
	OrderConfirmedExecutedChan chan []int
	InteruptChan chan os.Signal
}
type InternalMasterChannels struc {
	optimizationInitChan chan Master
	optimizationTriggerChan chan bool
	optimizationReturnChan chan [][]int

}
func Slave_external_chans_init() {
	ExSlaveChans.ToCommSlaveChan = make(chan Slave) 					//"sla"
	ExSlaveChans.ToCommOrderReceivedChan = make(chan []int) 			//"ore"
	ExSlaveChans.ToCommOrderExecutedChan = make(chan []int) 			//"oex"
	ExSlaveChans.ToCommOrderListReceivedChan = make(chan []int) 		//"ocr"
	ExSlaveChans.ToCommOrderConfirmedExecutionChan = make(chan []int)	//"oce"
	ExSlaveChans.ToCommExternalButtonPushedChan = make(chan []int)		//"ebp"
}
func Master_external_chans_init() {
	ExMasterChans.ToCommOrderListChan = make(chan [][]int) 			//"exo"
	ExMasterChans.ToCommReceivedConfirmationChan = make(chan []int) //"rco"
	ExMasterChans.ToCommExecutedConfirmationChan = make(chan []int)	//"eco"
}
func Slave_internal_chans_init(){
	InSlaveChans.OrderConfirmedExecutedChan = make(chan []int)
	InSlaveChans.InteruptChan = make(chan os.Signal,1) //must be buffered see package declaration
}
func Master_internal_chans_init() {
	InMasterChans.OptimizationInitChan = make(chan Master)
	InMasterChans.OptimizationTriggerChan = make(chan []int)
	InMasterChans.OptimizationReturnChan = make(chan [][]int)
}

func Slave_init() {
	buf := make([]byte,1024)
	s Slave{}
	connSend, connReceive := Network_init()
	_, _, err := connReceive.SetReadDeadline(time.Now().Add(Random_init(10, 100) * MilliSecond))
	_, _, err := connReceive.ReadFromUDP(buf) //n contanis numbers of used bytes
	if err != nil { //run master if connection fails
		go Master_init(s)
	} else {
		s Slave{}
		go Select_send_slave(connSend)
		go Select_receive()

		//communicate with statemachine
		go Recive_externalList()

		// ##1 Send order received <- update external list in struct

		// ##2 send floor reached <- get order from statemachine
		// <- receive confirmation on order executed and set lights
		

	}
}

func Master_init(s Slave) {
	connSend, _ := Network_init()
	ExNetChans.connSend <- connSend
	
	//intial sending contianing: ipadress, initialization, 
	//just listen to this
	m Master{s}
	go Select_send_master(c)
	go Select_receive()
	//if call for new optimalization
	//recive it and  iterate slaves laves
	go Interuption_killer()

	go ()

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
	<- InSlaveChans.InteruptChan
	//what should be done when ctrl-c is pressed????? 
	//<- goes here.
	//SystemInit()
	fmt.Println("Got ctrl-c signal")
	Exit(0)
}

func Error() { // handle error here?

}
func (m Master) Master_communication() {
	
	for {
		select {
			//triggers new optimization when new order received
			case order := <- ExCommChans.ToMasterExternalButtonPushed   
				InMasterChans.OptimizationTriggerChan <- order
			//receives new optimized orderList
			case orderList := <- InMasterChans.OptimizationReturnChan
				//send to slaves master
				ExMasterChans.ToCommOrderListChan: <- orderList

			case order := <-ExCommChans.ToMasterOrderExecutedChan://to spesific IP
				ExMasterChans.ToCommExecutedConfirmationChan <- order
			
			//Resopnd on orderList received
			case order := <-ExCommChans.ToMasterOrderListReceivedChan://with spesific IP
				InMasterChans.OrderReceivedMangerChan <- order


			case slave <- ExCommChans.ToMasterSlaveChan:
				m[slave.nr] = slave
		}
	}
}
func (s Slave) Slave_communication() {
	
	for {
		select {
			case slave := <-ExStateChans.DirectionUpdate:
				ExCommChans.ToMasterSlaveChan <- slave
			case slave := <-ExStateChans.CurrentFloorUpdate:
				ExCommChans.ToMasterSlaveChan <- slave

			//checks new button pressed
			case order := <-ExStateChans.ToSlaveExButtonPushedChan:
				ExSlaveChans.ToCommExternalButtonPushedChan <- order
			
			//receives order list
			case orderList := <- ExCommChan.ToSlaveOrderListChan:
				s.externalList = orderList
				//###### set lights on external floor
				ExSlaveChans.ToCommOrderListReceivedChan <- order 
				ExStateChans.NewOrderListChan <- orderList //set lights

			//Triggers if order is executed, sends confirmation further up
			case order := <- ExStateChans.ToSlaveExecutedOrderChan:
				ExSlaveChans.ToCommOrderExecutedChan <- order
		}
	}
}
func Order_received_manager(){
	for {
		order := <-InMasterChans.OrderReceivedManagerChan
		/*
		need to send a confirmation message to the ip that sent this every time we get the order.
		After some time without fun incoming on the channel we can assume that the confirmation has been received.
		*/

	}
}

func (s Slave) Send_slave_to_state() { //send next floor to statemachine
	if s.externalList[s.currentFloor][1] ==1 || s.internalList[s.currentFloor] == 1 {
		ExStateChan.slaveStateChan <- s.currentFloor

	} else if s.direction == 1 { //heading upwards -> can take higher orders
		for i := s.currentFloor; i<N_FLOORS; i++ {
			if s.externalList[i][0] == 1 || s.internalList[i] == 1 { // any orders on higher floors
				ExStateChan.slaveStateChan <- i 
				break
			}
		}
	} else if s.direction == -1 { // heading downwards -> can take lower orders
		for i := 0; i<s.currentFloor; i++ {//any orders on lower floors
			if s.externalList[i][0] == 1 || s.internalList[i] == 1{
				ExSateChan.slaveStateChan <- i 
				break
			} 
		}
	} else {
		Sleep(10 * Millisecond)
	}
}

