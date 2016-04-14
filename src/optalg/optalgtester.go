package optalg

import (
	. "../globalChans"
	. "../globalStructs"
	"../network"
	"../stateMachine"
	"fmt"
	"math"
	"time"
)

type ElevsToAgreeOnAssignedOrder struct {
	OrdAss   OrderAssigned
	ElevList []string
}

var allElevStates []ElevatorState

var unassignedOrders []Order

//INSIDE=GOOD
func GetUnassignedOrders() []Order {
	localCopy := make([]Order, len(unassignedOrders))
	copy(localCopy, unassignedOrders)
	return localCopy
}

//Har denne tre steder, gjør den global
func orderIsInQueue(orderQueue []Order, newOrder Order) bool {
	for _, v := range orderQueue {
		if v == newOrder {
			return true
		}
	}
	return false
}

//INSIDENAMING=GOOD
func removeOrder(orderList []Order, order Order) {
	for i, v := range orderList {
		if v == order {
			orderList = append(orderList[:i], orderList[i+1:]...)
			return
		}
	}
}

//INSIDENAMING=GOOD
func chooseOptimalElev(newOrder Order) string {
	elevStatesCopy := make([]ElevatorState, len(allElevStates))
	copy(elevStatesCopy, allElevStates)
	numOfElevs := len(elevStatesCopy)

	IPCostList := make([]int, numOfElevs)
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
		IPCostList[i] += int(math.Abs(float64(v.CurrentFloor - newOrder.Floor)))
		IPCostList[i] += len(v.OrderQueue) * 2
	}
	for k := 0; k < len(IPCostList); k++ {
		if IPCostList[k] < lowestCost {
			optimalIP = elevStatesCopy[k].IP
			lowestCost = IPCostList[k]
		} else if IPCostList[k] == lowestCost {
			if elevStatesCopy[k].IP > optimalIP {
				optimalIP = elevStatesCopy[k].IP
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
//INSIDENAMING=GOOD
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
//INSIDENAMING=GOOD
func UpdateElevatorStateList(newOrderToBeAssignedChan chan Order, resetAssignFuncChan chan bool) {

	go AssignOrdersAndWaitForAgreement(newOrderToBeAssignedChan, resetAssignFuncChan)
	for {
		select {
		case updatedElevState := <-FromNetworkNewElevStateChan:
			//fmt.Println("New state received!")
			elevInList := false
			for i, v := range allElevStates {
				if updatedElevState.IP == v.IP {
					var newestStateTime time.Time = allElevStates[i].Timestamp
					if updatedElevState.Timestamp.After(newestStateTime) {
						allElevStates[i] = updatedElevState
						for _, v := range allElevStates[i].OrderQueue {
							if orderIsInQueue(unassignedOrders, v) {
								removeOrder(unassignedOrders, v)
							}
						}
					}
					elevInList = true
					break
				}
			}
			if !elevInList {
				allElevStates = append(allElevStates, updatedElevState)
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

//INSIDENAMING=GOOD
func AssignOrdersAndWaitForAgreement(newOrderFromNetworkChan chan Order, resetAssignFuncChan chan bool) {

	var OrdersToBeAssignedByAll []ElevsToAgreeOnAssignedOrder
	localAddr := network.FindLocalIP()

	for {
		select {
		case newOrder := <-newOrderFromNetworkChan:
			fmt.Println("Order received from above:")
			orderIsRegistered := false
			fmt.Println("Iterating..")
			for _, v := range OrdersToBeAssignedByAll {
				if v.OrdAss.Order == newOrder {
					orderIsRegistered = true
				}
			}
			fmt.Println("done searching for order")
			if !orderIsRegistered {
				//ad external order to queue here
				//for now, the elevator states are all globally known. May send copy or something else later.
				assignedElevAddr := chooseOptimalElev(newOrder)
				fmt.Println("optalg complete")
				NewOrderToBeAssigned := OrderAssigned{Order: newOrder, AssignedTo: assignedElevAddr, SentFrom: localAddr}
				//Elevlist should be copied, global or maybe everyone that uses it should be in the same module
				elevList := network.ElevsSeen()
				OrdersToBeAssignedByAll = append(OrdersToBeAssignedByAll, ElevsToAgreeOnAssignedOrder{NewOrderToBeAssigned, elevList})
				if !orderIsInQueue(unassignedOrders, newOrder) {

					unassignedOrders = append(unassignedOrders, newOrder)
				}

				time.Sleep(time.Millisecond * 200) //This is to make sure you get to make the list before
				ToNetworkOrderAssignedToChan <- NewOrderToBeAssigned
			} else {
				fmt.Println("Order already registered. Discard message.")
			}
		case newOrdAss := <-FromNetworkOrderAssignedToChan:
			fmt.Println("newOrdAss := <-FromNetworkOrderAssignedToChan")
			posInSlice := -1
			for i, v := range OrdersToBeAssignedByAll {
				if newOrdAss.Order == v.OrdAss.Order {
					posInSlice = i
					break
				}
			}
			if posInSlice == -1 {
				fmt.Println("Order already assigned, throw awayyyyyyyy")
				stateMachine.PrintOrder(newOrdAss.Order)
			} else {
				stateMachine.PrintOrder(newOrdAss.Order)
				/*
					fmt.Println("                    This order is assigned to: ", newOrdAss.AssignedTo)
					fmt.Println("                            my ip is :", localAddr)
					fmt.Println("                   elevator that decided this is: ", newOrdAss.SentFrom)
				*/
				if newOrdAss.AssignedTo != OrdersToBeAssignedByAll[posInSlice].OrdAss.AssignedTo {
					fmt.Println("                             DISAGREEMENT, RECALCULATE")
					OrdersToBeAssignedByAll = append(OrdersToBeAssignedByAll[:posInSlice], OrdersToBeAssignedByAll[posInSlice+1:]...) //slicetricks
					assignedElevAddr := chooseOptimalElev(newOrdAss.Order)
					NewOrderToBeAssigned := OrderAssigned{Order: newOrdAss.Order, AssignedTo: assignedElevAddr, SentFrom: localAddr}
					elevList := network.ElevsSeen()
					OrdersToBeAssignedByAll = append(OrdersToBeAssignedByAll, ElevsToAgreeOnAssignedOrder{NewOrderToBeAssigned, elevList})
					time.Sleep(time.Millisecond * 200) //This is to make sure you get to make the list before
					ToNetworkOrderAssignedToChan <- NewOrderToBeAssigned
				} else {
					for i, v := range OrdersToBeAssignedByAll[posInSlice].ElevList {
						if newOrdAss.SentFrom == v {
							OrdersToBeAssignedByAll[posInSlice].ElevList = append(OrdersToBeAssignedByAll[posInSlice].ElevList[:i], OrdersToBeAssignedByAll[posInSlice].ElevList[i+1:]...)
							if len(OrdersToBeAssignedByAll[posInSlice].ElevList) == 0 {
								if newOrdAss.AssignedTo == localAddr {
									fmt.Println("X")
									fmt.Println("                                    im taking this order!")
									fmt.Println("X")
									NewOrderToLocalElevChan <- newOrdAss.Order
								}
								OrdersToBeAssignedByAll = append(OrdersToBeAssignedByAll[:posInSlice], OrdersToBeAssignedByAll[posInSlice+1:]...)
							}
						}
					}
				}
				fmt.Println("ENDOF:    newOrdAss := <-FromNetworkOrderAssignedToChan ")
			}
		case <-resetAssignFuncChan:
			fmt.Println("<-networkchangechan")
			OrdersToBeAssignedByAll = nil
		}
	}
}
