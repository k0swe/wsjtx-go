// Simple driver binary for wsjtx-go library

package main

import (
	"fmt"
	"github.com/xylo04/wsjtx-go/wsjtx"
	"reflect"
)

func main() {
	fmt.Println("Listening for WSJT-X...")
	c := make(chan interface{}, 5)
	go wsjtx.ListenToWsjtx(c)
	for {
		message := <-c
		switch message.(type) {
		case wsjtx.HeartbeatMessage:
			fmt.Println("Heartbeat:", message)
		case wsjtx.StatusMessage:
			fmt.Println("Status:", message)
		case wsjtx.DecodeMessage:
			fmt.Println("Decode:", message)
		case wsjtx.ClearMessage:
			fmt.Println("Clear:", message)
		default:
			fmt.Println("Other:", reflect.TypeOf(message), message)
		}
	}
}
