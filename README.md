# CSC464-Assignment2
Disclaimer: the code used in this assignment is not meant for use in any sort of production environment. It is intended to be used in an academic setting to demonstrate the functionality of vector clocks and the Byzantine generals algorithm.

## Vector clocks

This implementation creates three functions `a`, `b`, and `c` (or processes 1, 2, and 3 respectively) that are run as goroutines, call these functions processes to be consistent with Lamport's terminology. Each process can have simulate a local event by calling `localEvent()`, or they can use `sendMessage()` and `receiveMessage()` to communicate between each other. These simulate dependent events - or happens before relationships. Whenever an event happens, the process's clock is printed out with a letter: `L` means that this was a local process, `S` means a message was sent, and `R` means a message was received.

### Example run:
Sample output
```
2 -L- [0 0 1]
2 -L- [0 0 2]
2 -L- [0 0 3]
2 -L- [0 0 4]
0 -S- [1 0 0]
1 -R- [1 1 0]
1 -S- [1 2 0]
0 -R- [2 2 0]
```

From this output we can deduce several things.

1. Events in process 3 were entirely parallel to events in the other processes as the clock times were not less than or equal to any of the other clocks for _all_ processes.
2. The first sent message must happen before the first received message. This is because [1, 0, 0] is less than [1, 1, 0] at index 1 (and equal at the other indices), therefore the former clock shows an event that happened before the latter clock
3. Similarly, the second send happens before teh second receive as [1, 2, 0] is less than [2, 2, 0] at index 0 (and equal at the other indices).

This implementation demonstrates how the ordering of events can be inferred by analyzing vector clocks in each process after each event. It also demonstrates the transitive nature of the ordering of events. For example the first send happens before the first receive which happens before the second send. The clocks of the first [1, 0, 0] and second [1, 2, 0] send messages confirm this as the time at index 1 is less in the first process than in the second (and equal everywhere else).

## Byzantine Generals

Use: `go run ByzantineGenerals.go <m> <generals> <order>` where m is the number of traitors, generals is a list of letters indicading the generals' allegiances. The letter A means that the general is an ally and the letter T means that the general is a traitor. Lastly order is the initial order given by the commmander, it can be either ATTACK or RETREAT.

NOTE: `m` must be less than one third of the total number of generals for the algorithm to consistently work.

This implementation first creates the appropriate amount of `general` structs and assigns their allegiance. Then it creates a channel between each pair of generals, this channel will be used by each general to send order to each other. The commander general is called `commander ` and the others are in an array called `lieutenants`. The implementation then follows the three steps of the Byzantine general algorithm. Between each step there is a one second pause to allow all the goroutines to finish. A necessary improvement to this implementation would be to use waitGroups or some other synchronization mechanism, but for the purposes of this exercise this is not an important consideration as it still denomstrates an implementation of the algorithm. 

During the first step the commander sends each lieutenant their initial order who then set it as their active `order` and will use this value to pass on to the other lieutenants to attempt to reach consensus. In the second step, lieutenants send orders to each other and gather them in the `receivedOrders` array. In the third step, each general will pick the majority of their orders and set than as their active `order`. Finally when the algorithm is completed the received orders of the active generals are printed out.

## Tests

The file `ByzantineGeneralsTestCases.txt` contains commands to simulate test cases as well as their expected output, the more interesting test cases will be discussed here.

### Trivial example 1
```
$go run ByzantineGenerals.go 0 A A A A A ATTACK

The loyal lieutenant 1's order: ATTACK
The loyal lieutenant 2's order: ATTACK
The loyal lieutenant 3's order: ATTACK
The loyal lieutenant 4's order: ATTACK
```

This example demonstrates the base case of m=0 where there are no traitors. Essentially only the general sends out the given order and all the lieutenants just use it.

### Example 2

```
$go run ByzantineGenerals.go 2 A T A A A A ATTACK

The loyal lieutenant 1's order: ATTACK
The loyal lieutenant 2's order: ATTACK
The loyal lieutenant 3's order: ATTACK
The loyal lieutenant 5's order: ATTACK
```

This example demonstrates the case where one of the lieutenants is a traitor. The algorithm runs through one iteration of step two where all lieutenants talk to each other. Some lieutenants receive the wrong order to RETREAT from the traitor lieutenant, however the majority remains ATTACK and the lieutenants reach the correct decision. There are two things worth noting here, the first is that since the commander was loyal, all loyal generals must follow his order which they did. The second is that the same result would have been achieved without executing the second step of the algorithm, since all generals would have received the correct order. This brings us to the interesting example 3:

### Example 3

```
$go run ByzantineGenerals.go 1 T A A A A ATTACK

The loyal lieutenant 1's order: RETREAT
The loyal lieutenant 2's order: RETREAT
The loyal lieutenant 3's order: RETREAT
The loyal lieutenant 4's order: RETREAT
```

This is really the case that is solved by the algorithm. Notice how all the loyal genrals agreed on the same order, however it was not the order that the traitorous commander gave but this is ok. This example really demonstrates that the algorithm works.

The other test cases mostly serve to demonstrate that the solution scales well and still has correct results.