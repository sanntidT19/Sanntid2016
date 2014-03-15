package statemachine

import (
	"driver"
	"time"
)


const MAX_SPEED_UP = 300
const MAX_SPEED_DOWN = -300
const SPEED_STOP = 0

//if we get errors, this bool might be the bad guy
type Button struct {
	floor      int
	buttonType int
	turnOn     bool
}

//Functions used when running the elevator, find out better name and add prefix


func state_machine(goToFloorChan chan int, currentFloorChan chan int, currentDirChan chan int, buttonSliceChan [] chan int, lightSliceChan [] chan int, servedOrderChan chan bool) {
	driver.Elev_init() //Initiates the 
	motorChan := make(chan int)  //This channel is only between the statemachine and its functions
	privateSensorChan := make(chan int)
	var currentFloor =: driver.Elev_get_floor_sensor_signal() //We know that we are in a floor at this point.
	var previousFloor := -1
	var gtf int 
	go motor_control(speedChan)
	go button_updater(buttonSliceChan)
	go light_updater(lightSliceChan)
	go floor_reached(privateSensorChan)
	// This is where the statemachine is implemented, should it be a select case?
	


	go func() {
		for{
			select{
			//Slave could send a new command while the statemachine is serving another command, but it should fix the logic by itself
			// New order
			case gtf <-goToFloorChan:
				//You are in the floor, order served immediatly, maybe this if can be implemented in another case, but its here for now.
				if gtf == driver.Elev_get_floor_sensor_signal() {
					speedChan <- SPEED_STOP
					servedOrderChan <- true
					open_door() //Dont think we want this select loop to do anything else while the door is open. Solve with go open_door() if its not the case
				//You know you are under/above the current floor
				}else if gtf < currentFloor {
					speedChan <- MAX_SPEED_DOWN
					currentDirChan <- -1
				}else if gtf > currentFloor{
					speedChan <- MAX_SPEED_UP
					currentDirChan <- 1
				//Your last floor was the current floor, but something may have been pulled, so you dont know where you lie relative to it. Cant use direction.
				}else if gtf == currentFloor{
					//Using previousfloor can give you an idea in some cases.
					if previousFloor > currentFloor{
						speedchan <- MAX_SPEED_UP
						currentDirChan <- 1
					}else if previousFloor < currentFloor{  //This will also be the case if prevFloor is undefined (-1) They can never be the same.
						speedchan <- MAX_SPEED_DOWN
						currentDirChan <- -1
					//If someone has dragged
				}
			//New floor is reached and therefore shit is updated
			case cf := <-privateSensorChan:
				previousFloor = currentFloor
				currentFloor = cf
			case gtf == currentFloor{ //This will go immediatly after the case above has done its job, hopefully. Maybe need some extra care here.
				speedchan <- SPEED_STOP
				servedOrderChan <- true
				open_door()

			}




			}
			
		}

	}
}

//This one creates the basic button slice for our friends
func create_button_chan_slice()  [] chan Button //A little unsure of this
{
	chanSlice := make([] chan Button, N_FLOORS)
	var i = 0
	for i; i < N_FLOORS; i++ {
		chanSlice[i] = make(chan Button)
	}
	return chanSlice
}



func button_updater(buttonSlice [] chan Button) { //Sending the struct a level up, to the state machine setting and turning off lights.
	var buttonMatrix [][]int
	buttonMatrix = make([][]int, driver.N_FLOORS)
	var i int = 0
	var j int
	for i; i < driver.N_FLOORS; i++ {
		buttonMatrix[i] = make([]int, driver.N_BUTTONS) //Golang creates a slice of zeros by default
	}
	//Continious checking of buttons, buttonChan is a buffered channel who can fit N_FLOORS*N_BUTTONS elements
	for {
		sleep(time.Millisecond*20)  //Need a proper time to wait.
		i = 0
		for i; i < driver.N_FLOORS; i++ {
			j = 0
			for j; j < driver.N_BUTTONS; j++ {
				if buttonVar = driver.Elev_get_button_signal(i, j); buttonVar != buttonMatrix[i][j] { //Sending the struct if its pushed and hasnt been sent already
					if buttonVar == 1 {
						buttonSlice[i] <- Button(i, j, true)  //This might give an error
					}
					buttonMatrix[i][j] = buttonVar
				}
			}
		}
	}
}

//Try having a dedicated channel for each floor. Light updater will receive a floor command, and set the light on or off
//This will receive commands from two different holds, and only one will be served. I dont think this will be a problem

//If we get time: see how we can make this more dynamic.
func light_updater(buttonSlice [] chan Button){
	for{
		select{
			case butt := <- buttonSlice[0]:
				driver.Elev_set_button_lamp(butt.floor, butt.ButtonType, butt.TurnOn)
			case butt := <- buttonSlice[1]:
				driver.Elev_set_button_lamp(butt.floor, butt.ButtonType, butt.TurnOn)
			case butt := <- buttonSlice[2]:
				driver.Elev_set_button_lamp(butt.floor, butt.ButtonType, butt.TurnOn)
			case butt := <- buttonSlice[3]:
				driver.Elev_set_button_lamp(butt.floor, butt.ButtonType, butt.TurnOn)
		}
		
	}
}


func motor_control(speedChan chan float64) { //I think speedchan should not be buffered
	for {
		speedVal := <-speedChan
		driver.Elev_set_speed(speedVal)
	}
}



// Gets sensor signal and tells which floor is the current
func floor_reached(sensorChan chan int) {
	var previousFloor int = -1
	for {
		if currentFloor := driver.Elev_get_floor_sensor_signal(); currentFloor != -1 && currentFloor != previousFloor {
			driver.Elev_set_floor_indicator(currentFloor)
			previousFloor = currentFloor
			sensorChan <- currentFloor

		}
		time.Sleep(time.Millisecond * 100)
	}
}

func open_door(){
	driver.Elev_set_door_open_lamp(true)
	time.Sleep(time.Second*3)
	drive.Elev_set_door_open_lamp(false)
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

