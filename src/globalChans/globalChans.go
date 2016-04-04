package globalChans

import (
	"../globalStructs"
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


func InitChans() {
	ButtonPressedChan = make(chan globalStructs.Button)
	SetButtonLightChan = make(chan globalStructs.Button)
	
	FromNetworkNewOrderChan = make(chan globalStructs.Order)
	FromNetworkOrderAssignedToChan = make(chan globalStructs.OrderAssigned)
	FromNetworkOrderServedChan = make(chan globalStructs.Order)
	FromNetworkNewElevStateChan = make(chan globalStructs.ElevatorState)

	FromNetworkNewElevChan = make(chan string)
	FromNetworkElevGoneChan =make(chan string)
	FromNetworkNetworkDownChan = make(chan bool)
	FromNetworkNetworkUpChan = make(chan bool)

	ToNetworkNewOrderChan = make(chan globalStructs.Order)
	ToNetworkOrderAssignedToChan = make(chan globalStructs.OrderAssigned)
	ToNetworkOrderServedChan = make(chan globalStructs.Order)
	ToNetworkNewElevStateChan = make(chan globalStructs.ElevatorState)

	NewOrderToLocalElevChan = make(chan globalStructs.Order)
	
	//To optalg for now
	ToOptAlgDeleteElevChan = make(chan string)
	AddOrderAssignedToElevStateChan = make(chan globalStructs.OrderAssigned)
}