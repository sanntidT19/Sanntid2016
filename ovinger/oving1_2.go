package main

import (
	."fmt"
	"time"
	//"runtime"
)

var i = 0

func thread1(){
	for x:= 0; x < 1000000; x++{
		i++
	}
}

func thread2(){
	for x:= 0; x < 1000000; x++{
		i--
		
	}
}


func main(){
	//runtime.GOMAXPROCS(runtime.NumCPU()) 
	go thread1()
	go thread2()
	time.Sleep(1000*time.Millisecond)
	Println("Done:", i);
}
