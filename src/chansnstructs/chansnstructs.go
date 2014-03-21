package chansnstruct

import (
	. "net"
	"os"
	"time"
)

const (
	MAXWAIT  = time.Second
	PORT     = ":20019"
	N_FLOORS = 4

	IP1 = "129.241.187.147"
	//IP2 = 129.241.187.xxx
	//IP3 = 129.241.187.xxx
)

var ExNetChans NetworkExternalChannels
var ExSlaveChans ExternalSlaveChannels
var ExMasterChans ExternalMasterChannels
var ExCommChans ExternalCommunicationChannels
var ExStateMChans ExternalStateMachineChannels
var ExOptimalChans ExternalOptimalizationChannels

var InteruptChan chan os.Signal

type Master struct {
	S map[*UDPAddr]Slave
}

type Slave struct {
	AllExternalsOrder map[*UDPAddr][]Order
	InternalList      []int
	CurrentFloor      int
	Direction         int
}
type IpSlave struct {
	Ip *UDPAddr
	S  Slave
}

type IpOrderMessage struct {
	Ip    *UDPAddr
	Order Order
}

type Order struct {
	Floor      int
	ButtonType int
	TurnOn     bool
}

type State struct {
	Direction    int
	CurrentFloor int
}

type NetworkExternalChannels struct {
	ToNetwork  chan []byte
	ToComm     chan []byte
	ConnChan   chan Conn
	ToCommAddr chan *UDPAddr
	StopWrite  chan bool
}
type ExternalOptimalizationChannels struct {
	//InMasterChans.OptimizationInitChan = make(chan Master)
	OptimizationTriggerChan chan Master
	OptimizationReturnChan  chan [][]bool
}

type ExternalCommunicationChannels struct {
	//communication channels
	ToMasterSlaveChan                    chan IpSlave        //"sla"
	ToMasterOrderListReceivedChan        chan IpOrderMessage //"ore"
	ToMasterOrderExecutedChan            chan IpOrderMessage //"oex"
	ToMasterOrderExecutedReConfirmedChan chan IpOrderMessage //"oce"
	ToMasterExternalButtonPushedChan     chan IpOrderMessage //"ebp"
	ToSlaveOrderListChan                 chan [][]int        //"exo"
	//ToSlaveReceivedOrderListConfirmationChan chan ipOrderMessage //"rco"
	ToSlaveOrderExecutedConfirmedChan chan IpOrderMessage //"eco"
	ToSlaveImMasterChan               chan string         //"iam"

}
type ExternalSlaveChannels struct {
	ToCommSlaveChan                    chan Slave  //"sla"
	ToCommOrderListReceivedChan        chan []int  //"ore"
	ToCommOrderExecutedChan            chan []int  //"oex"
	ToCommOrderExecutedReConfirmedChan chan []int  //"oce"
	ToCommExternalButtonPushedChan     chan []int  //"ebp"
	ToCommImMasterChan                 chan string //"iam"

}
type ExternalMasterChannels struct {
	ToCommOrderListChan              chan [][]int //"exo"
	ToCommOrderExecutedConfirmedChan chan []int   //"eco"

}
type ExternalStateMachineChannels struct {
	ExternalButtonPressed chan []int
	OrderServed           chan []int
	CurrentState          chan []int
	GetSlaveStruct        chan bool
	ReturnSlaveStruct     chan Slave
	DirectionUpdate       chan int
}

func Channals_init() {
	network_external_chan_init()
	external_comm_channels_init()
	Slave_external_chans_init()
	Master_external_chans_init()
	Master_external_chans_init()
	External_state_machine_channels_init()
}

func network_external_chan_init() {
	ExNetChans.ToNetwork = make(chan []byte)
	ExNetChans.ToComm = make(chan []byte)
	ExNetChans.ConnChan = make(chan Conn)
	ExNetChans.ToCommAddr = make(chan *UDPAddr)
}

func external_comm_channels_init() {
	ExCommChans.ToSlaveOrderListChan = make(chan [][]int)                        //"ord"
	ExCommChans.ToMasterOrderListReceivedChan = make(chan IpOrderMessage)        //"ore"
	ExCommChans.ToMasterOrderExecutedChan = make(chan IpOrderMessage)            //"oex"
	ExCommChans.ToSlaveOrderExecutedConfirmedChan = make(chan IpOrderMessage)    //"eco"
	ExCommChans.ToMasterOrderExecutedReConfirmedChan = make(chan IpOrderMessage) //"oce"
	ExCommChans.ToMasterExternalButtonPushedChan = make(chan IpOrderMessage)     //"ebp"
	ExCommChans.ToMasterSlaveChan = make(chan IpSlave)                           //"sla"
	ExCommChans.ToSlaveImMasterChan = make(chan string)                          //"iam"
}

func Slave_external_chans_init() {
	ExSlaveChans.ToCommOrderListReceivedChan = make(chan []int) //"ore"
	//ExSlaveChans.ToCommOrderReceivedChan = make(chan ipOrderMessage)            //"oce"
	ExSlaveChans.ToCommOrderExecutedChan = make(chan []int)            //"oex"
	ExSlaveChans.ToCommOrderExecutedReConfirmedChan = make(chan []int) //"oce"
	ExSlaveChans.ToCommExternalButtonPushedChan = make(chan []int)     //"ebp"
	ExSlaveChans.ToCommSlaveChan = make(chan Slave)                    //"sla"

}
func Master_external_chans_init() {
	ExMasterChans.ToCommOrderListChan = make(chan [][]int)            //"ord"
	ExMasterChans.ToCommOrderExecutedConfirmedChan = make(chan []int) //"eco"
}

func External_state_machine_channels_init() {
	ExStateMChans.ExternalButtonPressed = make(chan []int)
	ExStateMChans.OrderServed = make(chan []int)
	ExStateMChans.CurrentState = make(chan []int)
	ExStateMChans.GetSlaveStruct = make(chan bool)
	ExStateMChans.ReturnSlaveStruct = make(chan Slave)
	ExStateMChans.DirectionUpdate = make(chan int)
}
func External_optimization_channel_init() {
	ExOptimalChans.OptimizationTriggerChan = make(chan Master)
	ExOptimalChans.OptimizationReturnChan = make(chan [][]bool)
}
