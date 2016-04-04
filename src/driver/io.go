package driver

// #cgo LDFLAGS: -lcomedi -lm
// #include <png.h>
// #include "io.h"
import "C"

//C.io_init returns =! 0 if init has been successfull, io_init returns true for the same case
func IoInit() bool {
	return bool(int(C.io_init()) != 0)
}

func IoSetBit(channel int) {
	C.io_set_bit(C.int(channel))
}

func IoClearBit(channel int) {
	C.io_clear_bit(C.int(channel))
}

func IoWriteAnalog(channel int, value int) {
	C.io_write_analog(C.int(channel), C.int(value))
}

//have changed this to return boolean values

func IoReadBit(channel int) bool {
	var i int = int(C.io_read_bit(C.int(channel)))
	if i == 1 {
		return true
	} else {
		return false
	}
}

func IoReadAnalog(channel int) int {
	return int(C.io_read_analog(C.int(channel)))
}
