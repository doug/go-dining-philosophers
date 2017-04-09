/**
* Toy Go implementation of Dining Philosophers problem
* http://en.wikipedia.org/wiki/Dining_philosophers_problem
* Author: Doug Fritz
* Date: 2011-01-05
**/
package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)


//Philosopher is astructure to store Philosopher details
type Philosopher struct {
	name        string
	mychopstick chan bool
	//mychopstick chan bool
	neighbor *Philosopher
	numEat   int
	timesEat int // A var to track the number of times Philosopher ate
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	rand := rand.Intn(max-min) + min
	// fmt.Printf("The random number selected is %v\n", rand)
	return rand
}

//makePhilosopher creates a Philospher and send message to the channel that he is available
func makePhilosopher(name string, neighbor *Philosopher, numEat int) *Philosopher {
	phil := &Philosopher{name, make(chan bool, 1), neighbor, numEat, 0}
	// phil.mychopstick <- "available"
	phil.mychopstick <- true
	return phil
}

func (phil *Philosopher) think() {
	fmt.Printf("%s is thinking.\n", phil.name)
	time.Sleep(time.Duration(rand.Int63n(1e9)))
}

func (phil *Philosopher) eat() {
	fmt.Printf("%s is eating.\n", phil.name)
	phil.timesEat++
	time.Sleep(time.Duration(rand.Int63n(1e9)))
}

func (phil *Philosopher) getChopsticks() {
	timeout := make(chan string, 1)
	go func() { time.Sleep(1e9); timeout <- "timesup!" }()
	<-phil.mychopstick
	fmt.Printf("%s has now acquired one chopstick.\n", phil.name)
	select {
	case <-phil.neighbor.mychopstick:
		fmt.Printf("%s has now acquired %s's chopstick.\n", phil.name, phil.neighbor.name)
		fmt.Printf("%s has two chopsticks now.\n", phil.name)
		return
	case <-timeout:
		// phil.mychopstick <- "available"
		phil.mychopstick <- true
		phil.think()
		phil.getChopsticks()
	}
}

func (phil *Philosopher) returnChopsticks() {
	// phil.mychopstick <- "available"
	phil.mychopstick <- false
	// phil.neighbor.mychopstick <- "available"
	phil.neighbor.mychopstick <- false
}

func (phil *Philosopher) dine(full chan *Philosopher, eatCount chan int, numEat int) {
	for i := 0; i < numEat; i++ {
		phil.think()
		phil.getChopsticks()
		phil.eat()
		phil.returnChopsticks()
	}
	full <- phil
	eatCount <- phil.timesEat
}

func main() {

	maxEat, _ := strconv.Atoi(os.Args[1])

	names := []string{"Plato", "Aristotle", "Augustine", "Aquinas", "Stein"}

	philosophers := make([]*Philosopher, len(names))

	var phil *Philosopher

	for i, name := range names {
		phil = makePhilosopher(name, phil, random(1, maxEat))
		philosophers[i] = phil
	}
	philosophers[0].neighbor = phil
	fmt.Printf("There are %d philosophers sitting at a dining table.\n", len(philosophers))
	// they each have one chopstick, and must borrow from their neighbor to eat.

	full := make(chan *Philosopher)
	eatCount := make(chan int)

	for _, phil := range philosophers {
		go phil.dine(full, eatCount, phil.numEat)
	}

	for i := 0; i < len(names); i++ {
		phil := <-full
		fmt.Printf("%s is done dining.\n", phil.name)
	}

	for _, phil := range philosophers {
		times := <-eatCount
		fmt.Printf("%v ate %v times\n", phil.name, times)
	}
}
