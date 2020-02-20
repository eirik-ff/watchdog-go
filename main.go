package main

import (
    "./watchdog"
)

func main() {
	go watchdog.Watchdog()

    for {
    }
}

