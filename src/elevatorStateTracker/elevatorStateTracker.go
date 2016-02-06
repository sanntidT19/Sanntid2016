import(
	"../globalStructs"
)

currentState ElevatorState


func initalize_state_tracker(){
	//read from file to check if system was killed
	//easy solution: if thats the case, set current state to that (may serve same order twice, but sverre wont die and its avoiding complicated solutions)
	//if not, initialize normally	
	// get floor and all that shit from other modules 
	//send the current state to everybody
}

func send_updated_elevator_state(){
	//call this whenever its updated 
}

func write_elevator_state_to_file(){
	//update this whenever the local elevator gets an order/command
}

func read_elevator_state_from_file(){
}

func detect_system_killed(){
	//maybe a different name
}