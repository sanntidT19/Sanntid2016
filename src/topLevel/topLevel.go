package topLevel

import (
	"../driver"
	. "../globalChans"
	. "../globalStructs"
	"../network"
	"../optalg"
	"../stateMachine"
	"encoding/gob"
	"fmt"
	"os"
	"time"
)

/*
TOP LEVEL FUNCTIONS GO HERE. WE WILL RENAME AND MOVE FUNCTIONS ONCE WE HAVE AN OVERVIEW OF HOW MUCH WE WILL END UP WITH
*/

var externalOrdersNotTaken []Order
var commonExternalArray [NUM_FLOORS][NUM_BUTTONS - 1]int

//ElevatorAssignedToNetworkChan OrderAssigned chan , ElevatorAssignedFromNetworkChan OrderAssigned chan

func orderIsInQueue(orderQueue []Order, newOrder Order) bool {
	for _, v := range orderQueue {
		if v == newOrder {
			return true
		}
	}
	return false
}

//Kanskje dele denne opp i to hvis man er sikker på at de ikke deler en variabel av no slag
func TopLogicNeedBetterName(newOrderToBeAssignedChan chan Order, resetAssignFuncChan chan bool) {
	var internalArray [NUM_FLOORS]int

	go ResendOrdersWhenError(resetAssignFuncChan)

	fmt.Println(commonExternalArray)
	for {
		select {
		case newButton := <-InternalButtonPressedChan:
			if internalArray[newButton.Floor] == 0 {
				internalArray[newButton.Floor] = 1 //Reset when order served. Tell that total state is changed?
				WriteCurrentOrdersToFile(AllOrders{InternalOrders: internalArray, ExternalOrders: commonExternalArray})
				driver.SetButtonLight(newButton, true)
				NewOrderToLocalElevChan <- newButton
			}
		case newOrder := <-FromNetworkNewOrderChan:
			fmt.Println("new order incoming")
			dir := UP
			if newOrder.Direction == DOWN {
				dir = 0
			}
			if commonExternalArray[newOrder.Floor][dir] == 0 {
				commonExternalArray[newOrder.Floor][dir] = 1
			}
			WriteCurrentOrdersToFile(AllOrders{InternalOrders: internalArray, ExternalOrders: commonExternalArray})

			if !orderIsInQueue(externalOrdersNotTaken, newOrder) {
				externalOrdersNotTaken = append(externalOrdersNotTaken, newOrder)
			}
			driver.SetButtonLight(newOrder, true)

			fmt.Println("TOPLOGIC: to assign func")
			newOrderToBeAssignedChan <- newOrder //the receiver will handle duplicates
			driver.SetButtonLight(newOrder, true)
			/*
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
			*/

		case servedOrder := <-OrderServedLocallyChan:
			//Denne tar imot alle ordre. Hvis nettet er oppe, skal man også sende videre til nett.
			if servedOrder.Direction == 0 {
				internalArray[servedOrder.Floor] = 0
			} else {
				dir := UP
				if servedOrder.Direction == DOWN {
					dir = 0
				}
				commonExternalArray[servedOrder.Floor][dir] = 0
				ToNetworkOrderServedChan <- servedOrder
			}
			WriteCurrentOrdersToFile(AllOrders{InternalOrders: internalArray, ExternalOrders: commonExternalArray})
			driver.SetButtonLight(servedOrder, false)

		case servedOrder := <-FromNetworkOrderServedChan:
			fmt.Println("servedOrder := <-FromNetworkOrderServedChan:")
			dir := UP
			if servedOrder.Direction == DOWN {
				dir = 0
			}
			commonExternalArray[servedOrder.Floor][dir] = 0
			WriteCurrentOrdersToFile(AllOrders{InternalOrders: internalArray, ExternalOrders: commonExternalArray})
			driver.SetButtonLight(servedOrder, false)

		case <-FromNetworkNewElevChan:
			stateMachine.SendElevStateToNetwork(ToNetworkNewElevStateChan)
			ToNetworkExternalArrayChan <- commonExternalArray
		case newExternalArray := <-FromNetworkExternalArrayChan:
			for i := 0; i < NUM_FLOORS; i++ {
				for j := 0; j < NUM_BUTTONS-1; j++ {
					if newExternalArray[i][j] == 1 {
						commonExternalArray[i][j] = 1
						dir := UP
						if j == 0 {
							dir = DOWN
						}
						driver.SetButtonLight(Order{Floor: i, Direction: dir}, true)
					}
				}
			}
			WriteCurrentOrdersToFile(AllOrders{InternalOrders: internalArray, ExternalOrders: commonExternalArray})
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
			for i := 0; i < NUM_FLOORS; i++ {
				for j := 0; j < NUM_BUTTONS-1; j++ {
					if commonExternalArray[i][j] == 1 {
						dir := UP
						if j == 0 {
							dir = DOWN
						}
						NewOrderToLocalElevChan <- Order{Floor: i, Direction: dir}
						time.Sleep(time.Millisecond * 200)
					}
				}
			}
		case elevGone := <-FromNetworkElevGoneChan:
			fmt.Println("ResendOrdersWhenError: elevGone")
			deadElevOrders := optalg.GetOrderQueueOfDeadElev(elevGone)
			fmt.Println("Orders of dead elev:")
			for _, v := range deadElevOrders {
				stateMachine.PrintOrder(v)
			}
			ToOptAlgDeleteElevChan <- elevGone

			resetAssignFuncChan <- true

			externalOrdersNotTaken := optalg.GetUnassignedOrders()
			for _, v := range deadElevOrders {
				ExternalButtonPressedChan <- v
				fmt.Println("                              dead order sent to network ")
				time.Sleep(time.Millisecond * 300)
			}
			for _, v := range externalOrdersNotTaken {
				ExternalButtonPressedChan <- v
				fmt.Println("                             unassigned order sent to network")
				time.Sleep(time.Millisecond * 300)
			}
		}
	}
}

const PATH_OF_SAVED_ORDER_STATE = "elevState.gob"

/*
func initalize_state_tracker(){
	//read from file to check if system was killed
	//easy solution: if thats the case, set current state to that (may serve same order twice, but sverre wont die and its avoiding complicated solutions)
	//if not, initialize normally
	// get floor and all that shit from other modules
	//send the current state to everybody
}*/

/*
func send_updated_elevator_state(){
	//call this whenever its updated. write to channels
}
*/
//endre navn
func WriteCurrentOrdersToFile(currentState AllOrders) {
	//temp for testing
	//update this whenever the local elevator gets an order/command
	dataFile, err := os.Create(PATH_OF_SAVED_ORDER_STATE)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dataEncoder := gob.NewEncoder(dataFile)
	dataEncoder.Encode(currentState)
	dataFile.Close()
}

func ReadOrdersStateBeforeShutdown() AllOrders {
	//start with reading it
	var formerState AllOrders

	if _, err := os.Stat(PATH_OF_SAVED_ORDER_STATE); os.IsNotExist(err) {
		fmt.Println("Local save of elevator state not detected. It has been cleared/this is the first run on current PC")
		return formerState
	}
	dataFile, err := os.Open(PATH_OF_SAVED_ORDER_STATE)

	dataDecoder := gob.NewDecoder(dataFile)
	err = dataDecoder.Decode(&formerState)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dataFile.Close()

	return formerState
}

func PrematureShutdownOccured(formerState AllOrders) bool {
	for i := 0; i < NUM_FLOORS; i++ {
		if formerState.InternalOrders[i] != 0 {
			return true
		}
		for j := 0; j < NUM_BUTTONS-1; j++ {
			if formerState.ExternalOrders[i][j] != 0 {
				return true
			}
		}
	}
	return false
}

//CHANNEL NAMES MIGHT BE WRONG. MIGHT NEED TO SWAP UP AND DOWN VALUES IN INNER FOR LOOP
//Make sure to update networkisup
func ReassignOrdersAfterShutdown(formerState AllOrders, networkIsUp bool) {
	for i := 0; i < NUM_FLOORS; i++ {
		if formerState.InternalOrders[i] == 1 {
			NewOrderToLocalElevChan <- Order{Floor: i, Direction: COMMAND}
			time.Sleep(time.Millisecond * 100)
		}
	}
	for i := 0; i < NUM_FLOORS; i++ {
		for j := 0; j < NUM_BUTTONS-1; j++ {
			if formerState.ExternalOrders[i][j] == 1 {
				direction := UP
				if j == 0 {
					direction = DOWN
				}
				if networkIsUp {
					ExternalButtonPressedChan <- Order{Floor: i, Direction: direction}
				} else {
					NewOrderToLocalElevChan <- Order{Floor: i, Direction: direction}
				}
				time.Sleep(time.Millisecond * 100)
			}
		}
	}

}
func StartupDraft() {
	formerState := ReadOrdersStateBeforeShutdown()
	var emptyState AllOrders = AllOrders{} //sleep to make sure network reads if its up

	if formerState != emptyState {
		if PrematureShutdownOccured(formerState) {

			//Set lights!
			for i := 0; i < NUM_FLOORS; i++ {
				if formerState.InternalOrders[i] == 1 {
					driver.SetButtonLight(Order{Floor: i, Direction: COMMAND}, true)
				}
			}
			for i := 0; i < NUM_FLOORS; i++ {
				for j := 0; j < NUM_BUTTONS-1; j++ {
					if formerState.ExternalOrders[i][j] == 1 {
						direction := UP
						if j == 0 {
							direction = DOWN
						}
						driver.SetButtonLight(Order{Floor: i, Direction: direction}, true)
					}
				}
			}
			//networkIsUp := readNetwork()//something like this           CHECK IF NETWORK IS UP HERE
			time.Sleep(time.Millisecond * 300)
			networkIsUp := false
			currentElevList := network.ElevsSeen()
			if len(currentElevList) > 0 {
				networkIsUp = true
			}
			ReassignOrdersAfterShutdown(formerState, networkIsUp) //FOR NOW TEMPFIX
		}
	}
}
