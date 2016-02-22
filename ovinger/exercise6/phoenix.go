package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
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
	addr := net.UDPAddr{
		Port: 20058,
		IP:   net.ParseIP("129.241.187.151"),
	}
	conn, _ := net.ListenUDP("129.241.187.151", &addr)
	fmt.Println("I am backup")
	masterChan := make(chan bool)
	go detect_precense(conn, masterChan)

	<-masterChan
	fmt.Println("I am master")
	spawn_backup()
	counterChan := make(chan int)
	go read_from_file(counterChan)
	go spam_precense()
	go counter_and_write_to_file(counterChan)
	time.Sleep(time.Second * 10000)
}

func counter_and_write_to_file(counterChan chan int) {
	counter := <-counterChan
	for {
		counter++
		ioutil.WriteFile(FILENAME, counter, os.ModeExclusive)
		fmt.Println(counter)
		time.Sleep(time.Millisecond * 200)
	}
}

func read_from_file(counterChan chan int) {
	fromFile, err := ioutil.ReadFile(FILENAME)
	if err != nil {
		panic(err)
	}
	counter := int(fromFile)
	counterChan <- counter
}

func spawn_backup() {
	cmd := exec.Command(gnome-terminal-x[""], "go run phoenix.go")
	err := exec.Run(cmd)
	if err != nil {
		fmt.Println("You fucked up in spawn_backup")
		panic(err)
	}
}

func spam_precense() {
	conn, err := net.Dial("udp", LOCAL_ADDR)
	if err != nil {
		panic(err)
	}
	for {
		_, err := conn.WriteToUDP([]byte("I am the master motherfucker"))
		if err != nil {
			fmt.Println("You fucked up in spam_precense")
			panic(err)
		}
		time.Sleep(time.Millisecond * 500)
	}
}

func detect_precense(conn *UDPConn, masterChan chan bool) {
	timerChan := time.After(time.Second * 3)
	p := make([]byte, 2048)
L:
	for {
		select {
		case <-timerChan:
			fmt.Println("Master dead")
			masterChan <- true
			break L
		default:
			_, _, err := conn.ReadFromUDP(p)
			if err != nil {
				fmt.Println("You fucked up in detect_precense.")
				panic(err)
			}
		}
	}
}
