package main

import (
	"fmt"
	"time"
)

func spinner(delay time.Duration, done <-chan struct{}) {

	isDone := false

	go func() {
		<-done
		isDone = true
	}()
	for {
		fmt.Print("\rdownload")
		for i := 0; i <= 5 ; i++  {
			if isDone {
				return
			}
			fmt.Print(".")
			time.Sleep(delay)
		}
	}
}

func testSpinner() {

	done := make(chan struct{})
	go spinner(1000 * time.Millisecond, done)
	//Do something
	time.Sleep(30 * time.Second)
	done <- struct{}{}
	close(done)
}

func main()  {
	testSpinner()
}
