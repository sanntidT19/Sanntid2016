package modules

import (
	"fmt"
	. "json"
	. "strings"
)

func Channels_init() {
	stateElevChan := make(chan int) //int?
	elevIOChannel := make(chan int) //int?

	slaveStateChan := make(chan int)

	stateElevChan := make(chan int)
	slaveCommChan := make(chan string)

	masterSlaveOderChan := make(chan [][]int)
	masterSlaveStringChan := make(chan string)

	masterCommOrderChan := make(chan [][]int)
	masterCommStringChan := make(chan string)

	commNetworkChan := make(chan []byte)
}

func Send_exteral_list(pack [][]int) []byte {
	pack := Marshall(pack)
	return "ord" + pack
}

func Send_i_am_slave(pack string) []byte {
	pack := Marshall(pack)
	return "iah" + pack
}

func Send_i_am_master(pack string) []byte {
	pack := Marshall(pack)
	return "iam" + pack
}

func Send_order_performed(pack string) []byte {
	pack := Marshall(pack)
	return "per" + pack
}

func Send_confirmation_from_slave(pack string) []byte {
	pack := Marshall(pack)
	return "sre" + pack
}

func Send_confirmation_from_master(pack string) []byte {
	pack := Marshall(pack)
	return "mre" + pack
}

func Send_state_changed(pack string) []byte {
	pack := Marshall(pack)
	return "sch" + pack
}

func Decrypt_message(message []byte) []byte {
	switch {
	case HasPrefix(message, "ias"): //I am slave
		message = TrimPrefix(message, "ias")
		str := Unmarshall(message)
		iasChan <- str

	case HasPrefix(message, "iam"): //I am master
		message = TrimPrefix(message, "iam")
		str := Unmarshall(message)
		iamChan <- str

	case HasPrefix(message, "per"): //Order performed
		message = TrimPrefix(message, "per")
		str := Unmarshall(message)
		perChan <- str

	case HasPrefix(message, "sre"): //recived order from slave
		message = TrimPrefix(message, "sec")
		str := Unmarshall(message)
		secChan <- str

	case HasPrefix(message, "ord"): //externalorderlist
		message = TrimPrefix(message, "ord")
		ordrlist := Unmarshall(message)
		ordChan <- ordrlist

	case HasPrefix(message, "mre"): //recived order from master
		message = TrimPrefix(message, "mre")
		str := Unmarshall(message)
		mreChan <- str

	case HasPrefix(message, "sch"): // state changed??????????what cind of variable is this
		message = TrimPrefix(message, "sch")
		str := Unmarshall(message)
		schChan <- str
	}
	return message
}
