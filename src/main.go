package main

import (
	"./driver"
	"./elevatorStateTracker"
	"./globalChans"
	//. "./globalStructs"
	"./optalg"
	"./stateMachine"
	"./topLevel"
	"./network"
	"./messages"
	"fmt"
	"time"
	//"encoding/gob"
	//"os"
)

func main() {
	globalChans.InitChans()
	driver.Init()
	go elevatorStateTracker.StartupDraft()
	go topLevel.TopLogicNeedBetterName()
	go optalg.UpdateElevatorStateList()
	go network.InitNetworkAndAlertChanges()
	go messages.MessagesTopAndWaitForNetworkChanges()
	go stateMachine.NewTopLoop()
	fmt.Println("Main: sleeping")
	time.Sleep(time.Second * 1000)
}
