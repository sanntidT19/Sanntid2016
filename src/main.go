package main

import (
	"./driver"
	"./globalChans"
	. "./globalStructs"
	"./messages"
	"./network"
	"./optalg"
	"./stateMachine"
	"./topLevel"
	"fmt"
	"time"
)

func main() {
	globalChans.InitChans()
	driver.Init()

	//global chans here
	newOrderToBeAssignedChan := make(chan Order)
	resetAssignFuncChan := make(chan bool)

	go topLevel.StartupDraft()
	go topLevel.TopLogicNeedBetterName(newOrderToBeAssignedChan, resetAssignFuncChan)
	go optalg.UpdateElevatorStateList(newOrderToBeAssignedChan, resetAssignFuncChan)
	go network.InitNetworkAndAlertChanges()
	go messages.MessagesTopAndWaitForNetworkChanges()
	go stateMachine.NewTopLoop()
	fmt.Println("Main: sleeping")
	time.Sleep(time.Second * 1000)
}
