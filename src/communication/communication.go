package main //communication for når denne skal kjøres, utelukkende

import (
	"bytes"
	"fmt"
	"net"
	"reflect"
	"time"
)

var broadcastAddr string = "255.255.255.255:20059" //"129.241.187.255:20059"
var commonPort string = "20059"
var broadcastPort string = "30059"

func main() {
	localAddr := getLocalIP()

	if localAddr == "" {
		fmt.Println("Problem using getLocalIP")
	}
	fmt.Println(localAddr)

	go broadcastPrecense(broadcastPort)
	go listenForElevators(broadcastPort)

	//connection := establish_connection(listenBroadCast, commonPort)
	//go readMessages(connection)
	//defer connection.Close()
	fmt.Println("End of main")
	time.Sleep(time.Second * 1000)
}

func establish_connection(remoteIp string, remotePort string) *net.UDPConn {
	fullAddr := remoteIp + ":" + remotePort
	remoteUDPAddr, _ := net.ResolveUDPAddr("udp4", fullAddr)
	connection, _ := net.ListenUDP("udp4", remoteUDPAddr)

	return connection

}

func getLocalIP() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Whoops..")
		panic(err)
	}
	for _, i := range ifaces {
		addrs, _ := i.Addrs()

		ipLow := net.ParseIP("129.241.187.000")
		ipHigh := net.ParseIP("129.241.187.255")
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if bytes.Compare(ip, ipLow) >= 0 && bytes.Compare(ip, ipHigh) <= 0 {
				fmt.Println("Printing type: ", reflect.TypeOf(ip))
				return ip.String()
			}
		}
	}
	return ""
}

func broadcastPrecense(broadcastPort string) {
	broadcastAddr := "255.255.255.255" + ":" + broadcastPort
	broadcastUDPAddr, _ := net.ResolveUDPAddr("udp", broadcastAddr)
	connection, err := net.DialUDP("udp", nil, broadcastUDPAddr)
	if err != nil {
		fmt.Println("You messed up in spam presence")
	}
	defer connection.Close()
	for {
		_, err := connection.Write([]byte("I am here"))
		if err != nil {
			fmt.Println("BroadcastPrecense: msg not sent")
		}
		time.Sleep(time.Millisecond * 500)
	}
}

func listenForElevators(broadcastPort string) {
	buffer := make([]byte, 2048)
	listenBroadcastAddress := "0.0.0.0" + ":" + broadcastPort
	broadcastUDPAddr, _ := net.ResolveUDPAddr("udp4", listenBroadcastAddress)
	connection, _ := net.ListenUDP("udp4", broadcastUDPAddr)
	defer connection.Close()
	for {
		_, senderAddr, err := connection.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading msg in listenForElevators, discard message")
		} else {
			fmt.Println("Sender: ", senderAddr.IP.String())
			//SEND TO OTHER GOROUTINE HANDLING ALL ELEVATORS PRESENT
		}
	}

}

func readMessages(remoteUDPConn *net.UDPConn) {
	buffer := make([]byte, 2048)
	for {
		_, senderAddr, err := remoteUDPConn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading message, do nothing")
		}
		fmt.Println("Sender: ", senderAddr.IP.String())
	}
}

func sumAllIncomingData(intsIncomingChan chan int) {
	var sum int = 0
	for {
		newInt := <-intsIncomingChan
		sum += newInt
	}
	fmt.Println("Current sum received ", sum)
}

func countBoolsIncoming(boolsIncomingChan chan bool) {
	var falseCounter int = 0
	var trueCounter int = 0
	for {
		newBool := <-boolsIncomingChan
		if newBool {
			trueCounter++
		} else {
			falseCounter++
		}
		fmt.Println("Number of true received:", trueCounter)
		fmt.Println("Number of false received: ", falseCounter)
	}
}
