	package optimal

	import (
	."chansnstructs"
	"fmt"
	)

	const (
		
	)

	type Cabin struct {
		Name         int
		Position     int    // current Position
		Direction    bool   // true if elevator is going up
		Buttons      []bool // list of current orders from cabin
		Timeload     int    // current load of elevator
		LowestFloor  int
		HighestFloor int
		Stops        int
		HasIntOrders bool
	}

	

	func Optimization(){

		

		var Input Master


		var ELEV int
		
		ELEV = len(Input.SlaveElev Type) int //??????

		//Y Matrix of order from hallway
		var Y [2][N_FLOORS]bool

		// E is output matrix sent to channel
		var E [ELEV][2][N_FLOORS]bool

		// T is array of current time load for each cabin
		var T [ELEV]int

		// C is array of parameters for each cabin
		var C [ELEV]Cabin



		for{

			// input from master - blocking function until something is there
			Input <- ExOptimalChans.OptimizationTriggerChan
			// Y is filled with orders from hallway
			Y = Input.SlaveElev[0].AllExternalsOrder

			//filling of cabins for each elevator
			for i := 0; i < ELEV; i++ {
				// fill recieved data to C array
				C[i] = Cabin{Input.SlaveElev[i].IP, Input.SlaveElev[i].CurrentFloor, Input.SlaveElev[i].Direction, Input.SlaveElev[i].InternalList[0:ELEV], 0, 0, 0, 0, false} 

				//Check if elevator needs to change direction - it happens if there is no order in original direction
				C[i].Direction = DirectionSwitch(C[i].Direction, C[i].HighestFloor, C[i].LowestFloor, C[i].Position)
				// Create bases of output matrices for each elevator
				E[i] = Transform(C[i].Buttons[0:N_FLOORS], C[i].Direction, C[i].Position)

				// Gets rid of duplicitous N_FLOORS allready served by cabins
				Y = (Already(Y, C[i].Buttons[0:N_FLOORS], C[i].Direction, C[i].Position))

			}

			// giving orders to elevator with current least load
			for j := 0; j < N_FLOORS; j++ {
				for n := 0; n < 2; n++ {
					if Y[n][j] {
						E, T = loop(C[0:ELEV], Y, E)
						E[Leastloadfu(T[0:ELEV])][n][j] = true
					}
				}
			}

			ExOptimalChans.OptimizationReturnChan <- E
		}
	}

	func loop(C []Cabin, H [2][N_FLOORS]bool, ExecutionList [ELEV][2][N_FLOORS]bool) ([ELEV][2][N_FLOORS]bool, [ELEV]int) {
		var TimeArr [ELEV]int
		// For each elevator do:
		for i := 0; i < ELEV; i++ {

			// Find Highest floor, Lowest floor, Number of stops, and if it has any orders.
			C[i].HighestFloor, C[i].LowestFloor, C[i].Stops, C[i].HasIntOrders = FindHighLow(ExecutionList[i])

			C[i].Timeload = Counttimeload(C[i].HasIntOrders, C[i].Direction, C[i].HighestFloor, C[i].LowestFloor, C[i].Position, C[i].Stops)

			// creates array of timeloads for
			TimeArr[i] = C[i].Timeload
			}

			return ExecutionList, TimeArr
		}

	//selects the elevator with least load
		func Leastloadfu(P []int) int {
			Chosen := 0
			Minimum := P[0]
			for i := 0; i < ELEV; i++ {
				if P[i] < Minimum {
					Minimum = P[i]
					Chosen = i
				}
			}
			return Chosen
		}

	// if direction is up and nothing is waiting there - switch direction down
		func DirectionSwitch(Di bool, Hi int, Lo int, Po int) bool {
			V := Di
			if Di && Hi <= Po {
				V = false
			}

			if !Di && Lo >= Po {
				V = true
			}

			return V
		}

	// count timeload
		func Counttimeload(Ha bool, Di bool, Hi int, Lo int, Po int, St int) int {
			var Ti int
			if Ha {
				if Di {
					Ti = ((2*Hi-Po-Lo)*TIMETRAVEL + St*TIMEOFSTOP)
				} else {
					Ti = ((Po+Hi-2*Lo)*TIMETRAVEL + St*TIMEOFSTOP)
				}
			} else {
				Ti = 0
			}
			return Ti
		}

	// get rid of floors alredy handled by internal orders
		func Already(P [2][N_FLOORS]bool, Bu []bool, Di bool, Po int) [2][N_FLOORS]bool {
			for n := 0; n < 2; n++ {
				for j := 0; j < N_FLOORS; j++ {

					if P[n][j] && Bu[j] && Po != n*(N_FLOORS-1) {
						P[n][j] = false
					}

				}
			}
			return P
		}

	// Finds Highest and Lowest number, Counts number of stops
	// and says if elevator has internal orders at all
		func FindHighLow(Ex [2][N_FLOORS]bool) (int, int, int, bool) {
		counter := false // counts if lowest level was allready found
		Ha := false // true cabin has some inputs
		St := 0 // number of stops
		Hi := 0 // highest floor
		Lo := 0 // lowest floor

		for j := 0; j < N_FLOORS; j++ {
			// find highest odered floor 
			if Ex[0][j] || Ex[1][j] { //doesn't matter which direction
				Hi = j 

				// find lowest odered floor
				if !counter {
					Lo = j
					counter = true
				}

				// counts number of stops
				St++

			}
			// false if there are no inside orders
			if Ex[0][j] || Ex[1][j] {
				Ha = true
			}
		}
		return Hi, Lo, St, Ha
	}

	// Transforms cabin buttons [array of bool for each floor]
	// into execution list [array of two bools for each floor] - distinguished up and down directions
	func Transform(Bu []bool, Di bool, Po int) [2][N_FLOORS]bool {
		var P [2][N_FLOORS]bool
		for j := 0; j < N_FLOORS; j++ {
			if Bu[j] {
				// Solves if elevator should handle floor on the way up or down.
				switch {
				case Po < j:
					P[0][j] = true
				case Po > j:
					P[1][j] = true
				case Po == j:
					// if button is pressed when cabin is in the same floor, it is considered as button was pressed after elevator left the floor
					if Di {
						P[1][j] = true
					} else {
						P[0][j] = true
					}
				}
			}
		}
		return P
	}