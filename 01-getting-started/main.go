package main

import (
	"fmt"
	"sync"
)

func printSomething(detail string, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Println(detail);
}

var msg string

func setMessage(message string, wg *sync.WaitGroup){
	defer wg.Done()
	msg = message
}

func printMessage(){
	fmt.Println(msg)
}


func main(){
	// go printSomething("this will not be printed")

	// printSomething("this will be printed")

	/*
	//----------------------introduction part 
	var wg sync.WaitGroup;

	words := [] string{
		"One",
		"Two",
		"three",
		"four",
		"five",
		"six",
		"seven",
		"eight",
	}

	wg.Add(len(words));

	for i, word := range words {
		go printSomething(fmt.Sprintf("%d : %s", i, word), &wg);
	}

	wg.Wait()
	// time.Sleep(time.Second)
	wg.Add(1)
	printSomething("atlast", &wg)

	*/

	//-----------------challange part

	var wg sync.WaitGroup

	wg.Add(1);
	go setMessage("Hello universe", &wg)
	wg.Wait()
	printMessage()

	wg.Add(1)
	go setMessage("Hello cosmos", &wg)
	wg.Wait()
	printMessage()

	wg.Add(1)
	go setMessage("Hello world", &wg)
	wg.Wait()
	printMessage()

	


}
