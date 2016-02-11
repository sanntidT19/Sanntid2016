package stateMachine

import(
	"../driver"
	"../globalChans"
	"time"

)
//up = 1, down = -1, stop = 0
var current_floor int = -1 
var last_floor int  = -1 
var desired_floor int = -1



func Stop_at_desired_floor(order_served_chan chan bool){
	for{
		if driver.Elev_get_floor_sensor_signal() == desired_floor{
			driver.Elev_drive_elevator(0)
			go driver.Open_door()
			order_served_chan <- true
		}
		time.Sleep(200*time.Millisecond)	
	}
}

func Execute_order(next_order_chan chan int){
    for{
		next_order := <-next_order_chan
		desired_floor = next_order
		if next_order > current_floor {
			driver.Elev_drive_elevator(1)	
		}else if next_order < current_floor{
			driver.Elev_drive_elevator(-1)	
		}
    }    
}

func Get_current_floor() {
	for{
		if sensor_result := driver.Elev_get_floor_sensor_signal() ;sensor_result != -1 && sensor_result != current_floor{
			last_floor = current_floor
			current_floor = sensor_result
		}
		time.Sleep(time.Millisecond *50)
	}
}

func write_to_matrix(button_pressed_chan chan Button){
	new_order := <-button_pressed_chan

	 
} // registrerer tastetrykk og legger dette inn i en matrise med oversikt over hvor det er bestillinger

func sort_order(order Button, current_order_queue [] Button) [] Button  {
	sim_dir := current_state.direction


	place_in_queue := -1
	//if queue is empty
	if len(current_order_queue == 0){
		sorted_order_queue := []Button{order}
		return sorted_order_queue
	}
	else if order.button_type == sim_dir && order.floor == driver.Elev_get_floor_sensor_signal() {
				//check if you are in the actual floor and should stop immediately
				place_in_queue = 0
	
	}else{
		for i:= 0; i < len(current_order_queue); i++{
				if order.button_type == sim_dir{
					}
					if sim_dir == UP && {
						
					}
					else if  current_order_queue[i].floor < order.floor && driver.Elev_get_floor_sensor_signal >= order.floor

				}   


				}
				else if order.button_type == UP && direction == 1 && last_floor > order.floor {
					

				}

		}
	}
}




//func send_next_order() // prioritere

/*
Structure of our matrix: 
*/
