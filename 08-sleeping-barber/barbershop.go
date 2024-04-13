package main

import (
	"time"

	"github.com/fatih/color"
)

type BarberShop struct {
	ShopCapacity    int
	HairCutDuration time.Duration
	NumberOfBarbers int
	BarberDoneChan  chan bool
	ClientsChan     chan string
	Open            bool
}

func (shop *BarberShop) addBarber(barber string) {
	shop.NumberOfBarbers++

	go func() {
		isSleeping := false
		color.Yellow("%s goes to the waiting room to check for clients", barber)

		for {
			if(len(shop.ClientsChan) == 0) {
				color.Yellow("There is nothing to do, so %s takes the nap", barber)
				isSleeping = true
			}

			client, isShopOpen := <-shop.ClientsChan

			if isShopOpen {
				if isSleeping {
					color.Yellow("%s wakes %s up", client, barber)
					isSleeping = false
				}

				//cut hair
				shop.cutHair(barber, client)
			} else {
				//shop is closed, so send the barber home
				shop.sendBarberHome(barber)
				return
				
			}

		}
	}()
}

func (shop *BarberShop) cutHair(barber string, client string) {
	color.Green("%s is cutting %s's hair", barber, client)
	time.Sleep(shop.HairCutDuration)
	color.Green("%s is finished cutting %s's hair", barber, client)
}

func (shop *BarberShop) sendBarberHome(barber string) {
	color.Cyan("%s is going home", barber)
	shop.BarberDoneChan <- true
}

func (shop *BarberShop) closeShopForDay() {
	color.Cyan("Shop closing for the day")

	close(shop.ClientsChan)
	shop.Open = false

	for a := 1; a <= shop.NumberOfBarbers; a++ {
		<- shop.BarberDoneChan
	}

	close(shop.BarberDoneChan)

	color.Green("The barbershop is closed and everyone is gone home")
}

func (shop *BarberShop) addClient(client string) {
	color.Green("--- %s arrives to the shop", client)

	if(shop.Open) {
		select{
		case shop.ClientsChan <- client:
			color.Yellow("%s takes the seat in the waiting room", client)
		default:
			color.Red("The waiting room is filled, %s leaves", client)
		}
	} else {
		color.Red("The shop is closed so %s leaves", client)
	}
}