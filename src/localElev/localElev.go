package localElev

import (
	"../driver"
	."../globalStructs"
	"../network"
	"fmt"
	"time"
)

/*
Functions controlling the local elevator, handling new orders and serving them
*/


var StateOfElev ElevatorState


func MoveElevAndOpenDoor(floorWithOrderReachedChan chan Order, orderServedChan chan Order, sendElevInDirectionChan chan int, sendElevStateChan chan ElevatorState) {
	for {
		select {
		case newestOrder := <-floorWithOrderReachedChan:
			if newestOrder.Direction != COMMAND {
				StateOfElev.Direction = newestOrder.Direction
			}
			driver.DriveElevator(0)
			driver.OpenDoor()
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

		}
	}
}

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

//Tells which direction the elevator should go, and when to stop
func feedDirectionCommandsToElev(orderQueueChangeChan chan bool, sendElevInDirectionChan chan int, targetFloorReachedChan chan Order) {
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
				if StateOfElev.CurrentFloor > StateOfElev.OrderQueue[0].Floor {
					sendElevInDirectionChan <- DOWN
				} else if StateOfElev.CurrentFloor < StateOfElev.OrderQueue[0].Floor {
					sendElevInDirectionChan <- UP
				}
			}
		}
	}
}

//
func RunElevAndManageOrderQueue(sendElevStateChan chan ElevatorState,newLocalOrderChan chan Order, localOrderServedChan chan Order) {
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
	go MoveElevAndOpenDoor(targetFloorReachedChan, orderServedChan, sendElevInDirectionChan,sendElevStateChan)
	go feedDirectionCommandsToElev(orderQueueChangeChan, sendElevInDirectionChan, targetFloorReachedChan)
	for {
		select {
		case servedOrder := <-orderServedChan:
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
			}
			SendElevStateToNetwork(sendElevStateChan)
		case newOrder := <-newLocalOrderChan: 
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


func insertOrderIntoQueue(newOrder Order, currentState ElevatorState) []Order {
	common_current_order_queue := currentState.OrderQueue
	orderQueueCopy := make([]Order, len(common_current_order_queue))
	copy(orderQueueCopy, common_current_order_queue)
	var placeInQueue int
	if len(orderQueueCopy) == 0 {
		newOrderQueue := []Order{newOrder}
		return newOrderQueue
	} else if (newOrder.Direction == currentState.Direction || newOrder.Direction == COMMAND) && newOrder.Floor == driver.GetFloorSensorSignal() {
		//check if you are in the actual floor in the correct direction and should stop immediately
		placeInQueue = 0
	} else {
		currentDirection := currentState.Direction
		firstFloorToVisit := currentState.CurrentFloor
		if driver.GetFloorSensorSignal() == -1 {
			fmt.Println("Elevator not in current floor, simulate past current floor")
			firstFloorToVisit = currentState.CurrentFloor + currentDirection 
		}
		placeInQueue = SimulateElevDrivingFindOrderIndex(firstFloorToVisit, currentDirection, orderQueueCopy, newOrder)
	}
	newOrderQueue := append(orderQueueCopy[:placeInQueue], append([]Order{newOrder}, orderQueueCopy[placeInQueue:]...)...)
	return newOrderQueue
}

/*
Simulate visiting all the floors in both directions. 
Serve orders in the current queue with the correct direction and floor.
Index of new order is found when it is in the correct direction and floor.
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
	simulatedFloor -= simulatedDirection //bounds were exceeded and you take one step back
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
			i++
			if i == len(orderQueue) {
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