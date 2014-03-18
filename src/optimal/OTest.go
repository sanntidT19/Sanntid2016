package optimal

import ( 
	"fmt"
	"os"
	"os/exec"
)

const (
	FLOORS = 4
	ELEV = 3 // It will be counted from recieved message
	TIMEOFSTOP = 1
	TIMETRAVEL = 2
)

type Cabin struct {
	Name int
	Position int // current Position
	Direction bool // true if elevator is going up
	Buttons[]bool // list of current orders from cabin
	Timeload int // current load of elevator
	LowestFloor int
	HighestFloor int
	Stops int
	HasIntOrders bool
}

// C is array of orders from cabin
var C[ELEV] Cabin
// H are orders from hallway

// Yet to serve - mirror of hallway, but will clear until every order from hallway is issued - I guess we should'n delete values in H
var H [2][FLOORS] bool
var Y [2][FLOORS] bool
var Outotal [2][FLOORS] bool
var TimeArr[ELEV] int

func main() {
	// clears the command line
	cmd := exec.Command("cmd", "/c", "cls")
    cmd.Stdout = os.Stdout
    cmd.Run()
	// starting conditions of elevators
	dire := []bool{true, false, true}
	posit := []int {3,2,2}
	
	// Input from inside cabins
	InBut := []bool{false, true, false, false,
			false, false, false, false,
			false, false, false, true,}
	
	//input from hallway - top and bottom floor can only go one direction
	OuUp := [FLOORS]bool{true, true, false, false}
	OuDo := [FLOORS]bool{false, true, false, true}		

	Outotal[0] = OuUp
	Outotal[1] = OuDo
	

	//input from hallway into struct	
	H = Outotal
	Y = H

	// For each elevator do:
	for i := 0; i < ELEV; i++ {
		
		C[i] = Cabin{i, posit[i], dire[i], InBut[i*FLOORS:i*FLOORS+FLOORS],0,0,0,0,false} // fill recieved data to C array
		
		// Find Highest floor, Lowest floor, Number of stops, and if it has any orders.
		C[i].HighestFloor, C[i].LowestFloor, C[i].Stops, C[i].HasIntOrders = FindHighLow(C[i].Buttons[0:FLOORS])

		//Check if elevator needs to change direction
    	C[i].Direction = DirectionSwitch(C[i].Direction, C[i].HighestFloor, C[i].LowestFloor, C[i].Position)
    	C[i].Timeload = Counttimeload(C[i].HasIntOrders, C[i].Direction, C[i].HighestFloor, C[i].LowestFloor, C[i].Position, C[i].Stops)	
    	fmt.Println("External list",H)
		H = (Allready(H, C[i].Buttons[0:FLOORS], C[i].Direction, C[i].Position))
		fmt.Println("External list",H)		
		// get rid of floors inside of elevator movements
		//
		//

	TimeArr[i] = C[i].Timeload


	// Just for testing
	fmt.Println("Name = ",C[i].Name)
	fmt.Println("Position",i, " = ",C[i].Position)
	fmt.Println("Direction",i, " = ",C[i].Direction)
	fmt.Println("Buttons",i, " = ",C[i].Buttons[0:4])
	fmt.Println("HighestFloor",i, " = ", C[i].HighestFloor)
	fmt.Println("LowestFloor",i, " = ", C[i].LowestFloor)
	fmt.Println("Stops",i, " = ", C[i].Stops)
	fmt.Println("Timeload",i, " = ", C[i].Timeload)
	fmt.Println()
	}
	// finds the number of elevator with least load
	Leastload := Leastloadfu(TimeArr[0:ELEV])	
	fmt.Println("Number of Leastload",Leastload)
	// = InsideMovement(C[Leastload] Cabin)
}

//selects the elevator with least load	
func Leastloadfu(P[] int) int {				
		Chosen := 0
		Minimum := P[0]
    	for i := 0; i < ELEV; i++ {
        	if P[i] < Minimum{
           		Minimum = P[i]
         		Chosen = i
			}	
		}
	return Chosen
}

// if direction is up and nothing is waiting there - switch direction down
func DirectionSwitch(Di bool, Hi int, Lo int, Po int) bool{
		V := Di
		if Di && Hi <= Po{
			V = false
			fmt.Println("switch down")
		}
		
		if !Di && Lo >= Po{
			V = true
			fmt.Println("switch UP")
		}
		
	return V
}
		// problem is when only button with actual possition is turned on - this should be solved elsewhere. Probably just run function open doors and unactivate.
		// Now it is considered that the button was activated just after the cabin left current floor.




// count timeload 
func Counttimeload(Ha bool, Di bool, Hi int, Lo int, Po int, St int) int{
		var Ti int
		if Ha {
			if Di {
				Ti = ((2 * Hi - Po - Lo) * TIMETRAVEL + St * TIMEOFSTOP)
			} else {
				Ti = ((Po + Hi - 2 * Lo) * TIMETRAVEL + St * TIMEOFSTOP)
			}
		} else {
				Ti = 0
			}
				return Ti
}

// get rid of floors alredy handled
func Allready(P[2][FLOORS] bool, Bu[] bool, Di bool, Po int) [2][FLOORS]bool{
		for n := 0; n < 2; n++ {
			for j := 0; j < FLOORS; j++ {
			
        		if P[n][j] && Bu[j] && Po != n * (FLOORS-1) {
					P[n][j] = false
				}

			}
		}
		return P	
}
// Finds Highest and Lowest number, Counts number of stops and says if elevator has internal orders at all
func FindHighLow(Bu[] bool) (int, int, int, bool){
		counter := false // counts if lowest level was allready found
		Ha := false
		St := 0
		Hi := 0
		Lo := 0
		for j := 0; j < FLOORS; j++ {
			// find highest odered floor
        	if Bu[j] {
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
			if Bu[j] {
				Ha = true
			}
    	}
    return Hi, Lo, St, Ha
}
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
