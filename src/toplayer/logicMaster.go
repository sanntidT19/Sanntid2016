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
func Slave_external_chans_init() {
	ExSlaveChans.ToCommSlaveChan = make(chan Slave) //"sla"
	ExSlaveChans.ToCommOrderReceivedChan = make(chan []int) //"ore"
	ExSlaveChans.ToCommOrderExecutedChan = make(chan []int) //"oex"
	ExSlaveChans.ToCommOrderConfirmedReceivedChan = make(chan []int) //"ocr"
	ExSlaveChans.ToCommOrderConfirmedExecutuinChan = make(chan []int) //"oce"
	ExSlaveChans.ToCommExternalButtonPushedChan = make(chan []int)	//"ebp"
}
func Master_external_chans_init() {
	ExMasterChans.ToCommOrderListChan = make(chan [][]int) //"exo"
	ExMasterChans.ToCommReceivedConfirmationChan = make(chan []int)   //"rco"
	ExMasterChans.ToCommExecutedConfirmationChan = make(chan []int)   //"eco"
}
func Slave_internal_chans_init(){
	InSlaveChans.OrderConfirmedExecutedChan = make(chan []int)
	InSlaveChans.InteruptChan = make(chan os.Signal,1) //must be buffered see package declaration
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
	}
func Interuption_killer() {
	signal.Notify(InSlaveChans.InteruptChan, os.Interrupt)
	signal.Notify(InSlaveChans.InteruptChan, syscall.SIGTERM)
	<- InSlaveChans.InteruptChan
	//what should be done when ctrl-c is pressed????? 
	//<- goes here.
	fmt.Println("Got ctrl-c signal")
	Exit(0)
}

func Error() { // handle error here?

}

func (s Slave) Recive_externalList() {
	for {
	s.externalList = <- ExComExMasterChanshans.ToSlaveOrderListChan 
	}
}

func Send_confirmation_to_master(message string) {
	//when recieved from state machine
	for {
		select {
			case order := <-ExStateChans.ToSlaveExButtonPushedChan
				ExSlaveChans.ToCommExternalButtonPushedChan <- order
			//case order := //?????????????
			//	ExSlaveChans.ToCommOrderReceivedChan <- order
			case order := <- ExStateChans.ToSlaveExecutedOrderChan:
				ExSlaveChans.ToCommOrderExecutedChan <- order
				////internl chans:
			//case order := : 
			//	order <- ExSlaveChans.ToCommOrderConfirmedReceivedChan 
			case order <- InSlaveChans.OrderConfirmedExecutedChan: 
				ExSlaveChans.ToCommOrderConfirmedExecutuinChan <- order

		}
	}
}
func Get_order_executed_from_slave() []int {
	order := <- ExCommChans.ToMasterOrderReceivedChan
	Send_order_executed_confirmation_to_slave(order)
	//delete from externalList
	return order
}
func Send_order_executed_to_master(order []int) {
	ExSlaveChans.ToCommOrderReceivedChan <- order
}

func Send_order_executed_confirmation_to_slave(order) {
	ExCommChans.ToSlaveConfirmedExecutionChan <- order
}
func Get_order_confirmation_from_master() {
	order := <- ExSlaveChans.ToSlaveConfirmedExecutionChan
	//turn off ligths

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
func (m Master) Get_slave_from_slave() {

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
