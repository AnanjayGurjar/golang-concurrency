package main

import (
	"fmt"
	"sync"
	"time"
)

type Philosopher struct {
	name      string
	rightFork int
	leftFork  int
}

var philosophers = []Philosopher{
	{
		name: "Plato", leftFork: 4, rightFork: 0,
	},
	{
		name: "Scorates", leftFork: 0, rightFork: 1,
	},
	{
		name: "Aristotle", leftFork: 1, rightFork: 2,
	},
	{
		name: "Pascal", leftFork: 2, rightFork: 3,
	},
	{
		name: "Locke", leftFork: 3, rightFork: 4,
	},
}

var hunger = 3 //how many times does a philosopher eats
var eatTime = 1 * time.Second
var thinkTime = 3 * time.Second
var sleepTime = 1 * time.Second

func main() {

	fmt.Printf("The Dining Philosophers problem\n")
	fmt.Printf("-------------------------------\n")

	fmt.Printf("The table is empty\n")

	dine()

	fmt.Printf("The table is empty\n")
}

func dine() {
	//to check if everyone is done eating
	wg := &sync.WaitGroup{}
	wg.Add(len(philosophers))

	//to check if everyone is seated
	seated := &sync.WaitGroup{}
	seated.Add(len(philosophers))

	forks := make(map[int]*sync.Mutex)

	for i := 0; i < len(philosophers); i++ {
		forks[i] = &sync.Mutex{}
	}

	for i := 0; i < len(philosophers); i++ {
		go diningProblem(philosophers[i], wg, forks, seated)
	}

	wg.Wait()

}

func diningProblem(philosopher Philosopher, wg *sync.WaitGroup, forks map[int]*sync.Mutex, seated *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("%s is seated at the table\n", philosopher.name)

	seated.Done()
	//wait till everyone is seated
	seated.Wait()

	for i := hunger; i > 0; i-- {

		if(philosopher.leftFork > philosopher.rightFork){
			forks[philosopher.rightFork].Lock()
			fmt.Printf("\t%s picked the right fork\n", philosopher.name)

			forks[philosopher.leftFork].Lock()
			fmt.Printf("\t%s picked the left fork\n", philosopher.name)
		}else{
			forks[philosopher.leftFork].Lock()
			fmt.Printf("\t%s picked the left fork\n", philosopher.name)

			forks[philosopher.rightFork].Lock()
			fmt.Printf("\t%s picked the right fork\n", philosopher.name)
		}
		

		fmt.Printf("\t%s has both the forks and is eating...\n", philosopher.name)
		time.Sleep(eatTime)

		fmt.Printf("\t%s is thinking...\n", philosopher.name)
		time.Sleep(thinkTime)

		forks[philosopher.leftFork].Unlock()
		forks[philosopher.rightFork].Unlock()

		fmt.Printf("\t%s put down the forks\n", philosopher.name)

	}
	fmt.Printf("%s is satisfied and left the table\n", philosopher.name)
}
