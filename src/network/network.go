package network

import (
	"fmt"
	"math/rand"
	. "net"
	"time"
)

const (
	MAXWAIT = time.Second
	PORT    = ":20019"
)
var ExNetChan NetworkExternalChannels
/*
func main() {
	nr := backup()

	primary(nr)

}
*/

type NetworkExternalChannels struct {
	ToNetwork chan []byte
	ToComm chan []byte
}
func network_external_chan_init() {
	ExNetChan.ToNetwork = make(chan []byte)
	ExNetChan.ToComm = make(chan []byte)
}


func Network_init() Conn {
	fmt.Println("gi")
	address, err := ResolveUDPAddr("udp", "129.241.187.255"+PORT) //leser bare fra porten generellt
	c, err := DialUDP("udp", nil, address)

	if err != nil {
		fmt.Println(err.Error())
	}
	return c

}

//use chan in outer loop, who fetches to_writing
//if chan is taken as parameter we are unable to run it in a for loop, because the channel is emptied at the first run
func Send(c Conn, to_writing []byte) { //Olav: does this need a buffer as paramter aswell??
	for {
		_, err := c.Write(to_writing)

		if err != nil {
			fmt.Println(err.Error())
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func Receive(networkToComm chan []byte) { //does only need a connection who listen to a port like ":20019" not the entire ip adress.
	buf := make([]byte, 1024)
	addr, _ := ResolveUDPAddr("udp", PORT)
	c, err := ListenUDP("udp", addr)

	defer c.Close()

	c.SetReadDeadline(time.Now().Add(1 * time.Second)) //returns error if deadline is reached
	_, _, err = c.ReadFromUDP(buf)                     //n contanis numbers of used bytes, fills buf with content on the connection
	if err == nil {                                    //if error is nil, read from buffer
		ExNetChan.ToComm <- buf
	} else {
		//break
	}

}

func Choose_master() {
	//all say i am slave with random timeout time
	//all will obey if one is master

	//when time runs out it calls out i am master

	/*PSUDO
	ALL: broadcast "i am slave"
	ALL: set readDeadline(random time) - listening for "i am master"
	MASTER: first one who times out broadcast "i am master".
	ALL - MASTER: will continue listening for "i am master"
	***We have a master***
	*/
	go Slave_elevator()
}

func Slave_elevator() {

	buf := make([]byte, 1024)
	addr, _ := ResolveUDPAddr("udp", PORT)
	c, _ := ListenUDP("udp", addr)

	defer c.Close()

	for {
		//rand := Random_int(600, 1000)
		c.SetReadDeadline(time.Now().Add(10 * time.Millisecond))
		_, _, err := c.ReadFromUDP(buf) //n contanis numbers of used bytes

		if err == nil { // of readdeadline dont kicks in
			//decrypt buf
			//if decryptet buf equals iam
			//keep on serching

		} else { // if readdeadline kicks in
			//first one here becomes master(?)

			//this will just be called in case of there is no master
			go Master_elevator()
			break
		}

	}
}

func Master_elevator() {
	address, _ := ResolveUDPAddr("udp", "129.241.187.255"+PORT)
	c, _ := DialUDP("udp", nil, address)

	fmt.Println("primary")
	for {
		_, err := c.Write([]byte("iam"))
		if err != nil {
			fmt.Println("fail")
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func Random_int(min int, max int) int { //gives a random int for waiting
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}
