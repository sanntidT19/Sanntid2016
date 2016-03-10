package main //communication for når denne skal kjøres, utelukkende

import (
	"bytes"
	"fmt"
	"net"
	"reflect"
	"time"
)

var broadcastAddr string = "129.241.187.255:20059"
var commonPort string = "20059"

func main() {
	localAddr := getLocalIP()
	if localAddr == "" {
		fmt.Println("Problem using getLocalIP")
	}
	fmt.Println(localAddr)
	spam_precense(broadcastAddr)
	connection := establish_connection(localAddr, commonPort)
	go readMessages(connection)
	defer connection.Close()
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
			fmt.Println("Printing type: ", reflect.TypeOf(ip))
			if bytes.Compare(ip, ipLow) >= 0 && bytes.Compare(ip, ipHigh) <= 0 {
				return ip.String()
			}
		}
	}
	return ""
}

func spam_precense(remoteAddr string) {
	udpRemote, _ := net.ResolveUDPAddr("udp", remoteAddr)

	connection, err := net.DialUDP("udp", nil, udpRemote)
	if err != nil {
		fmt.Println("You messed up in spam presence")
		panic(err)
	}
	defer connection.Close()
	for {
		_, err := connection.Write([]byte("I am here"))
		if err != nil {
			fmt.Println("You messed up in spam_precense")
			panic(err)
		}
		time.Sleep(time.Millisecond * 500)
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
