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
var listOfElevatorsInNetwork []string;


var connectionList map[string] *net.UDPConn



type addrAndTimer struct{
	DeadLine time.Time
	NetWorkAddr string
}

type Message struct{
	Tag string
	Ack bool
	Data []byte 	
}

func main() {
	newShoutFromElevatorChan := make(chan string)
	newElevatorChan := make(chan string)
	newConnectionChan := make(chan string)
	elevatorGoneChan := make(chan string)
	endConnectionChan := make(chan string)
	sendNetworkMessageChan := make(chan []byte)
	

	localAddr := getLocalIP()

	if localAddr == "" {
		fmt.Println("Problem using getLocalIP, expecting office-pc")
		localAddr = "129.241.154.78"
	}
	fmt.Println(localAddr)

	go broadcastPrecense(broadcastPort)
	go listenForBroadcast(broadcastPort, newShoutFromElevatorChan)
	go readUpdateElevatorOverview(newShoutFromElevatorChan, newElevatorChan, elevatorGoneChan)
	go SendMessagesToAllElevators(sendNetworkMessageChan, newConnectionChan, endConnectionChan)
	go readMessagesFromNetwork(localAddr, commonPort)

	go func(){
		for{
			sendNetworkMessageChan <- []byte("I am yolomaster")
			time.Sleep(time.Millisecond*500)
		}

	}()
	for{
		select{
		case elevGone := <- elevatorGoneChan:
			fmt.Println("Elevator gone, address: ", elevGone)
			endConnectionChan <- elevGone
		case newElev := <- newElevatorChan:
			fmt.Println("New elevator, address: ",newElev)
			newConnectionChan <- newElev
		}
	}

	//defer connection.Close()
	fmt.Println("End of main")
}
//Start with int to test. Then network and test this



func readUpdateElevatorOverview(newShoutFromElevatorChan chan string, newElevatorChan chan string, elevatorGoneChan chan string){
	elevatorTimerList := []addrAndTimer{}
	for{
		select{
		case newShout := <- newShoutFromElevatorChan:
			elevatorInList := false
				index := 0;
			//Index might be invalid if list is altered somewhere else during the for loop.
			for i, v := range elevatorTimerList{
				if v.NetWorkAddr == newShout{
					elevatorInList = true
					index = i;
					break;
				}
			}
			if elevatorInList{
				elevatorTimerList[index].DeadLine = time.Now().Add(time.Second * 3) 
			}else{
				elevatorTimerList = append(elevatorTimerList, addrAndTimer{DeadLine: time.Now().Add(time.Second *3), NetWorkAddr: newShout})
				newElevatorChan <- newShout
			}
		default:
			for i, v := range elevatorTimerList{
				if time.Now().After(v.DeadLine){
					elevatorGoneChan <- v.NetWorkAddr
					elevatorTimerList = append(elevatorTimerList[:i], elevatorTimerList[i+1:]...) //From slicetricks. Remove element.
				}
			}
		}
		time.Sleep(time.Millisecond*200)
	}
}

func connectToElevator(remoteIp string, remotePort string) *net.UDPConn {
	fullAddr := remoteIp + ":" + remotePort
	remoteUDPAddr, _ := net.ResolveUDPAddr("udp4", fullAddr)
	connection, _ := net.DialUDP("udp4", nil, remoteUDPAddr)

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

func listenForBroadcast(broadcastPort string, newShoutFromElevatorChan chan string) {
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
			newShoutFromElevatorChan <- senderAddr.IP.String()
		}
	}

}

/* I think you only need to listen on one port */
/* test to iterate through connections if this doesnt work*/ 

func readMessagesFromNetwork(localAddr string, commonPort string) {
	buffer := make([]byte, 2048)
	fullAddr := localAddr + ":" + commonPort
	listenConnAddress, _ := net.ResolveUDPAddr("udp4",fullAddr)
	listenConn, err := net.ListenUDP("udp4",listenConnAddress)
	if err != nil{
		fmt.Println("Error setting up listen-connection")
	}
	for {{
		_, senderAddr, err := listenConn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading message, do nothing")
		}else{
			messageIndex := bytes.IndexByte(buffer, 0)
			message := string(buffer[:messageIndex])
			fmt.Println("Reading message from: ", senderAddr.IP.String(), "Message: ", message)
		}
		}
	}
}

//
func SendMessagesToAllElevators(sendNetworkMessageChan chan []byte, newConnectionChan chan string, endConnectionChan chan string){
	connectionList := make(map[string] *net.UDPConn)
	for _,deferConn := range connectionList{
		defer deferConn.Close()
	}
	for{
		select{
		case newElevator := <-newConnectionChan:
			connectionList[newElevator] = connectToElevator(newElevator, commonPort)
			defer connectionList[newElevator].Close() //This might be enough
		case deadElevator := <-endConnectionChan:
			connectionList[deadElevator].Close() 
		case newMessage := <-sendNetworkMessageChan:
				for _,conn := range connectionList{
					_, err := conn.Write(newMessage)
					if err != nil {
						fmt.Println("Error in sending func, maybe some error handling later")
				}
			}
		}
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
