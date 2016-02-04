package main

import(
	"./driver"
	"fmt"
	"time"
	"./stateMachine"
	//"./globalChans"
	
)

func main(){
	driver.Elev_main_tester_function()
	next_order_chan := make(chan int)
	order_served_chan := make(chan bool)
	go stateMachine.Get_current_floor()
	go stateMachine.Execute_order(next_order_chan)
	go stateMachine.Stop_at_desired_floor(order_served_chan)
	next_order_chan <- 3
	<-order_served_chan
	fmt.Printf("Order_served\n")
	time.Sleep(4*time.Second)
	fmt.Printf("End of main \n")
}
