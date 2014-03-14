package modules

const (
	FLOORS = 4
)

type Slave struct {
	internalList [FLOORS]bool
	externalList [][]int
	currentFloor int //get from driver/IO
	direction    int // get from driver/IO

}

/*
Get orders from the optimalizaton algorithm
*/
func (s Slave) Run_elevator() {

	//communicate with driver

	//run externalList

	//internal list

}

func (s Slave) Recive_orders(incomingOrderChan chan [][]int) {
	select {
	case o := <-incomingOrderChan:
		//write to UDP, let master know
		s.externalList = o

	}
}

func (s Slave) Send_message_to_master(message string, outgoingOrderChan chan string) {
	//transfer to the communication module(put on correct tag)
	outgoingOrderChan <- message
}

func (s Slave) Update_current_floor(currentFloorChan chan int) {
	select {
	case o <- currenFloorChan:
		s.currentFloor = o
	}
}
