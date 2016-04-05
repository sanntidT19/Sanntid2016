package optalg

import (
	. "../globalChans"
	. "../globalStructs"
	"fmt"
	"math"
)

/*
var el_state1 ElevatorState = ElevatorState{MyIP: "123.123.123.123",
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
var all_elevs []ElevatorState

func Opt_alg(new_order Order) string {
	numOfElevs := len(all_elevs)
	IP_cost_list := make([]int, numOfElevs)
	Queue_len_list := make([]int, numOfElevs)
	var Ele_nmr int = -1
	var IP_score int = 100
	var Optimal_IP string = "0"
	fmt.Println("optalg: number of elevs seen: ", numOfElevs)
	for i, v := range all_elevs {
		Queue_len_list[i] = len(v.OrderQueue)
		if v.CurrentFloor < new_order.Floor {
			if v.Direction != 1 {
				IP_cost_list[i] += 1
			}
		} else if v.CurrentFloor > new_order.Floor {
			if v.Direction != -1 {
				IP_cost_list[i] += 1
			}
		}
		float_difference := float64(v.CurrentFloor - new_order.Floor)
		IP_cost_list[i] += int(math.Abs(float_difference))
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
	return Optimal_IP
}

/*
func main() {
	fmt.Print("optimal IP er:", opt_alg(new_order))
}
*/
func GetOrderQueueOfDeadElev(deadIP string) []Order {
	for _, v := range all_elevs {
		if v.IP == deadIP {
			listCopy := make([]Order, len(v.OrderQueue))
			copy(listCopy, v.OrderQueue)
			return listCopy
		}
	}
	return nil
}

//may not need channels, think about if its better to just call it from somewhere else
func UpdateElevatorStateList() {
	for {
		select {
		case updatedElevState := <-FromNetworkNewElevStateChan:
			elevInList := false
			for i, v := range all_elevs {
				if updatedElevState.IP == v.IP {
					all_elevs[i] = updatedElevState
					elevInList = true
					break
				}
			}
			if !elevInList {
				all_elevs = append(all_elevs, updatedElevState)
			}
		case elevatorTakesOrder := <-AddOrderAssignedToElevStateChan:
			for _, v := range all_elevs {
				if elevatorTakesOrder.AssignedTo == v.IP {
					v.OrderQueue = append(v.OrderQueue, elevatorTakesOrder.Order)
				}
			}
		case deadElev := <-ToOptAlgDeleteElevChan:
			for i, v := range all_elevs {
				if v.IP == deadElev {
					all_elevs = append(all_elevs[:i], all_elevs[i+1:]...)
					break
				}
			}
		}
	}
}
