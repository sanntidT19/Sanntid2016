package network

const BROADCAST_DEADLINE = 3

type addrAndDeadline struct {
	DeadLine    time.Time
	Addr string
}


var broadcastPort string = "30059"
var listOfElevatorsInNetwork []string

//Denne kan kalles initNetworkAndAlertOfChanges
func CommNeedBetterName() {
	newShoutFromElevatorChan := make(chan string)
	newElevatorChan := make(chan string)
	newConnectionChan := make(chan string)
	elevatorGoneChan := make(chan string)
	endConnectionChan := make(chan string)
	sendNetworkMessageChan := make(chan []byte)
	messageFromNetworkChan := make(chan []byte)
	networkDownRemoveAcksChan := make(chan bool)

	ackdByAllChan := make(chan MessageWithHeader)
	resendMessageChan := make(chan MessageWithHeader)
	newAckFromNetworkChan := make(chan MessageWithHeader)
	newAckStartChan := make(chan MessageWithHeader)
	elevatorListChangedChan := make(chan bool)
	LocalAddr = FindLocalIP()
	if LocalAddr == "" {
		fmt.Println("Problem using getLocalIP, expecting office-pc")
		LocalAddr = "129.241.154.78"
	}

	go broadcastPrecense(broadcastPort)
	go listenForBroadcast(broadcastPort, newShoutFromElevatorChan)
	go detectNewAndDeadElevs(newShoutFromElevatorChan, newElevatorChan, elevatorGoneChan)
	//go DistributeOrdersToNetwork()
	go func() {
		for {
			select {
			case elevGone := <-elevatorGoneChan:
				fmt.Println("ip list before going through the list. Case: elevgone", listOfElevatorsInNetwork)
				pos := -1
				for i, v := range listOfElevatorsInNetwork {
					if v == elevGone {
						pos = i
						break
					}
				}
				if pos == -1 {
					fmt.Println("main go func: elevator not found, pos == -1")
				} else {
					fmt.Println("position in list", pos)
					fmt.Println("list before change: ", listOfElevatorsInNetwork)
					listOfElevatorsInNetwork = append(listOfElevatorsInNetwork[:pos], listOfElevatorsInNetwork[pos+1:]...)
					fmt.Println("list after change: ", listOfElevatorsInNetwork)

				}
				fmt.Println("Elevator gone, address: ", elevGone)
				endConnectionChan <- elevGone
				elevatorListChangedChan <- true
				FromNetworkElevGoneChan <- elevGone
				if len(listOfElevatorsInNetwork) == 0 {
					fmt.Println("Network is gone!")
					FromNetworkNetworkDownChan <- true
					networkDownRemoveAcksChan <- true
				}

			case newElev := <-newElevatorChan:
				fmt.Println("New elevator, address: ", newElev)
				fmt.Println("list before change: ", listOfElevatorsInNetwork)
				listOfElevatorsInNetwork = append(listOfElevatorsInNetwork, newElev)
				fmt.Println("after appending ip to list : ", listOfElevatorsInNetwork)
				newConnectionChan <- newElev
				elevatorListChangedChan <- true
				FromNetworkNewElevChan <- newElev
				if len(listOfElevatorsInNetwork) == 1 {
					FromNetworkNetworkUpChan <- true
				}
			}
		}
	}()
	for {
		select {
		case <-ackdByAllChan:
			//fmt.Println("ackd by all!")
		case <-FromNetworkNetworkUpChan:
			fmt.Println("Network is up!")

		}
	}
}

func connectToElevator(remoteIp string, remotePort string) *net.UDPConn {
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
	broadcastAddr := "255.255.255.255" + ":" + broadcastPort
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
	listenBroadcastAddress := "0.0.0.0" + ":" + broadcastPort
	broadcastUDPAddr, _ := net.ResolveUDPAddr("udp4", listenBroadcastAddress)
	connection, _ := net.ListenUDP("udp4", broadcastUDPAddr)
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

func GetElevList() []string {
	copyElevList := make([]string, len(listOfElevatorsInNetwork))
	copy(copyElevList, listOfElevatorsInNetwork)
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
				newElevChan <- newShout
			}
		default:
			for i, v := range elevDeadlineList {
				if time.Now().After(v.DeadLine) {
					elevDeadChan <- v.Addr
					elevDeadlineList = append(elevDeadlineList[:i], elevDeadlineList[i+1:]...) //From slicetricks. Remove element.
				}
			}
		}
		time.Sleep(time.Millisecond * 200)
	}
}
