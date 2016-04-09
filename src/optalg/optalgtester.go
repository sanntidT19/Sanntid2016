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
var allElevs []ElevatorState

func Opt_alg(new_order Order) string {
	numOfElevs := len(allElevs)
	all_elevs := make([]ElevatorState, numOfElevs)
	copy(all_elevs, allElevs)

	IP_cost_list := make([]int, numOfElevs)
	Queue_len_list := make([]int, numOfElevs)
	var Ele_nmr int = -1
	var IP_score int = 100
	var Optimal_IP string = "0"
	fmt.Println("optalg: number of elevs seen: ", numOfElevs)
	for i, v := range all_elevs {
		Queue_len_list[i] = len(v.OrderQueue)
		if v.CurrentFloor < new_order.Floor {
			if v.Direction == DOWN {
				IP_cost_list[i] += 1
			}
		} else if v.CurrentFloor > new_order.Floor {
			if v.Direction == UP {
				IP_cost_list[i] += 1
				//Add distance to last order in queue and distance from last order to new order, only if there are orders
				if len(v.OrderQueue) > 0{
					lastFloorInOrderQueue := v.OrderQueue[len(v.OrderQueue) -1].Floor
					
					floorsToBeVisited := int(math.Abs(float64(v.CurrentFloor - lastFloorInOrderQueue)) + math.Abs(float64(new_order.Floor - lastFloorInOrderQueue)))


					IP_cost_list[i] += floorsToBeVisited
				}
			}
		}
		float_difference := float64(v.CurrentFloor - new_order.Floor)
		IP_cost_list[i] += int(math.Abs(float_difference))
		IP_cost_list[i] += len(v.OrderQueue)
	}
	for k := 0; k < len(IP_cost_list); k += 1 {
		if IP_cost_list[k] < IP_score {
			Optimal_IP = all_elevs[k].IP
			IP_score = IP_cost_list[k]
			Ele_nmr = k
		} else if IP_cost_list[k] == IP_score {
			if len(all_elevs[k].OrderQueue) < len(all_elevs[Ele_nmr].OrderQueue) {
				Optimal_IP = all_elevs[k].IP
				Ele_nmr = k
			} else if len(all_elevs[k].OrderQueue) == len(all_elevs[Ele_nmr].OrderQueue) {
				if all_elevs[k].IP > Optimal_IP {
					Optimal_IP = all_elevs[k].IP
					Ele_nmr = k
				}
			}
		}
	}
	fmt.Println("For this order: ")
	stateMachine.PrintOrder(new_order)
	fmt.Println("My choice: ", Optimal_IP)
	return Optimal_IP
}

/*
func main() {
	fmt.Print("optimal IP er:", opt_alg(new_order))
}
*/
func GetOrderQueueOfDeadElev(deadIP string) []Order {
	for _, v := range allElevs {
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
			fmt.Println("New state received!")
			elevInList := false
			for i, v := range allElevs {
				if updatedElevState.IP == v.IP {
					allElevs[i] = updatedElevState
					elevInList = true
					break
				}
			}
			if !elevInList {
				allElevs = append(allElevs, updatedElevState)
			}
		case elevatorTakesOrder := <-AddOrderAssignedToElevStateChan:
			for _, v := range allElevs {
				if elevatorTakesOrder.AssignedTo == v.IP {
					v.OrderQueue = append(v.OrderQueue, elevatorTakesOrder.Order)
				}
			}
		case deadElev := <-ToOptAlgDeleteElevChan:
			for i, v := range allElevs {
				if v.IP == deadElev {
					allElevs = append(allElevs[:i], allElevs[i+1:]...)
					break
				}
			}
		}
	}
}
