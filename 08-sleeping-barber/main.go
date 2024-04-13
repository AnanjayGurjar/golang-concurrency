// problem descritpion : "https://en.wikipedia.org/wiki/Sleeping_barber_problem"
package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

var seatingCapacity = 10
var arrivalRate = 100
var cutDuration = 1000 * time.Millisecond
var timeOpen = 10 * time.Second


func main() {
	//seed our random number generator
	rand.Seed(time.Now().UnixNano())

	//print welcome message
	color.Yellow("The Sleeping barber problem")
	color.Yellow("---------------------------")

	//create channels
	clientChan := make(chan string, seatingCapacity)
	doneChan := make(chan bool)
	
	//create barbershop
	shop := BarberShop {
		ShopCapacity: seatingCapacity,
		HairCutDuration: cutDuration,
		NumberOfBarbers: 0,
		ClientsChan: clientChan,
		BarberDoneChan: doneChan,
		Open: true,
	}

	color.Green("The Shop is open for the day!")

	//add barbers
	shop.addBarber("Frank")
	shop.addBarber("John")
	shop.addBarber("Steve")
	shop.addBarber("Paul")
	//start barbershop as goroutine
	shopClosing := make(chan bool)
	closed := make(chan bool)

	go func(){
		<- time.After(timeOpen)
		shopClosing <- true
		shop.closeShopForDay()
		closed <- true
	}()

	//add clients
	i := 1
	go func() {
		for {
			randomMillisecond := rand.Int() % (2 * arrivalRate)

			select {
			case <- shopClosing:
				return
			case <- time.After(time.Millisecond * time.Duration(randomMillisecond)):
				shop.addClient(fmt.Sprintf("Client #%d", i))
				i++
			}
		}
	}()

	//block until barbershop is closed
	<- closed
}