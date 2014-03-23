package chansnstruct

import (
	. "net"
	"os"
	"time"
)

const (
	UP         = 0 // is this correct?
	DOWN       = 1 // si this correct?
	N_BUTTONS  = 8 // @Yngve is this correct??
	MAXWAIT    = time.Second
	PORT       = ":20019"
	N_FLOORS   = 4
	TIMEOFSTOP = 1 // Time spent in one floor
	TIMETRAVEL = 2 // Time of travel between floors
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

type IpState struct {
	Ip           *UDPAddr
	Direction    int
	CurrentFloor int
}

type Master struct {
	SlaveIp      []*UDPAddr
	ExternalList map[*UDPAddr]*[N_FLOORS][2]bool
	Statelist    map[*UDPAddr]IpState
}

type Slave struct {
	IP           *UDPAddr
	ExternalList map[*UDPAddr]*[N_FLOORS][2]bool
	InternalList []bool
	CurrentFloor int
	Direction    int
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

type IpOrderList struct {
	Ip           *UDPAddr
	ExternalList map[*UDPAddr]*[N_FLOORS][2]bool
}

type NetworkExternalChannels struct {
	ToNetwork  chan []byte
	ToComm     chan []byte
	ToCommAddr chan *UDPAddr
	StopWrite  chan bool
}
type ExternalOptimalizationChannels struct {
	//InMasterChans.OptimizationInitChan = make(chan Master)
	OptimizationTriggerChan chan IpOrderMessage
	OptimizationReturnChan  chan IpOrderMessage
}

type ExternalCommunicationChannels struct {
	//communication channels
	ToMasterSlaveChan                    chan IpSlave        //"sla"
	ToMasterOrderListReceivedChan        chan IpOrderList    //"ore"
	ToMasterOrderExecutedChan            chan IpOrderMessage //"oex"
	ToMasterOrderExecutedReConfirmedChan chan IpOrderMessage //"oce"
	ToMasterExternalButtonPushedChan     chan IpOrderMessage //"ebp"
	ToMasterImSlaveChan                  chan IpOrderMessage //"ias"
	ToMasterUpdateState                  chan IpState        //"ust"
	ToSlaveOrderListChan                 chan IpOrderList    //"exo"
	ToSlaveOrderExecutedConfirmedChan    chan IpOrderMessage //"eco"
	ToSlaveImMasterChan                  chan string         //"iam"
	ToSlaveUpdateStateReceivedChan       chan IpState        //"sus"
	ToSlaveButtonPressedConfirmedChan    chan IpOrderMessage //"bpc"

}
type ExternalSlaveChannels struct {
	ToCommSlaveChan                    chan Slave          //"sla"
	ToCommOrderListReceivedChan        chan IpOrderList    //"ore"
	ToCommOrderExecutedChan            chan IpOrderMessage //"oex"
	ToCommOrderExecutedReConfirmedChan chan IpOrderMessage //"oce"
	ToCommExternalButtonPushedChan     chan IpOrderMessage //"ebp"
	ToCommImSlaveChan                  chan IpOrderMessage //"ias"
	ToCommUpdatedStateChan             chan IpState        //"ust"

}
type ExternalMasterChannels struct {
	ToCommOrderListChan              chan map[*UDPAddr]*[N_FLOORS][2]bool //"exo"
	ToCommOrderExecutedConfirmedChan chan IpOrderMessage                  //"eco"
	ToCommImMasterChan               chan string                          //"iam"
	ToCommUpdateStateReceivedChan    chan IpState                         // "sus"

}
type ExternalStateMachineChannels struct {
	ButtonPressed      chan IpOrderMessage //State machine needs to keep track of internal and external
	OrderServed        chan Order
	CurrentState       chan IpState
	GetSlaveStruct     chan bool
	ReturnSlaveStruct  chan Slave
	DirectionUpdate    chan int
	SingleExternalList chan [N_FLOORS][2]bool
	LightChan          chan [N_FLOORS][2]bool
}

func Channels_init() {
	network_external_chan_init()
	Communication_external_channels_init()
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

func Communication_external_channels_init() {
	ExCommChans.ToSlaveOrderListChan = make(chan IpOrderList)                    //"ord"
	ExCommChans.ToMasterOrderListReceivedChan = make(chan IpOrderList)           //"ore"
	ExCommChans.ToMasterOrderExecutedChan = make(chan IpOrderMessage)            //"oex"
	ExCommChans.ToSlaveOrderExecutedConfirmedChan = make(chan IpOrderMessage)    //"eco"
	ExCommChans.ToMasterOrderExecutedReConfirmedChan = make(chan IpOrderMessage) //"oce"
	ExCommChans.ToMasterExternalButtonPushedChan = make(chan IpOrderMessage)     //"ebp"
	ExCommChans.ToMasterSlaveChan = make(chan IpSlave)                           //"sla"
	ExCommChans.ToSlaveImMasterChan = make(chan string)                          //"iam"
	ExCommChans.ToMasterImSlaveChan = make(chan IpOrderMessage)                  //"ims"
	ExCommChans.ToMasterUpdateState = make(chan IpState)
	ExCommChans.ToSlaveButtonPressedConfirmedChan = make(chan IpOrderMessage)
}

func Slave_external_chans_init() {
	ExSlaveChans.ToCommOrderListReceivedChan = make(chan IpOrderList) //"ore"
	//ExSlaveChans.ToCommOrderReceivedChan = make(chan ipOrderMessage)            //"oce"
	ExSlaveChans.ToCommOrderExecutedChan = make(chan IpOrderMessage)            //"oex"
	ExSlaveChans.ToCommOrderExecutedReConfirmedChan = make(chan IpOrderMessage) //"oce"
	ExSlaveChans.ToCommExternalButtonPushedChan = make(chan IpOrderMessage)     //"ebp"
	ExSlaveChans.ToCommSlaveChan = make(chan Slave)                             //"sla"
	ExSlaveChans.ToCommImSlaveChan = make(chan IpOrderMessage)
	ExSlaveChans.ToCommUpdatedStateChan = make(chan IpState) //"ust" //CHECK USAGE IN LINE 60 IN FUNCSUGG

}
func Master_external_chans_init() {
	ExMasterChans.ToCommOrderListChan = make(chan map[*UDPAddr]*[N_FLOORS][2]bool) //"ord"
	ExMasterChans.ToCommOrderExecutedConfirmedChan = make(chan IpOrderMessage)     //"eco"
	ExMasterChans.ToCommImMasterChan = make(chan string)
}

func External_state_machine_channels_init() {
	ExStateMChans.ButtonPressed = make(chan IpOrderMessage)
	ExStateMChans.OrderServed = make(chan Order)
	ExStateMChans.CurrentState = make(chan IpState)
	ExStateMChans.GetSlaveStruct = make(chan bool)
	ExStateMChans.ReturnSlaveStruct = make(chan Slave)
	ExStateMChans.DirectionUpdate = make(chan int)
	ExStateMChans.SingleExternalList = make(chan [N_FLOORS][2]bool)
	ExStateMChans.LightChan = make(chan [N_FLOORS][2]bool)
}
func External_optimization_channel_init() {
	ExOptimalChans.OptimizationTriggerChan = make(chan IpOrderMessage)
	ExOptimalChans.OptimizationReturnChan = make(chan IpOrderMessage)
}

//MEMBER FUNCTIONS
func (m Master) Set_external_list_order(ip *UDPAddr, floor int, buttonType int, ipOrder IpOrderMessage) {
	m.ExternalList[ip][floor][buttonType] = ipOrder.Order.TurnOn
}
func (m Master) Get_external_list() map[*UDPAddr]*[N_FLOORS][2]bool {
	return m.ExternalList
}
func (s Slave) Overwrite_external_list(newExternalList map[*UDPAddr]*[N_FLOORS][2]bool) {
	s.ExternalList = newExternalList
}
func (s Slave) Get_ip() *UDPAddr {
	return s.IP
}
