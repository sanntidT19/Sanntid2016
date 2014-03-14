package statemachine

import (
	"driver"
	"time"
)

//Functions used when running the elevator, find out better name and add prefix



//This one creates the basic button array for our friends
func create_button_chan_slice()  [] chan Button //A little unsure of this
{
	chanSlice := make([] chan Button, N_FLOORS)
	var i = 0
	for i; i < N_FLOORS; i++ {
		chanArray[i] = make(chan Button)
	}
	return chanSlice
}

//if we get errors, this bool might be the bad guy
type Button struct {
	floor      int
	buttonType int
	turnOn     bool
}

func button_updater(buttonSlice [] chan Button) { //Sending the struct a level up, to the state machine setting and turning off lights.
	var button_matrix [][]int
	button_matrix = make([][]int, driver.N_FLOORS)
	var i int = 0
	var j int
	for i; i < driver.N_FLOORS; i++ {
		button_matrix[i] = make([]int, driver.N_BUTTONS) //Golang creates a slice of zeros by default
	}
	//Continious checking of buttons, buttonChan is a buffered channel who can fit N_FLOORS*N_BUTTONS elements
	for {
		sleep(time.Millisecond*20)  //Need a proper time to wait.
		i = 0
		for i; i < driver.N_FLOORS; i++ {
			j = 0
			for j; j < driver.N_BUTTONS; j++ {
				if buttonVar = driver.Elev_get_button_signal(i, j); buttonVar != button_matrix[i][j] { //Sending the struct if its pushed and hasnt been sent already
					if buttonVar == 1 {
						buttonSlice[i] <- Button(i, j, true)  //This might give an error
					}
					button_matrix[i][j] = buttonVar
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


func floor_getter(sensorlightchan chan int) {
	for {
		if c := driver.Elev_get_floor_sensor_signal(); c != -1 {
			driver.Elev_set_floor_indicator(c)
			sensorlightchan <- c
		}
		time.Sleep(time.Millisecond * 100)
	}
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