package main

import "fmt"

const (
	FLOORS = 4
	ELEV = 3
	TIMEOFSTOP = 3
	TIMETRAVEL = 4
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
}


// Hallway orders need to be split for each direction
type Hallway struct {
	HallUP []bool // list of current orders from hallway Upwards
	HallDO []bool // list of current orders from hallway Downwards
}	


// C is array of orders from cabin
var C[ELEV] Cabin
// H are orders from hallway
var H Hallway
// Yet to serve - mirror of hallway, but will clear until every order from hallway is issued - I guess we should'n delete values in H
var Y Hallway

var Chosen int
	var Minimum int

func main() {

	// starting conditions of elevators
	dire := []bool{true, false, false}
	move := []bool{true, true, true}
	posit := []int {0,2,2}
	
	// Input from inside cabins
	InBut := []bool{true, false, false, true,
			true, false, false, true,
			true, true, false, true,}
	
	//input from hallway - top and bottom floor can only go one direction

	OutUP := []bool{true, true, false, false}
	OutDO := []bool{false, true, false, true,}
	
	//input from hallway into struct	
	H := Hallway{OutUP[0:FLOORS],OutDO[0:FLOORS]}
	Y := H  // I am working with just reference and i want to copy H to Y instead

	// For each elevator do:
	for i := 0; i < ELEV; i++ {
		
		C[i] = Cabin{i, posit[i], dire[i], move[i], InBut[i*FLOORS:i*FLOORS+FLOORS],0,0,0,0} // fill recieved data to C array
		counter := false // counts if lowest level was allready found

		//For each floor do:
		for j := 0; j < FLOORS; j++ {
			// find highest odered floor
        		if C[i].Buttons[j] {
				C[i].HighestFloor = j

				// find lowest odered floor
				if !counter {
				C[i].LowestFloor = j
				counter = true
				}

				// counts number of stops
				C[i].Stops++ 			
		
			}
						
    		}
		
		// problem is when only button with actual possition is turned on - this should be solved elsewhere. Probably just run function open doors and unactivate.
		// Now it is considered that the button was activated just after the cabin left current floor.
		

		// if direction is up and nothing is waiting there - switch direction down
		if C[i].Direction && C[i].HighestFloor <= C[i].Position{
			C[i].Direction = false
		}
		// the same as above but for oposite direction
		if !C[i].Direction && C[i].LowestFloor >= C[i].Position{
			C[i].Direction = true
		}
		

		// count timeload
		if C[i].Direction {
			C[i].Timeload = ((2 * C[i].HighestFloor - C[i].Position - C[i].LowestFloor) * TIMETRAVEL + C[i].Stops * TIMEOFSTOP)
		} else {
				C[i].Timeload = ((C[i].Position + C[i].HighestFloor - 2 * C[i].LowestFloor) * TIMETRAVEL + C[i].Stops * TIMEOFSTOP)
			}
		

		// get rid of floors alredy handled
		for j := 0; j < FLOORS; j++ {
			// if going up
        		if Y.HallUP[j] && C[i].Buttons[j] && C[i].Direction && C[i].Position != 0 {
				Y.HallUP[j] = false
			}
			// if going down
			if Y.HallDO[j] && C[i].Buttons[j] && !C[i].Direction && C[i].Position != FLOORS-1 {
				Y.HallDO[j] = false
			}
		}	
		
				
		// get rid of floors inside of elevator movements
		//
		//


	// Just for testing
	fmt.Println("Name = ",C[i].Name)
	fmt.Println("Position",i, " = ",C[i].Position)
	fmt.Println(C[i].Direction)
	fmt.Println(C[i].Buttons[0:4])
	fmt.Println("HighestFloor",i, " = ", C[i].HighestFloor)
	fmt.Println("LowestFloor",i, " = ", C[i].LowestFloor)
	fmt.Println("Stops",i, " = ", C[i].Stops)
	fmt.Println("Timeload",i, " = ", C[i].Timeload)
	fmt.Println()
	fmt.Println(H.HallUP [0:4])
	fmt.Println(H.HallDO [0:4])
	fmt.Println()
	fmt.Println(Y.HallUP [0:4])
	fmt.Println(Y.HallDO [0:4])
	fmt.Println()
	fmt.Println()
	}

	//prints the number of elevator with the least load		
	fmt.Println("from leastload ",Leastload(C[0:ELEV]))	
	
}
//selects the elevator with least load	
func Leastload(P[] Cabin) int {			
   	
	

	Minimum = P[0].Timeload;
 
    	for i := 0; i < ELEV; i++ {
        	if P[i].Timeload < Minimum{
           		Minimum = P[i].Timeload
         		Chosen = i
		}	
	}
	return Chosen
}


/*
// obsolete version
// For each elevator do:
	for i := 0; i < ELEV; i++ {
		C[i] = Cabin{i, posit[i], dire[i], move[i], InUP[i*FLOORS:i*FLOORS+FLOORS],InDO[i*FLOORS:i*FLOORS+FLOORS],0,0,0,0,0,0}
		
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
