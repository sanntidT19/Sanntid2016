package globalChans

import (
	"../globalStructs"
)

var ButtonPressedChan chan globalStructs.Button
var SetButtonLightChan chan globalStructs.Button

func Init_chans() {
	ButtonPressedChan = make(chan globalStructs.Button)
	SetButtonLightChan = make(chan globalStructs.Button)
}
/*
global chans som må fikses:
Externalbuttonpressedchan(driver og network)
Internalbuttonpressedchan(driver og toplevel)



ResetExternalOrdersInQueueChan (statemachine og toplevel)
/*