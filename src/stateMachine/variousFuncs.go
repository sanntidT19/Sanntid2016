package stateMachine

import(
	"../driver"
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
