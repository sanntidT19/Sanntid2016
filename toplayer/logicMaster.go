package modules

const (
	FLOORS = 4
)

type Slave struct {
	nr int
	internalList []bool
	externalList [][]int
	currentFloor int //get from driver/IO
	direction    int // get from driver/IO

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
	s.externalList = <- commToSlaveOrderChan
	}
}

func Send_confirmation_to_master_or_comm(message string, commToSlaveMastersBackConfirmChan chan string) {
	//transfer to the communication module(put on correct tag)
	commToSlaveMastersBackConfirmChan <- message
}

func (s Slave) Update_current_floor_and_direction(currentFloorChan chan int, directionChan chan int, slaveToCommSlaveStructChan chan int) {
	select {
	case floor := <-currenFloorChan:
		s.currentFloor = floor
	case dir := <-directionChan:
		s.driection = dir
	}


	s.Send_slave_to_master(slaveToCommSlaveStructChan chan Slave)


}

func (s Slave) Send_slave_to_master(slaveToCommSlaveStructChan chan Slave) {
	slaveToCommSlaveStructChan <- s
}

func (s Slave) Send_slave_to_state(slaveStateChan chan int) { //send next floor to statemachine
	var nextFloor int

	//if heading up, pick up others who is going up, both external and internal


	s.externalList[][]

	s.currentFloor

	s.direction

	s.internalList
	//if heading down pick up othes who iss going down, both external and internal



	 slaveStateChan <- nextFloor
}


type Master struct {
	nr int 
	externalList  [][]int
	currentFloors []int
	directions    []int
	internalList  []int
}

func Send_external_list_to_slaves(masterToCommOrderChan chan [][]int) {
	//sends new external list to communication module 
	masterToCommOrderChan <- Get_optimal_externalList()
}
func (m Master) Get_optimal_externalList(newExternalList chan [][]int) [][]int{
	//gets new OptimnalExternalList from module
	m.externalList
	m.nr 
	m.currentFloors
	m.directions
	m.internalList
}
