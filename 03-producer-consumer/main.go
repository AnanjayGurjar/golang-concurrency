package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

const NumberOfPizzas = 10

var pizzasMade, pizzasFailed, total int

type Producer struct {
	data chan PizzaOrder
	quit chan chan error
}

type PizzaOrder struct {
	pizzaNumber int
	message     string
	success     bool
}

func (p *Producer) Close() error {
	ch := make(chan error)
	p.quit <- ch
	return <-ch
}
func pizzeria(pizzaMaker *Producer) {
	//keep track of pizza we are making
	var i = 0
	//run forever until we recieve a quit notification
	//try to make pizzas

	for {
		//try to make a pizza
		currentPizza := makePizza(i)

		if currentPizza != nil {
			i = currentPizza.pizzaNumber
			select {
			//we tried making the pizza and sent it to the data channel
			case pizzaMaker.data <- *currentPizza:

			case quitChan := <-pizzaMaker.quit:
				close(pizzaMaker.data)
				close(quitChan)
				return
			}
		}
	}

}

func makePizza(pizzaNumber int) *PizzaOrder {
	pizzaNumber++

	if pizzaNumber <= NumberOfPizzas {
		delay := rand.Intn(5) + 1

		fmt.Printf("Recieved order #%d\n", pizzaNumber)

		rnd := rand.Intn(12) + 1
		msg := ""
		success := false

		if rnd < 5 {
			pizzasFailed++
		} else {
			pizzasMade++
		}

		total++
		fmt.Printf("Making pizza #%d, will take %d seconds...\n", pizzaNumber, delay)

		//delay for a bit
		time.Sleep(time.Duration(delay) * time.Second)

		if rnd <= 2 {
			msg = fmt.Sprintf("*** We ran out of ingredients for pizza #%d", pizzaNumber)
		} else if rnd <= 4 {
			msg = fmt.Sprintf("***The cook quit while making pizza #%d", pizzaNumber)
		} else {
			success = true
			msg = fmt.Sprintf("Pizza order %d is ready", pizzaNumber)
		}

		p := PizzaOrder{
			pizzaNumber: pizzaNumber,
			message:     msg,
			success:     success,
		}

		return &p
	}

	return &PizzaOrder{
		pizzaNumber: pizzaNumber,
	}

}
func main() {
	//seed the random number generator
	rand.Seed(time.Now().UnixNano())

	//print out a message
	color.Cyan("Pizzeria is open for buisness")
	color.Cyan("-----------------------------")

	//create a producer
	pizzaJob := &Producer{
		data: make(chan PizzaOrder),
		quit: make(chan chan error),
	}

	//run the producer in the bg
	go pizzeria(pizzaJob)

	//create and run consumer
	for i := range pizzaJob.data {
		if i.pizzaNumber <= NumberOfPizzas {
			if i.success {
				color.Green(i.message)
			} else {
				color.Red(i.message)
			}
		} else {
			color.Cyan("Done making pizzas....")
			err := pizzaJob.Close()
			if err != nil {
				color.Red("*** Error closing channel***", err)
			}
		}
	}

	//print out the ending message

	color.Cyan("Done for the day!!!")
	color.Cyan("We made %d pizzas, but failed to make %d, with total %d attempts", pizzasMade, pizzasFailed, total)
}
