package statemachine

import (
	"driver"
)

//Functions used when running the elevator, find out better name and add prefix

type Button struct {
	floor      int
	buttonType int
	pushed     int
}

func button_updater(buttonChan chan Button) { //Sending the struct a level up, to the state machine setting and turning off lights.
	var button_matrix [][]int
	button_matrix = make([][]int, driver.N_FLOORS)
	var i int = 0
	var j int
	for i; i < driver.N_FLOORS; i++ {
		button_matrix[i] = make([]int, driver.N_BUTTONS) //Golang creates a slice of zeros by default
	}
	//Continious checking of buttons, buttonChan is a buffered channel who can fit N_FLOORS*N_BUTTONS elements

	for {
		i = 0
		for i; i < driver.N_FLOORS; i++ {
			j = 0
			for j; j < driver.N_BUTTONS; j++ {
				if buttonVar = driver.Elev_get_button_signal(i, j); buttonVar != button_matrix[i][j] { //Sending the struct if its pushed and hasnt been sent already
					if buttonVar == 1 {
						buttonChan <- Button(i, j, buttonVar)
					}
					button_matrix[i][j] = buttonVar
				}
			}
		}
	}
}

func motor_control(speedChan chan float64) { //I think speedchan should not be buffered
	for {
		speedVal := <-speedChan
		driver.Elev_set_speed(speedVal)
	}
}
