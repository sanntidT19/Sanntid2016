package globalChans

import(
	"../globalStructs"
)

func Initalize_chans(){
	buttonPressedChan := make(chan globalStructs.Button)
	setButtonLightChan := make (chan globalStructs.Button)
}


