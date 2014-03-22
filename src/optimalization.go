package optimalization

import (
	"fmt"
)
const (

)

func Optimalization_init(m Master) {
	for {
		order := <-ExOptimalChans.OptimizationTriggerChan
		
		workloadMatrix := make([]int, N_ELEV)
		for i:=0; i<N_FLOORS; i++ {
			if m.S[IP].AllExternalOrders[IP].TurnOn == true {
				workLoadMatrix[i]++
			}
			if  m.S[IP].InternalList[i] == true {
				workLoadMatrix[i]++
			}
			if m.S[IP].Direction < (order.Floor-)
		}

		//Chooses elevator with lead workload
		NrMinWorkloadElevator := 0
		tempLoad :=	9999
		for i:=0; i<N_FLOORS; i++ {
			if tempLoad < workLoadMatrix[i]
				NrMinWorkloadElevator = i
		}
		return NrMinWorkloadElevator




	}
}