package watchdog

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"../network/bcast"
)

const startupTimeLimit time.Duration = 2 * time.Second

// Watchdog supervises a process sending on port wdPort. If no message has
// arrived in a specified time interval, Watchdog starts the process again.
//
// * Parameter port is which port the watchdog will listen to.
// * Parameter timeout is how long between messages are allowed. If no message is
// received after timeout amount of time, the process is restarted.
// * Parameter message is the message that needs to be received for the message
// to be accepted as a "still alive" message.
// * Parameter exePath is the path of the executable to be executed on time out.
// * Parameter args are arguments that should be passed to the executable.
func Watchdog(port int, timeout time.Duration, message string, exePath string, args []string) {
	wdChan := make(chan string)
	bcast.InitLogger()
	go bcast.Receiver(port, wdChan)

	wdTimer := time.NewTimer(timeout)
	respawn := false

	for {
		select {
		case msg := <-wdChan:
			fmt.Printf("Received message: \"%s\"\n", msg)
			if msg == message {
				wdTimer.Stop()
				wdTimer.Reset(timeout)
			}

		case <-wdTimer.C:
			// process did not respond in time, respawn
			fmt.Println("No message received after timeout")
			respawn = true
		}

		if respawn {
			respawn = false

			fmt.Println("Spawing process")
			cmd := exec.Command(exePath, args...)
			// TODO: figure out how to prefix output with a string
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err := cmd.Start()
			if err != nil {
				panic(fmt.Sprintf("Couldn't respawn process: %#v", cmd))
			}

			disown := exec.Command("disown", "-a")
			disown.Run()

			// Wait a while for process to start again
			<-time.After(startupTimeLimit)
			wdTimer.Reset(timeout)
			fmt.Println("Watchdog timer started again")
		}
	}
}
