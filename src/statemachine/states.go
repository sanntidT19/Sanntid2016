package statemachine

import (
	"./driver"
	"fmt"
	"time"
)

const MAX_SPEED_UP = 300
const MAX_SPEED_DOWN = -300
const SPEED_STOP = 0

// Get it confirmed that this is the case
const DIR_UP = 0
const DIR_DOWN = 1

var ExStateMChans ExternalStateMachineChannels
var InStateMChans InternalStateMachineChannels

//if we get errors, this bool might be the bad guy
type Button struct {
	floor      int
	buttonType int
	turnOn     bool
}

type State struct {
	direction    int
	currentFloor int
}

//Functions used when running the elevator, find out better name and add prefix
/*
The elevator manager
	- Takes in an array, and updates buttons/lights to be cleared/set, comparing with its own private array, only sends to a channel if something is different
	- Has a logic function sending the correct order to the elevator_worker
	- Sends these updates to the mothership: state(struct) updated, order (button struct) served, has/has not command (bool?)



*/
type ExternalStateMachineChannels struct {

}
type InternalStateMachineChannels struct {

}

func External_state_machine_channels_init() {

}
func Internal_state_machine_channels_init() {

}


func Elevator_manager( /*orderArrayChan chan [][] int   <- this will come from above, but not in this test program */ /* this may need to have a different name ->  ,currentStateChan chan State*/) {
	driver.Elev_init()

	var orderArray [driver.N_FLOORS][2]int // Up and down goes here
	var commandOrderList [driver.N_FLOORS]int
	var managersCurrentOrder Button
	externalButtPressChan := make(chan Button)
	internalButtPressChan := make(chan Button)
	buttUpdatedChan := make(chan Button)
	orderArrayChan := make(chan [driver.N_FLOORS][2]int)
	setLightChan := make(chan Button)
	commandOrderChan := make(chan [driver.N_FLOORS]int)
	currentStateChan := make(chan State)
	ordersCalculatedChan := make(chan []Button)
	orderToWorkerChan := make(chan int)
	orderServedChan := make(chan bool)
	orderToManagerChan := make(chan Button)
	tempOrderServedChan := make(chan bool)

	//buttonSliceChan := Create_button_chan_slice()
	//lightSliceChan := Create_button_chan_slice()

	go Button_updater(externalButtPressChan, internalButtPressChan, buttUpdatedChan)
	go Light_updater(setLightChan, buttUpdatedChan)

	go Choose_next_order(orderArrayChan, commandOrderChan, currentStateChan, ordersCalculatedChan)

	go Elevator_worker(orderToWorkerChan, currentStateChan, orderServedChan)
	go Send_orders_to_worker(orderToManagerChan, ordersCalculatedChan, tempOrderServedChan)

	//This logic below should be done somewhere else, I think
	for {
		fmt.Println("I've come this far!")
		select {
		//If something is pressed, the channels are updated
		case orderButt := <-internalButtPressChan:

			commandOrderList[orderButt.floor] = 1
			commandOrderChan <- commandOrderList
			setLightChan <- orderButt
		case orderButt := <-externalButtPressChan:

			orderArray[orderButt.floor][orderButt.buttonType] = 1
			orderArrayChan <- orderArray
			setLightChan <- orderButt
		case managersCurrentOrder := <-orderToManagerChan: //Just so the manager can keep up with the current order

			orderToWorkerChan <- managersCurrentOrder.floor
		case <-orderServedChan:
			fmt.Println("I'm in this case!")

			if managersCurrentOrder.buttonType == driver.COMMAND {
				commandOrderList[managersCurrentOrder.floor] = 0
			} else {
				orderArray[managersCurrentOrder.floor][managersCurrentOrder.buttonType] = 0
			}
			tempOrderServedChan <- true
			setLightChan <- Button{managersCurrentOrder.floor, managersCurrentOrder.buttonType, false}
		}
	}

}

//Going where it is told to go, based on information on where it is.   Gotofloor needs to be buttonized.
func Elevator_worker(goToFloorChan chan int, currentStateChan chan State, orderServedChan chan bool) {
	speedChan := make(chan float64) //This channel is only between the statemachine and its functions
	privateSensorChan := make(chan int)
	currentDir := driver.Elev_get_direction()
	currentFloor := driver.Elev_get_floor_sensor_signal() //We know that we are in a floor at this point.
	previousFloor := -1
	orderedFloor := -1
	go Motor_control(speedChan)
	go Is_floor_reached(privateSensorChan)
	// This is where the statemachine is implemented, should it be a select case?
	for {
		select {
		//Slave could send a new command while the statemachine is serving another command, but it should fix the logic by itself
		// New order
		case gtf := <-goToFloorChan:
			orderedFloor = gtf
			//You are in the floor, order served immediatly, maybe this if can be implemented in another case, but its here for now.
			if gtf == driver.Elev_get_floor_sensor_signal() {
				speedChan <- SPEED_STOP
				orderServedChan <- true
				Open_door() //Dont think we want this select loop to do anything else while the door is open. Solve with go open_door() if its not the case
				//You know you are under/above the current floor
			} else if gtf < currentFloor {
				speedChan <- MAX_SPEED_DOWN
				currentDir = DIR_DOWN
				currentStateChan <- State{currentDir, currentFloor} //Do this for all? Should we send if its already going down (no state changed then)

			} else if gtf > currentFloor {
				speedChan <- MAX_SPEED_UP
				currentDir = DIR_UP
				currentStateChan <- State{currentDir, currentFloor}
				//Your last floor was the current floor, but something may have been pulled, so you dont know where you lie relative to it. Cant use direction.
			} else if gtf == currentFloor { //this may now go all the time
				//Using previousfloor can give you an idea in some cases.
				if previousFloor > currentFloor {
					speedChan <- MAX_SPEED_UP
					currentDir = DIR_UP
					currentStateChan <- State{currentDir, currentFloor}
				} else if previousFloor < currentFloor { //This will also be the case if prevFloor is undefined (-1) They can never be the same.
					speedChan <- MAX_SPEED_DOWN
					currentDir = DIR_DOWN
					currentStateChan <- State{currentDir, currentFloor}
				}
			}
		//New floor is reached and therefore shit is updated
		case cf := <-privateSensorChan:
			previousFloor = currentFloor
			currentFloor = cf
			currentStateChan <- State{currentDir, currentFloor}
			if orderedFloor == currentFloor {
				speedChan <- SPEED_STOP
				orderServedChan <- true
				Open_door()
			}
		}
	}
}

//This function has control over the orderlist and speaks directly to the worker.
func Send_orders_to_worker(orderToManagerChan chan Button, ordersCalculatedChan chan []Button, tempOrderServedChan chan bool) {
	var currentOrderList []Button
	var currentOrderIter int
	emptyButton := Button{}
	for {
		select {
		//New orders sorted, picked. scrap the old one
		case currentOrderList := <-ordersCalculatedChan:
			currentOrderIter = 0
			orderToManagerChan <- currentOrderList[currentOrderIter]
			currentOrderIter++
		case <-tempOrderServedChan: //This also needs to be sent upwards somehow
			if currentOrderList[currentOrderIter] != emptyButton /*this need to be checked*/ {
				orderToManagerChan <- currentOrderList[currentOrderIter]
				currentOrderIter++
			}
		}
	}
}

//Swapping : x,y = y,x
// logic here, no problemo, motherfucker
func Choose_next_order(orderArrayChan chan [driver.N_FLOORS][2]int, commandOrderChan chan [driver.N_FLOORS]int, currentStateChan chan State, ordersCalculatedChan chan []Button) {
	var currentFloor, currentDir int
	var dirIter int //Deciding where to iterate first
	var orderArray [driver.N_FLOORS][2]int
	var commandList [driver.N_FLOORS]int
	var firstPriority, secondPriority int
	incomingUpdateChan := make(chan bool)
	for {
		select {
		case currentState := <-currentStateChan:
			currentFloor = currentState.currentFloor
			currentDir = currentState.direction
		case orderArray = <-orderArrayChan:
			fmt.Println("I am also in this choosenextorder-case")
			incomingUpdateChan <- true
		case commandList = <-commandOrderChan:
			fmt.Println("hope im not here")
			incomingUpdateChan <- true
		case fucker := <-incomingUpdateChan:
			fmt.Println("which means im here aswell")
			//Makin a slice of sorted orders and sending it to the elevator manager.
			resultOrderSlice := make([]Button, driver.N_FLOORS*driver.N_BUTTONS) //This needs to be printed
			resultIter := 0
			if fucker {
				fmt.Println("fuck")
			}
			dirIter = 1
			firstPriority = driver.UP
			secondPriority = driver.DOWN
			if currentDir == DIR_DOWN {
				dirIter = -1
				firstPriority = driver.DOWN
				secondPriority = driver.UP
			}
			//If we are above/below the floor we need to prioritize that as one of the less attractive ones
			if driver.Elev_get_floor_sensor_signal() != currentFloor { //Maybe assign a value here
				currentFloor += dirIter //Setting the startingfloor to iterate from
			}
			i := currentFloor
			// Iterating in the most desirable direction
			for ; i < driver.N_FLOORS && i >= 0; i += dirIter {
				if isCommandOrder := commandList[i]; isCommandOrder == 1 {
					resultOrderSlice[resultIter] = Button{i, driver.COMMAND, true}
					resultIter++
				}
				if isOrder := orderArray[i][firstPriority]; isOrder == 1 {
					resultOrderSlice[resultIter] = Button{i, firstPriority, true}
					resultIter++
				}
			}
			//Now going from top/bottom and checking the orders in the other direction, as well as commands not yet checked
			dirIter = dirIter * -1
			firstPriority, secondPriority = secondPriority, firstPriority
			canCheckCommand := false
			for ; i < driver.N_FLOORS && i >= 0; i += dirIter {
				if canCheckCommand {
					if isCommandOrder := commandList[i]; isCommandOrder == 1 {
						resultOrderSlice[resultIter] = Button{i, driver.COMMAND, true}
						resultIter++
					}
				}
				if i == currentFloor {
					canCheckCommand = true
				}
				if isOrder := orderArray[i][firstPriority]; isOrder == 1 {
					resultOrderSlice[resultIter] = Button{i, firstPriority, true}
					resultIter++
				}
			}
			//Lowest priority: checking from bottom/top to current floor if there are any bastards wanting an elevated experience all commands have been checked
			dirIter = dirIter * -1
			firstPriority, secondPriority = secondPriority, firstPriority
			for ; i != currentFloor; i += dirIter {
				if isOrder := orderArray[i][firstPriority]; isOrder == 1 {
					resultOrderSlice[resultIter] = Button{i, firstPriority, true}
					resultIter++
				}

				ordersCalculatedChan <- resultOrderSlice
			}
		}
	}
}

//This one creates the basic button slice for our friends

func Create_button_chan_slice() []chan Button { //A little unsure of this
	//fmt.Println("making slice for you, sir")
	chanSlice := make([]chan Button, driver.N_FLOORS)
	for i := 0; i < driver.N_FLOORS; i++ {
		chanSlice[i] = make(chan Button)
	}
	return chanSlice
}

/*
func old_Button_updater(buttonSlice []chan Button) { //Sending the struct a level up, to the state machine setting and turning off lights.
	var buttonMatrix [][]int
	buttonMatrix = make([][]int, driver.N_FLOORS)
	for i := 0; i < driver.N_FLOORS; i++ {
		buttonMatrix[i] = make([]int, driver.N_BUTTONS) //Golang creates a slice of zeros by default
	}
	fmt.Print(buttonMatrix)
	go func
	//Continious checking of buttons, buttonChan is a buffered channel who can fit N_FLOORS*N_BUTTONS elements
	for {
		time.Sleep(time.Millisecond * 20) //Need a proper time to wait.
		for i := 0; i < driver.N_FLOORS; i++ {
			for j := 0; j < driver.N_BUTTONS; j++ {
				if buttonVar := driver.Elev_get_button_signal(i, j); buttonVar != buttonMatrix[i][j] { //Sending the struct if its pushed and hasnt been sent already
					fmt.Println("Here is drivers version of button pressed: ", buttonVar)
					fmt.Println("floor and button:", i, j)
					if buttonVar == 1 {

						buttonSlice[i] <- Button{i, j, true} //This might give an error
					}
					buttonMatrix[i][j] = buttonVar
				}
			}
		}
	}
}
*/
//Slave should only send to buttonUpdated if something comes from above, ie not from button_updater
func Button_updater(externalButtPressChan chan Button, internalButtPressChan chan Button, buttonUpdatedChan chan Button) { //Sending the struct a level up, to the state machine setting and turning off lights.
	buttonMatrix := make([][]int, driver.N_FLOORS)
	for i := 0; i < driver.N_FLOORS; i++ {
		buttonMatrix[i] = make([]int, driver.N_BUTTONS) //Golang creates a slice of zeros by default
	}
	//fmt.Print(buttonMatrix)
	go func() {
		for {
			butt := <-buttonUpdatedChan //Word from above that some button is updated
			if butt.turnOn {
				buttonMatrix[butt.floor][butt.buttonType] = 1
			} else {
				buttonMatrix[butt.floor][butt.buttonType] = 0
			}
		}
	}()
	//Continious checking of buttons, buttonChan is a buffered channel who can fit N_FLOORS*N_BUTTONS elements
	for {
		time.Sleep(time.Millisecond * 40) //Need a proper time to wait.
		for i := 0; i < driver.N_FLOORS; i++ {
			for j := 0; j < driver.N_BUTTONS; j++ {
				if buttonVar := driver.Elev_get_button_signal(i, j); buttonVar != buttonMatrix[i][j] { //Sending the struct if its pushed and hasnt been sent already
					//fmt.Println("Here is drivers version of button pressed: ", buttonVar)
					//fmt.Println("floor and button:", i, j)
					if buttonVar == 1 && j != driver.COMMAND {

						externalButtPressChan <- Button{i, j, true} //   YO!   Need to make this one sexier. Maybe one channel for each button
					} else if buttonVar == 1 && j == driver.COMMAND {
						internalButtPressChan <- Button{i, j, true}
					}
					buttonMatrix[i][j] = 1 //Confirm press to avoid spamming
				}
			}
		}
	}
}

//Try having a dedicated channel for each floor. Light updater will receive a floor command, and set the light on or off
//This will receive commands from two different holds, and only one will be served. I dont think this will be a problem

//If we get time: see how we can make this more dynamic. Yhis also turns off if the bool is false.
func Light_updater(setLightChan chan Button, buttUpdatedChan chan Button) {
	for {
		butt := <-setLightChan
		_ = driver.Elev_set_button_lamp(butt.floor, butt.buttonType, butt.turnOn)
		buttUpdatedChan <- butt
	}
}

func Motor_control(speedChan chan float64) { //I think speedchan should not be buffered
	for {
		speedVal := <-speedChan
		if speedVal == 0 {
			driver.Elev_stop_elevator()
		} else {
			driver.Elev_set_speed(speedVal)
		}
	}
}

// Gets sensor signal and tells which floor is the current
func Is_floor_reached(sensorChan chan int) {
	var previousFloor int = -1
	for {
		if currentFloor := driver.Elev_get_floor_sensor_signal(); currentFloor != -1 && currentFloor != previousFloor {
			driver.Elev_set_floor_indicator(currentFloor)
			previousFloor = currentFloor
			sensorChan <- currentFloor

		}
		time.Sleep(time.Millisecond * 250)
	}
}

func Open_door() {
	driver.Elev_set_door_open_lamp(true)
	time.Sleep(time.Second * 3)
	driver.Elev_set_door_open_lamp(false)
}

/*func main() {
	var err error = driver.Elev_init()
	sensorlightchan := make(chan int)
	buttonSliceChan := make
	if err != nil {
		fmt.Println(err)
	}
	go light_setter(sensorlightchan)
	go light_getter(sensorlightchan)
	driver.Elev_set_speed(-300)

	time.Sleep(time.Second * 8)
	driver.Elev_set_speed(0)

}
*/
