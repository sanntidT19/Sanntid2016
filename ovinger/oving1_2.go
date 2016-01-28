package main

import (
	."fmt"
	"time"
	"runtime"
)

var i = 0


func thread1(ch1 chan int){
	for x:= 0; x < 1000000; x++{
		i :=<-ch1
		i++
		ch1 <-i
	}
}

func thread2(ch1 chan int){
	for x:= 0; x < 1000001; x++{
		i:=<-ch1
		i--
		ch1<-i
	}
}


func main(){
	runtime.GOMAXPROCS(runtime.NumCPU()) 
	chan_sea := make(chan int,1)
	chan_sea<-i
	go thread1(chan_sea)
	go thread2(chan_sea)
	time.Sleep(1000*time.Millisecond)
	i:= <-chan_sea
	Println("Done:", i);
}
