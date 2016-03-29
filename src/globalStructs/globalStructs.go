package globalStructs

type Button struct {
	Floor          int
	Button_type    int
	Button_pressed bool
}

type OrderAssigned struct{
	AssignedTo string
	SentFrom string
	Order ExternalOrder
}

type ExternalOrder struct {
	Floor int
	Direction int 
}

//This is a struct you could send to all other elevators
type ElevatorState struct {
	IP            string //not an int, always useful to have
	CurrentFloor  int
	PreviousFloor int
	Direction     int
	OrderQueue    []Button //This is an array or something of all orders currently active for this elevator.

}

const (
	UP          = 1
	DOWN        = -1
	COMMAND     = 0
	NUM_BUTTONS = 3
	NUM_FLOORS  = 4
	MOTOR_SPEED = 2800
)
