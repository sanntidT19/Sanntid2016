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

	//Make all global channels
	newOrderToBeAssignedChan := make(chan Order) //rename slightly
	resetAssignFuncChan := make(chan bool)

	internalButtonChan := make(chan Order)
	newLocalOrderChan := make(chan Order)
	localOrderServedChan := make(chan Order)

	newConnectionChan :=  make(chan string)
	endConnectionChan := make(chan string)
	resendUnackdMessagesChan := make(chan bool)

	newElevChan := make(chan string)
	deadElevChan := make(chan string)
	networkDownChan := make(chan bool)

	removeElevChan := make(chan string)

	receiveNewOrderChan := make(chan Order)
	receiveOrderAssChan := make(chan OrderAssigned)
	receiveOrderServedChan := make(chan Order)
	receiveElevStateChan := make(chan ElevatorState)
	receiveExternalArrayChan := make(chan [NUM_FLOORS][NUM_BUTTONS -1]int)

	fromDecodeChans := MessageChans{
		NewOrderChan : receiveNewOrderChan,
		OrderAssChan : receiveOrderAssChan,
		OrderServedChan : receiveOrderServedChan,
		ElevStateChan : receiveElevStateChan,
		ExternalArrayChan : receiveExternalArrayChan}

	sendNewOrderChan := make(chan Order)
	sendOrderAssChan := make(chan OrderAssigned)
	sendOrderServedChan := make(chan Order)
	sendElevStateChan := make(chan ElevatorState)
	sendExternalArrayChan := make(chan [NUM_FLOORS][NUM_BUTTONS -1]int)

	toEncodeChans := MessageChans{
		NewOrderChan : sendNewOrderChan,
		OrderAssChan : sendOrderAssChan,
		OrderServedChan : sendOrderServedChan,
		ElevStateChan : sendElevStateChan,
		ExternalArrayChan : sendExternalArrayChan}

	driver.Init(sendNewOrderChan,internalButtonChan)
	go topLevel.StartupDraft(sendNewOrderChan,newLocalOrderChan)
	//hvis tid: samle de i toplogic i en struct (newOrderToBeAssignedChan, sendNewOrderChan, receiveNewOrderChan, sendOrderServedChan, receiveOrderServedChan, sendElevStateChan, sendExternalArrayChan, receiveExternalArrayChan, internalButtonChan, newLocalOrderChan, localOrderServedChan, newElevChan)
	go topLevel.TopLogicNeedBetterName(newOrderToBeAssignedChan, sendNewOrderChan, receiveNewOrderChan, sendOrderServedChan, receiveOrderServedChan, sendElevStateChan, sendExternalArrayChan, receiveExternalArrayChan, internalButtonChan, newLocalOrderChan,localOrderServedChan, newElevChan)
	go topLevel.ResendOrdersIfNetworkError(resetAssignFuncChan, sendNewOrderChan,newLocalOrderChan,deadElevChan,networkDownChan,removeElevChan)
	go optalg.UpdateElevatorStateList(newOrderToBeAssignedChan, resetAssignFuncChan, sendOrderAssChan, receiveOrderAssChan, receiveElevStateChan,newLocalOrderChan,removeElevChan)
	go network.InitNetworkAndAlertChanges(newConnectionChan, endConnectionChan, resendUnackdMessagesChan, newElevChan, deadElevChan, networkDownChan)
	go messages.MessagesTopAndWaitForNetworkChanges(fromDecodeChans,toEncodeChans, newConnectionChan, endConnectionChan,resendUnackdMessagesChan)
	go stateMachine.NewTopLoop(sendElevStateChan, newLocalOrderChan, localOrderServedChan)
	fmt.Println("Main: sleeping")
	time.Sleep(time.Second * 1000)
}
