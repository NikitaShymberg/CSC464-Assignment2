package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type general struct {
	ID             int
	allegiance     string
	channels       []chan string
	order          string
	receivedOrders []string
}

func (g *general) receiveOrder(n int) {
	// fmt.Printf("%v is waiting for order from %v\n", g.ID, n)
	g.order = <-g.channels[n]
	// fmt.Printf("%v RECEIVED order from %v\n", g.ID, n)
	g.receivedOrders = append(g.receivedOrders, g.order)
}

func (g *general) sendOrder() {
	if g.allegiance == "A" {
		for i, c := range g.channels {
			go func(c chan string, i int) {
				// fmt.Printf("%v sending to %v\n", g.ID, i)
				c <- g.order
				// fmt.Printf("%v SENT to %v\n", g.ID, i)
			}(c, i)
		}
	} else {
		for i, c := range g.channels {
			var order string
			if i%2 != 0 {
				order = g.order
			} else {
				if g.order == "ATTACK" {
					order = "RETREAT"
				} else {
					order = "ATTACK"
				}
			}
			go func(c chan string, order string) {
				c <- order
			}(c, order)
		}
	}
}

func (g *general) finalizeOrder() {
	var attacks int
	var retreats int
	fmt.Printf("Received orders by %v : %v\n", g.ID, g.receivedOrders)
	for _, v := range g.receivedOrders {
		if v == "ATTACK" {
			attacks++
		} else {
			retreats++
		}
	}
	if attacks > retreats {
		g.order = "ATTACK"
	} else if retreats > attacks {
		g.order = "RETREAT"
	} else {
		g.order = "TIE"
	}
}

func main() {
	var m int
	var commanderAllegiance string
	var generalAlliances []string // TODO: rename me I'm only the input
	var order string

	m, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	commanderAllegiance = os.Args[2]

	for _, v := range os.Args[3 : len(os.Args)-1] {
		if v != "A" && v != "T" {
			log.Fatal("Generals can only be Allies (A) or Traitors (T)")
		}
		generalAlliances = append(generalAlliances, v)
	}

	order = os.Args[len(os.Args)-1]
	if order != "ATTACK" && order != "RETREAT" {
		log.Fatal("Order can only be ATTACK or RETREAT")
	}

	fmt.Printf("M: %v\n", m)
	fmt.Printf("COMMANDER: %v\n", commanderAllegiance)
	fmt.Printf("GENENERALS: %v\n", generalAlliances)
	fmt.Printf("ORDER: %v\n", order)

	// Create generals
	var commanderToGenerals []chan string
	for range generalAlliances {
		commanderToGenerals = append(commanderToGenerals, make(chan string))
	}
	commander := general{allegiance: commanderAllegiance, channels: commanderToGenerals, order: order}

	var lieutenants []general
	for i, v := range generalAlliances {
		var channels = make([]chan string, len(generalAlliances)+1)
		channels[0] = commanderToGenerals[i]
		lieutenants = append(lieutenants, general{ID: i + 1, allegiance: v, channels: channels})
	}

	// Add channels to the other lieutenants
	for i := range lieutenants {
		for j := i; j < len(lieutenants); j++ {
			newChan := make(chan string)
			// lieutenants[i].channels = append(lieutenants[i].channels, newChan)
			// lieutenants[j].channels = append(lieutenants[j].channels, newChan)
			lieutenants[i].channels[lieutenants[j].ID] = newChan
			lieutenants[j].channels[lieutenants[i].ID] = newChan
		}
	}

	// Step 1
	go commander.sendOrder()
	for i := range lieutenants {
		updatedGeneral := make(chan general)
		go func(g general) {
			g.receiveOrder(0)
			updatedGeneral <- g
		}(lieutenants[i])
		go func(lieutenants []general, i int) {
			lieutenants[i] = <-updatedGeneral
		}(lieutenants, i)
	}

	time.Sleep(time.Second)

	// Step 2
	for i := 0; i < m; i++ {
		for j := range lieutenants {
			go func(g general) {
				lieutenants[j].sendOrder()
			}(lieutenants[j])

			for k := range lieutenants {
				if k != j {
					updatedGeneral := make(chan general)
					go func(g general, from int) {
						// fmt.Printf("S %v .... %v\n", g.ID, from)
						g.receiveOrder(from)
						// fmt.Printf("F %v .... %v\n", g.ID, from)
						updatedGeneral <- g
					}(lieutenants[k], lieutenants[j].ID)
					// go func(lieutenants []general, k int) {
					lieutenants[k] = <-updatedGeneral
					// }(lieutenants, k)
				}
			}
		}
	}

	time.Sleep(time.Second)

	// Step 3
	for _, g := range lieutenants {
		g.finalizeOrder()
	}

	time.Sleep(time.Second)

	for _, g := range lieutenants {
		if g.allegiance == "A" {
			fmt.Printf("The loyal lieutenant %v's order: %v\n", g.ID, g.order)
		}
	}

}
