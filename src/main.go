package main

import (
	"./driver"
	. "./globalChans"
	//. "./globalStructs"
	"./stateMachine"
	"fmt"
	"time"
	//"encoding/gob"
	//"os"
)

/*
const PATH_OF_SAVED_STATE = "elevState.gob"

func write_elevator_state_to_file() {
	//temp for testing
	test_struct := ElevatorState{255, 19, 1, 1, 100}
	//update this whenever the local elevator gets an order/command
	dataFile, err := os.Create(PATH_OF_SAVED_STATE)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dataEncoder := gob.NewEncoder(dataFile)
	dataEncoder.Encode(test_struct)
	dataFile.Close()
}

func read_elevator_state_from_file() {
	//start with reading it
	data := ElevatorState{0, 0, 0, 0, 0}

	if _, err := os.Stat(PATH_OF_SAVED_STATE); os.IsNotExist(err) {
		fmt.Println("Local save of elevator state not detected.\n It has been cleared/this is the first run on current PC")
		return
	}
	dataFile, err := os.Open(PATH_OF_SAVED_STATE)

	dataDecoder := gob.NewDecoder(dataFile)
	err = dataDecoder.Decode(&data)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dataFile.Close()
	//currentState = data
}
*/

/*
func main() {
	/*
		//First test of basics
		driver.Elev_main_tester_function()
		next_order_chan := make(chan int)
		order_served_chan := make(chan bool)
		go stateMachine.Get_current_floor()
		go stateMachine.Execute_order(next_order_chan)
		go stateMachine.Stop_at_desired_floor(order_served_chan)
		next_order_chan <- 3
		<-order_served_chan
		fmt.Printf("Order_served\n")
		time.Sleep(4*time.Second)

		//Write to file test
		fmt.Printf("Test of file writing and reading/n")
		write_elevator_state_to_file()
		time.Sleep(time.Second * 2)
		read_elevator_state_from_file()
		fmt.Printf("End of main \n")

	Init_chans()
	driver.Elev_main_tester_function()
	stateMachine.State_machine_top_loop()
	for {
		time.Sleep(5 * time.Second)
		fmt.Printf("Main waiting woop woop...\n")
	}
}
*/

func main() {
	InitChans()
	driver.ElevMainTesterFunction()
	go stateMachine.NewTopLoop()
	/*floorReachedChan := make(chan Button)
	orderServedChan := make(chan Button)
	goUpDownChan := make(chan int)
	go stateMachine.MoveElevatorAndOpenDoor(floorReachedChan, orderServedChan, goUpDownChan)
	*/
	fmt.Println("Main: sleeping")
	time.Sleep(time.Second * 1000)
}
