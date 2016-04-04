package globalStructs

const(
	UP          = 1
	DOWN        = -1
	COMMAND     = 0
	NUM_BUTTONS = 3
	NUM_FLOORS  = 4
	MOTOR_SPEED = 2800
)


type Button struct {
	Floor          int
	Button_type    int
	Button_pressed bool
}

type OrderAssigned struct{
	AssignedTo string
	SentFrom string
	Order Order
}

type Order struct {
	Floor int
	Direction int 
}

//This is a struct you could send to all other elevators
type ElevatorState struct {
	IP            string //not an int, always useful to have
	CurrentFloor  int
	PreviousFloor int
	Direction     int
	OrderQueue    []Order //This is an array or something of all orders currently active for this elevator.

}

type AllOrders struct{
	ExternalOrders [NUM_FLOORS][NUM_BUTTONS-1] int
	InternalOrders [NUM_FLOORS] int
}