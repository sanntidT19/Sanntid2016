package optalg

import (
	. "../globalChans"
	. "../globalStructs"
	"../stateMachine"
	"fmt"
	"math"
)

/*0
,State{MyIP: "123.123.123.123",
	CurrentFloor: 3,
	LastFloor:    2,
	Direction:    -1}
var el_state2 ElevatorState = ElevatorState{MyIP: "123.123.123.124",
	CurrentFloor: 2,
	LastFloor:    1,
	Direction:    1,
	OrderQueue:       []Button{Button{3, -1, false}, {2, 0, false}}}

var el_state3 ElevatorState = ElevatorState{MyIP: "123.123.123.125",
	CurrentFloor: 3,
	LastFloor:    3,
	Direction:    1,
	OrderQueue:       []Button{Button{2, -1, false}}}
*/
var allElevStates []ElevatorState

func OptAlg(newOrder Order) string {
	elevStatesCopy := make([]ElevatorState, len(allElevStates))
	copy(elevStatesCopy, allElevStates)
	numOfElevs := len(elevStatesCopy)

	IPCostList := make([]int, numOfElevs)
	var elevNmr int = -1
	var lowestCost int = 100
	var optimalIP string = "0"
	fmt.Println("optalg: number of elevs seen: ", numOfElevs)
	for i, v := range elevStatesCopy {
			if v.CurrentFloor < newOrder.Floor {
				if v.Direction == DOWN {
					if len(v.OrderQueue) > 0 {
						lastFloorInOrderQueue := v.OrderQueue[len(v.OrderQueue)-1].Floor

						floorsToBeVisited := int(math.Abs(float64(v.CurrentFloor-lastFloorInOrderQueue)) + math.Abs(float64(newOrder.Floor-lastFloorInOrderQueue)))

						IPCostList[i] += floorsToBeVisited
						fmt.Println("LENGTH OF QUEUE IS LONGER THAN 0")
					}

				}
			} else if v.CurrentFloor > newOrder.Floor {
				if v.Direction == UP {
					//Add distance to last order in queue and distance from last order to new order, only if there are orders
					if len(v.OrderQueue) > 0 {
						lastFloorInOrderQueue := v.OrderQueue[len(v.OrderQueue)-1].Floor

						floorsToBeVisited := int(math.Abs(float64(v.CurrentFloor-lastFloorInOrderQueue)) + math.Abs(float64(newOrder.Floor-lastFloorInOrderQueue)))

						IPCostList[i] += floorsToBeVisited
						fmt.Println("LENGTH OF QUEUE IS LONGER THAN 0")
					}
				}
			}
		floatDifference := float64(v.CurrentFloor - newOrder.Floor)

		IPCostList[i] += int(math.Abs(floatDifference))
		IPCostList[i] += len(v.OrderQueue)

	}
	for k := 0; k < len(IPCostList); k++ {
		if IPCostList[k] < lowestCost {
			optimalIP = elevStatesCopy[k].IP
			lowestCost = IPCostList[k]
			elevNmr = k
		} else if IPCostList[k] == lowestCost {
				if elevStatesCopy[k].IP > optimalIP {
					optimalIP = elevStatesCopy[k].IP
					elevNmr = k
				}
			
		}
	}

	fmt.Println("For this order: ")
	stateMachine.PrintOrder(newOrder)
	fmt.Println("My choice: ", optimalIP)
	return optimalIP
}

/*
func main() {
	fmt.Print("optimal IP er:", opt_alg(new_order))
}
*/
func GetOrderQueueOfDeadElev(deadIP string) []Order {
	for _, v := range allElevStates {
		if v.IP == deadIP {
			listCopy := make([]Order, len(v.OrderQueue))
			copy(listCopy, v.OrderQueue)
			return listCopy
		}
	}
	fmt.Println("Elevator not found")
	return nil
}

//may not need channels, think about if its better to just call it from somewhere else
func UpdateElevatorStateList() {
	for {
		select {
		case updatedElevState := <-FromNetworkNewElevStateChan:
			//fmt.Println("New state received!")
			elevInList := false
			for i, v := range allElevStates {
				if updatedElevState.IP == v.IP {
					allElevStates[i] = updatedElevState
					elevInList = true
					break
				}
			}
			if !elevInList {
				allElevStates = append(allElevStates, updatedElevState)
			}
		case elevatorTakesOrder := <-AddOrderAssignedToElevStateChan:
			for _, v := range allElevStates {
				if elevatorTakesOrder.AssignedTo == v.IP {
					v.OrderQueue = append(v.OrderQueue, elevatorTakesOrder.Order)
				}
			}
		case deadElev := <-ToOptAlgDeleteElevChan:
			for i, v := range allElevStates {
				if v.IP == deadElev {
					allElevStates = append(allElevStates[:i], allElevStates[i+1:]...)
					break
				}
			}
		}
	}
}
