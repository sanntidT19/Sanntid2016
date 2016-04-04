import(
	."../globalStructs"
	"encoding/gob"
	"os"
	"fmt"
)

currentState ElevatorState

//this will show 
const PATH_OF_SAVED_ORDER_STATE = "elevState.gob"



type AllOrders type{
	ExternalOrders [NUM_FLOORS][NUM_BUTTONS]int
	InternalOrders [NUM_FLOORS] int
}

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

func WriteCurrentOrdersToFile(currentState AllOrders){
	//temp for testing
	//update this whenever the local elevator gets an order/command
	dataFile, err := os.Create(PATH_OF_SAVED_ORDER_STATE)
	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}

	dataEncoder := gob.NewEncoder(dataFile)
	dataEncoder.Encode(currentState)
	dataFile.Close()
}

func ReadOrdersStateBeforeShutdown() AllOrders{
	//start with reading it
	var formerState AllOrders

	if _, err := os.Stat(PATH_OF_SAVED_STATE); os.IsNotExist(err){
		fmt.Println("Local save of elevator state not detected. It has been cleared/this is the first run on current PC")
		return nil
	}
	dataFile, err := os.Open(PATH_OF_SAVED_STATE)

	dataDecoder := gob.NewDecoder(dataFile)
	err = dataDecoder.Decode(&formerState)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dataFile.Close()

	return formerState
}


func PrematureShutdownOccured(formerState) bool{
	//only call when initializing, REDO WHEN WE KNOW HOW THE ORDERS WILL LOOK
	for i:= 0; i < NUM_FLOORS; i++{
		if formerState.InternalOrders[i] != 0{
			return true
		}
		for j:= 0; j < NUM_BUTTONS-1; j++{
			 if formerState.ExternalOrders[i][j] != 0{
				return true
			}
		}
	}
	return false
}

//CHANNEL NAMES MIGHT BE WRONG. MIGHT NEED TO SWAP UP AND DOWN VALUES IN INNER FOR LOOP
//Make sure to update networkisup
func ReassignOrdersAfterShutdown(formerState AllOrders, networkIsUp bool){
	for i := 0; i < NUM_FLOORS; i++{
		if formerState.InternalOrders[i] == 1{
			NewOrderToLocalElevChan <- Order{Floor: i, Direction: COMMAND}
		}
	}
	for i := 0; i< NUM_FLOORS; i++{
		for j := 0; j < NUM_BUTTONS-1; j++{
			if formerState.ExternalOrders[i][j] == 1{
				direction := UP
				if j == 0{
					direction = DOWN
				}
				if networkIsUp{
					ToNetWorkNewOrderChan <- Order{Floor: i, Direction: direction}
				}else{
					NewOrderToLocalElevChan <- Order{Floor:i, Direction: direction}
				}
			}
		}
	}

}
func StartUpDraft(){
	formerState := ReadOrdersStateBeforeShutdown()
	if formerState != nil{
		if PrematureShutdownOccured(formerState){
			networkIsUp := readNetwork()//something like this
			ReassignOrdersAfterShutdown(formerState,networkIsUp)
			//SET UP LIGHTS HERE
		}		
	}
}
//If network disappears
func SendAllExternalOrdersToLocalElev(currentState AllOrders){
	for i:= 0; i < NUM_FLOORS; i++{
		for j := 0; j < NUM_BUTTONS-1; j++{
			if currentState[i][j] == 1{
				direction := UP
				if j == 0{
					direction = DOWN
				}
				NewOrderToLocalElevChan<-Order{Floor:i, Direction: direction}
			}
		}
	}
}


//Assume network is up. If its not, it will be detected and a different function will be called
func ResendOrdersOfLostElev(orderQueue []Order, sendOrderToNetworkChan chan Order){
	for _, v := range orderQueue{
		if v.Direction != COMMAND{
			sendOrderToNetworkChan <- v
		} 
	}
}

/*
What is needed to save to make sure no orders are lost.
ExternalArray and Internalarray.
There is no need to know direction, current floor, or any such thing.
Not even if other elevators have known stuff.

HAVE AN ARRAY OF UNASSIGNED ORDERS. WHEN ORDER COMES. ADD THIS TO THE QUEUE.
WHEN ELEVATOR DISAPPEARS, ALSO SEND THESE AND RESET STRUCT WAITING FOR AGREEMENT.


Case: unserved orders are up, network detected.
Send all orders over network. If network then suddenly shuts down,
one will resend all external orders anyway.

network not present:
send all orders to local statemachine.


If network suddenly appears: local statemachine will have taken all.
Too bad. T_T
*/

/* 
after init:

elevator appears: just include it in structure which optalg uses.
If this one already has orders, this needs to be checked, only to set lights
do not interfere with its current queue. reset orders to be assigned by all.
resend all unassigned orders

elevator disappears:
both other should notice this at the same time (ish). Go through this 
elevators current queue. Send these over the network. 
Maybe: Have a list of unassigned orders from AssignordersAndwait for agreement
All unassigned orders should be resent.

network disappears:
This is detected. All external orders should be sent to statemachine.
All unassigned orders should be registered before sent to optalg.
So no need to check with optalg-algorithm

save and write to file BEFORE lights are turned on/off.

*/
/*

comm:
need functions/channels that lets the other modules know when a new elevator is present,
an elevator is gone, network i down, and network is up
OK

toplevel: 
assignordersandwaitforagreement needs to reset queue when elevlist changed
or network gone
OK


optalg:
need a structure that contains all elevstructures.
need to update these whenever it is received from network.
OK

remove when elevlist change/network gone.
ok

somewhere:
need to maintain the list of elevstructs whenever its sent over network.
OK

link everything together

*/