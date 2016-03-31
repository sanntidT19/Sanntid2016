package globalChans

import (
	"../globalStructs"
)

var ButtonPressedChan chan globalStructs.Button
var SetButtonLightChan chan globalStructs.Button

var FromNetworkNewOrderChan chan globalStructs.Order
var FromNetworkOrderAssignedToChan chan globalStructs.OrderAssigned
var FromNetworkOrderServedChan chan globalStructs.Order
var FromNetworkElevlistChangedChan chan //TO BE FIXED
var FromNetworkNewElevStateChan chan globalStructs.ElevatorState


var ToNetworkNewOrderChan chan globalStructs.Order
var ToNetworkOrderAssignedToChan chan globalStructs.OrderAssigned
var ToNetworkOrderServedChan chan globalStructs.Order //Maybe external
var ToNetworkNewElevStateChan chan globalStructs.ElevatorState

func Init_chans() {
	ButtonPressedChan = make(chan globalStructs.Button)
	SetButtonLightChan = make(chan globalStructs.Button)
	
	FromNetworkNewOrderChan = make(chan globalStructs.Order)
	FromNetworkOrderAssignedToChan = make(chan globalStructs.OrderAssigned)
	FromNetworkOrderServedChan = make(chan globalStructs.Order)
	FromNetworkNewElevStateChan = make(chan globalStructs.ElevatorState)

	/*
	NewOrderFromNetWorkChan*/
	ToNetworkNewOrderChan = make(chan globalStructs.Order)
	ToNetworkOrderAssignedToChan = make(chan globalStructs.OrderAssigned)
	ToNetworkOrderServedChan = make(chan globalStructs.Order)
	ToNetworkNewElevStateChan = make(chan globalStructs.ElevatorState)

}

/*
global chans som må fikses:
Externalbuttonpressedchan(driver og network)
Internalbuttonpressedchan(driver og toplevel)



ResetExternalOrdersInQueueChan (statemachine og toplevel)
/*