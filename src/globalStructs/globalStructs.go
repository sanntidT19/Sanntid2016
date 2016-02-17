package globalStructs

type Button struct {
	Floor          int
	Button_type    int
	Button_pressed bool
}

//This is a struct you could send to all other elevators
type ElevatorState struct {
	MyIP         int //not an int, always useful to have
	CurrentFloor int
	LastFloor    int
	Direction    int
	Orders       int //This is an array or something of all orders currently active for this elevator.

}

const (
	UP          = 1
	DOWN        = -1
	COMMAND     = 0
	NUM_BUTTONS = 3
	NUM_FLOORS  = 4
	MOTOR_SPEED = 2800
)
