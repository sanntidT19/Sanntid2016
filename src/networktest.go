package main

import (
	"fmt"
	//. "net"
	//"os"
	//"os/signal"
	//"syscall"
	//"time"
	"os/exec"
)

const (
	PORT = ":20019"
)

func main() {

	cmd := exec.Command("mate-terminal", "-x", "go", "run", "tempmain.go")
	cmd.Start()

	fmt.Println("test")
}

/*
	fmt.Println("hei")
	//buf := make([]byte, 1024)
	Network_init()
	//err := connReceive.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	//_, _, err = connReceive.ReadFromUDP(buf) //n contanis numbers of used bytes
	time.Sleep(time.Second)
	fmt.Println("hei")
}
func Network_init() {
	InteruptChan = make(chan os.Signal, 1)
	go Interuption_killer()
	buf := make([]byte, 1024)
	fmt.Println("gi")
	addr, err := ResolveUDPAddr("udp", "129.241.187.255"+PORT) //leser bare fra porten generellt
	c1, err := DialUDP("udp", nil, addr)
	go send(c1)
	if err != nil {
		fmt.Println(err.Error())
	}

	addr2, _ := ResolveUDPAddr("udp", PORT)
	c2, err := ListenUDP("udp", addr2)

	c2.SetReadDeadline(time.Now().Add(300 * time.Millisecond)) //returns error if deadline is reached
	n, address, err := c2.ReadFromUDP(buf)
	str := ipMessage{address, string(buf)}
	fmt.Println("structip", str.ip)
	fmt.Println(string(buf))
	fmt.Println(n)
	//defer c1.Close()
	//defer c2.Close()
}
func send(c1 Conn) {
	for {
		c1.Write([]byte("neooooooo"))
		time.Sleep(1000 * time.Millisecond)
	}
}

type ipOrderMessage struct {
	ip    *UDPAddr
	order string
}

func Interuption_killer() {
	signal.Notify(InteruptChan, os.Interrupt)
	signal.Notify(InteruptChan, syscall.SIGTERM)
	<-InteruptChan
	//what should be done when ctrl-c is pressed?????
	//<- goes here.
	//SystemInit()
	fmt.Println("Got ctrl-c signal")
	os.Exit(0)
}

/*CRTL_C FUNCTION
func cleanup() {
    fmt.Println("cleanup")
}


func main() {
    c := make(chan os.Signal)

    go func() {
    	signal.Notify(c, os.Interrupt)
    	signal.Notify(c, syscall.SIGTERM)
        <-c
        cleanup()
        os.Exit(1)
    }()

    for {
        fmt.Println("sleeping...")
        time.Sleep(10 * time.Second) // or runtime.Gosched() or similar per @misterbee
    }
}
*/
