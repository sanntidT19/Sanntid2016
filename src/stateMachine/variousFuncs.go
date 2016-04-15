package stateMachine

import (
	"../driver"
	."../globalStructs"
	"../network"
	"fmt"
	"time"
)

//up = 1, down = -1, stop = 0
//We need better names for these to avoid confusion. The whole floor-thing needs to be thought through
var current_floor int = -1
var last_floor int = -1
var desired_floor int = -1
var door_open bool = false

// For new suggestion
var StateOfElev ElevatorState

func PrintOrder(order Order) {
	var x string
	if order.Direction == COMMAND {
		x = "Command"
	} else if order.Direction == UP {
		x = "Up"
	} else {
		x = "Down"
	}
	fmt.Println("Floor: ", order.Floor+1, " Type :", x)
}

//UTKAST TIL REFORMATERT HEIS ER UNDER HER. KOM MED GODE NAVN PÅ TING ALLEREDE NÅ
func MoveElevatorAndOpenDoor(floorWithOrderReachedChan chan Order, orderServedChan chan Order, sendElevInDirectionChan chan int, sendElevStateChan chan ElevatorState) {
	//CAN BE: IDLE, OPEN_DOOR, GO_TO_FLOOR (evt can G_T_F be go_up, and go_down)
	for {
		select {
		case newestOrder := <-floorWithOrderReachedChan:
			if newestOrder.Direction != COMMAND {
				StateOfElev.Direction = newestOrder.Direction
			}
			driver.DriveElevator(0)
			driver.OpenDoor()
			//After door is closed. Tell someone above
			orderServedChan <- newestOrder
		case goInDirection := <-sendElevInDirectionChan:
			if goInDirection == UP {
				driver.DriveElevator(UP)
				StateOfElev.Direction = UP
			} else {
				driver.DriveElevator(DOWN)
				StateOfElev.Direction = DOWN
			}
			SendElevStateToNetwork(sendElevStateChan)
			//go to floor, if below, go up. if above, go down. This is only needed when the elevator is idle.
		}
	}
}

//
func UpdateNewFloorReached(sendElevStateChan chan ElevatorState) {
	mostRecentFloorVisited := driver.GetFloorSensorSignal()
	for {
		if sensor_result := driver.GetFloorSensorSignal(); sensor_result != -1 && sensor_result != mostRecentFloorVisited {
			mostRecentFloorVisited = sensor_result
			StateOfElev.PreviousFloor = StateOfElev.CurrentFloor
			StateOfElev.CurrentFloor = sensor_result
			SendElevStateToNetwork(sendElevStateChan)
		}
		time.Sleep(time.Millisecond * 50)
	}
}

//split into two: need to control shared var stopSignalSent
func feedDirectionCommandsToElev(orderQueueChangeChan chan bool, sendElevInDirectionChan chan int, targetFloorReachedChan chan Order) {
	//State of elev is known globally for now
	//Send til hit når ordre er served også
	stopSignalSent := false
	go func() {
		for {
			if len(StateOfElev.OrderQueue) > 0 {
				if StateOfElev.OrderQueue[0].Floor == driver.GetFloorSensorSignal() {
					if !stopSignalSent {
						targetFloorReachedChan <- StateOfElev.OrderQueue[0]
						stopSignalSent = true
					}
				}
			}
			time.Sleep(time.Millisecond * 50)
		}
	}()
	for {
		<-orderQueueChangeChan
		stopSignalSent = false
		if len(StateOfElev.OrderQueue) > 0 {
			if driver.ElevNotMoving() {
				fmt.Println("current floor state", StateOfElev.CurrentFloor)
				fmt.Println("floor of order", StateOfElev.OrderQueue[0].Floor)
				if StateOfElev.CurrentFloor > StateOfElev.OrderQueue[0].Floor {
					sendElevInDirectionChan <- DOWN
				} else if StateOfElev.CurrentFloor < StateOfElev.OrderQueue[0].Floor {
					sendElevInDirectionChan <- UP
				}
			}
		}
	}
}

func NewTopLoop(sendElevStateChan chan ElevatorState,newLocalOrderChan chan Order, localOrderServedChan chan Order) {
	//Her kan man anta at man står stille i en etasje
	StateOfElev.CurrentFloor = driver.GetFloorSensorSignal()
	StateOfElev.Direction = DOWN
	StateOfElev.IP = network.FindLocalIP()
	SendElevStateToNetwork(sendElevStateChan)

	targetFloorReachedChan := make(chan Order)
	orderServedChan := make(chan Order)
	sendElevInDirectionChan := make(chan int)
	orderQueueChangeChan := make(chan bool)
	//go SpamCurrentQueue()
	go UpdateNewFloorReached(sendElevStateChan)
	go MoveElevatorAndOpenDoor(targetFloorReachedChan, orderServedChan, sendElevInDirectionChan,sendElevStateChan)
	go feedDirectionCommandsToElev(orderQueueChangeChan, sendElevInDirectionChan, targetFloorReachedChan)
	fmt.Println("Statemachine: ready")
	for {
		select {
		//Might need to put the first case in its on goroutine
		//Send up about change
		case servedOrder := <-orderServedChan:
			fmt.Println("This order is served")
			PrintOrder(servedOrder)
			indexOfServedOrder := -1
			for i, v := range StateOfElev.OrderQueue {
				if v == servedOrder {
					indexOfServedOrder = i
				}

			}
			if indexOfServedOrder == -1 {
				fmt.Println("Error! Couldnt find served order in local queue")
			} else {
				StateOfElev.OrderQueue = append(StateOfElev.OrderQueue[:indexOfServedOrder], StateOfElev.OrderQueue[indexOfServedOrder+1:]...)

				orderQueueChangeChan <- true

				localOrderServedChan <- servedOrder

				fmt.Println("order served")
			}
			SendElevStateToNetwork(sendElevStateChan)
		case newOrder := <-newLocalOrderChan: //newOrderToElevChan when we have mothership ready
			if orderIsInList(StateOfElev.OrderQueue, newOrder) {
				break
			} else {
				StateOfElev.OrderQueue = insertOrderIntoQueue(newOrder, StateOfElev)
				orderQueueChangeChan <- true
			}
			SendElevStateToNetwork(sendElevStateChan)
		}
	}

}

func SendElevStateToNetwork(sendElevStateChan chan ElevatorState) {
	copyOfElevState := StateOfElev
	copyOfElevState.Timestamp = time.Now()
	copyOfOrderList := make([]Order, len(StateOfElev.OrderQueue))
	copy(copyOfOrderList, StateOfElev.OrderQueue)
	copyOfElevState.OrderQueue = copyOfOrderList
	sendElevStateChan <- copyOfElevState
}

func SpamCurrentQueue() {
	for {
		for _, v := range StateOfElev.OrderQueue {
			PrintOrder(v)
		}
		fmt.Println()
		time.Sleep(time.Millisecond * 200)
	}
}

func insertOrderIntoQueue(newOrder Order, currentState ElevatorState) []Order {
	common_current_order_queue := currentState.OrderQueue
	orderQueueCopy := make([]Order, len(common_current_order_queue))
	copy(orderQueueCopy, common_current_order_queue)
	var placeInQueue int
	//if queue is empty
	if len(orderQueueCopy) == 0 {
		newOrderQueue := []Order{newOrder}
		return newOrderQueue
	} else if (newOrder.Direction == currentState.Direction || newOrder.Direction == COMMAND) && newOrder.Floor == driver.GetFloorSensorSignal() {
		//check if you are in the actual floor in correct direction and should stop immediately
		placeInQueue = 0
	} else {
		//Perform a simulation of what the elevator will do with current queue and find out where the new order belongs
		currentDirection := currentState.Direction
		firstFloorToVisit := currentState.CurrentFloor
		if driver.GetFloorSensorSignal() == -1 {
			fmt.Println("Elevator not in current floor, simulate past current floor")
			firstFloorToVisit = currentState.CurrentFloor + currentDirection //er dette riktig?
		}
		placeInQueue = SimulateElevDrivingFindOrderIndex(firstFloorToVisit, currentDirection, orderQueueCopy, newOrder)
	}
	newOrderQueue := append(orderQueueCopy[:placeInQueue], append([]Order{newOrder}, orderQueueCopy[placeInQueue:]...)...)
	return newOrderQueue
}

/*
Simulate driving the elevator from current floor and serve all possible orders.
*/
func SimulateElevDrivingFindOrderIndex(startingFloor int, startingDirection int, orderQueue []Order, newOrder Order) int {
	simulatedFloor := startingFloor
	simulatedDirection := startingDirection
	placeInQueue := 0
	i := 0
	for simulatedFloor >= 0 && simulatedFloor < NUM_FLOORS {
		if newOrder.Floor == simulatedFloor && (simulatedDirection == newOrder.Direction || newOrder.Direction == COMMAND) {
			placeInQueue = i
			return placeInQueue
		} else if orderQueue[i].Floor == simulatedFloor && (simulatedDirection == orderQueue[i].Direction || orderQueue[i].Direction == COMMAND) {
			i++
			if i == len(orderQueue) {
				placeInQueue = i
				return placeInQueue
			}
		}
		simulatedFloor += simulatedDirection
	}

	simulatedFloor -= simulatedDirection
	if simulatedDirection == UP {
		simulatedDirection = DOWN
	} else {
		simulatedDirection = UP
	}
	for simulatedFloor >= 0 && simulatedFloor < NUM_FLOORS {
		if newOrder.Floor == simulatedFloor && (simulatedDirection == newOrder.Direction || newOrder.Direction == COMMAND) {
			placeInQueue = i
			return placeInQueue
		} else if orderQueue[i].Floor == simulatedFloor && (simulatedDirection == orderQueue[i].Direction || orderQueue[i].Direction == COMMAND) {
			i++
			if i == len(orderQueue) {
				placeInQueue = i
				return placeInQueue
			}
		}
		simulatedFloor += simulatedDirection
	}

	simulatedFloor -= simulatedDirection //bounds were exceeded and you take one step back

	if simulatedDirection == UP {
		simulatedDirection = DOWN
	} else {
		simulatedDirection = UP
	}

	for simulatedFloor != (startingFloor + simulatedDirection) {
		if newOrder.Floor == simulatedFloor && (simulatedDirection == newOrder.Direction || newOrder.Direction == COMMAND) {
			placeInQueue = i
			return placeInQueue
		} else if orderQueue[i].Floor == simulatedFloor && (simulatedDirection == orderQueue[i].Direction || orderQueue[i].Direction == COMMAND) {
			//order was theoretically served and next element in the current queue is up for evaluation
			i++
			if i == len(orderQueue) {
				//all orders in queue evaluated. put new order at end
				placeInQueue = i
			}
		}
		simulatedFloor += simulatedDirection
	}
	return placeInQueue
}


func orderIsInList(orderQueue []Order, newOrder Order) bool {
	for _, v := range orderQueue {
		if v == newOrder {
			return true
		}
	}
	return false
}