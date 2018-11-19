package main

import (
	"fmt"
)

// TODO: write tests

// NumThreads is the number of goroutines to be run
const NumThreads = 3

type vectorClock struct {
	owner int
	times [NumThreads]int
}

type message struct {
	clock vectorClock
	data  string
}

func (vc *vectorClock) localEvent() {
	vc.times[vc.owner]++
}

func (vc *vectorClock) sendMessage(data string, dest chan message) {
	vc.times[vc.owner]++
	var m message
	m.data = data
	m.clock = *vc

	fmt.Printf("Sending message message %v ...\n", m.data)
	dest <- m
}

func (vc *vectorClock) receiveMessage(src chan message) {
	m := <-src
	for i, v := range m.clock.times {
		if vc.times[i] < v {
			vc.times[i] = v
		}
	}
	fmt.Printf("Received message %v\n", m.data)
}

func a(toB chan message, toC chan message, out chan vectorClock) {
	var myClock vectorClock
	myClock.owner = 0

	myClock.localEvent()
	myClock.sendMessage("Hello!", toB)
	myClock.localEvent()
	myClock.receiveMessage(toC)

	out <- myClock
}

func b(toA chan message, toC chan message, out chan vectorClock) {
	var myClock vectorClock
	myClock.owner = 1

	myClock.receiveMessage(toA)
	myClock.sendMessage("World!", toC)
	myClock.localEvent()
	myClock.localEvent()

	out <- myClock
}

func c(toA chan message, toB chan message, out chan vectorClock) {
	var myClock vectorClock
	myClock.owner = 2

	myClock.localEvent()
	myClock.localEvent()
	myClock.receiveMessage(toB)
	myClock.localEvent()
	myClock.sendMessage("Star", toA)
	myClock.localEvent()

	out <- myClock
}

func main() {
	results := make(chan vectorClock)
	aToB := make(chan message)
	aToC := make(chan message)
	bToC := make(chan message)

	go a(aToB, aToC, results)
	go b(aToB, bToC, results)
	go c(aToC, bToC, results)

	var finalClocks []vectorClock

	for i := 0; i < NumThreads; i++ {
		finalClocks = append(finalClocks, <-results)
	}

	for _, clock := range finalClocks {
		fmt.Printf("Owner %v, clock %v\n", clock.owner, clock.times)
	}
}
