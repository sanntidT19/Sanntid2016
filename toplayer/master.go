package modules

import (
	"fmt"
)

const ()

type Master struct {

	externalList  [][]int
	currentFloors []int
	directions    []int
}

func (m Master) getOptimalExternalList(optimalityChan chan [][]int) {
	//gets new OptimnalExternalList from module
}

func (m Master) sendExternalListToSlaves(masterToCommunicationChan chan ) {
	//sends new external list to communication module 
}