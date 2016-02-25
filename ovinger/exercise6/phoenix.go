package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	//"os/exec"
	"encoding/binary"
	"time"
)

const FILENAME = "counterSave.txt"

/*
Need:
Listen function that writes to "master" channel when master has been silent for 3 seconds
When "master" is written to, start a backup
simple counter that writes to file.
Simple reader that reads from file when it is master
*/
func main() {
	/*addr := net.UDPAddr{
		Port: 20058,
		IP:   net.ParseIP("129.241.187.151"),
	}*/

	fmt.Println("I am backup")

	localIP := "129.241.187.151"
	localPort := "20058"
	localAddr := localIP + ":" + localPort

	udpAddr, err := net.ResolveUDPAddr("udp4", localAddr)

	if err != nil {
		panic(err)
	}

	connection, err := net.ListenUDP("udp4", udpAddr)
	if err != nil {
		panic(err)
	}
	fmt.Println("Listening on port ", localPort)
	defer connection.Close()

	masterChan := make(chan bool)
	go detect_precense(connection, masterChan)

	<-masterChan
	fmt.Println("I am master")
	//spawn_backup()
	counterChan := make(chan int64)
	go read_from_file(counterChan)
	go spam_precense(localAddr)
	go counter_and_write_to_file(counterChan)
	time.Sleep(time.Second * 10000)
}

func counter_and_write_to_file(counterChan chan int64) {
	counter := <-counterChan
	var binary_counter []byte
	for {
		counter++
		fmt.Println("counter1", counter)
		_ = binary.PutVarint(binary_counter, counter)
		ioutil.WriteFile(FILENAME, binary_counter, os.ModeExclusive)
		fmt.Println("counter2", counter)
		time.Sleep(time.Millisecond * 200)
	}
}

//Need to check for the case when there's nothing there.
func read_from_file(counterChan chan int64) {
	var counter int64
	if _, err := os.Stat("FILENAME"); os.IsNotExist(err) {
		fmt.Println("no save exists")
		counter = 0
	} else {
		fromFile, err := ioutil.ReadFile(FILENAME)
		if err != nil {
			fmt.Println(err)
		}
		counter, _ = binary.Varint(fromFile)
	}
	counterChan <- counter
}

/*
func spawn_backup() {
	cmd := exec.Command(gnome-terminal-x[""], "go run phoenix.go")
	err := exec.Run(cmd)
	if err != nil {
		fmt.Println("You messed up in spawn_backup")
		panic(err)
	}
}
*/
func spam_precense(localAddr string) {
	udpRemote, _ := net.ResolveUDPAddr("udp", localAddr)

	connection, err := net.DialUDP("udp", nil, udpRemote)
	if err != nil {
		fmt.Println("You messed up in spam presenceclear")
		panic(err)
	}
	defer connection.Close()
	for {
		_, err := connection.Write([]byte("I am the master"))
		if err != nil {
			fmt.Println("You messed up in spam_precense")
			panic(err)
		}
		time.Sleep(time.Millisecond * 500)
	}
}

func detect_precense(connection *net.UDPConn, masterChan chan bool) {
	//timerChan := time.After(time.Second * 3)
	buffer := make([]byte, 2048)
	for {
		t := time.Now()
		connection.SetDeadline(t.Add(3 * time.Second))
		_, _, err := connection.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("UDP timeout: ", err)
			masterChan <- true
			break
		}
		/*
			fmt.Println("I'm getting into this for")
			select {
			case <-time.After(time.Second * 3):
				fmt.Println("Master dead")
				masterChan <- true
				break L
				/*
					default:
						_, _, err := connection.ReadFromUDP(buffer)
						if err != nil {
							fmt.Println("You messed up in detect_precense.")
							panic(err)
						}
			}(*/
	}
}
