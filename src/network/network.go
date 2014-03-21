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
	c1, err := DialUDP("udp", nil, addr)

	if err != nil {
		fmt.Println(err.Error())
	}

	defer c1.Close()
	addr2, _ := ResolveUDPAddr("udp", PORT)
	c2, err := ListenUDP("udp", addr2)

	//c2.SetReadDeadline(time.Now().Add(300 * time.Millisecond)) //returns error if deadline is reached
	//_, _, err = c2.ReadFromUDP(buf)

	return c1, c2

}

///////We need to select push to network or send
func Push_to_network(to_writing []byte, c Conn) {
	for {
		err := c.SetWriteDeadline(time.Now().Add(10 * time.Millisecond))
		_, err = c.Write(to_writing)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			//break
		}
	}
}

func Send() {
	byteArr := <-ExNetChans.ToNetwork
	c := <-ExNetChans.ConnChan

	go Write_to_network(byteArr, c)
}
func Write_to_network(to_writing []byte, c Conn) {
	for {
		err := c.SetWriteDeadline(time.Now().Add(10 * time.Millisecond))
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

	defer c.Close()
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

/*
func Choose_master() {
	go Slave_elevator()
}

func Slave_elevator() {

	buf := make([]byte, 1024)
	addr, _ := ResolveUDPAddr("udp", PORT)
	c, _ := ListenUDP("udp", addr)

	defer c.Close()

	for {
		//rand := Random_int(600, 1000)
		c.SetReadDeadline(time.Now().Add(300 * Millisecond))
		_, _, err := c.ReadFromUDP(buf) //n contanis numbers of used bytes

		if err == nil { // of readdeadline dont kicks in
			//decrypt buf
			//if decryptet buf equals iam
			//keep on serching
			Decrypt_message(buf)
			<-ExComChans.ToSlaveImMasterChan

		} else { // if readdeadline kicks in
			//first one here becomes master(?)

			//this will just be called in case of there is no master
			go Master_elevator()
			break
		}

	}
}

func Master_elevator() {

	for {
		MC.ToCommImMasterChan <- true
		time.Sleep(50 * time.Millisecond)
	}
}
*/
func Random_init(min int, max int) int { //gives a random int for waiting
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}
