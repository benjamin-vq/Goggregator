package main

import "log"

func assert(condition bool, msg string) {
	if !condition {
		log.Fatal(msg)
	}
}
