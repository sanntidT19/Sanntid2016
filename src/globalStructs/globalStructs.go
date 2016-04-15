package globalStructs
import(
	"time"
)
const (
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

type OrderAssigned struct {
	AssignedTo string
	SentFrom   string
	Order      Order
}

type Order struct {
	Floor     int
	Direction int
}

//This is a struct you could send to all other elevators
//Remember to set timestamp. And check the timestamp when a new state comes to optalgtester
type ElevatorState struct {
	IP            string //not an int, always useful to have
	CurrentFloor  int
	PreviousFloor int
	Direction     int
	OrderQueue    []Order //This is an array or something of all orders currently active for this elevator
	Timestamp time.Time 
}

type AllOrders struct {
	ExternalOrders [NUM_FLOORS][NUM_BUTTONS - 1]int
	InternalOrders [NUM_FLOORS]int
}

type MessageChans struct{
	NewOrderChan chan Order
	OrderAssChan chan OrderAssigned
	OrderServedChan chan Order
	ElevStateChan chan ElevatorState
	ExternalArrayChan chan [NUM_FLOORS][NUM_BUTTONS -1] int
}

func orderIsInList(orderQueue []Order, newOrder Order) bool {
	for _, v := range orderQueue {
		if v == newOrder {
			return true
		}
	}
	return false
}
