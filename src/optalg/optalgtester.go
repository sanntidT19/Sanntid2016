package main 

type ElevatorState struct {
	MyIP         string //not an int, always useful to have
	CurrentFloor int
	LastFloor    int
	Direction    int
	Orders       []Button //This is an array or something of all orders currently active for this elevator.

}


type Button struct {
	Floor          int
	Button_type    int
	Button_pressed bool
}


var el_state1 ElevatorState = ElevatorState{MyIP: "123.123.123.123",
			CurrentFloor: 2,
			LastFloor: 3,
			Direction: -1}
var el_state2 ElevatorState = ElevatorState{MyIP: "123.123.123.124",
			CurrentFloor: 2,
			LastFloor: 3,
			Direction: -1,
			Orders: []Button{Button{3, -1, false}, {2, 0, false}}}
			
var el_state3 ElevatorState =ElevatorState{MyIP: "123.123.123.125",
			CurrentFloor: 2,
			LastFloor: 3,
			Direction: -1,
			Orders: []Button{Button{2, -1, false}, Button{1, 1, false}}}
			
var all_elevs  = []ElevatorState{el_state1, el_state2, el_state3}
var new_order Button = Button{Floor: 0,
		Button_type: -1,
		Button_pressed: false}

func opt_alg(new_order Button) string {
	numOfElevs := len(all_elevs)
	var IP_cost_list [numOfElevs]int;
	var IP_score string = "100";
	var Optimal_IP string = "0"; 
	for i_v := range all_elevs{
		IP_cost_list[i] += len(v.Orders);
		if v.Direction != new_order.Button_type{
			IP_cost_list[i] += 1;
			}
		IP_cost_list[i] += abs(v.CurrentFloor - new_order.Floor);
		}	
	for k := 1; k < len(IP_cost_list); k+=1{
		if IP_cost_list[k] < IP_score{
			Optimal_IP = all_elevs[k].MyIP;
		}else if IP_cost_list[k] == IP_score{
			if all_elevs[k].MyIP > Optimal_IP{
				Optimal_IP := all_elevs[k].MyIP;
				}
			}
	} 		
}

func main(){
	optalgtest.opt_alg(new_order); 
}



















