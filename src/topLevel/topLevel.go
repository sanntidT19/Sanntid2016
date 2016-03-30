package topLevel

import(
	"../communication"
	"../optalg"
	"../globalStructs"
	"../globalChans"
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


func AssignOrdersAndWaitForAgreement(newOrderFromNetworkChan chan Button, ElevatorAssignedToNetworkChan chan , ElevatorAssignedFromNetworkChan chan){
	var OrdersToBeAssignedByAll []AssignedOrderAndElevList
	localAddr := communication.GetMyIP()
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
				OrdersToBeAssignedByAll = append(OrdersToBeAssignedByAll, AssignedOrderAndElevList{NewOrderToBeAssigned,ElevList})
				time.Sleep(time.Millisecond * 20) //This is to make sure you get to make the list before
				ElevatorAssignedToNetworkChan <- NewOrderToBeAssigned
			}else{
				fmt.Println("Order already registered. Discard message.")
			}
		case newOrdAss := <- ElevatorAssignedFromNetworkChan:
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
					ElevatorAssignedToNetworkChan <- NewOrderToBeAssigned
				}else{
					for i,v := range OrdersToBeAssignedByAll[posInSlice].ElevList{
						if newOrdAss.SentFrom == v{
							OrdersToBeAssignedByAll[posInSlice].ElevList = append(OrdersToBeAssignedByAll[posInSlice].ElevList[:i], OrdersToBeAssignedByAll[posInSlice].ElevList[i+1:]...)
							if len(OrdersToBeAssignedByAll[posInSlice].ElevList) == 0{
								OrdersToBeAssignedByAll = append(OrdersToBeAssignedByAll[:posInSlice], OrdersToBeAssignedByAll[posInSlice+1:]...)
								if newOrdAss.AssignedTo == localAddr{
									IWillTakeOrderChan <-newOrdAss
									NewOrderToLocalElevChan <- newOrdAss.Order
								}	
							}
						}
					}
				}
			}	
		case <-elevlistchange:
			OrdersToBeAssignedByAll = nil
			//dump?
		}
	}
}

func ReassignAllOrders(commonExternalArray [][] int){
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


func TopLogicNeedBetterName(){
	commonExternalArray [NUM_FLOORS][NUM_BUTTONS -1] int // UP DOWN
	internalArray[NUM_FLOORS] int
	for{
		select{
			case newButton <-internalButtonPressedChan:
				
				if internalArray[newButton.Floor] == 0{
					newOrderToStateMachineChan <- newButton
					internalArray[newButton.Floor] = 1  //Reset when order served. Tell that total state is changed?
					SetButtonLight(newButton)
				}
			case newOrder := <-ExternalOrderFromNetWorkChan:
				if commonExternalArray[newOrder.Floor][newOrder.Direction] == 0{
					newOrderToOptAlgChan <- newOrder
					commonExternalArray[newOrder.Floor][newOrder.Direction] = 1
					SetButtonLight(newOrder,true)
				}
			case <-ElevListChange:
				ResetExternalOrdersInQueueChan<-true
				ReassignAllOrders(newOrder,true)
			case servedOrder<-InternalOrderServedChan: //Here or directly from statemachine
				internalArray[servedOrder.Floor] = 0
				SetButtonLight(servedOrder,false)
			case servedOrder <-externalOrderServedfromnetwork:
				commonExternalArray[servedOrder.Floor][servedOrder.Direction]
				SetButtonLight(servedOrder,false)
		}
	}
}
