import
(
	"sync"
	"time"
)


//Masters functions, pick good names.

func order()){
	for {
		order := <-InMasterChans.OrderReceivedManagerChan
		/*
		need to send a confirmation message to the ip that sent this every time we get the order.
		After some time without fun incoming on the channel we can assume that the confirmation has been received.
		*/

	}
}

func Incoming_Order_executed(){
	var countdownChan := make(chan Order)
	var timerMap := make(map[IpOrderMessage] Time)
	var incomingOrderMap = struct{
		sync.RWMutex
    	m map[Order] IP}{m: make(map[Order] IP)}
	go func () {
		for{
			//Updates the queue, if the same kind of messages are sent simultaneously
			orderExe := <- ExCommChans.ToMasterOrderExecutedChan  //This is on IP-form
			incomingOrderMap.RLock()
			inQueue := incomingOrderMap.m[orderExe.Order] //It will be nil if its not in the map
			incomingOrderMap.RUnlock()
			if inQueue == nil{ //If its not in queue we should 
				toExternalOrderChan <- orderExe
				incomingOrderMap.Lock()
				incomingOrderMap.m[orderExe.Order] = orderExe.IP
				incomingOrderMap.Unlock()
				countdownChan <- orderExe.Order
			}
			ExMasterChans.ToCommExecutedConfirmationChan <- orderExe
		}
	}
	//Timer function, that concurrently deletes orders/renews the timer.
	go func(){
		for{
			select{
				case orderReceived := <-countdownChan:
					countdownChan[orderReceived] = time.Now()
			}
			default {
				for order, timestamp := timerMap.range(){
					if time.Since(timestamp) > time.Millisecond* 500{ //temp 
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