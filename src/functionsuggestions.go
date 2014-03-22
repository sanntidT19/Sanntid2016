package test


import
(
	"sync"
	"time"
)

WE MAY NEED TO MAKE A COPY, I DO NOT KNOW IF WE CAN ITERATE THROUGH THE MAPS WHILE WE ARE WRITING TO IT. HOW DO WE RLOCK AN IF-STATEMENT IN THIS CASE?
//Masters functions, pick good names.


func Master_updated_state_incoming(){
	for {
		updatedState := <- ToMasterUpdateStateChan 					THIS NEEDS TO BE MADE/IS IT MADE??? //This is Ip and everything
		MasterUpdateCentralChan <- updatedState				THIS NEEDS TO BE MADE/IS IT MADE???
		ToCommUpdateStateReceivedChan <-updatedState

	}

}

func Master_updated_arrays_outgoing(){
	localSlaveMap := make(map[IP] Slave)  //think about this one
	timerMap := make(map[IP] time.Time)  //timers for each IP
	startCountdownChan := make(chan bool)
	countdownFinishedChan := make(chan bool)



	//Get the total shit if it has been updated, and the slaves need to know
	go func(){
		for{
			select{
				case localSlaveMap = <- arraysHasBeenUpdatedChan:   THIS NEEDS TO BE MADE/IS IT MADE???   ALSO: BAD NAME
					ToCommAllArrays <- localSlaveMap


					startCountdownChan <- true
				case <-countdownFinishedChan:
					ToCommAllArrays <- localSlaveMap
			}	

		}
	}()
	//Timer goroutine, receives answer and removes from 
	//Also sends again if no answer
	go func(){
		for{
			<-startCountdownChan
			timer := time.After(500 * time.Millisecond)  // If all answers, find a way to stop timer
			<-timer
			countdownFinishedChan <- true




		}
	}()
}

//Either copy-paste this or send it to optimization-module in the code where it is handled. May just have a goroutine in this function as well  
func Master_incoming_order_executed(){
	countdownChan := make(chan Order)
	timerMap := make(map[IpOrderMessage] Time)
	var incomingOrderMap = struct{
		sync.RWMutex
    	m map[Order] IP}{m: make(map[Order] IP)}
	go func () {
		for{
			//Updates the queue, if the same kind of messages are sent simultaneously
			orderExe := <- ExCommChans.ToMasterOrderExecutedChan  //This is on IP-message-form



			ToOrderHandlerChan <- orderExe     /*The code that receives isnt made yet. Should handle optimization module there.*/ THIS NEEDS TO BE MADE/IS IT MADE??? 



			incomingOrderMap.RLock()
			inQueue := incomingOrderMap.m[orderExe.Order] //It will be nil if its not in the map
			incomingOrderMap.RUnlock()
			if inQueue == nil{ //If its not in queue we should
				toExternalOrderChan <- orderExe
				incomingOrderMap.Lock()
				incomingOrderMap.m[orderExe.Order] = orderExe.IP
				incomingOrderMap.Unlock()
			}
			countdownChan <- orderExe.Order
			ExMasterChans.ToCommExecutedConfirmationChan <- orderExe
		}}()
	//Timer function, that concurrently deletes orders/renews the timer.
	go func(){
		for{
			select{
				case orderReceived := <-countdownChan:
					timerMap[orderReceived] = time.Now()
			}
			default {
				for order, timestamp := timerMap.range(){
					if time.Since(timestamp) > time.Millisecond* 500{ //RLock before and after this if?
						delete(timerMap, order)
						incomingOrderMap.Lock()
						delete(incomingOrderMap.m,order)
						incomingOrderMap.Unlock()
					}
				}
				time.Sleep(25*Millisecond) // Change to optimize
			}
		}
		}()
	}
}





	/*
	Knows where this came from.
	Update the array that the optimization algorithm uses, but dont send it to the algorithm.
	I guess you need to send this to all the slaves and stop when they have confirmed that they have received it.
	They need to know this in case the master goes down. This will also trigger the lights of all elevators.
	Just keep sending this to all the slaves and mark the slaves as they confirm that they have received it.
	When all slaves are marked this will stop sending.

	If multiple different orders are incoming, how do we spawn to threads? Does goroutines work this way?
*/



//Slaves functions

func Slave_Order_Outgoing(){
	countdownChan := make(chan Order)
	outgoingOrderMap := struct{sync.RWMutex ,m map[Order] IP}
    	{m: make(map[Order] time.Time)}

	go func () {
		for{
			//Updates the queue, if the same kind of messages are sent simultaneously
			orderOut := <- ToSlaveOrderOutChan           THIS NEEDS TO BE MADE/IS IT MADE???
			outgoingOrderMap.Lock()
			outgoingOrderMap.m[orderOut] = time.Now()
			outgoingOrderMap.Unlock()
			ToCommOrderOutChan <- orderOut              THIS NEEDS TO BE MADE/IS IT MADE???				
		}
	}()
	//Timer function, that concurrently renews the timer and sends another message if the old one has timed out.
	go func(){
		for{
			for orderOut, timestamp := outgoingOrderMap.m.range(){
					if time.Since(timestamp) > time.Millisecond* 500{ //temp
						outgoingOrderMap.Lock()
						outgoingOrderMap[orderOut] = time.Now()
						outgoingOrderMap.Unlock()
						ToCommOrderOutChan <- orderOut
					}
				}
				time.Sleep(25*Millisecond) // Change to optimize
			}
		}
		}()
	}
	//Waits for a response and removes the element from map if it has been confirmed by master.
	go func(){
		for{
			orderOut := <- ToSlaveOrderOutConfirmedChan  	THIS NEEDS TO BE MADE/IS IT MADE???
			outgoingOrderMap.Lock()
			delete(outgoingOrderMap,orderOut)
			outgoingOrderMap.Unlock()		
		}
	}
}

//No queueing necessary
func Slave_order_arrays_incoming(){
	var orderArray Master  //Think sending the master could be good, but master isnt a good name
	go func() {
		for{
			orderArray := <- ToSlaveArraysChan
			//Here we need to save all the information about the other slaves, and send our own to the statemachine
			FromSlaveArraysReceivedChan <- Message back
		}
	}
}
//Sends state if timer expires or state changes.
func Slave_state_updated(){
	var sendAgainTimer time.Time    THESE TWO NEEDS TO BE VERIFIED
	var currentState State
	for{
		select{
		case currentState = <- ToSlaveFromStateChan:
			FromSlaveToCommStateChan <- currentState                 		THIS NEEDS TO BE MADE/IS IT MADE???
			sendAgainTimer = time.After(500* time.Millisecond)
		case currentStateReceived := <-ToSlaveFromCommChan:
			if currentStateReceived == currentState {
				sendAgainTimer = nil                  //Not sure if this is legal, will this send to channel if its set to nil??
			}
		case <- sendAgainTimer:  //This will be sent when time runs out, I think.
			FromSlaveToCommStateChan <- currentState    		THIS NEEDS TO BE MADE/IS IT MADE???
			sendAgainTimer = time.After(500* time.Millisecond)
		}
	}
}