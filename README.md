Elevator-progress-
==================
--Notes for what needs to be done--

--The driver module: Need control of elevator using go-esque techniques
******MASTER********
*The behaviour of the master:
Construncts the externalOrderList with help from the optimalization algorithm(OA). 
also needs the functionality of the slave.

******SLAVE*******
*The behaviour of oue slave: 
It should take an order from the master and take the elevator to the desired floor. When the elevator has arrived it should alert its master. 

If the elevator is going up and have an order from the command panel, this should be prioritized.If the command order is in the other way of the current direction and the elevator gets an external order which is in the same direction as the current direction this order should be prioritized. We have this, since we assume that there is no evil passengers saying that they want up, then ordering a lower floor.

The module should control all lights and keep the door open for three seconds if its to stop on a floor. It should also tell its master about new orders and current floor and direction (only telling when there's a change, so this doesnt get spammed. Need to use channels smartly here, I think) 


Think we need some channel use for the checking of elevator sensors and buttons, cant check all the time. 

******NETWORK********
--The network module:

Yngve: Should we have a format on the messages sent in the form of that the first sign in the message specifies what type of message it is?

If multiple slaves says "I am slave" at the same time, we will use net.setReadDeadLine hopefully with randim waittime, the first one times out will call out; i am master,
If one elevator detects multiple new elevators, where 1 is saying "I'm master", then it should assume that its connected to an already established group and subdue to its master. 

--Proposition to message types sent over the network:
(some of these can also be merged/combined)
nr. component: name; format
1. Slave: I am slave;"ias"
2. Master: I am master;"iam"
3. Slave: Order performed; "per"
4. Slave: Order received;"sre"
5. Slave: 
6. Slave: 
7. Master: New external list of orders;"ord"
8. Master: Order received;"mre"
9. Master: Confirm slave excecuted order;"msi"
10. Slave/Master: State changed (direction/current floor/stopped);"sch"
11. (Should a slave tell the master that its not doing anything and is available?);


is #10 enough to get all the elevators to be up to date about the complete elevator-structure that the master needs?



--How do we choose a master?
	- Can have the slaves wait a random seeded time to listen 	if other slaves also want to be master
	- If a master suddenly gets messages from more than one 	other elevator and one of the elevators claim to be 	master, then it should assume that its connected to an 	established network and submit as a slave.    


--Does master need to know about the internal orders?
	- Slaves can tell that they have an order, and combined 	with the state information(direction/current floor), the 	optimization algorithm can mark some of the floors as
	"low priority" or something. If all else fails, pick
	these unattractive floors.  


--All elevators need to have a copy of the current elevators and its orders/position/


--The order optimalization module: Is this going to adapt itself to a bigger number of elevators? I reckon yes!





We also need to find out the interfaces between the different modules. Expect tons of work the next two weeks.
