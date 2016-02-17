package driver

// #cgo LDFLAGS: -lcomedi -lm
// #include <png.h>
// #include "elev.h"
import "C"

import (
	//"errors"
	. "../globalChans"
	. "../globalStructs"
	"fmt"
	"time"
)

//We need global channels if we are to communicate between modules

//This part is copied from last years project and needs to be thorougly confirmed working
var lamp_channel_matrix [][]int
var button_channel_matrix [][]int

func Elev_make_std_l_matrix() [][]int { //stupid name? we need to agree on a stanard name convention
	std_matrix := make([][]int, 4, 8)
	std_matrix[0] = make([]int, NUM_BUTTONS)
	std_matrix[0] = []int{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1}
	std_matrix[1] = make([]int, NUM_BUTTONS)
	std_matrix[1] = []int{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2}
	std_matrix[2] = make([]int, NUM_BUTTONS)
	std_matrix[2] = []int{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3}
	std_matrix[3] = make([]int, NUM_BUTTONS)
	std_matrix[3] = []int{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4}
	return std_matrix
}

func Elev_make_std_b_matrix() [][]int {
	std_matrix := make([][]int, 4, 8) //Find out if the capacity here is valid
	std_matrix[0] = make([]int, NUM_BUTTONS)
	std_matrix[0] = []int{BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1}
	std_matrix[1] = make([]int, NUM_BUTTONS)
	std_matrix[1] = []int{BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2}
	std_matrix[2] = make([]int, NUM_BUTTONS)
	std_matrix[2] = []int{BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3}
	std_matrix[3] = make([]int, NUM_BUTTONS)
	std_matrix[3] = []int{BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4}
	return std_matrix
}

// evnt for Elev_get_floor_sensor_signal() == -1 {kjør ned}

func Elev_init() {
	C.elev_init()
	lamp_channel_matrix = Elev_make_std_l_matrix()
	button_channel_matrix = Elev_make_std_b_matrix()
	if Elev_get_floor_sensor_signal() == -1 {
		Elev_drive_elevator(-1)
		for Elev_get_floor_sensor_signal() == -1 {
			time.Sleep(50 * time.Millisecond)
		}
		Elev_drive_elevator(0)
	}
	fmt.Printf("Initialization of elevator complete.\n")

}

func Elev_drive_elevator(dirn int) {
	if dirn == 0 {
		io_write_analog(MOTOR, 0)
	} else if dirn > 0 {
		io_clear_bit(MOTORDIR)
		io_write_analog(MOTOR, MOTOR_SPEED)
	} else if dirn < 0 {
		io_set_bit(MOTORDIR)
		io_write_analog(MOTOR, MOTOR_SPEED)
	}
}

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

func elev_set_button_light(button int, floor int, value bool) {
	if value {
		io_set_bit(lamp_channel_matrix[floor][button])
	} else {
		io_clear_bit(lamp_channel_matrix[floor][button])
	}
	return
}

//Vurder senere om bool eller int er best her
func elev_get_button_signal(button int, floor int) bool {
	if io_read_bit(button_channel_matrix[floor][button]) {
		return true
	} else {
		return false
	}
}

func Elev_get_floor_sensor_signal() int {
	return int(C.elev_get_floor_sensor_signal())

}

func elev_set_floor_light(floor int) {
	C.elev_set_floor_indicator(C.int(floor))
	return
}

func Open_door() {
	io_set_bit(LIGHT_DOOR_OPEN)
	time.Sleep(3 * time.Second)
	io_clear_bit(LIGHT_DOOR_OPEN)
}

func Check_for_buttons_pressed(button_pressed_chan chan Button) {
	for {
		for i := 0; i < NUM_FLOORS; i++ {
			for j := 0; j < NUM_BUTTONS; j++ {
				if io_read_bit(button_channel_matrix[i][j]) {
					button_pressed_chan <- Button{i, j, true}
				}
			}

		}
		time.Sleep(50 * time.Millisecond)
	}
}

//Should make this general for turning on and off
func Set_button_lights(button_pressed_chan chan Button) {
	for {
		change_button := <-button_pressed_chan
		if change_button.Button_pressed {
			io_set_bit(lamp_channel_matrix[change_button.Floor][change_button.Button_type])
		} else {
			io_clear_bit(lamp_channel_matrix[change_button.Floor][change_button.Button_type])
		}

	}
}

//Lag alle simple funksjoner først. Bruker drivere som vi allerede har. Gjør det simpelt.
//Heller mer komplekst og "go-ete" når funksjoner skal settes sammen i loops og whatever

func Elev_main_tester_function() {
	io_init()
	Elev_init()
	go Elev_floor_light_updater()
	go Check_for_buttons_pressed(ButtonPressedChan)
	go Set_button_lights(SetButtonLightChan)
	return
}

//Use of last_floor may need to be exported og gotten somewhere else
func Elev_floor_light_updater() {
	current_floor := -1
	for {
		time.Sleep(200 * time.Millisecond)
		current_floor = Elev_get_floor_sensor_signal()
		if current_floor != -1 {
			elev_set_floor_light(current_floor)
		}
	}
}
