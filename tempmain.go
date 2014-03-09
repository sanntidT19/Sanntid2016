package main 

import 
(
	"driver"
	"errors"
	"fmt"
)


func main() {
	err = driver.elev_init()
	if err != nil{
		fmt.Println(err)
	}
}
