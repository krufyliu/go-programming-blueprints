package main

import "time"

type message struct {
	From    string
	Message string
	When    time.Time
}
