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
func Slave_external_chans_init() {
	ExSlaveChans.ToCommSlaveChan = make(chan Slave) //"sla"
	ExSlaveChans.ToCommOrderReceivedChan = make(chan []int) //"ore"
	ExSlaveChans.ToCommOrderExecutedChan = make(chan []int) //"oex"
	ExSlaveChans.ToCommOrderListReceivedChan = make(chan []int) //"ocr"
	ExSlaveChans.ToCommOrderConfirmedExecutionChan = make(chan []int) //"oce"
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
func Send_to_slave() {
	ExMasterChans.ToCommOrderListChan 
	ExMasterChans.ToCommReceivedConfirmationChan 
	ExMasterChans.ToCommExecutedConfirmationChan 
	for {
		select {
			case order := <-ExCommChans.ToSlaveReceivedOrderListConfirmationChan:

			case order := <-ExCommChans.ToSlaveExecutedConfirmationChan:
				ExToCommChans.ToSlaveExecutedConfirmationCHan
			case order := <-ExCommChans.ToSlaveOrderListChan:
				ExToCommChans.ToSlaveOrderListChan <- order
		}
	}
}
func Send_confirmation_to_master() {
	//when recieved from state machine
	for {
		select {
			case order := <-ExStateChans.ToSlaveExButtonPushedChan
				ExSlaveChans.ToCommExternalButtonPushedChan <- order

			case order := <- ExStateChans.ToSlaveOrderListReceivedChan:
				ExSlaveChans.ToCommOrderListReceivedChan <- order 

			case order := <- ExStateChans.ToSlaveExecutedOrderChan:
				ExSlaveChans.ToCommOrderExecutedChan <- order
				////internl chans:
			//case order := : 
			//	order <- ExSlaveChans.ToCommOrderListReceivedChan 
			case order <- InSlaveChans.OrderConfirmedExecutedChan: 
				ExSlaveChans.ToCommOrderConfirmedExecutionChan <- order

		}
	}
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

