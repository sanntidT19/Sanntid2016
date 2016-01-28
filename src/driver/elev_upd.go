package driver

import(
	"errors"
	"fmt"
	"C"
)

//This looks good
func Elev_set_door_open_lamp(turn_on bool) {
	if turn_on {
		io_set_bit(LIGHT_DOOR_OPEN)
	} else {
		io_clear_bit(LIGHT_DOOR_OPEN)
	}
}

func Elev_set_stop_lamp(turn_on bool) {
	if turn_on {
		io_set_bit(STOP)
	} else {
		io_clear_bit(STOP)
	}
}


func Elev_init() /*error*/ {
	fmt.Printf("Kommer da hit\n")
	C.elev_init()
	/*if !io_init() {
		return errors.New("IO initialization failed")
	}
	
	Elev_set_stop_lamp(true)
	Elev_set_door_open_lamp(true)
	fmt.Printf("Kommer da hit\n")
	
	lamp_channel_matrix = Elev_make_std_l_matrix()
	button_channel_matrix = Elev_make_std_b_matrix()
	//Zero all floor button lamps
	//var i int = 0
	for i := 0; i < N_FLOORS; i++ {
		if i != 0 {
			Elev_set_button_lamp(i, DOWN, false)
		}
		if i != N_FLOORS-1 {
			Elev_set_button_lamp(i, UP, false)
		}
		Elev_set_button_lamp(i, COMMAND, false)
	}
	//Clear stop lamp, foor open lamp, and set floor indicatior and ground floor.
	Elev_set_stop_lamp(false)
	Elev_set_door_open_lamp(false)
	Elev_set_floor_indicator(0)

	//Running down until it reaches a floor
	currentFloor := Elev_get_floor_sensor_signal()
	if currentFloor == -1 {
		Elev_set_speed(-300) //Can we write constants here?
		fmt.Println("current direction: ", Elev_get_direction())
		for currentFloor == -1 {
			currentFloor = Elev_get_floor_sensor_signal()
		}
		Elev_stop_elevator()
	}

	//Should current floo
	*/
	return nil
}
