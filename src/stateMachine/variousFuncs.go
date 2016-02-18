package stateMachine

import (
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

//Dummy struct, for now
var current_state ElevatorState

//Should say it opens door or we need to move that one
//Will this stop when we are idle as well? Desired floor should be set to -1 somewhere
func Stop_at_desired_floor(order_served_chan chan bool) {
	for {
		if driver.Elev_get_floor_sensor_signal() == desired_floor {
			desired_floor = -1000 //tempfix
			driver.Elev_drive_elevator(0)
			fmt.Println("I GET HERE BEFORE OPENING DOOR")
			go driver.Open_door()
			order_served_chan <- true
			fmt.Println("I get here at least")
		}
		time.Sleep(200 * time.Millisecond)
	}
}

//avoid the word "execute". Better name needed
func Execute_order(next_order_chan chan int) {
	for {
		next_order := <-next_order_chan
		desired_floor = next_order
		if next_order > current_state.CurrentFloor {
			driver.Elev_drive_elevator(1)
			current_state.Direction = 1
		} else if next_order < current_state.CurrentFloor {
			driver.Elev_drive_elevator(-1)
			current_state.Direction = -1
		}
	}
}

func Get_current_floor() {
	for {
		if sensor_result := driver.Elev_get_floor_sensor_signal(); sensor_result != -1 && sensor_result != current_floor {
			current_state.LastFloor = current_state.CurrentFloor
			current_state.CurrentFloor = sensor_result
		}
		time.Sleep(time.Millisecond * 50)
	}
}

//LOL WHAT IS THIS? PAST YNGVE/SIGURD YOU SO CRAZY xD
/*func write_to_matrix(button_pressed_chan chan Button) {
	new_order := <-button_pressed_chan

} // registrerer tastetrykk og legger dette inn i en matrise med oversikt over hvor det er bestillinger
*/
func sort_order_queue(new_order Button, current_order_queue []Button) []Button {
	var simulated_direction int
	var simulated_floor int
	place_in_queue := 0
	fmt.Println("Current order queue at start of sort", current_order_queue)
	//if queue is empty
	if len(current_order_queue) == 0 {
		sorted_order_queue := []Button{new_order}
		return sorted_order_queue
	} else if (new_order.Button_type == current_state.Direction || new_order.Button_type == COMMAND) && new_order.Floor == driver.Elev_get_floor_sensor_signal() {
		//check if you are in the actual floor in correct direction and should stop immediately
		place_in_queue = 0
	} else {
		fmt.Println("Start of simulation")
		i := 0
		place_found := false
		//Perform a simulation of what the elevator will do with current queue and find out where the new order belongs
		simulated_direction = current_state.Direction
		simulated_floor = current_state.CurrentFloor + simulated_direction
		fmt.Println("Sim_floor loop 1: ", simulated_floor)
		for simulated_floor >= 0 && simulated_floor < NUM_FLOORS {
			if new_order.Floor == simulated_floor && (simulated_direction == new_order.Button_type || new_order.Button_type == COMMAND) {
				place_in_queue = i
				place_found = true
				break
			} else if current_order_queue[i].Floor == simulated_floor && (simulated_direction == current_order_queue[i].Button_type || current_order_queue[i].Button_type == COMMAND) {
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

				if new_order.Floor == simulated_floor && (simulated_direction == new_order.Button_type || new_order.Button_type == COMMAND) {
					place_in_queue = i
					place_found = true
					break
				} else if current_order_queue[i].Floor == simulated_floor && (simulated_direction == current_order_queue[i].Button_type || current_order_queue[i].Button_type == COMMAND) {
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

				if new_order.Floor == simulated_floor && (simulated_direction == new_order.Button_type || new_order.Button_type == COMMAND) {
					place_in_queue = i
					place_found = true
					break
				} else if current_order_queue[i].Floor == simulated_floor && (simulated_direction == current_order_queue[i].Button_type || current_order_queue[i].Button_type == COMMAND) {
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

		/*for i:= 0; i < len(current_order_queue); i++{
				if order.button_type == sim_dir{
					}
					if sim_dir == UP && {

					}
					else if  current_order_queue[i].floor < order.floor && driver.Elev_get_floor_sensor_signal >= order.floor

				}


				}
				else if order.button_type == UP && direction == 1 && last_floor > order.floor {


				}

		}*/
	}
	fmt.Println("Place in queue:", place_in_queue)
	fmt.Println(current_order_queue)
	fmt.Println("capacity: ", cap(current_order_queue))
	//current_order_queue = current_order_queue[0 : len(current_order_queue)+1]
	fmt.Println("Before copying ", current_order_queue)
	//sorted_order_queue := make([]Button,len(current_order_queue)+1)
	//sorted_order_queue[:place_in_queue] = current_order_queue[:place_in_queue]
	//sorted_order_queue[place_in_queue+1:] = current_order_queue[place_in_queue]

	current_order_queue = append(current_order_queue[:place_in_queue], append([]Button{new_order}, current_order_queue[place_in_queue:]...)...)

	/*
		sorted_order_queue := []Button{}
		sorted_order_queue = append(sorted_order_queue, current_order_queue[:place_in_queue])
		sorted_order_queue = append(sorted_order_queue, new_order)
		sorted_order_queue = append(sorted_order_queue, current_order_queue[place_in_queue:])
	*/
	//copy(current_order_queue[place_in_queue+1:], current_order_queue[place_in_queue:])
	//current_order_queue[place_in_queue] = new_order
	//found place in queue, make a new and updated queue
	//current_order_queue = append(current_order_queue, Button{})
	//copy(current_order_queue[:place_in_queue+1], current_order_queue[place_in_queue:])
	//current_order_queue[place_in_queue] = new_order
	fmt.Println("End of sorting queue: ", current_order_queue)
	return current_order_queue
}

//For select that runs it all together goes here. Better name must be found.
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
		case new_order := <-ButtonPressedChan:
			if isOrderInQueue(order_queue, new_order) {
				break
			}
			/*
				sort order in queue, if queue is empty, execute?(doesnt seem like an optimal structure)
				it should only sort. Maybe even in an own goroutine
			*/
			SetButtonLightChan <- new_order
			order_queue = sort_order_queue(new_order, order_queue)
			next_order_chan <- order_queue[0].Floor //this needs to be separated from the others
		default:
			//fmt.Println("I am not stuck anywhere here")
			time.Sleep(time.Millisecond * 100)

		}

	}
}

//func send_next_order() // prioritere

/*
Structure of our matrix:
*/
func isOrderInQueue(order_queue []Button, new_order Button) bool {
	for _, queueElements := range order_queue {
		if queueElements == new_order {
			return true
		}
	}
	return false
}
