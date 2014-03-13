package main

//https://groups.google.com/forum/#!topic/golang-china/KG0Bgf0CAQc
import (
	"fmt"
	. "net"
	"os/exec"
	"strconv"
	"time"
)

const (
	MAXWAIT = time.Second
	PORT = ":20019"
)
/*
func main() {
	nr := backup()

	primary(nr)

}
*/
func Network_init() {
	
	address, err := ResolveUDPAddr("udp", PORT) //leser bare fra porten generellt
	conn, err := ListenUDP("udp", address)

	if err != nil {
		fmt.Println(err)
	}

}
func Write(c Conn, to_writing ) bool { //Olav: does this need a buffer as paramter aswell??
	buf := make([]byte, 1024)
	//Keep running
	buf = codeMaster(externalList)//json: packs down externallist
	_, err := conn.Write(buf)

	if err != nil {
		fmt.Println(err)
	}
	time.Sleep(300 * time.Millisecond)
}

func Read(c Conn) [][]int { //does only need a connection who listen to a port like ":20019" not the entire ip adress.
	buf := make([]byte, 1024)
	//addr, err := ResolveUDPAddr("udp", PORT)
	//c, err := ListenUDP("udp", addr)

	c.SetReadDeadline(time.Now().Add(1 * time.Second)) //returns error if deadline is reached
	n, _, err := c.ReadFromUDP(buf) //n contanis numbers of used bytes, fills buf with content on the connection
	if err == nil { //if error is nil, read from buffer
		newExternalList := decodeSlave(buf)
		//nr, _ = strconv.Atoi(string(buf[0:n]))
	} else {
		break
	}
}
c.Close()
return newExternalList
}

