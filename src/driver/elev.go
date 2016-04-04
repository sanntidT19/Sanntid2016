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

func ElevMakeStdLMatrix() [][]int { //stupid name? we need to agree on a stanard name convention
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

func ElevMakeStdBMatrix() [][]int {
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

func ElevInit() {
	C.elev_init()
	lamp_channel_matrix = ElevMakeStdLMatrix()
	button_channel_matrix = ElevMakeStdBMatrix()
	if ElevGetFloorSensorSignal() == -1 {
		ElevDriveElevator(-1)
		for ElevGetFloorSensorSignal() == -1 {
			time.Sleep(50 * time.Millisecond)
		}
		ElevDriveElevator(0)
	}
	fmt.Printf("Initialization of elevator complete.\n")

}


func ElevDriveElevator(dirn int) {
	if dirn == 0 {
		IoWriteAnalog(MOTOR, 0)
	} else if dirn > 0 {
		IoClearBit(MOTORDIR)
		IoWriteAnalog(MOTOR, MOTOR_SPEED)
	} else if dirn < 0 {
		IoSetBit(MOTORDIR)
		IoWriteAnalog(MOTOR, MOTOR_SPEED)
	}
}

func ElevSetDoorOpenLamp(turn_on bool) {
	if turn_on {
		IoSetBit(LIGHT_DOOR_OPEN)
	} else {
		IoClearBit(LIGHT_DOOR_OPEN)
	}
}


func ElevSetStopLamp(turn_on bool) {
	if turn_on {
		IoSetBit(STOP)
	} else {
		IoClearBit(STOP)
	}
}

func ElevSetButtonLight(button int, floor int, value bool) {
	if value {
		IoSetBit(lamp_channel_matrix[floor][button])
	} else {
		IoClearBit(lamp_channel_matrix[floor][button])
	}
	return
}

//Vurder senere om bool eller int er best her
func ElevGetButtonSignal(button int, floor int) bool {
	if IoReadBit(button_channel_matrix[floor][button]) {
		return true
	} else {
		return false
	}
}

func ElevGetFloorSensorSignal() int {
	return int(C.elev_get_floor_sensor_signal())

}

func ElevSetFloorLight(floor int) {
	C.elev_set_floor_indicator(C.int(floor))
	return
}


func OpenDoor() {
	IoSetBit(LIGHT_DOOR_OPEN)
	time.Sleep(3 * time.Second)
	IoClearBit(LIGHT_DOOR_OPEN)
}

func ElevNotMoving() bool {
	if IoReadAnalog(MOTOR) == 0 {
		return true
	} else {
		return false
	}
}
func CheckForButtonsPressed() {
	for {
		for i := 0; i < NUM_FLOORS; i++ {
			for j := 0; j < NUM_BUTTONS; j++ {
				if IoReadBit(button_channel_matrix[i][j]) {
					var button_type int
					if j == 0 {
						button_type = UP
					} else if j == 1 {
						button_type = DOWN
					} else {
						button_type = COMMAND
					}
					if button_type == COMMAND{
						InternalButtonPressedChan <- Order{i,button_type}
					}else{
						ExternalButtonPressedChan <- Order{i, button_type}
					}
				}
			}

		}
		time.Sleep(100 * time.Millisecond) //Doblet sleep. Se hvordan det går. Kanskje også en sleep etter 
	}
}


func SetButtonLight(ButtonLight Order, turnOn bool){
	if ButtonLight.Button_type == UP {
		ButtonLight.Button_type = 0
	} else if ButtonLight.Button_type == DOWN {
		ButtonLight.Button_type = 1
	} else {
		ButtonLight.Button_type = 2
	}
	
	if turnOn {
		io_set_bit(lamp_channel_matrix[ButtonLight.Floor][ButtonLight.Button_type])
	} else {
		io_clear_bit(lamp_channel_matrix[ButtonLight.Floor][ButtonLight.Button_type])
	}

}


//Rename
func ElevMainTesterFunction() {
	IoInit()
	ElevInit()
	go ElevFloorLightUpdater()
	go CheckForButtonsPressed(ButtonPressedChan)
	return
}
//Use of last_floor may need to be exported og gotten somewhere else
func ElevFloorLightUpdater() {
	current_floor := -1
	for {
		time.Sleep(200 * time.Millisecond)
		current_floor = ElevGetFloorSensorSignal()
		if current_floor != -1 {
			ElevSetFloorLight(current_floor)
		}
	}
}
