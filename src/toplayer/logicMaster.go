package toplayer

//export GOPATH =$HOME/Elevator-progress-/
import(
	."time"
)
const (
	N_FLOORS = 4
)


//Global chans
var ExSlaveChans ExternalSlaveChannels
var ExMasterChans ExternalMasterChannels

type Slave struct {
	nr int
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
	ToCommOrderConfirmedReceivedChan  chan []int //"ocr"
	ToCommOrderConfirmedExecutuinChan chan []int //"oce"
}
type ExternalMasterChannels struct {
	ToCommOrderListChan            chan [][]int //"exo"
	ToCommReceivedConfirmationChan chan []int   //"rco"
	ToCommExecutedConfirmationChan chan []int   //"eco"
}
func Slave_chans_init() {
	ExSlaveChans.ToCommSlaveChan = make(chan Slave) //"sla"
	ExSlaveChans.ToCommOrderReceivedChan = make(chan []int) //"ore"
	ExSlaveChans.ToCommOrderExecutedChan = make(chan []int) //"oex"
	ExSlaveChans.ToCommOrderConfirmedReceivedChan = make(chan []int) //"ocr"
	ExSlaveChans.ToCommOrderConfirmedExecutuinChan = make(chan []int) //"oce"
}
func Master_chans_init() {
	ExMasterChans.ToCommOrderListChan = make(chan [][]int) //"exo"
	ExMasterChans.ToCommReceivedConfirmationChan = make(chan []int)   //"rco"
	ExMasterChans.ToCommExecutedConfirmationChan = make(chan []int)   //"eco"
}

func Slave_init() {
	connSend, connReceive := Network_init()
	.SetReadDeadline(time.Now().Add(Random_init(10, 100) * MilliSecond))
	
	s Slave{}
	go Select_send()
	go Select_receive()


	//communicate with statemachine
	go Recive_externalList()

}

func Master_init() {

	//intial sending contianing: ipadress, initialization, 
	//just listen to this
	m Master{}
	go Select_send()
	go Select_receive()
	//if call for new optimalization
	//send it futher

	//recive it and  iterate slaves slaves
	go Send_im_master()
}

func (s Slave) Recive_externalList() {
	for {
	s.externalList = <- ExComExMasterChanshans.ToSlaveOrderListChan 
	}
}

func Send_confirmation_to_master(message string) {
	//when recieved from state machine

	<- ExSlaveChans.ToCommOrderReceivedChan 
	<- ExSlaveChans.ToCommOrderExecutedChan 
	<- ExSlaveChans.ToCommOrderConfirmedReceivedChan 
	<- ExSlaveChans.ToCommOrderConfirmedExecutuinChan 

	
}

func (s Slave) Update_current_floor_and_direction(currentFloorChan chan int, directionChan chan int, ToCommSlaveStructChan chan int) {
	/*
	select {
	case floor := <-currenFloorChan:
		s.currentFloor = floor
	case dir := <-directionChan:
		s.driection = dir
	}

	s.Send_slave_to_master()
*/

}

func (s Slave) Send_slave_to_master() {
	ExSlaveChans.ToCommSlaveChan <- s
}

func (s Slave) Send_slave_to_state(slaveStateChan chan int) { //send next floor to statemachine
	if s.externalList[s.currentFloor][1] ==1 || s.internalList[s.currentFloor] == 1 {
		//slaveStateChan <- s.currentFloor

	} else if s.direction == 1 { //heading upwards -> can take higher orders
		for i := s.currentFloor; i<N_FLOORS; i++ {
			if s.externalList[i][0] == 1 || s.internalList[i] == 1 { // any orders on higher floors
				slaveStateChan <- i 
				break
			}
		}
	} else if s.direction == -1 { // heading downwards -> can take lower orders
		for i := 0; i<s.currentFloor; i++ {//any orders on lower floors
			if s.externalList[i][0] == 1 || s.internalList[i] == 1{
				slaveStateChan <- i 
				break
			} 
		}
	} else {
		Sleep(10 * Millisecond)
	}
}



func Distribute_orders() {
	//sends new external list to communication module 
	ExMasterChans.ToCommOrderChan <- Get_optimal_externalList()
}
func (m Master) Get_optimal_externalList() {
	//newExternalList := updateExternalList()
	//m.externalList = newExternalList 
	//ExMasterChans.ToCommOrderChan <- newExternalList
}
