package main

import (
	"time"

	"./watchdog"
)

const wdPort int = 57005

func main() {

	go watchdog.Watchdog(57005, 5*time.Second, nil)

	for {
	}
}
