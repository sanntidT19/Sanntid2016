package modules

const (
	FLOORS = 4
)


//Global chans
var SC SlaveChans
var MC MasterChan

type Slave struct {
	nr int
	internalList []int
	externalList [][]int
	currentFloor int //get from driver/IO
	direction    int // get from driver/IO

}

type Master struct {
	s []Slaves
}
type SlaveChans {
	ToCommSlaveChan                   chan Slave //"sla"
	ToCommOrderReceivedChan           chan []int //"ore"
	ToCommOrderExecutedChan           chan []int //"oex"
	ToCommOrderConfirmedReceivedChan  chan []int //"ocr"
	ToCommOrderConfirmedExecutuinChan chan []int //"oce"
}
type MasterChans {
	ToCommOrderListChan            chan [][]int //"exo"
	ToCommImMasterChan             chan string  //"iam"
	ToCommReceivedConfirmationChan chan []int   //"rco"
	ToCommExecutedConfirmationChan chan []int   //"eco"
}
func Slave_chans_init() {
	SC.ToCommSlaveChan := make(chan Slave) //"sla"
	SC.ToCommOrderReceivedChan := make(chan []int) //"ore"
	SC.ToCommOrderExecutedChan := make(chan []int) //"oex"
	SC.ToCommOrderConfirmedReceivedChan := make(chan []int) //"ocr"
	SC.ToCommOrderConfirmedExecutuinChan := make(chan []int) //"oce"
}
func Master_chans_init() {
	MC.ToCommOrderListChan := make(chan [][]int) //"exo"
	MC.ToCommImMasterChan := make(chan string)  //"iam"
	MC.ToCommReceivedConfirmationChan := make(chan []int)   //"rco"
	MC.ToCommExecutedConfirmationChan := make(chan []int)   //"eco"
}

/*
Get orders from the optimalizaton algorithm
*/
func Slave_init() {

	//communicate with driver

	//run externalList

	//internal list

}
func Master_init() {

}

func (s Slave) Recive_externalList_from(externalListChan chan [][]int) {
	s.externalList = <- CO.ToSlaveOrderChan
	}
}

func Send_confirmation_to_master_or_comm(message string) {
	//transfer to the communication module(put on correct tag)
	CC.ToSlaveMastersBackConfirmChan <- message
}

func (s Slave) Update_current_floor_and_direction(currentFloorChan chan int, directionChan chan int, ToCommSlaveStructChan chan int) {
	select {
	case floor := <-currenFloorChan:
		s.currentFloor = floor
	case dir := <-directionChan:
		s.driection = dir
	}


	s.Send_slave_to_master(SC.ToCommSlaveStructChan chan Slave)


}

func (s Slave) Send_slave_to_master() {
	SC.ToCommSlaveStructChan <- s
}

func (s Slave) Send_slave_to_state(slaveStateChan chan int) { //send next floor to statemachine
	if (s.externalList[currentFlorr][1] || s.internalList[currentFloor]) == 1 {
		slaveStateChan <- s.currentFloor

	} else if s.dir == 1 { //heading upwards -> can take higher orders
		for i := s.currentFloor; i<driver.N_FLOORS; i++ {
			if (s.externalList[i][0] || s.internalList[i]) == 1{ // any orders on higher floors
				slaveStateChan <- i 
				break
		}
	} else if s.dir == -1 { // heading downwards -> can take lower orders
		for i := 0; i<s.currentFloor; i++ {//any orders on lower floors
			if (s.externalList[i][0] || s.internalList[i]) == 1{
				slaveStateChan <- i 
				break
			} 
		}
	} else {
		Sleep(10 * Millisecound)
	}
}



func Send_external_list_to_slaves() {
	//sends new external list to communication module 
	MC.ToCommOrderChan <- Get_optimal_externalList()
}
func (m Master) Get_optimal_externalList() {
	newExternalList := updateExternalList()
	m.externalList = newExternalList 
	MC.ToCommOrderChan <- newExternalList
}
