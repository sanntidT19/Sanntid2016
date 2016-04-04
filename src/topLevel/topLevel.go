package topLevel

import(
	"../communication"
	"../optalg"
	."../globalStructs"
	."../globalChans"
	"time"
	"fmt"
)


/*
TOP LEVEL FUNCTIONS GO HERE. WE WILL RENAME AND MOVE FUNCTIONS ONCE WE HAVE AN OVERVIEW OF HOW MUCH WE WILL END UP WITH
*/

type AssignedOrderAndElevList struct{
	OrdAss OrderAssigned
	ElevList []string
}

var externalOrdersNotTaken []Order
var commonExternalArray [NUM_FLOORS][NUM_BUTTONS -1] int
//ElevatorAssignedToNetworkChan OrderAssigned chan , ElevatorAssignedFromNetworkChan OrderAssigned chan

func AssignOrdersAndWaitForAgreement(newOrderFromNetworkChan chan Button, resetAssignFuncChan chan bool, newOrderAssignedChan chan OrderAssigned){

	var OrdersToBeAssignedByAll []AssignedOrderAndElevList
	localAddr := communication.GetMyIP()
	elevList [] string
	for {
		select{
		case newOrder := <- newOrderFromNetworkChan:
			orderIsRegistered := false
			for _,v := range OrdersToBeAssignedByAll{
				if v.OrdAss.Order == newOrder{
					orderIsRegistered = true
				}
			}
			if !orderIsRegistered{
				//for now, the elevator states are all globally known. May send copy or something else later.
				assignedElevAddr := optalg.Opt_alg(newOrder)
				NewOrderToBeAssigned := OrderAssigned{Order: newOrder, AssignedTo: assignedElevAddr, SentFrom: localAddr}
				//Elevlist should be copied, global or maybe everyone that uses it should be in the same module
				OrdersToBeAssignedByAll = append(OrdersToBeAssignedByAll, AssignedOrderAndElevList{NewOrderToBeAssigned,elevList})
				time.Sleep(time.Millisecond * 20) //This is to make sure you get to make the list before
				ToNetWorkOrderAssignedToChan <- NewOrderToBeAssigned
			}else{
				fmt.Println("Order already registered. Discard message.")
			}
		case newOrdAss := <- FromNetworkOrderAssignedToChan:
			posInSlice := -1
			for i, v := range OrdersToBeAssignedByAll {
				if newAssOrd.Order == v.OrdAss.Order{
					posInSlice = i
					break;
				}
			}
			if posInSlice == -1{
				fmt.Println("Old/garbage order")
			}else{
				if newOrdAss.AssignedTo != OrdersToBeAssignedByAll.OrdAss.AssignedTo{
					fmt.Println("Disagreement, recalculate with optalg")
					OrdersToBeAssignedByAll = append(OrdersToBeAssignedByAll[:posInSlice], OrdersToBeAssignedByAll[posInSlice+1:]...) //slicetricks
					
					assignedElevAddr := optalg.Opt_alg(newOrdAss.Order)
					NewOrderToBeAssigned := OrderAssigned{Order: newOrdAss.Order, AssignedTo: assignedElevAddr, SentFrom: localAddr}
					OrdersToBeAssignedByAll = append(OrdersToBeAssignedByAll, AssignedOrderAndElevList{NewOrderToBeAssigned,ElevList})
					time.Sleep(time.Millisecond * 20) //This is to make sure you get to make the list before
					ToNetWorkOrderAssignedToChan <- NewOrderToBeAssigned
				}else{
					for i,v := range OrdersToBeAssignedByAll[posInSlice].ElevList{
						if newOrdAss.SentFrom == v{
							OrdersToBeAssignedByAll[posInSlice].ElevList = append(OrdersToBeAssignedByAll[posInSlice].ElevList[:i], OrdersToBeAssignedByAll[posInSlice].ElevList[i+1:]...)
							if len(OrdersToBeAssignedByAll[posInSlice].ElevList) == 0{
								OrdersToBeAssignedByAll = append(OrdersToBeAssignedByAll[:posInSlice], OrdersToBeAssignedByAll[posInSlice+1:]...)
								newOrderAssignedChan <-newOrdAss
								}	
							}
						}
					}
				} 
			}	
		case <-resetAssignFuncChan:
			elevList = communication.GetElevList()
			OrdersToBeAssignedByAll = nil
		}
	}
}
//This may not be used, yo
func SendAllOrdersOfQueueToNetwork(commonExternalArray [][] int){
	for i := 0; i < NUM_FLOORS; i++{
		for j := 0; j < NUM_BUTTONS -1; j++{
			if commonExternalArray[i][j] == 1 {
				direction = UP
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

func orderIsInQueue(order_queue []Order, newOrder Order) bool {
	for _, queueElements := range order_queue {
		if queueElements == new_order {
			return true
		}
	}
	return false
}

func TopLogicNeedBetterName(){
	commonExternalArray [NUM_FLOORS][NUM_BUTTONS -1] int // UP DOWN
	internalArray[NUM_FLOORS] int

	resetAssignFuncChan := make(chan bool)
	go ResendOrdersWhenError(resetAssignFuncChan)
	go AssignOrdersAndWaitForAgreement(newOrderToBeAssignedChan, resetAssignFuncChan, newOrderAssignedChan)
	for{
		select{
			case newButton <-internalButtonPressedChan:
				
				if internalArray[newButton.Floor] == 0{
					internalArray[newButton.Floor] = 1  //Reset when order served. Tell that total state is changed?
					elevatorStateTracker.WriteCurrentOrdersToFile(Allorders{InternalOrders: internalArray, ExternalOrders: commonExternalArray})
					SetButtonLight(newButton)
					NewOrderToLocalElevChan <- newButton
				}
			case newOrder := <-ExternalOrderFromNetWorkChan:
				dir := UP
				if newOrder.Direction == DOWN {
					dir = 0
				}
				if commonExternalArray[newOrder.Floor][dir] == 0{
					commonExternalArray[newOrder.Floor][dir] = 1
				}

				elevatorStateTracker.WriteCurrentOrdersToFile(Allorders{InternalOrders: internalArray, ExternalOrders: commonExternalArray})

				if !orderIsInQueue(ordersNotConfirmedTaken, newOrder){
					externalOrdersNotTaken = append(externalOrdersNotTaken,newOrder)
				}
				newOrderToBeAssignedChan <- newOrder //the receiver will handle duplicates
				SetButtonLight(newOrder,true)

			case newOrdAss := <-newOrderAssignedChan:
				AddOrderAssignedToElevStateChan <- newOrdAss
				for i,v := range externalOrdersNotTaken{
					if v == newOrdAss.Order{
						externalOrdersNotTaken = append(externalOrdersNotTaken[:i], externalOrdersNotTaken[i+1:]...)
					}
				}
				if newOrdass.AssignedTo == communication.GetMyIP(){
					NewOrderToLocalElevChan <- newOrdAss.Order
				}


			case servedOrder<-InternalOrderServedChan: //Here or directly from statemachine
				internalArray[servedOrder.Floor] = 0

				elevatorStateTracker.WriteCurrentOrdersToFile(Allorders{InternalOrders: internalArray, ExternalOrders: commonExternalArray})

				SetButtonLight(servedOrder,false)
			case servedOrder <-externalOrderServedfromnetwork:
				dir := UP
				if servedOrder.Direction == DOWN{
					dir = 0
				}
				commonExternalArray[servedOrder.Floor][dir] = 0
				elevatorStateTracker.WriteCurrentOrdersToFile(Allorders{InternalOrders: internalArray, ExternalOrders: commonExternalArray})
				SetButtonLight(servedOrder,false)

 			case newElev := <-FromNetworkNewElevChan:
				fmt.Println("DONT KNOW IF WE NEED THIS CASE")

		}
	}
}

//MULIG VI MÅ FLYTTE DE SOM TAR IMOT NYE ORDRE OG DE SOM SENDER NYE ORDRE I FORSKJELLIGE LOOPS
//PROBLEM NÅR MANGE ORDRE BLIR SENDT PÅ NYTT

func ResendOrdersWhenError(resetAssignFuncChan chan bool){
	for{
		select{
			case <-FromNetworkNetworkDownChan:
				resetAssignFuncChan<-true
				externalOrdersNotTaken = nil
				for i:= 0; i < NUM_FLOORS; i++{
					for j := 0; j < NUM_BUTTONS -1; j++{
						if commonExternalArray[i][j] == 1{
							dir = UP
							if j == 0{
								dir = DOWN
							}
							NewOrderToLocalElevChan <- Order{Floor: i, Direction: dir}
						}
					}
				}
			case elevGone := <-FromNetworkElevGoneChan:
				deadElevOrders := optalgtester.GetOrderQueueOfDeadElev(elevGone)
				ToOptAlgDeleteElevChan <-elevGone
				resetAssignFuncChan<-true
				for _, v := range deadElevOrders}{
					ToNetworkNewOrderChan <- v
				}
				for _, v := range externalOrdersNotTaken{
					ToNetworkNewOrderChan <- v 
				}
		}
	}
}