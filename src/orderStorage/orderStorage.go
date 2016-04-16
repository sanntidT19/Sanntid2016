package orderStorage

import (
	"../driver"
	. "../globalStructs"
	"../network"
	"../orderAssignment"
	"../localElev"
	"encoding/gob"
	"fmt"
	"os"
	"time"
)

/*
Functions assuring no orders are lost when errors occur locally or somewhere else in the network
*/

const PATH_OF_SAVED_ORDER_STATE = "mostRecentOrderState.gob"
var externalOrdersNotTaken []Order
var commonExternalArray [NUM_FLOORS][NUM_BUTTONS - 1]int

//Helper function
func orderIsInList(orderQueue []Order, newOrder Order) bool {
	for _, v := range orderQueue {
		if v == newOrder {
			return true
		}
	}
	return false
}

//Save the current order state to file every time it changes
func RegisterNewAndServedOrders(newOrderToBeAssignedChan chan Order, sendNewOrderChan chan Order, receiveNewOrderChan chan Order, sendOrderServedChan chan Order, receiveOrderServedChan chan Order, sendElevStateChan chan ElevatorState, sendExternalArrayChan chan [NUM_FLOORS][NUM_BUTTONS-1]int , receiveExternalArrayChan chan [NUM_FLOORS][NUM_BUTTONS-1]int, internalButtonChan chan Order, newLocalOrderChan chan Order, localOrderServedChan chan Order, newElevChan chan string) {
	var internalArray [NUM_FLOORS]int
	fmt.Println(commonExternalArray)
	for {
		select {
		case newOrder := <-internalButtonChan:
			if internalArray[newOrder.Floor] == 0 {
				internalArray[newOrder.Floor] = 1
				WriteCurrentOrdersToFile(AllOrders{InternalOrders: internalArray, ExternalOrders: commonExternalArray})
				driver.SetButtonLight(newOrder, true)
				newLocalOrderChan <- newOrder
			}
		case newOrder := <-receiveNewOrderChan:
			dir := UP
			if newOrder.Direction == DOWN {
				dir = 0
			}
			if commonExternalArray[newOrder.Floor][dir] == 0 {
				commonExternalArray[newOrder.Floor][dir] = 1
			}
			WriteCurrentOrdersToFile(AllOrders{InternalOrders: internalArray, ExternalOrders: commonExternalArray})
			driver.SetButtonLight(newOrder, true)
			newOrderToBeAssignedChan <- newOrder
			driver.SetButtonLight(newOrder, true)

		case servedOrder := <-localOrderServedChan:
			if servedOrder.Direction == 0 {
				internalArray[servedOrder.Floor] = 0
			} else {
				dir := UP
				if servedOrder.Direction == DOWN {
					dir = 0
				}
				commonExternalArray[servedOrder.Floor][dir] = 0
				sendOrderServedChan <- servedOrder
			}
			WriteCurrentOrdersToFile(AllOrders{InternalOrders: internalArray, ExternalOrders: commonExternalArray})
			driver.SetButtonLight(servedOrder, false)

		case servedOrder := <-receiveOrderServedChan:
			dir := UP
			if servedOrder.Direction == DOWN {
				dir = 0
			}
			commonExternalArray[servedOrder.Floor][dir] = 0
			WriteCurrentOrdersToFile(AllOrders{InternalOrders: internalArray, ExternalOrders: commonExternalArray})
			driver.SetButtonLight(servedOrder, false)

		case <-newElevChan:
			localElev.SendElevStateToNetwork(sendElevStateChan)
			sendExternalArrayChan <- commonExternalArray

		case newExternalArray := <-receiveExternalArrayChan:
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

func ResendOrdersIfNetworkError(resetAssignFuncChan chan bool, sendNewOrderChan chan Order,newLocalOrderChan chan Order, deadElevChan chan string, networkDownChan chan bool, removeElevChan chan string) {
	for {
		select {
		case <-networkDownChan:
			resetAssignFuncChan <- true
			for i := 0; i < NUM_FLOORS; i++ {
				for j := 0; j < NUM_BUTTONS-1; j++ {
					if commonExternalArray[i][j] == 1 {
						dir := UP
						if j == 0 {
							dir = DOWN
						}
						newLocalOrderChan <- Order{Floor: i, Direction: dir}
						time.Sleep(time.Millisecond * 200)
					}
				}
			}
		case elevGone := <-deadElevChan:
			deadElevOrders := orderAssignment.GetOrderQueueOfDeadElev(elevGone)
			removeElevChan <- elevGone
			resetAssignFuncChan <- true
			externalOrdersNotTaken := orderAssignment.GetUnassignedOrders()
			for _, v := range deadElevOrders {
				sendNewOrderChan <- v
				time.Sleep(time.Millisecond * 300)
			}
			for _, v := range externalOrdersNotTaken {
				sendNewOrderChan <- v
				time.Sleep(time.Millisecond * 300)
			}
		}
	}
}


func WriteCurrentOrdersToFile(currentState AllOrders) {
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

func ReassignOrdersIfPrematureShutdown(formerState AllOrders, networkIsUp bool, sendNewOrderChan chan Order, newLocalOrderChan chan Order) {
	for i := 0; i < NUM_FLOORS; i++ {
		if formerState.InternalOrders[i] == 1 {
			newLocalOrderChan <- Order{Floor: i, Direction: COMMAND}
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
					sendNewOrderChan <- Order{Floor: i, Direction: direction}
				} else {
					newLocalOrderChan <- Order{Floor: i, Direction: direction}
				}
				time.Sleep(time.Millisecond * 100)
			}
		}
	}

}
func RestoreOrderStateBeforeShutdown(sendNewOrderChan chan Order,newLocalOrderChan chan Order) {
	formerState := ReadOrdersStateBeforeShutdown()
	var emptyState AllOrders = AllOrders{}
	if formerState != emptyState {
		if PrematureShutdownOccured(formerState) {

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
			time.Sleep(time.Millisecond * 300)
			networkIsUp := false
			currentElevList := network.ElevsSeen()
			if len(currentElevList) > 0 {
				networkIsUp = true
			}
			ReassignOrdersIfPrematureShutdown(formerState, networkIsUp, sendNewOrderChan, newLocalOrderChan) //FOR NOW TEMPFIX
		}
	}
}
