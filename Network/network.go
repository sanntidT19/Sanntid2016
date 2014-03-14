package main

//https://groups.google.com/forum/#!topic/golang-china/KG0Bgf0CAQc
import (
	"bytes"
	"fmt"
	"math/rand"
	. "net"
	"time"
)

const (
	MAXWAIT = time.Second
	PORT    = ":20019"
)

/*
func main() {
	nr := backup()

	primary(nr)

}
*/
func Network_init() Conn {
	fmt.Println("gi")
	address, err := ResolveUDPAddr("udp", "129.241.187.255"+PORT) //leser bare fra porten generellt
	c, err := DialUDP("udp", nil, address)

	if err != nil {
		fmt.Println(err)
	}
	return c

}
func Send(c Conn, to_writing []byte) bool { //Olav: does this need a buffer as paramter aswell??
	_, err := c.Write(to_writing)

	if err != nil {
		fmt.Println(err)
		return false
	}
	time.Sleep(300 * time.Millisecond)
	return true
}

func Recive() []byte { //does only need a connection who listen to a port like ":20019" not the entire ip adress.
	buf := make([]byte, 1024)
	addr, _ := ResolveUDPAddr("udp", PORT)
	c, err := ListenUDP("udp", addr)

	//addr, err := ResolveUDPAddr("udp", PORT)
	//c, err := ListenUDP("udp", addr)
	defer c.Close()

	c.SetReadDeadline(time.Now().Add(1 * time.Second)) //returns error if deadline is reached
	_, _, err = c.ReadFromUDP(buf)                     //n contanis numbers of used bytes, fills buf with content on the connection
	if err == nil {                                    //if error is nil, read from buffer
		return buf
	} else {
		//break
	}
	return nil
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
	c.SetReadDeadline(time.Now().Add(Random_int(600, 1000) * time.MilliSecond))
	_, _, _ = c.ReadFromUDP(buf) //n contanis numbers of used bytes
		
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

func Master_elevator() {
	address, _ := ResolveUDPAddr("udp", "129.241.187.255"+PORT)
	c, _ := DialUDP("udp", nil, address)
	
	fmt.Println("primary")
	for {

		_, err := c.Write([]byte("iam")
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
