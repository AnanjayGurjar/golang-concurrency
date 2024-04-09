package main

import (
	"fmt"
	"sync"
)
/*
//----------- basic usage of mutex

var msg string;
var wg sync.WaitGroup

func updateMessage(updatedMessage string, mutex *sync.Mutex){
	defer wg.Done()
	mutex.Lock()
	msg = updatedMessage
	mutex.Unlock()
}
func main(){
	msg = "Hello world!"

	var mutex sync.Mutex
	

	wg.Add(2)
	go updateMessage("Hello Universe", &mutex)
	go updateMessage("Hello Cosmos", &mutex)
	wg.Wait()

	fmt.Println(msg);


}
*/


// ---------------------- sligtly complex example: you are given multiple income sources per week, need to print final income after year and current income after each week

type Income struct {
	source string
	amount int
}
func main(){
	bankBalance := 0
	var wg sync.WaitGroup;
	var mutex sync.Mutex;

	incomes := []Income{
		{"Job", 200},
		{"Gifts", 10},
		{"Part Time", 50},
		{source: "ROI", amount: 75},
	}

	wg.Add(len(incomes))

	for i, income := range incomes{
		go func (i int, income Income)  {

			defer wg.Done()

			for week := 1; week <= 52; week++ {
				mutex.Lock()
				currentBankBalance := bankBalance
				currentBankBalance += income.amount
				bankBalance = currentBankBalance
				mutex.Unlock()

				fmt.Printf("On Week %d got $%d.00 from %s and current bank balance is $%d.00\n", week, income.amount, income.source, currentBankBalance)

				
			}
		}(i, income)
	}

	wg.Wait()

	fmt.Printf("Final bank Balance at the end of year is $%d.00\n", bankBalance)


}