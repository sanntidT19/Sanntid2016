package network

import (
	. "chansnstructs"
	"fmt"
	"math/rand"
	. "net"
	"time"
)

func Network_init() (Conn, Conn) {
	fmt.Println("gi")
	addr, err := ResolveUDPAddr("udp", "129.241.187.255"+PORT) //leser bare fra porten generellt
	SendConn, err := DialUDP("udp", nil, addr)

	if err != nil {
		fmt.Println(err.Error())
	}

	return SendConn

}

func Write_to_network(to_writing []byte, c Conn) {
	for {
		//err := c.SetWriteDeadline(time.Now().Add(50 * time.Millisecond))
		_, err = c.Write(to_writing)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			//break
		}
	}
}

func Receive() { //will error trigger if just read fails? or will it only go on deadline?
	buf := make([]byte, 1024)
	addr, _ := ResolveUDPAddr("udp", PORT)
	c, err := ListenUDP("udp", addr)

	//this will also check if the master is still there.
	c.SetReadDeadline(time.Now().Add(300 * time.Millisecond)) //returns error if deadline is reached
	n, sendingAddr, err := c.ReadFromUDP(buf)                 //n contanis numbers of used bytes, fills buf with content on the connection

	fmt.Println(sendingAddr)

	if err == nil { //if error is nil, read from buffer
		ExNetChans.ToComm <- buf[0:n]
		ExNetChans.ToCommAddr <- addr

		//ExSlaveChans.ToSlaveImMasterChan <- true
	} else {
		//ExSlaveChans.ToSlaveImMasterChan <- false
	}

}

func Random_init(min int, max int) int { //gives a random int for waiting
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}
