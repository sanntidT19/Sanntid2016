package network

import (
"net";
"bufio";
"fmt";
."time"
)

func initNetwork() {

	raddr, err := ResolveUDPAddr("udp", "129.241.187.255:20019")

	c ,err :=  net.DialUDP("udp","129.241.187.157:20019")
	if err != nil {
     	    fmt.Printf("unable to connect to server, code red: %s\n", err.Error())
     	}
	
	go Echo(c)
	go Broadcast(c)
	//close the connections when program is closed(?)
	defer c.Close()
	defer raddr.Close()
	/*
	for {
		
		Echo(raddr)
		Broadcast(c)
		Sleep(100*Millisecond)
	}
	*/
	
}

func Read(raddr net.UDPConn, networkChan chan []byte) {
	message, err := bufio.NewReader(raddr).ReadString('\x00')

	if err != nil {
		fmt.Printf("Failed to read: %s\n",err.Error())	
		return
	}
	networkChan <- message
	//fmt.Println(message)

}

func Broadcast(c net.UDPConn, networkChan chan []byte) {
	holdMessage <- networkChan
	//sould be sendt sevral times to be sure the reciver gets it??????
	_, err = c.WriteToUDP(holdMessage)

	if err != nil {
		fmt.Printf("Failed to write: %s\n", err.Error())
		return	
	}

}
	// defer c.Close() //Closes the TCP connection so you dont abuse it, defer is called when the function returns.
