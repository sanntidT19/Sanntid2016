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

//ElevatorAssignedToNetworkChan OrderAssigned chan , ElevatorAssignedFromNetworkChan OrderAssigned chan

func AssignOrdersAndWaitForAgreement(newOrderFromNetworkChan chan Button, ElevListChangedResetAssignFuncChan chan bool, NewOrderAssignedChan chan OrderAssigned ){

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
								NewOrderAssignedChan <-newOrdAss
								}	
							}
						}
					}
				} 
			}	
		case <-ElevListChangedResetAssignFuncChan:
			elevList = communication.GetElevList()
			OrdersToBeAssignedByAll = nil
		}
	}
}
//This may not be used, yo
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
					internalArray[newButton.Floor] = 1  //Reset when order served. Tell that total state is changed?
					SetButtonLight(newButton)
					newOrderToStateMachineChan <- newButton
				}
			case newOrder := <-ExternalOrderFromNetWorkChan:
				dir := UP
				if newOrder.Direction == DOWN {
					dir = 0
				}
				if commonExternalArray[newOrder.Floor][dir] == 0{
					commonExternalArray[newOrder.Floor][dir] = 1

					WRITE TO FILE

				}
				newOrderToOptAlgChan <- newOrder //the receiver will handle duplicates
				SetButtonLight(newOrder,true)

			case <-ElevListChanged://NEEDS TO CHANGE
				ResetExternalOrdersInQueueChan<-true
				ReassignAllOrders(newOrder,true)
			case servedOrder<-InternalOrderServedChan: //Here or directly from statemachine
				internalArray[servedOrder.Floor] = 0

					WRITE TO FILE

				SetButtonLight(servedOrder,false)
			case servedOrder <-externalOrderServedfromnetwork:
				dir := UP
				if servedOrder.Direction == DOWN{
					dir = 0
				}
				commonExternalArray[servedOrder.Floor][dir]
				SetButtonLight(servedOrder,false)
		}
	}
}
