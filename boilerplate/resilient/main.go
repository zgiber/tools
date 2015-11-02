package main

import (
	"fmt"
	"log"
	"time"
)

func main() {

	// the main communication channel
	c := make(chan string, 10)

	// delay mechanism for service restarts
	backoff := 1

	// testing (will close chan once after 5 seconds)
	go stopAfter5(c)

	// main loop
	for {

		// checking how long the service was running
		// if long enough, the backoff is reset to 1
		t := time.Now()
		err := runService(c)
		if err != nil {
			log.Println(err)
		}

		// new communication channel for the restarted
		// service.
		c = make(chan string, 10)

		// exponential backoff with maximum 32 sec delays.
		if time.Since(t) < time.Duration(2)*time.Second {
			log.Printf("Restarting service in %vs\n", backoff)
			time.Sleep(time.Duration(backoff) * time.Second)
			if backoff < 32 {
				backoff *= 2
			}
		}
	}
}

func runService(c chan string) (err error) {

	// channel for fatal errors
	// quit if receive from this.
	e := make(chan error)

	go doSomething1(c, e)
	go doSomething2(c, e)

	for {
		select {
		case err = <-e:
			if _, ok := <-c; ok {

				// closing the channel makes sure that
				// all goroutines writing to it will panic
				// panic must be recovered in each goroutine !!!
				close(c)
			}
			return

		// Important! always check at read if the channel is closed
		// reading from closed channel floods an empty string.
		case m, ok := <-c:

			// if the channel is closed, quit.
			if ok {
				fmt.Println(m)
			} else {
				return
			}
		}
	}
}

func stopAfter5(c chan string) {
	select {
	case <-time.After(5 * time.Second):
		close(c)
	}
}

func doSomething1(c chan string, e chan error) {

	defer handlePanic(e)

	for {
		select {
		case <-time.After(1 * time.Second):
			c <- "hello from 1"
		}
	}
}

// sending hello every 500 ms
func doSomething2(c chan string, e chan error) {

	time.Sleep(500 * time.Millisecond)
	defer handlePanic(e)

	for {
		select {
		case <-time.After(1 * time.Second):
			c <- "hello from 2"
		}
	}
}

func handlePanic(e chan error) {
	if r := recover(); r != nil {
		log.Println(r)
	}
	select {

	// sending the error on channel e makes the "parent" service stop.
	case e <- fmt.Errorf("Process exited"):
	default:
	}
}
