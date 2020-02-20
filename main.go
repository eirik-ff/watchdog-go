package main

import (
	"flag"
	"os/user"
	"strings"
	"time"

	"./watchdog"
)

const wdPortDefault int = 57005
const wdTimeoutDefault = 5000 // 5 seconds

func main() {
	usr, _ := user.Current()
	homeDir := usr.HomeDir

	port := flag.Int("port", wdPortDefault, "Port for communicating with the watchdog")
	timeout := flag.Int("timeout", wdTimeoutDefault, "Timeout in milliseconds")
	exe := flag.String("exec", homeDir+"/sanntid-heis-gr28/heis",
		strings.Join([]string{
			"Path of the executable that will be respawned. Supports arguments if ",
			"placed in quotes ('single' or \"double\"). Use \"go build\" to create",
			" an executable file in Go",
		}, ""))
	flag.Parse()

	args := strings.Split(*exe, " ")
	exePath := args[0]
	args = args[1:]

	go watchdog.Watchdog(*port, time.Duration(*timeout)*time.Millisecond, exePath, args)

	for {
	}
}
