package watchdog

import (
	"fmt"

    "../lib/network/network/bcast"
)

const wdPort int = 57005

// Watchdog supervises a process sending on port wdPort. If no message has arrived
// in a specified time interval, Watchdog starts the process again.
func Watchdog() {
	wdChan := make(chan string)

	go bcast.Receiver(wdPort, wdChan)

	for {
		select {
		case msg := <-wdChan:
			fmt.Printf("Received message: \"%s\"\n", msg)
		}
	}
}
