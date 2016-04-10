package topLevel

import (
	"../communication"
	"../driver"
	"../elevatorStateTracker"
	. "../globalChans"
	. "../globalStructs"
	"../optalg"
	"../stateMachine"
	"fmt"
	"time"
)

/*
TOP LEVEL FUNCTIONS GO HERE. WE WILL RENAME AND MOVE FUNCTIONS ONCE WE HAVE AN OVERVIEW OF HOW MUCH WE WILL END UP WITH
*/

type AssignedOrderAndElevList struct {
	OrdAss   OrderAssigned
	ElevList []string
}

var externalOrdersNotTaken []Order
var commonExternalArray [NUM_FLOORS][NUM_BUTTONS - 1]int

//ElevatorAssignedToNetworkChan OrderAssigned chan , ElevatorAssignedFromNetworkChan OrderAssigned chan

func AssignOrdersAndWaitForAgreement(newOrderFromNetworkChan chan Order, resetAssignFuncChan chan bool, newOrderAssignedChan chan OrderAssigned) {

	var OrdersToBeAssignedByAll []AssignedOrderAndElevList
	localAddr := communication.GetLocalIP()

	for {
		select {
		case newOrder := <-newOrderFromNetworkChan:
			fmt.Println("newOrder := <-newOrderFromNetworkChan:")
			orderIsRegistered := false
			fmt.Println("Iterating..")
			for _, v := range OrdersToBeAssignedByAll {
				if v.OrdAss.Order == newOrder {
					orderIsRegistered = true
				}
			}
			fmt.Println("done searching for order")
			if !orderIsRegistered {
				//for now, the elevator states are all globally known. May send copy or something else later.
				assignedElevAddr := optalg.OptAlg(newOrder)
				fmt.Println("optalg complete")
				NewOrderToBeAssigned := OrderAssigned{Order: newOrder, AssignedTo: assignedElevAddr, SentFrom: localAddr}
				//Elevlist should be copied, global or maybe everyone that uses it should be in the same module
				elevList := communication.GetElevList()
				OrdersToBeAssignedByAll = append(OrdersToBeAssignedByAll, AssignedOrderAndElevList{NewOrderToBeAssigned, elevList})
				time.Sleep(time.Millisecond * 100) //This is to make sure you get to make the list before
				ToNetworkOrderAssignedToChan <- NewOrderToBeAssigned
			} else {
				fmt.Println("Order already registered. Discard message.")
			}
		case newOrdAss := <-FromNetworkOrderAssignedToChan:
			fmt.Println("newOrdAss := <-FromNetworkOrderAssignedToChan")
			posInSlice := -1
			for i, v := range OrdersToBeAssignedByAll {
				if newOrdAss.Order == v.OrdAss.Order {
					posInSlice = i
					break
				}
			}
			if posInSlice == -1 {
				fmt.Println("Old/garbage order")
			} else {
				fmt.Println("posInSlice: ",posInSlice)
				if newOrdAss.AssignedTo != OrdersToBeAssignedByAll[posInSlice].OrdAss.AssignedTo {
					fmt.Println("Disagreement, recalculate with optalg")
					OrdersToBeAssignedByAll = append(OrdersToBeAssignedByAll[:posInSlice], OrdersToBeAssignedByAll[posInSlice+1:]...) //slicetricks
					assignedElevAddr := optalg.OptAlg(newOrdAss.Order)
					NewOrderToBeAssigned := OrderAssigned{Order: newOrdAss.Order, AssignedTo: assignedElevAddr, SentFrom: localAddr}
					elevList := communication.GetElevList()
					OrdersToBeAssignedByAll = append(OrdersToBeAssignedByAll, AssignedOrderAndElevList{NewOrderToBeAssigned, elevList})
					time.Sleep(time.Millisecond * 20) //This is to make sure you get to make the list before
					ToNetworkOrderAssignedToChan <- NewOrderToBeAssigned
				} else {
					for i, v := range OrdersToBeAssignedByAll[posInSlice].ElevList {
						if newOrdAss.SentFrom == v {
							OrdersToBeAssignedByAll[posInSlice].ElevList = append(OrdersToBeAssignedByAll[posInSlice].ElevList[:i], OrdersToBeAssignedByAll[posInSlice].ElevList[i+1:]...)
							if len(OrdersToBeAssignedByAll[posInSlice].ElevList) == 0 {
								OrdersToBeAssignedByAll = append(OrdersToBeAssignedByAll[:posInSlice], OrdersToBeAssignedByAll[posInSlice+1:]...)
								newOrderAssignedChan <- newOrdAss
							}
						}
					}
				}
				fmt.Println("ENDOF:    newOrdAss := <-FromNetworkOrderAssignedToChan ")
			}
		case <-resetAssignFuncChan:
			OrdersToBeAssignedByAll = nil
		}
	}
}

/*
func UpdateNetworkCondition() {
	networkIsUp := false
	for {
		select {
		case <-FromNetworkNetworkUpChan:
			networkIsUp = true
		case <-FromNetworkNetworkDownChan:
			networkIsUp = false
		}

	}
}/*

//This may not be used, yo
/*
func SendAllOrdersOfQueueToNetwork(commonExternalArray [][] int){
	for i := 0; i < NUM_FLOORS; i++{
		for j := 0; j < NUM_BUTTONS -1; j++{
			if commonExternalArray[i][j] == 1 {
				direction := UP
				if j != 0{
					direction = DOWN
				}
				newOrder := ExternalOrder{Floor: i, Direction: direction}
				//Rename channel
				newExternalOrderToNetwork <- newOrder
				//Maybe a tiny sleep here, maybe not
			}
			commonExternalArray[i][j] = 0
		}
	}
}
*/

func orderIsInQueue(orderQueue []Order, newOrder Order) bool {
	for _, v := range orderQueue {
		if v == newOrder {
			return true
		}
	}
	return false
}

func TopLogicNeedBetterName() {
	var internalArray [NUM_FLOORS]int
	resetAssignFuncChan := make(chan bool)
	newOrderToBeAssignedChan := make(chan Order)
	orderDoneAssignedChan := make(chan OrderAssigned)
	go ResendOrdersWhenError(resetAssignFuncChan)
	go AssignOrdersAndWaitForAgreement(newOrderToBeAssignedChan, resetAssignFuncChan, orderDoneAssignedChan)

	fmt.Println(commonExternalArray)
	for {
		select {
		case newButton := <-InternalButtonPressedChan:
			if internalArray[newButton.Floor] == 0 {
				internalArray[newButton.Floor] = 1 //Reset when order served. Tell that total state is changed?
				elevatorStateTracker.WriteCurrentOrdersToFile(AllOrders{InternalOrders: internalArray, ExternalOrders: commonExternalArray})
				driver.SetButtonLight(newButton, true)
				NewOrderToLocalElevChan <- newButton
			}
		case newOrder := <-FromNetworkNewOrderChan:
			dir := UP
			if newOrder.Direction == DOWN {
				dir = 0
			}
			if commonExternalArray[newOrder.Floor][dir] == 0 {
				commonExternalArray[newOrder.Floor][dir] = 1
			}
			elevatorStateTracker.WriteCurrentOrdersToFile(AllOrders{InternalOrders: internalArray, ExternalOrders: commonExternalArray})

			if !orderIsInQueue(externalOrdersNotTaken, newOrder) {
				externalOrdersNotTaken = append(externalOrdersNotTaken, newOrder)
			}
			driver.SetButtonLight(newOrder, true)

			fmt.Println("TOPLOGIC: to assign func")
			newOrderToBeAssignedChan <- newOrder //the receiver will handle duplicates
			driver.SetButtonLight(newOrder, true)

		case newOrdAss := <-orderDoneAssignedChan:
			fmt.Println("newOrdAss := <-orderDoneAssignedChan:")
			AddOrderAssignedToElevStateChan <- newOrdAss
			for i, v := range externalOrdersNotTaken {
				if v == newOrdAss.Order {
					externalOrdersNotTaken = append(externalOrdersNotTaken[:i], externalOrdersNotTaken[i+1:]...)
				}
			}
			if newOrdAss.AssignedTo == communication.GetLocalIP() {
				NewOrderToLocalElevChan <- newOrdAss.Order
			}

		case servedOrder := <-OrderServedLocallyChan:
			//Denne tar imot alle ordre. Hvis nettet er oppe, skal man også sende videre til nett.
			if(servedOrder.Direction == 0){
				internalArray[servedOrder.Floor] = 0
			}else{
				dir := UP
				if servedOrder.Direction == DOWN {
					dir = 0
				}
				commonExternalArray[servedOrder.Floor][dir] = 0
				ToNetworkOrderServedChan <- servedOrder
			}
			elevatorStateTracker.WriteCurrentOrdersToFile(AllOrders{InternalOrders: internalArray, ExternalOrders: commonExternalArray})
			driver.SetButtonLight(servedOrder, false)

		case servedOrder := <-FromNetworkOrderServedChan:
			fmt.Println("servedOrder := <-FromNetworkOrderServedChan:")
			dir := UP
			if servedOrder.Direction == DOWN {
				dir = 0
			}
			commonExternalArray[servedOrder.Floor][dir] = 0
			elevatorStateTracker.WriteCurrentOrdersToFile(AllOrders{InternalOrders: internalArray, ExternalOrders: commonExternalArray})
			driver.SetButtonLight(servedOrder, false)

		case <-FromNetworkNewElevChan:
			resetAssignFuncChan <- true
			stateMachine.SendElevStateToNetwork(ToNetworkNewElevStateChan)
			ToNetworkExternalArrayChan<-commonExternalArray
		case newExternalArray := <-FromNetworkExternalArrayChan:
			for i := 0; i < NUM_FLOORS; i++{
				for j:=0; j < NUM_BUTTONS-1; j++{
					if newExternalArray[i][j] == 1{
						commonExternalArray[i][j] = 1
						dir := UP
						if j == 0{
							dir = DOWN
						}
						driver.SetButtonLight(Order{Floor: i, Direction: dir},true)
					}
				}
			}
			elevatorStateTracker.WriteCurrentOrdersToFile(AllOrders{InternalOrders: internalArray, ExternalOrders: commonExternalArray})
		}
	}
}

//MULIG VI MÅ FLYTTE DE SOM TAR IMOT NYE ORDRE OG DE SOM SENDER NYE ORDRE I FORSKJELLIGE LOOPS
//PROBLEM NÅR MANGE ORDRE BLIR SENDT PÅ NYTT


//CHANGE THIS FUCKING NAME
func ResendOrdersWhenError(resetAssignFuncChan chan bool) {
	for {
		select {
		case <-FromNetworkNetworkDownChan:
			resetAssignFuncChan <- true
			externalOrdersNotTaken = nil
			for i := 0; i < NUM_FLOORS; i++ {
				for j := 0; j < NUM_BUTTONS-1; j++ {
					if commonExternalArray[i][j] == 1 {
						dir := UP
						if j == 0 {
							dir = DOWN
						}
						NewOrderToLocalElevChan <- Order{Floor: i, Direction: dir}
						time.Sleep(time.Millisecond*200)
					}
				}
			}
		case elevGone := <-FromNetworkElevGoneChan:
			fmt.Println("ResendOrdersWhenError: elevGone")
			deadElevOrders := optalg.GetOrderQueueOfDeadElev(elevGone)
			fmt.Println("Orders of dead elev:")
			for _, v := range deadElevOrders{
				stateMachine.PrintOrder(v)
			}
			ToOptAlgDeleteElevChan <- elevGone
			resetAssignFuncChan <- true
			for _, v := range deadElevOrders {
				ExternalButtonPressedChan <- v
				time.Sleep(time.Millisecond*40)
			}
			for _, v := range externalOrdersNotTaken {
				ExternalButtonPressedChan <- v
				time.Sleep(time.Millisecond*40)
			}
		}
	}
}
