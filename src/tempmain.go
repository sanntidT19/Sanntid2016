package main

import (
	. "statemachine"
	//"fmt"
)

//Making a test function to see how the state machine works
func main() {
	Elevator_manager()

	block := make(chan bool)
	<-block
	/*
		goToFloorChan := make(chan int)
		currentFloorChan := make(chan int)
		currentDirChan := make(chan int)
		buttonSliceChan := statemachine.Create_button_chan_slice()
		lightSliceChan := statemachine.Create_button_chan_slice()
		servedOrderChan := make(chan bool)
		go statemachine.Elevator_worker(goToFloorChan, currentFloorChan, currentDirChan, buttonSliceChan, servedOrderChan)
		go statemachine.Light_updater(lightSliceChan)
		for {
			select {
			case button_Pushed := <-buttonSliceChan[0]:
				goToFloorChan <- 0
				lightSliceChan[0] <- button_Pushed
			case button_Pushed := <-buttonSliceChan[1]:
				goToFloorChan <- 1
				lightSliceChan[1] <- button_Pushed
			case button_Pushed := <-buttonSliceChan[2]:
				goToFloorChan <- 2
				lightSliceChan[2] <- button_Pushed
			case button_Pushed := <-buttonSliceChan[3]:
				goToFloorChan <- 3
				lightSliceChan[3] <- button_Pushed
			case dir := <-currentDirChan:
				fmt.Println("Current direction: %v", dir)
			case <-servedOrderChan:
				fmt.Println("Order served")
			case currFl := <-currentFloorChan:
				fmt.Println("Current floor: ", currFl)
			}
		}
	*/
}
