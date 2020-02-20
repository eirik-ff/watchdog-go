package watchdog

import (
	"fmt"
	"os/exec"
	"time"

	"../lib/network/network/bcast"
)

const startupTimeLimit time.Duration = 2 * time.Second

// Watchdog supervises a process sending on port wdPort. If no message has
// arrived in a specified time interval, Watchdog starts the process again.
// Parameter port is which port the watchdog will listen to.
// Parameter timeout is how long between messages are allowed. If no message is
// received after timeout amount of time, the process is restarted.
// Parameter cmd is the command that should be executed after timeout duration.
func Watchdog(port int, timeout time.Duration, exePath string, args []string) {
	wdChan := make(chan string)
	go bcast.Receiver(port, wdChan)

	wdTimer := time.NewTimer(timeout)
	respawn := false

	for {
		select {
		case msg := <-wdChan:
			fmt.Printf("Received message: \"%s\"\n", msg)
			wdTimer.Stop()
			wdTimer.Reset(timeout)

		case <-wdTimer.C:
			// process did not respond in time, respawn
			fmt.Println("No message received after timeout")
			respawn = true
		}

		if respawn {
			respawn = false

			fmt.Println("Spawing process")
			cmd := exec.Command(exePath, args...)

			err := cmd.Start()
			if err != nil {
				panic(fmt.Sprintf("Couldn't respawn process: %#v", cmd))
			}

			// Wait a while for process to start again
			<-time.After(startupTimeLimit)
			wdTimer.Reset(timeout)
			fmt.Println("Watchdog timer started again")
		}
	}
}
