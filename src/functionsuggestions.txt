
func Order_received_manager(){
	for {
		order := <-InMasterChans.OrderReceivedManagerChan
		/*
		need to send a confirmation message to the ip that sent this every time we get the order.
		After some time without fun incoming on the channel we can assume that the confirmation has been received.
		*/

	}
}