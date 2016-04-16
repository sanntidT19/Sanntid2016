package orderAssignment

import (
	. "../globalStructs"
	"../network"
	"fmt"
	"math"
	"time"
)

/*
Functions deciding which elevator in the network should get the order,
based on their states which are also maintained.
*/

type ElevsToAgreeOnAssignedOrder struct {
	OrdAss   OrderAssigned
	ElevList []string
}

var allElevStates []ElevatorState

var unassignedOrders []Order


func GetUnassignedOrders() []Order {
	localCopy := make([]Order, len(unassignedOrders))
	copy(localCopy, unassignedOrders)
	return localCopy
}

//Helper function
func orderIsInList(orderQueue []Order, newOrder Order) bool {
	for _, v := range orderQueue {
		if v == newOrder {
			return true
		}
	}
	return false
}


func removeOrderFromUnassigned(order Order) {
	for i, v := range unassignedOrders {
		if v == order {
			unassignedOrders = append(unassignedOrders[:i], unassignedOrders[i+1:]...)
			return
		}
	}
}


func chooseOptimalElev(newOrder Order) string {
	elevStatesCopy := make([]ElevatorState, len(allElevStates))
	copy(elevStatesCopy, allElevStates)
	numOfElevs := len(elevStatesCopy)

	IPCostList := make([]int, numOfElevs)
	var lowestCost int = 100
	var optimalIP string = "0"
	for i, v := range elevStatesCopy {
		if v.CurrentFloor < newOrder.Floor {
			if v.Direction == DOWN {
				if len(v.OrderQueue) > 0 {
					lastFloorInOrderQueue := v.OrderQueue[len(v.OrderQueue)-1].Floor

					floorsToBeVisited := int(math.Abs(float64(v.CurrentFloor-lastFloorInOrderQueue)) + math.Abs(float64(newOrder.Floor-lastFloorInOrderQueue)))

					IPCostList[i] += floorsToBeVisited
				}
			}
		} else if v.CurrentFloor > newOrder.Floor {
			if v.Direction == UP {
				//Add distance to last order in queue and distance from last order to new order, only if there are orders
				if len(v.OrderQueue) > 0 {
					lastFloorInOrderQueue := v.OrderQueue[len(v.OrderQueue)-1].Floor

					floorsToBeVisited := int(math.Abs(float64(v.CurrentFloor-lastFloorInOrderQueue)) + math.Abs(float64(newOrder.Floor-lastFloorInOrderQueue)))

					IPCostList[i] += floorsToBeVisited
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
	return optimalIP
}


func GetOrderQueueOfDeadElev(deadIP string) []Order {
	for _, v := range allElevStates {
		if v.IP == deadIP {
			listCopy := make([]Order, len(v.OrderQueue))
			copy(listCopy, v.OrderQueue)
			return listCopy
		}
	}
	return nil
}


func UpdateElevatorStateList(newOrderToBeAssignedChan chan Order, resetAssignFuncChan chan bool, sendOrderAssChan chan OrderAssigned, receiveOrderAssChan chan OrderAssigned, receiveElevStateChan chan ElevatorState, newLocalOrderChan chan Order,removeElevChan chan string) {

	go AssignOrdersAndWaitForAgreement(newOrderToBeAssignedChan, resetAssignFuncChan, sendOrderAssChan, receiveOrderAssChan, newLocalOrderChan)
	for {
		select {
		case updatedElevState := <-receiveElevStateChan:
			//fmt.Println("New state received!")
			elevInList := false
			for i, v := range allElevStates {
				if updatedElevState.IP == v.IP {
					var newestStateTime time.Time = allElevStates[i].Timestamp
					if updatedElevState.Timestamp.After(newestStateTime) {
						allElevStates[i] = updatedElevState
						for _, v := range allElevStates[i].OrderQueue {
							if orderIsInList(unassignedOrders, v) {
								removeOrderFromUnassigned(v)
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
		case deadElev := <-removeElevChan:
			for i, v := range allElevStates {
				if v.IP == deadElev {
					allElevStates = append(allElevStates[:i], allElevStates[i+1:]...)
					break
				}
			}
		}
	}
}

func AssignOrdersAndWaitForAgreement(newOrderFromNetworkChan chan Order, resetAssignFuncChan chan bool, sendOrderAssChan chan OrderAssigned, receiveOrderAssChan chan OrderAssigned,newLocalOrderChan chan Order) {
	var OrdersToBeAssignedByAll []ElevsToAgreeOnAssignedOrder
	localAddr := network.FindLocalIP()
	for {
		select {
		case newOrder := <-newOrderFromNetworkChan:
			orderIsRegistered := false
			for _, v := range OrdersToBeAssignedByAll {
				if v.OrdAss.Order == newOrder {
					orderIsRegistered = true
				}
			}
			if !orderIsRegistered {
				assignedElevAddr := chooseOptimalElev(newOrder)
				NewOrderToBeAssigned := OrderAssigned{Order: newOrder, AssignedTo: assignedElevAddr, SentFrom: localAddr}
				elevList := network.ElevsSeen()
				OrdersToBeAssignedByAll = append(OrdersToBeAssignedByAll, ElevsToAgreeOnAssignedOrder{NewOrderToBeAssigned, elevList})
				if !orderIsInList(unassignedOrders, newOrder) {
					unassignedOrders = append(unassignedOrders, newOrder)
				}
				time.Sleep(time.Millisecond * 200) 
				sendOrderAssChan <- NewOrderToBeAssigned
			} else {
				fmt.Println("Order already registered,discard message.")
			}
		case newOrdAss := <-receiveOrderAssChan:
			posInSlice := -1
			for i, v := range OrdersToBeAssignedByAll {
				if newOrdAss.Order == v.OrdAss.Order {
					posInSlice = i
					break
				}
			}
			if posInSlice == -1 {
				fmt.Println("Assignment registered, discard message.")
			} else {

				if newOrdAss.AssignedTo != OrdersToBeAssignedByAll[posInSlice].OrdAss.AssignedTo {
					fmt.Println("Disagreement, recalculate.")
					OrdersToBeAssignedByAll = append(OrdersToBeAssignedByAll[:posInSlice], OrdersToBeAssignedByAll[posInSlice+1:]...)
					assignedElevAddr := chooseOptimalElev(newOrdAss.Order)
					NewOrderToBeAssigned := OrderAssigned{Order: newOrdAss.Order, AssignedTo: assignedElevAddr, SentFrom: localAddr}
					elevList := network.ElevsSeen()
					OrdersToBeAssignedByAll = append(OrdersToBeAssignedByAll, ElevsToAgreeOnAssignedOrder{NewOrderToBeAssigned, elevList})
					time.Sleep(time.Millisecond * 200) 
					sendOrderAssChan <- NewOrderToBeAssigned
				} else {
					for i, v := range OrdersToBeAssignedByAll[posInSlice].ElevList {
						if newOrdAss.SentFrom == v {
							OrdersToBeAssignedByAll[posInSlice].ElevList = append(OrdersToBeAssignedByAll[posInSlice].ElevList[:i], OrdersToBeAssignedByAll[posInSlice].ElevList[i+1:]...)
							if len(OrdersToBeAssignedByAll[posInSlice].ElevList) == 0 {
								if newOrdAss.AssignedTo == localAddr {
									newLocalOrderChan <- newOrdAss.Order
								}
								OrdersToBeAssignedByAll = append(OrdersToBeAssignedByAll[:posInSlice], OrdersToBeAssignedByAll[posInSlice+1:]...)
							}
						}
					}
				}
				
			}
		case <-resetAssignFuncChan:
			OrdersToBeAssignedByAll = nil
		}
	}
}
