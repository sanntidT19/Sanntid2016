package chansnstruct

import (
	. "net"
	"os"
	"time"
)

const (
	UP        = 0 // is this correct?
	DOWN      = 1 // si this correct?
	N_BUTTONS = 8 // @Yngve is this correct??
	MAXWAIT   = time.Second
	PORT      = ":20019"
	N_FLOORS  = 4

	//IP1 = "129.241.187.147"
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
	SlaveElev map[*UDPAddr]Slave
}

type Slave struct {
	IP                *UDPAddr
	AllExternalsOrder map[*UDPAddr][]Order
	InternalList      []bool
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
	ToCommAddr chan *UDPAddr
	StopWrite  chan bool
}
type ExternalOptimalizationChannels struct {
	//InMasterChans.OptimizationInitChan = make(chan Master)
	OptimizationTriggerChan chan Order
	OptimizationReturnChan  chan []Order
}

type ExternalCommunicationChannels struct {
	//communication channels
	ToMasterSlaveChan                    chan IpSlave        //"sla"
	ToMasterOrderListReceivedChan        chan IpOrderMessage //"ore"
	ToMasterOrderExecutedChan            chan IpOrderMessage //"oex"
	ToMasterOrderExecutedReConfirmedChan chan IpOrderMessage //"oce"
	ToMasterExternalButtonPushedChan     chan IpOrderMessage //"ebp"
	ToSlaveOrderListChan                 chan []Order        //"exo"
	ToSlaveOrderExecutedConfirmedChan    chan IpOrderMessage //"eco"
	ToSlaveImMasterChan                  chan string         //"iam"
	ToMasterImSlaveChan                  chan IpOrderMessage //"ias"

}
type ExternalSlaveChannels struct {
	ToCommSlaveChan                    chan Slave          //"sla"
	ToCommOrderListReceivedChan        chan Order          //"ore"
	ToCommOrderExecutedChan            chan Order          //"oex"
	ToCommOrderExecutedReConfirmedChan chan Order          //"oce"
	ToCommExternalButtonPushedChan     chan Order          //"ebp"
	ToCommImSlaveChan                  chan IpOrderMessage //"ias"

}
type ExternalMasterChannels struct {
	ToCommOrderListChan              chan []Order //"exo"
	ToCommOrderExecutedConfirmedChan chan Order   //"eco"
	ToCommImMasterChan               chan string  //"iam"
}
type ExternalStateMachineChannels struct {
	ExternalButtonPressed chan Order
	OrderServed           chan Order
	CurrentState          chan Order
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
	ExNetChans.ToCommAddr = make(chan *UDPAddr)
}

func external_comm_channels_init() {
	ExCommChans.ToSlaveOrderListChan = make(chan []Order)                        //"ord"
	ExCommChans.ToMasterOrderListReceivedChan = make(chan IpOrderMessage)        //"ore"
	ExCommChans.ToMasterOrderExecutedChan = make(chan IpOrderMessage)            //"oex"
	ExCommChans.ToSlaveOrderExecutedConfirmedChan = make(chan IpOrderMessage)    //"eco"
	ExCommChans.ToMasterOrderExecutedReConfirmedChan = make(chan IpOrderMessage) //"oce"
	ExCommChans.ToMasterExternalButtonPushedChan = make(chan IpOrderMessage)     //"ebp"
	ExCommChans.ToMasterSlaveChan = make(chan IpSlave)                           //"sla"
	ExCommChans.ToSlaveImMasterChan = make(chan string)                          //"iam"
	ExCommChans.ToMasterImSlaveChan = make(chan IpOrderMessage)                  //"ims"
}

func Slave_external_chans_init() {
	ExSlaveChans.ToCommOrderListReceivedChan = make(chan Order) //"ore"
	//ExSlaveChans.ToCommOrderReceivedChan = make(chan ipOrderMessage)            //"oce"
	ExSlaveChans.ToCommOrderExecutedChan = make(chan Order)            //"oex"
	ExSlaveChans.ToCommOrderExecutedReConfirmedChan = make(chan Order) //"oce"
	ExSlaveChans.ToCommExternalButtonPushedChan = make(chan Order)     //"ebp"
	ExSlaveChans.ToCommSlaveChan = make(chan Slave)                    //"sla"
	ExSlaveChans.ToCommImSlaveChan = make(chan IpOrderMessage)

}
func Master_external_chans_init() {
	ExMasterChans.ToCommOrderListChan = make(chan []Order)            //"ord"
	ExMasterChans.ToCommOrderExecutedConfirmedChan = make(chan Order) //"eco"
	ExMasterChans.ToCommImMasterChan = make(chan string)
}

func External_state_machine_channels_init() {
	ExStateMChans.ExternalButtonPressed = make(chan Order)
	ExStateMChans.OrderServed = make(chan Order)
	ExStateMChans.CurrentState = make(chan Order)
	ExStateMChans.GetSlaveStruct = make(chan bool)
	ExStateMChans.ReturnSlaveStruct = make(chan Slave)
	ExStateMChans.DirectionUpdate = make(chan int)
}
func External_optimization_channel_init() {
	ExOptimalChans.OptimizationTriggerChan = make(chan Order)
	ExOptimalChans.OptimizationReturnChan = make(chan []Order)
}
