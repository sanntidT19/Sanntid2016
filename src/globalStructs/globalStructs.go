package globalStructs


type Button struct {
	floor int
	button_type int
	pressed bool
}


//This is a struct you could send to all other elevators
type ElevatorState struct {
	myIP int //not an int, always useful to have
	currentFloor int
	lastFloor int
	direction int
	orders int //This is an array or something of all orders currently active for this elevator. 
}


