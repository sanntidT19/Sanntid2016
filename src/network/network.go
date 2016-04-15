package network

import (
	"bytes"
	"fmt"
	"net"
	"time"
)

type addrAndDeadline struct {
	DeadLine time.Time
	Addr     string
}

const BROADCAST_DEADLINE = 2

var broadcastPort string = "30059"
var listOfElevsInNetwork []string

//Denne kan kalles initNetworkAndAlertOfChanges

/*
SØRG FOR AT ALLE KANALER SOM TRENGER Å VITE AT LISTEN HAR ENDRET SEG, FÅR VITE DETTE.
FOR ØYEBLIKKET ER DET KUN Å GJØRE DEN ENE KANALEN TIL MESSAGES GLOBAL
SPØR OM EN LUR MÅTE Å BRUKE GLOBALE KANALER PÅ I MORRA
*/
func InitNetworkAndAlertChanges(newConnectionChan chan string, endConnectionChan chan string, resendUnackdMessagesChan chan bool, alertNewElevChan chan string, alertDeadElevChan chan string, alertNetworkDownChan chan bool) {
	newShoutFromElevatorChan := make(chan string)
	newElevChan := make(chan string)
	elevDeadChan := make(chan string)
	go broadcastPrecense(broadcastPort)
	go listenForBroadcast(broadcastPort, newShoutFromElevatorChan)
	go detectNewAndDeadElevs(newShoutFromElevatorChan, newElevChan, elevDeadChan)
	for {
		select {
		case elevGone := <-elevDeadChan:
			fmt.Println("ip list before going through the list. Case: elevgone", listOfElevsInNetwork)
			pos := -1
			for i, v := range listOfElevsInNetwork {
				if v == elevGone {
					pos = i
					break
				}
			}
			if pos == -1 {
				fmt.Println("main go func: elevator not found, pos == -1")
			} else {
				fmt.Println("position in list", pos)
				fmt.Println("list before change: ", listOfElevsInNetwork)
				listOfElevsInNetwork = append(listOfElevsInNetwork[:pos], listOfElevsInNetwork[pos+1:]...)
				fmt.Println("list after change: ", listOfElevsInNetwork)

			}
			fmt.Println("Elevator gone, address: ", elevGone)
			endConnectionChan <- elevGone
			resendUnackdMessagesChan <- true
			alertDeadElevChan <- elevGone
			if len(listOfElevsInNetwork) == 0 {
				fmt.Println("Network is gone!")
				alertNetworkDownChan <- true
				resendUnackdMessagesChan <-false
			}

		case newElev := <-newElevChan:
			fmt.Println("New elevator, address: ", newElev)
			fmt.Println("list before change: ", listOfElevsInNetwork)
			listOfElevsInNetwork = append(listOfElevsInNetwork, newElev)
			fmt.Println("after appending ip to list : ", listOfElevsInNetwork)
			fmt.Println("                           NewMessages is receiving")
			newConnectionChan <- newElev
			resendUnackdMessagesChan <-true
			fmt.Println("                           NewMessages is reading")
			alertNewElevChan <- newElev //send to somewhere else
		}
		fmt.Println("end of network-select")
	}
}

//tenk litt gjennom om/hvordan denne bør lages
/*
func alertSystemOfNetworkChanges(newElev string, deadElev string, networkIsDown bool){
	if newElev != ""{

	}else if deadElev != ""{

	}
}
*/

func ConnectToElevator(remoteIp string, remotePort string) *net.UDPConn {
	fullAddr := remoteIp + ":" + remotePort
	remoteUDPAddr, _ := net.ResolveUDPAddr("udp4", fullAddr)
	connection, _ := net.DialUDP("udp4", nil, remoteUDPAddr)
	return connection

}

func FindLocalIP() string {
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
				//fmt.Println("Printing type: ", reflect.TypeOf(ip))
				return ip.String()
			}
		}
	}
	return ""
}

//Broadcasts to network "i am here"
func broadcastPrecense(broadcastPort string) {
	broadcastAddr := "129.241.187.255" + ":" + broadcastPort
	broadcastUDPAddr, _ := net.ResolveUDPAddr("udp", broadcastAddr)
	connection, err := net.DialUDP("udp", nil, broadcastUDPAddr)
	if err != nil {
		fmt.Println(err)
	}
	defer connection.Close()
	for {
		_, err := connection.Write([]byte("I am here"))
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(time.Millisecond * 500)
	}
}

func listenForBroadcast(broadcastPort string, newShoutFromElevatorChan chan string) {
	buffer := make([]byte, 2048)
	listenBroadcastAddress := ":" + broadcastPort
	broadcastUDPAddr, _ := net.ResolveUDPAddr("udp", listenBroadcastAddress)
	connection, _ := net.ListenUDP("udp", broadcastUDPAddr)
	defer connection.Close()
	for {
		_, senderAddr, err := connection.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error receiving broadcast, discard message")
		} else {
			newShoutFromElevatorChan <- senderAddr.IP.String()
		}
	}

}

//Navnet på denne må endres, slik at man skjønner hva man får av den
func ElevsSeen() []string {
	copyElevList := make([]string, len(listOfElevsInNetwork))
	copy(copyElevList, listOfElevsInNetwork)
	return copyElevList
}

func detectNewAndDeadElevs(newShoutFromElevChan chan string, newElevChan chan string, elevDeadChan chan string) {
	elevDeadlineList := []addrAndDeadline{}
	for {
		select {
		case newShout := <-newShoutFromElevChan:
			elevIsInList := false
			placeInList := 0
			for i, v := range elevDeadlineList {
				if v.Addr == newShout {
					elevIsInList = true
					placeInList = i
					break
				}
			}
			if elevIsInList {
				elevDeadlineList[placeInList].DeadLine = time.Now().Add(time.Second * BROADCAST_DEADLINE)
			} else {
				elevDeadlineList = append(elevDeadlineList, addrAndDeadline{DeadLine: time.Now().Add(time.Second * BROADCAST_DEADLINE), Addr: newShout})
				fmt.Println("newElevChan <- newShout")
				newElevChan <- newShout
			}
		default:
			for i, v := range elevDeadlineList {
				fmt.Println("Now: ", time.Now())
				fmt.Println("Deadline: ", v.DeadLine)
				if time.Now().After(v.DeadLine) {
					fmt.Println("                            elevDeadChan <- v.Addr")
					elevDeadChan <- v.Addr
					elevDeadlineList = append(elevDeadlineList[:i], elevDeadlineList[i+1:]...) //From slicetricks. Remove element.
					break
				}
			}
		}
		time.Sleep(time.Millisecond * 200)
	}
}
