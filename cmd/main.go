package main

import (
	"bufio"
	"fmt"
	"github.com/k0swe/wsjtx-go"
	"os"
	"reflect"
	"strings"
)

// Simple driver binary for wsjtx-go library.
func main() {
	fmt.Println("Listening for WSJT-X...")
	wsjtxChannel := make(chan interface{}, 5)
	wsjtxServer := wsjtx.MakeServer()
	go wsjtxServer.ListenToWsjtx(wsjtxChannel)

	stdinChannel := make(chan string, 5)
	go stdinCmd(stdinChannel)

	for {
		select {
		case message := <-wsjtxChannel:
			handleServerMessage(message)
		case command := <-stdinChannel:
			command = strings.ToLower(command)
			handleCommand(command, wsjtxServer)
		}
	}
}

// Goroutine to listen to stdin.
func stdinCmd(c chan string) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		for scanner.Scan() {
			input := scanner.Text()
			c <- input
		}
		if err := scanner.Err(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

// When we receive WSJT-X messages, display them.
func handleServerMessage(message interface{}) {
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
	case wsjtx.WSPRDecodeMessage:
		fmt.Println("WSPR Decode:", message)
	case wsjtx.LoggedAdifMessage:
		fmt.Println("Logged Adif:", message)
	default:
		fmt.Println("Other:", reflect.TypeOf(message), message)
	}
}

// When we get a command from stdin, send WSJT-X a message.
func handleCommand(command string, wsjtxServer wsjtx.Server) {
	var err error
	switch command {
	case "clear":
		fmt.Println("Sending Clear")
		err = wsjtxServer.Clear(wsjtx.ClearMessage{Id: "WSJT-X", Window: 2})

	case "close":
		fmt.Println("Sending Close")
		err = wsjtxServer.Close(wsjtx.CloseMessage{Id: "WSJT-X"})
	}
	if err != nil {
		fmt.Println(err)
	}
}
