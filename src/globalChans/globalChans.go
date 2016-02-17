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
