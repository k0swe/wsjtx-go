// Simple driver binary for wsjtx-go library

package main

import (
	"fmt"
	"github.com/k0swe/wsjtx-go/wsjtx"
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
		case wsjtx.QsoLoggedMessage:
			fmt.Println("QSO Logged:", message)
		case wsjtx.CloseMessage:
			fmt.Println("Close:", message)
		case wsjtx.LoggedAdifMessage:
			fmt.Println("Logged Adif:", message)
		default:
			fmt.Println("Other:", reflect.TypeOf(message), message)
		}
	}
}
