import(
	"../globalStructs"
	"encoding/gob"
	"os"
	"fmt"
)

currentState ElevatorState

const PATH_OF_SAVED_STATE = "elevState.gob"

func initalize_state_tracker(){
	//read from file to check if system was killed
	//easy solution: if thats the case, set current state to that (may serve same order twice, but sverre wont die and its avoiding complicated solutions)
	//if not, initialize normally	
	// get floor and all that shit from other modules 
	//send the current state to everybody
}

func send_updated_elevator_state(){
	//call this whenever its updated. write to channels
}

func write_elevator_state_to_file(){
	//temp for testing
	test_struct := ElevatorState{255,2,1,1,0}
	//update this whenever the local elevator gets an order/command
	dataFile, err := os.Create(PATH_OF_SAVED_STATE)
	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}

	dataEncoder := gob.NewEncoder(dataFile)
	dataEncoder.Encode(test_struct)
	dataFile.Close()
}

func read_elevator_state_from_file(){
	//start with reading it
	var data ElevatorState

	if _, err := os.Stat("PATH_OF_SAVED_STATE"); os.IsNotExist(err){
		fmt.Println("Local save of elevator state not detected. It has been cleared/this is the first run on current PC")
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

	currentState = data
}

func detect_system_killed() bool{
	//only call when initializing, REDO WHEN WE KNOW HOW THE ORDERS WILL LOOK
	for i:= 0; i < NUM_FLOORS; i++{
		for j:= 0; j < NUM_BUTTONS; j++{
			/*if currentState.orders[i][j] != 0{
				return true
			}
			*/
			if currentState.orders == 1{
				return true
			}
		}
	}
	return false
}