package stateMachine

import (
	"../communication"
	"../driver"
	. "../globalChans"
	. "../globalStructs"
	"fmt"
	"time"
)

//up = 1, down = -1, stop = 0
//We need better names for these to avoid confusion. The whole floor-thing needs to be thought through
var current_floor int = -1
var last_floor int = -1
var desired_floor int = -1
var door_open bool = false

//Dummy struct, for now, must be set to have dir as -1 or 1 for sort function to work properly
var current_state ElevatorState

// For new suggestion
var StateOfElev ElevatorState

//Should say it opens door or we need to move that one
//Will this stop when we are idle as well? Desired floor should be set to -1 somewhere
/*func Stop_at_desired_floor(order_served_chan chan bool) {
	for {
		if driver.ElevGetFloorSensorSignal() == desired_floor {
			desired_floor = -1000 //tempfix
			driver.ElevDriveElevator(0)
			fmt.Println("I GET HERE BEFORE OPENING DOOR")
			door_open = true
			driver.Open_door()
			door_open = false
			order_served_chan <- true
			fmt.Println("I get here at least")
		}
		time.Sleep(200 * time.Millisecond)
	}
}
*/
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

/*
//avoid the word "execute". Better name needed
func Execute_order(next_order_chan chan int) {
	for {
		next_order := <-next_order_chan
		desired_floor = next_order
		fmt.Printf("desired floor is:", desired_floor)
		if next_order > current_state.CurrentFloor && !door_open {
			driver.ElevDriveElevator(1)
			current_state.Direction = 1
		} else if next_order < current_state.CurrentFloor && !door_open {
			driver.ElevDriveElevator(-1)
			current_state.Direction = -1
		}
	}
}
*/
/*
func Get_current_floor() {
	for {
		if sensor_result := driver.ElevGetFloorSensorSignal(); sensor_result != -1 && sensor_result != current_state.CurrentFloor {
			current_state.LastFloor = current_state.CurrentFloor
			current_state.CurrentFloor = sensor_result
			fmt.Println("Sensor result:", sensor_result)
		}
		time.Sleep(time.Millisecond * 50)
	}
}*/
/*
 */
func sort_order_queue(new_order Order, current_state ElevatorState) []Order {
	common_current_order_queue := current_state.OrderQueue
	current_order_queue := make([]Order, len(common_current_order_queue))
	copy(current_order_queue, common_current_order_queue)

	var simulated_direction int
	var simulated_floor int
	place_in_queue := 0
	fmt.Println("Current order queue at start of sort", current_order_queue)
	//if queue is empty
	if len(current_order_queue) == 0 {
		sorted_order_queue := []Order{new_order}
		fmt.Println("Queue is empty")
		return sorted_order_queue
	} else if (new_order.Direction == current_state.Direction || new_order.Direction == COMMAND) && new_order.Floor == driver.ElevGetFloorSensorSignal() {
		//check if you are in the actual floor in correct direction and should stop immediately
		place_in_queue = 0
	} else {
		fmt.Println("Start of simulation")
		i := 0
		place_found := false
		//Perform a simulation of what the elevator will do with current queue and find out where the new order belongs
		simulated_direction = current_state.Direction
		simulated_floor = current_state.CurrentFloor
		if driver.ElevGetFloorSensorSignal() == -1 {
			fmt.Println("Elevator not in current floor, simulate past current floor")
			simulated_floor = current_state.CurrentFloor + simulated_direction //er dette riktig?
		}
		fmt.Println("Sim_floor loop 1: ", simulated_floor)
		fmt.Println("Sim_dir loop 1: ", simulated_direction)
		fmt.Println("new order floor: ", new_order.Floor)
		fmt.Println("new order direction: ", new_order.Direction)
		for simulated_floor >= 0 && simulated_floor < NUM_FLOORS {
			if new_order.Floor == simulated_floor && (simulated_direction == new_order.Direction || new_order.Direction == COMMAND) {
				place_in_queue = i
				place_found = true
				break
			} else if current_order_queue[i].Floor == simulated_floor && (simulated_direction == current_order_queue[i].Direction || current_order_queue[i].Direction == COMMAND) {
				//order was theoretically served and next element in the current queue is up for evaluation
				i++
				fmt.Println("i, sim1: ", i)
				if i == len(current_order_queue) {
					//all orders in queue evaluated. put new order at end
					place_in_queue = i
					place_found = true
					break
				}
			}
			simulated_floor += simulated_direction
		}
		simulated_floor -= simulated_direction //bounds were exceeded and you take one step back
		simulated_direction *= -1
		if !place_found {
			for simulated_floor >= 0 && simulated_floor < NUM_FLOORS {

				if new_order.Floor == simulated_floor && (simulated_direction == new_order.Direction || new_order.Direction == COMMAND) {
					place_in_queue = i
					place_found = true
					break
				} else if current_order_queue[i].Floor == simulated_floor && (simulated_direction == current_order_queue[i].Direction || current_order_queue[i].Direction == COMMAND) {
					//order was theoretically served and next element in the current queue is up for evaluation
					i++
					fmt.Println("i, sim2: ", i)
					if i == len(current_order_queue) {
						//all orders in queue evaluated. put new order at end
						place_in_queue = i
						place_found = true
						break
					}
				}
				simulated_floor += simulated_direction
			}
		}
		simulated_floor -= simulated_direction //bounds were exceeded and you take one step back
		simulated_direction *= -1

		if !place_found {
			for simulated_floor != (current_state.CurrentFloor + simulated_direction) {
				if new_order.Floor == simulated_floor && (simulated_direction == new_order.Direction || new_order.Direction == COMMAND) {
					place_in_queue = i
					place_found = true
					break
				} else if current_order_queue[i].Floor == simulated_floor && (simulated_direction == current_order_queue[i].Direction || current_order_queue[i].Direction == COMMAND) {
					//order was theoretically served and next element in the current queue is up for evaluation
					i++
					fmt.Println("i, sim3: ", i)
					if i == len(current_order_queue) {
						//all orders in queue evaluated. put new order at end
						place_in_queue = i
						place_found = true
						break
					}
				}
				simulated_floor += simulated_direction
			}
		}
	}
	sorted_order_queue := append(current_order_queue[:place_in_queue], append([]Order{new_order}, current_order_queue[place_in_queue:]...)...)
	/*fmt.Println("Order queue after sort")
	for _, v := range sorted_order_queue {
		Print_order(v)
	}*/
	return sorted_order_queue
}

//For select that runs it all together goes here. Better name must be found.
/*
func State_machine_top_loop() {
	next_order_chan := make(chan int)
	order_served_chan := make(chan bool)
	order_queue := []Button{}

	go Stop_at_desired_floor(order_served_chan)
	go Get_current_floor()
	go Execute_order(next_order_chan)
	for {
		select {
		case <-order_served_chan:
			fmt.Println("order queue when order is served: ", order_queue)
			localCopyOfOrderServed := order_queue[0]
			localCopyOfOrderServed.Button_pressed = false
			SetButtonLightChan <- localCopyOfOrderServed
			if len(order_queue) > 1 {
				order_queue = order_queue[1:]
				next_order_chan <- order_queue[0].Floor
			} else {
				fmt.Printf("Here everything is ok\n")
				order_queue = order_queue[1:]
				fmt.Printf("I also get here\n")
				//next line is just so door isnt always open
			}
			//Denne casen kan trolig gjøres parallelt med de andre casene.
		case new_order := <-ButtonPressedChan:
			if isOrderInQueue(order_queue, new_order) {
				break
			}
			/*
				sort order in queue, if queue is empty, execute?(doesnt seem like an optimal structure)
				it should only sort. Maybe even in an own goroutine

			SetButtonLightChan <- new_order

			order_queue = sort_order_queue(new_order, order_queue)

			next_order_chan <- order_queue[0].Floor
			//this needs to be separated from the others
		default:
			//fmt.Println("I am not stuck anywhere here")
			time.Sleep(time.Millisecond * 100)

		}
	}
}
*/
/*
 */
func isOrderInQueue(order_queue []Order, new_order Order) bool {
	for _, queueElements := range order_queue {
		if queueElements == new_order {
			return true
		}
	}
	return false
}

//UTKAST TIL REFORMATERT HEIS ER UNDER HER. KOM MED GODE NAVN PÅ TING ALLEREDE NÅ
func MoveElevatorAndOpenDoor(floorWithOrderReachedChan chan Order, orderServedChan chan Order, sendElevInDirectionChan chan int) {
	//CAN BE: IDLE, OPEN_DOOR, GO_TO_FLOOR (evt can G_T_F be go_up, and go_down)
	for {
		select {
		case newestOrder := <-floorWithOrderReachedChan:
			if newestOrder.Direction != COMMAND {
				StateOfElev.Direction = newestOrder.Direction
			}
			driver.ElevDriveElevator(0)
			driver.OpenDoor()
			//After door is closed. Tell someone above
			orderServedChan <- newestOrder
		case goInDirection := <-sendElevInDirectionChan:
			if goInDirection == UP {
				driver.ElevDriveElevator(UP)
				StateOfElev.Direction = UP
			} else {
				driver.ElevDriveElevator(DOWN)
				StateOfElev.Direction = DOWN
			}
			ToNetworkNewElevStateChan <- StateOfElev
			//go to floor, if below, go up. if above, go down. This is only needed when the elevator is idle.
		}
	}
}

func UpdateNewFloorReached() {
	mostRecentFloorVisited := driver.ElevGetFloorSensorSignal()
	for {
		if sensor_result := driver.ElevGetFloorSensorSignal(); sensor_result != -1 && sensor_result != mostRecentFloorVisited {
			mostRecentFloorVisited = sensor_result
			StateOfElev.PreviousFloor = StateOfElev.CurrentFloor
			StateOfElev.CurrentFloor = sensor_result
			ToNetworkNewElevStateChan <- StateOfElev
		}
		time.Sleep(time.Millisecond * 50)
	}
}

/*
func SetButtonLights(){
	Need to set lights for all elevators.
	Maybe separate one for internal.
}
*/

//Vurder omformulere denne, slik at den kalles istedenfor å kjøres i en egen goroutine
func feedDirectionCommandsToElev(orderQueueChangeChan chan bool, sendElevInDirectionChan chan int, targetFloorReachedChan chan Order) {
	//State of elev is known globally for now
	//Send til hit når ordre er served også
	stopSignalSent := false
	go func() {
		for {
			if len(StateOfElev.OrderQueue) > 0 {
				if StateOfElev.OrderQueue[0].Floor == driver.ElevGetFloorSensorSignal() {
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

func NewTopLoop() {
	//Her kan man anta at man står stille i en etasje
	StateOfElev.CurrentFloor = driver.ElevGetFloorSensorSignal()
	StateOfElev.Direction = DOWN
	StateOfElev.IP = communication.GetLocalIP()
	ToNetworkNewElevStateChan <- StateOfElev

	targetFloorReachedChan := make(chan Order)
	orderServedChan := make(chan Order)
	sendElevInDirectionChan := make(chan int)
	orderQueueChangeChan := make(chan bool)
	go SpamCurrentQueue()
	go UpdateNewFloorReached()
	go MoveElevatorAndOpenDoor(targetFloorReachedChan, orderServedChan, sendElevInDirectionChan)
	go feedDirectionCommandsToElev(orderQueueChangeChan, sendElevInDirectionChan, targetFloorReachedChan)
	fmt.Println("Statemachine: ready")
	defer driver.ElevDriveElevator(0)
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

			
				OrderServedLocallyChan <- servedOrder

				fmt.Println("order served")
			}
			ToNetworkNewElevStateChan <- StateOfElev
		case newOrder := <-NewOrderToLocalElevChan: //newOrderToElevChan when we have mothership ready
			if isOrderInQueue(StateOfElev.OrderQueue, newOrder) {
				break
			} else {
				StateOfElev.OrderQueue = sort_order_queue(newOrder, StateOfElev)
				copyOfOrderList := make([]Order, len(StateOfElev.OrderQueue))
				copy(copyOfOrderList, StateOfElev.OrderQueue)
				orderQueueChangeChan <- true
			}
			ToNetworkNewElevStateChan <- StateOfElev
		}
	}

}

func SendElevStateToNetwork(toNetworkChan chan ElevatorState){
	copyOfElevState := StateOfElev
	copyOfOrderList := make([]Order, len(StateOfElev.OrderQueue))
	copy(copyOfOrderList,StateOfElev.OrderQueue)
	copyOfElevState.OrderQueue = copyOfOrderList
	toNetworkChan <-copyOfElevState
}

func SpamCurrentQueue(){
	for{
		for _, v := range StateOfElev.OrderQueue{
			PrintOrder(v)
		}
		fmt.Println()
		time.Sleep(time.Second*1)
	}
}



func insertOrderIntoQueue(newOrder Order, currentState ElevatorState) []Order {
	common_current_order_queue := currentState.OrderQueue
	orderQueueCopy := make([]Button, len(common_current_order_queue))
	copy(orderQueueCopy, common_current_order_queue)
	var placeInQueue int
	fmt.Println("Current order queue at start of sort", orderQueueCopy)
	//if queue is empty
	if len(orderQueueCopy) == 0 {
		newOrderQueue := []Order{newOrder}
		return newOrderQueue
	} else if (newOrder.Direction == currentState.Direction || newOrder.Direction == COMMAND) && newOrder.Floor == driver.Elev_get_floor_sensor_signal() {
		//check if you are in the actual floor in correct direction and should stop immediately
		placeInQueue = 0
	} else {
		//Perform a simulation of what the elevator will do with current queue and find out where the new order belongs
		currentDirection = currentState.Direction
		firstFloorToVisit = currentState.CurrentFloor
		if driver.Elev_get_floor_sensor_signal() == -1 {
			fmt.Println("Elevator not in current floor, simulate past current floor")
			firstFloorToVisit = currentState.CurrentFloor + simulatedDirection //er dette riktig?
		}
		placeInQueue = SimulateElevDrivingFindOrderIndex(firstFloorToVisit,currentDirection, orderQueueCopy, newOrder)
	}
	newOrderQueue := append(orderQueueCopy[:placeInQueue], append([]Button{newOrder}, orderQueueCopy[placeInQueue:]...)...)
	fmt.Println("Order queue after sort")
	for _, v := range newOrderQueue {
		Print_order(v)
	}
	return newOrderQueue
}



/*
Simulate driving the elevator from current floor and serve all possible orders.
*/
func SimulateElevDrivingFindOrderIndex(startingFloor int,startingDirection int,orderQueue []Order, newOrder Order) int{
	simulatedFloor := startingFloor
	simulatedDirection := startingDirection
	placeInQueue := 0
	i := 0
	for simulatedFloor >= 0 && simulatedFloor < NUM_FLOORS {
		if new_order.Floor == simulatedFloor && (simulatedDirection == new_order.Direction || new_order.Direction == COMMAND) {
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
	if simulatedDirection == UP{
		simulatedDirection = DOWN
	}else{
		simulatedDirection = UP
	}
	for simulatedFloor >= 0 && simulatedFloor < NUM_FLOORS {
			if new_order.Floor == simulatedFloor && (simulatedDirection == new_order.Direction || new_order.Direction == COMMAND) {
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
	}

	simulatedFloor -= simulatedDirection //bounds were exceeded and you take one step back
	if simulatedDirection == UP{
		simulatedDirection = DOWN
	}else{
		simulatedDirection = UP
	}
	for simulatedFloor != (current_state.CurrentFloor + simulatedDirection) {
		if new_order.Floor == simulatedFloor && (simulatedDirection == new_order.Direction || new_order.Direction == COMMAND) {
			placeInQueue = i
				return placeInQueue
			} else if orderQueue[i].Floor == simulatedFloor && (simulatedDirection == orderQueue[i].Direction || orderQueue[i].Direction == COMMAND) {
				//order was theoretically served and next element in the current queue is up for evaluation
				i++
				fmt.Println("i, sim3: ", i)
				if i == len(orderQueue) {
					//all orders in queue evaluated. put new order at end
					place_in_queue = i
				}
			}
		simulatedFloor += simulatedDirection
	}
	return placeInQueue
}
func isOrderInQueue(order_queue []Button, new_order Button) bool {
	for _, queueElements := range order_queue {
		if queueElements == new_order {
			return true
		}
	}
	return false
}
