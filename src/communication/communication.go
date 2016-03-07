package communication //main for når denne skal kjøres, utelukkende

import (
	"fmt"
	"net"
	"time"
	//"reflect"
	"bytes"
)

/*
func main() {
	fmt.Println("Main started")
	localIP := getLocalIP()
	fmt.Println(localIP)
	time.Sleep(time.Second * 1)
}
*/
func getLocalIP() net.IP {
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
			//fmt.Println("Printing type: ", reflect.TypeOf(ip))
			if bytes.Compare(ip, ipLow) >= 0 && bytes.Compare(ip, ipHigh) <= 0 {
				return ip
			}
		}
	}
	return nil
}
