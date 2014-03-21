package optimal

import (
	. "chansnstructs"
	"fmt"
)

const (
	FLOORS     = 4
	ELEV       = 3 // It will be counted from recieved message
	TIMEOFSTOP = 1
	TIMETRAVEL = 2
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
m.s[]

func Optimization(X [ELEV]Slave) [ELEV][2][FLOORS]bool {
	// input in format:
	/*
	   	type Slave struct {
	   	Nr int
	   	InternalList [FLOORS]bool
	   	ExternalList [2][FLOORS]bool
	   	CurrentFloor int
	   	Direction    bool
	   }
	*/

	// H are orders from hallway
	//Y is mirror of hallway
	var Y [2][FLOORS]bool
	// C is array of parameters for each cabin
	var E [ELEV][2][FLOORS]bool
	var T [ELEV]int
	var C [ELEV]Cabin

	// make mirror of hallway - we will not overwrite hallway
	Y = X[0].ExternalList
	fmt.Println("External list", Y)
	fmt.Println()
	//filling of cabins
	for i := 0; i < ELEV; i++ {

		C[i] = Cabin{X[i].Nr, X[i].CurrentFloor, X[i].Direction, X[i].InternalList[0:ELEV], 0, 0, 0, 0, false} // fill recieved data to C array

		//Check if elevator needs to change direction - it happens if there is no order in original direction
		C[i].Direction = DirectionSwitch(C[i].Direction, C[i].HighestFloor, C[i].LowestFloor, C[i].Position)

		E[i] = Transform(C[i].Buttons[0:FLOORS], C[i].Direction, C[i].Position)

		// Gets rid of duplicitous floors allready served by cabins
		Y = (Allready(Y, C[i].Buttons[0:FLOORS], C[i].Direction, C[i].Position))

	}

	// giving orders to elevator with current least load
	for j := 0; j < FLOORS; j++ {
		for n := 0; n < 2; n++ {
			if Y[n][j] {
				E, T = loop(C[0:ELEV], Y, E)
				E[Leastloadfu(T[0:ELEV])][n][j] = true

			}
		}
	}

	fmt.Println("Toto je E0", E[0])
	fmt.Println("Toto je E1", E[1])
	fmt.Println("Toto je E2", E[2])

	fmt.Println()

	return E

}

func loop(C []Cabin, H [2][FLOORS]bool, ExecutionList [ELEV][2][FLOORS]bool) ([ELEV][2][FLOORS]bool, [ELEV]int) {
	var TimeArr [ELEV]int
	// For each elevator do:
	for i := 0; i < ELEV; i++ {

		// Find Highest floor, Lowest floor, Number of stops, and if it has any orders.
		C[i].HighestFloor, C[i].LowestFloor, C[i].Stops, C[i].HasIntOrders = FindHighLow(ExecutionList[i])

		C[i].Timeload = Counttimeload(C[i].HasIntOrders, C[i].Direction, C[i].HighestFloor, C[i].LowestFloor, C[i].Position, C[i].Stops)

		// creates array of timeloads for
		TimeArr[i] = C[i].Timeload

		// Just for testing
		/*	fmt.Println("Name = ",C[i].Name)
			fmt.Println("Position",i, " = ",C[i].Position)
			fmt.Println("Direction",i, " = ",C[i].Direction)
			fmt.Println("Buttons",i, " = ",C[i].Buttons[0:4])
			fmt.Println("HighestFloor",i, " = ", C[i].HighestFloor)
			fmt.Println("LowestFloor",i, " = ", C[i].LowestFloor)
			fmt.Println("Stops",i, " = ", C[i].Stops)
			fmt.Println("Timeload",i, " = ", C[i].Timeload)
			fmt.Println()
		*/
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

// problem is when only button with actual possition is turned on - this should be solved elsewhere. Probably just run function open doors and unactivate.
// Now it is considered that the button was activated just after the cabin left current floor.

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

// get rid of floors alredy handled
func Allready(P [2][FLOORS]bool, Bu []bool, Di bool, Po int) [2][FLOORS]bool {
	for n := 0; n < 2; n++ {
		for j := 0; j < FLOORS; j++ {

			if P[n][j] && Bu[j] && Po != n*(FLOORS-1) {
				P[n][j] = false
			}

		}
	}
	return P
}

// Finds Highest and Lowest number, Counts number of stops and says if elevator has internal orders at all
func FindHighLow(Ex [2][FLOORS]bool) (int, int, int, bool) {
	counter := false // counts if lowest level was allready found
	Ha := false
	St := 0
	Hi := 0
	Lo := 0

	for j := 0; j < FLOORS; j++ {
		// find highest odered floor
		if Ex[0][j] || Ex[1][j] {
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

// Transforms cabin buttons into execution list
func Transform(Bu []bool, Di bool, Po int) [2][FLOORS]bool {
	var P [2][FLOORS]bool
	for j := 0; j < FLOORS; j++ {
		if Bu[j] {
			// Solves if elevator should handle floor on the way up or down.
			switch {
			case Po < j:
				P[0][j] = true
			case Po > j:
				P[1][j] = true
			case Po == j:
				// if button is pressed when cabin is in the same floor, it is considered as elevator already left
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

/*
func HallwayToExec(P[][][] bool,) [2][FLOORS]bool{

		return P
}
*/

/*
func InsideMovement(P[]) {

}

func BoolToFloorNum(P[][]) {
	for j := 0; j < FLOORS; j++ {
	if P[0][i]
}
*/
/*
// obsolete version
// For each elevator do:
	for i := 0; i < ELEV; i++ {
		C[i] = Cabin{i, posiTimeArr[i], dire[i], move[i], InUP[i*FLOORS:i*FLOORS+FLOORS],InDO[i*FLOORS:i*FLOORS+FLOORS],0,0,0,0,0,0}

	//For each floor do:
		for j := 0; j < FLOORS; j++ {
			// find highest odered floor
        		if C[i].ServiceUp[j] {
				C[i].HighestFloor = j
			}
			// find lowest ordered floor
			if C[i].ServiceDO[FLOORS - 1 - j] {
				C[i].LowestFloor = FLOORS - j
			}
			// count number of stops (both Directions)
			if C[i].ServiceUp[j] || C[i].ServiceDO[j]{
				C[i].Stops++
			}

    		}

		if C[i].Direction {
			C[i].Timeload = ((2 * C[i].HighestFloor - C[i].Position - C[i].LowestFloor) * TIMETRAVEL + C[i].Stops * TIMEOFSTOP)
		} else {
			C[i].Timeload = ((C[i].Position + C[i].HighestFloor - 2 * C[i].LowestFloor) * TIMETRAVEL + C[i].Stops * TIMEOFSTOP)
		}
*/

// What if only button with current flore is swithced on?

/*		if !C[i]Movement {

			if C[i].LowestFloor && C[i].HighestFloor == C[i].Position{
				C[i].Timeload = 0
			}else	{

				if C[i].Direction && C[i].HighestFloor <= C[i].Position{
				C[i].Direction = false
				}
				if !C[i].Direction && C[i].LowestFloor >= C[i].Position{
				C[i].Direction = true
				}
			}
		}

*/
