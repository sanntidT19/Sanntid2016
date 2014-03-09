package driver // where "driver" is the folder that contains io.go, io.c, io.h, channels.go, channels.c and driver.go
/*
//package main //I think this should be here!!!!!!!!!!!!!!!!!!!!!!!-Olav

/*
Translated to GO
I am not totaly sure if the matrices is equal to C but the are on the correct form, maybe they need to be transposed(?)
I left the C-code commented out after the go functions


Yngve: I think slices is the idiomatic way to go in go.

-Olav
*/



/*
#cgo LDFLAGS: -lcomedi -lm
#include "io.h"
#include "io.c"
*/
import (
"C"
"errors"
)



// Wrapper for libComedi Elevator control.
// These functions provides an interface to the elevators in the real time lab
//
// 2007, Martin Korsgaard


//#include "channels.h"
//#include "elev.h"
//#include "io.h"

//#include <assert.h>
//#include <stdlib.h>

// Number of signals and lamps on a per-floor basis (excl sensor)
const N_BUTTONS = 3
var N_FLOORS := 4 //is this required here? should this be all caps?


//constants making it easier to read the code, will prob be used
const UP = 0
const DOWN = 1
const COMMAND = 2

var lamp_channel_matrix [][] int := nil
var button_channel_matrix [][] int := nil

//Making the standard matrix for our project
func elev_make_std_l_matrix() [][]int {  //stupid name? we need to agree on a stanard name convention
	std_matrix := make([][] int, 4,8)
	std_matrix[0] = make([]int,N_BUTTONS)
	std_matrix[0] = []int{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1}
	std_matrix[1] = make([]int,N_BUTTONS)
	std_matrix[1] = []int{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2}
	std_matrix[2] = make([]int,N_BUTTONS)
	std_matrix[2] = []int{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3}
	std_matrix[3] = make([]int,N_BUTTONS)
	std_matrix[3] = []int{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4}
	return std_matrix
	}		
}

func elev_make_std_b_matrix() [][] int{
	std_matrix := make([][] int, 4,8) //Find out if the capacity here is valid
	std_matrix[0] = make([]int,N_BUTTONS)
	std_matrix[0] = []int{FLOOR_UP1, FLOOR_DOWN1, FLOOR_COMMAND1}
	std_matrix[1] = make([]int,N_BUTTONS)
	std_matrix[1] = []int{FLOOR_UP2, FLOOR_DOWN2, FLOOR_COMMAND2}
	std_matrix[2] = make([]int,N_BUTTONS)
	std_matrix[2] = []int{FLOOR_UP3, FLOOR_DOWN3, FLOOR_COMMAND3}
	std_matrix[3] = make([]int,N_BUTTONS)
	std_matrix[3] = []int{FLOOR_UP4, FLOOR_DOWN4, FLOOR_COMMAND4}
	return std_matrix
}

//Possibility for extension of floors, but not used in this project. capacity of the matrix is extended if necessary
//It is assumed that the added "floor" is either on the top or the bottom of the elevator-shaft.
func elev_extend_matrix(matrix [][] int,light_up int ,light_down int,light_command int, int floor) ([][] int, int) {
	if len(matrix) == cap(matrix) {
		temp_matrix := make([][]int, len(matrix), (cap(matrix)+1)*2)
		copy(temp_matrix,matrix)
		matrix := temp_matrix
		}
		
	//insert "add floor"-code in here	
	}
	return new_matrix, floors + 1 // is N_FLOORS available here? 
}


/*
type lamp_channel_matrix := [N_FLOORS][N_BUTTONS]int{
[N_BUTTONS]int{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
[N_BUTTONS]int{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
[N_BUTTONS]int{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
[N_BUTTONS]int{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4}}

static const int lamp_channel_matrix[N_FLOORS][N_BUTTONS] = {
{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
};

button_channel_matrix := [N_FLOORS][N_BUTTONS]int

{[N_BUTTONS]int{FLOOR_UP1, FLOOR_DOWN1, FLOOR_COMMAND1},
[N_BUTTONS]int{FLOOR_UP2, FLOOR_DOWN2, FLOOR_COMMAND2},
[N_BUTTONS]int{FLOOR_UP3, FLOOR_DOWN3, FLOOR_COMMAND3},
[N_BUTTONS]int{FLOOR_UP4, FLOOR_DOWN4, FLOOR_COMMAND4}}


static const int button_channel_matrix[N_FLOORS][N_BUTTONS] = {
{FLOOR_UP1, FLOOR_DOWN1, FLOOR_COMMAND1},
{FLOOR_UP2, FLOOR_DOWN2, FLOOR_COMMAND2},
{FLOOR_UP3, FLOOR_DOWN3, FLOOR_COMMAND3},
{FLOOR_UP4, FLOOR_DOWN4, FLOOR_COMMAND4},
};
*/


//We havent used goroutines nor channels here, since this is a one time event for the life of an elevator
func elev_init() error {
//Initializes hardware
//(?) correct function call?
	if !io_init()) {
	return errors.New("IO initialization failed")
	}
	lamp_channel_matrix = elev_make_std_l_matrix()
	button_channel_matrix = elev_make_std_b_matrix()
	//Zero all floor button lamps
	for i:=0; i<N_FLOORS;++i {
		if i != 0 {
			elev_set_button_lamp(DOWN, i , 0)
		}
		if i != N_FLOORS-1 {
			elev_set_button_lamp(UP, i, 0)
		}
		elev_set_button_lamp(COMMAND, i, 0)
	}
	//Clear stop lamp, foor open lamp, and set floor indicatior and ground floor.
	elev_set_stop_lamp(0)
	elev_set_door_open_lamp(0)
	elev_set_floor_indicator(0)
	return nil
}
/*
int elev_init(void){
// Init hardware
if (!io_init())
return 0;

// Zero all floor button lamps
for (int i = 0; i < N_FLOORS; ++i) {
if (i != 0)
elev_set_button_lamp(BUTTON_CALL_DOWN, i, 0);

if (i != N_FLOORS-1)
elev_set_button_lamp(BUTTON_CALL_UP, i, 0);

elev_set_button_lamp(BUTTON_COMMAND, i, 0);
}

// Clear stop lamp, door open lamp, and set floor indicator to ground floor.
elev_set_stop_lamp(0);
elev_set_door_open_lamp(0);
elev_set_floor_indicator(0);

// Return success.
return 1;
}

*/

 // This looks good, isnt +-300 the desired speeds?
func elev_set_speed(speed int) {
//In order to sharply stop the elevator, the direction bit is toggled, before setting speed to zero
	last_speed := 0
//if to start (speed > 0)
	if speed > 0 {
		io_clear_bit(MOTORDIR)
	} else if speed < 0 {
		io_set_bit(MOTORDIR)  // if to stop (speed == 0)
	} else if last_speed < 0 { 
		io_clear_bit(MOTORDIR)
	} else if last_speed > 0 {
		io_set_bit(MOTORDIR)
	}
	last_speed = speed

	//Write new setting to motor
	io_write_analog(MOTOR, 2048 + 4*abs(speed))

}
/*
void elev_set_speed(int speed){
// In order to sharply stop the elevator, the direction bit is toggled
// before setting speed to zero.
static int last_speed = 0;
// If to start (speed > 0)
if (speed > 0)
io_clear_bit(MOTORDIR);
else if (speed < 0)
io_set_bit(MOTORDIR);

// If to stop (speed == 0)
else if (last_speed < 0)
io_clear_bit(MOTORDIR);
else if (last_speed > 0)
io_set_bit(MOTORDIR);

last_speed = speed ;

// Write new setting to motor.
io_write_analog(MOTOR, 2048 + 4*abs(speed));
}
*/

    // This looks good. I suggest we find a smart use of channels to read and write light-bits
func elev_get_floor_sensor_signal(){
	if io_read_bit(SENSOR1) {
		return 0
	} else if io_read_bit(SENSOR2) {
		return 1
	} else if io_read_bit(SENSOR3) {
		return 2
	} else if io_read_bit(SENSOR4) {
		return 3
	}
//by convention in go no else is used, the function will not continue if one of the previous returns is called(?)  Yngve: I think else is used plenty
	return -1 
}
/*

int elev_get_floor_sensor_signal(void){
if (io_read_bit(SENSOR1))
return 0;
else if (io_read_bit(SENSOR2))
return 1;
else if (io_read_bit(SENSOR3))
return 2;
else if (io_read_bit(SENSOR4))
return 3;
else
return -1;
}
*/

   //This will return -1 if something fails  Looks good
func elev_get_button_signal(button int, floor int) (int, error) {
	if floor < 0 || floor >= N_FLOORS {
		return -1, errors.New("Floor value not in valid region")
	}
	if button < 0 || button >= N_BUTTONS {
		return -1, errors.New("Button value not in valid region") 
	}
	if button == UP && floor ==  N_FLOORS -1{
		return -1, errors.New("This floor has no defined up button")
	}
	if button == DOWN && floor ==  0{
		return -1, errors.New("This floor has no defined down button")
	}
	if io_read_bit(button_channel_matrix[floor][button]) {
		return 1,nil
	}
	return 0,nil
}
/*

int elev_get_button_signal(elev_button_type_t button, int floor){
assert(floor >= 0);
assert(floor < N_FLOORS);
assert(!(button == BUTTON_CALL_UP && floor == N_FLOORS-1));
assert(!(button == BUTTON_CALL_DOWN && floor == 0));
assert( button == BUTTON_CALL_UP ||
button == BUTTON_CALL_DOWN ||
button == BUTTON_COMMAND);

if (io_read_bit(button_channel_matrix[floor][button]))
return 1;
else
return 0;
}
*/



	// This looks good
func elev_set_floor_indicator(floor int) error {	
	if floor < 0 || floor >= N_FLOORS {
	return errors.New("Floor value not in valid region")
	}
    // Binary encoding. One light must always be on
    if floor && 0x02 {
     io_set_bit(FLOOR_IND1)
    } else {
     io_clear_bit(FLOOR_IND1)
    }
    if floor && 0x01 {
     io_set_bit(FLOOR_IND2)
    } else {
     io_clear_bit(FLOOR_IND2)
    }
    return nil
}

/*
void elev_set_floor_indicator(int floor){
assert(floor >= 0);
assert(floor < N_FLOORS);

// Binary encoding. One light must always be on.
if (floor & 0x02)
io_set_bit(FLOOR_IND1);
else
io_clear_bit(FLOOR_IND1);
if (floor & 0x01)
io_set_bit(FLOOR_IND2);
else
io_clear_bit(FLOOR_IND2);
}
*/

// There is no assert in go, Looks good
func elev_set_button_lamp(button int, floor int, value int) error {
	
	if floor < 0 || floor >= N_FLOORS {
		return errors.New("Floor value not in valid region")
	}
	if button < 0 || button >= N_BUTTONS {
		return errors.New("Button value not in valid region") 
	}
	if button == UP && floor ==  N_FLOORS -1{
		return errors.New("This floor has no defined up button")
	}
	if button == DOWN && floor ==  0{
		return errors.New("This floor has no defined down button")
	}
    if value == 1 {
    	io_set_bit(lamp_channel_matrix[floor][button])
    } else {
    	io_clear_bit(lamp_channel_matrix[floor][button])
    }
    return nil
}


	//This looks good 
func elev_set_door_open_lamp(value int) {
	if value {
		io_set_bit(DOOR_OPEN)
	}
	else {
		io_clear_bit(DOOR_OPEN)
	}
}
/*
void elev_set_door_open_lamp(int value){
if (value)
io_set_bit(DOOR_OPEN);
else
io_clear_bit(DOOR_OPEN);
}
*/



/*
void elev_set_button_lamp(elev_button_type_t button, int floor, int value){
assert(floor >= 0);
assert(floor < N_FLOORS);
assert(!(button == BUTTON_CALL_UP && floor == N_FLOORS-1));
assert(!(button == BUTTON_CALL_DOWN && floor == 0));
assert( button == BUTTON_CALL_UP ||
button == BUTTON_CALL_DOWN ||
button == BUTTON_COMMAND);

if (value == 1)
io_set_bit(lamp_channel_matrix[floor][button]);
else
io_clear_bit(lamp_channel_matrix[floor][button]);
}
*/



















// Not needed
func elev_set_stop_lamp(value int) {
if value {
io_set_bit(LIGHT_STOP)
}
else {
io_clear_bit(LIGHT_STOP)
}
}
/*
void elev_set_stop_lamp(int value){
if (value)
io_set_bit(LIGHT_STOP);
else
io_clear_bit(LIGHT_STOP);
}
*/

// Dont need this shit.
func elev_get_stop_signal() {
return io_read_bit(STOP)
}
/*
int elev_get_stop_signal(void){
return io_read_bit(STOP);
}
*/

// Dont need this shit?
func elev_get_obstruction_signal() {
return io_read_bit(OBSTRUCTION)
}
/*
int elev_get_obstruction_signal(void){
return io_read_bit(OBSTRUCTION);
}
*/



