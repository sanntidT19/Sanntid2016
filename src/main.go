package main

import (
	"./communication"
	"./driver"
	"./elevatorStateTracker"
	. "./globalChans"
	//. "./globalStructs"
	"./optalg"
	"./stateMachine"
	"./topLevel"

	"fmt"
	"time"
	//"encoding/gob"
	//"os"
)

func main() {
	InitChans()
	driver.ElevMainTesterFunction()
	go elevatorStateTracker.StartupDraft()
	go topLevel.TopLogicNeedBetterName()
	go optalg.UpdateElevatorStateList()
	go communication.CommNeedBetterName()
	go stateMachine.NewTopLoop()
	fmt.Println("Main: sleeping")
	time.Sleep(time.Second * 1000)
}
