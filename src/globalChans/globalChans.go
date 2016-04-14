package globalChans

import (
	"../globalStructs"
	"fmt"
)

var ExternalButtonPressedChan chan globalStructs.Order
var InternalButtonPressedChan chan globalStructs.Order

var FromNetworkNewOrderChan chan globalStructs.Order
var FromNetworkOrderAssignedToChan chan globalStructs.OrderAssigned
var FromNetworkOrderServedChan chan globalStructs.Order

//var FromNetworkElevlistChangedChan chan //TO BE FIXED
var FromNetworkNewElevStateChan chan globalStructs.ElevatorState

var FromNetworkNewElevChan chan string
var FromNetworkElevGoneChan chan string
var FromNetworkNetworkDownChan chan bool
var FromNetworkNetworkUpChan chan bool

var ToNetworkNewOrderChan chan globalStructs.Order
var ToNetworkOrderAssignedToChan chan globalStructs.OrderAssigned
var ToNetworkOrderServedChan chan globalStructs.Order //Maybe external
var ToNetworkNewElevStateChan chan globalStructs.ElevatorState

var NewOrderToLocalElevChan chan globalStructs.Order

var ToOptAlgDeleteElevChan chan string

var AddOrderAssignedToElevStateChan chan globalStructs.OrderAssigned

var OrderServedLocallyChan chan globalStructs.Order //Cover this in statemachine too

var FromNetworkExternalArrayChan chan [globalStructs.NUM_FLOORS][globalStructs.NUM_BUTTONS - 1]int
var ToNetworkExternalArrayChan chan [globalStructs.NUM_FLOORS][globalStructs.NUM_BUTTONS - 1]int

//var ToNetworkExternalButtonPressedChan chan globalStructs.Order

var ToMessagesNewElevChan chan string
var ToMessagesDeadElevChan chan string
var ToMessagesNetworkDownChan chan bool

func InitChans() {
	fmt.Println("initializing chans")
	FromNetworkNewOrderChan = make(chan globalStructs.Order)
	FromNetworkOrderAssignedToChan = make(chan globalStructs.OrderAssigned)
	FromNetworkOrderServedChan = make(chan globalStructs.Order)
	FromNetworkNewElevStateChan = make(chan globalStructs.ElevatorState)
	FromNetworkExternalArrayChan = make(chan [globalStructs.NUM_FLOORS][globalStructs.NUM_BUTTONS - 1]int)
	FromNetworkNewElevChan = make(chan string)
	FromNetworkElevGoneChan = make(chan string)
	FromNetworkNetworkDownChan = make(chan bool)
	FromNetworkNetworkUpChan = make(chan bool)

	ToMessagesNewElevChan = make(chan string)
	ToMessagesDeadElevChan = make(chan string)
	ToMessagesNetworkDownChan = make(chan bool)

	ExternalButtonPressedChan = make(chan globalStructs.Order)

	ToNetworkNewOrderChan = make(chan globalStructs.Order)
	ToNetworkOrderAssignedToChan = make(chan globalStructs.OrderAssigned)
	ToNetworkOrderServedChan = make(chan globalStructs.Order)
	ToNetworkNewElevStateChan = make(chan globalStructs.ElevatorState)
	ToNetworkExternalArrayChan = make(chan [globalStructs.NUM_FLOORS][globalStructs.NUM_BUTTONS - 1]int)

	ExternalButtonPressedChan = make(chan globalStructs.Order)

	NewOrderToLocalElevChan = make(chan globalStructs.Order)
	OrderServedLocallyChan = make(chan globalStructs.Order)
	ToOptAlgDeleteElevChan = make(chan string)
	AddOrderAssignedToElevStateChan = make(chan globalStructs.OrderAssigned)

	InternalButtonPressedChan = make(chan globalStructs.Order)

	//ToNetworkExternalButtonPressedChan = make(chan globalStructs.Order)

	ToMessagesNewElevChan = make(chan string)
	ToMessagesDeadElevChan = make(chan string)
	ToMessagesNetworkDownChan = make(chan bool)

	fmt.Println("finished initializing chans")
}
