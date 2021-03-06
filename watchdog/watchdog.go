package watchdog

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
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
	bcast.InitLogger(fmt.Sprintf("watchdog%d", port))
	wdChan := make(chan string)
	go bcast.Receiver(port, wdChan)

	wdTimer := time.NewTimer(timeout)
	respawn := false
	heispid := 0

	args = append(args, "&") // start in background (needed for nohup)
	args = append([]string{exePath}, args...)

	for {
		select {
		case msg := <-wdChan:
			if strings.HasPrefix(msg, message) {
				localHeispid := 0
				fmt.Sscanf(msg, message+":%d", &localHeispid)
				if localHeispid != heispid {
					fmt.Printf("PID of elevator process: %d\n", localHeispid)
				}
				heispid = localHeispid
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

			if heispid > 0 {
				kill := exec.Command("kill", "-9", strconv.Itoa(heispid))
				fmt.Println("Running kill")
				kill.Start()
			}

			cmd := exec.Command("nohup", args...)
			fmt.Println("Spawing process:", cmd)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err := cmd.Start()
			if err != nil {
				panic(fmt.Sprintf("Couldn't respawn process: %#v", cmd))
			}

			// Wait a while for process to start again
			<-time.After(startupTimeLimit)
			wdTimer.Reset(timeout)
			fmt.Println("Watchdog timer started again")

			// rename output file
			mv := exec.Command("mv", "nohup.out", fmt.Sprintf("logs/heis%d.log", port))
			mv.Run()
		}
	}
}
