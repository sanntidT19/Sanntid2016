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
func (s Slave) runElevator() {

	//communicate with driver

}
func (s Slave) alertMaster(alertMasterChan chan string) {
	select {
	case <-alertMasterChan:
		//write to UDP, let master know
	}

}

func reciveOrders(incomingOrderChan chan [][]int) {
	select {
	case <-incomingOrderChan:
		//write to UDP, let master know
	}
}
