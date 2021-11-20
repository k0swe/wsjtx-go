package integration

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"testing"
	"time"
)

// WsjtxFake is a test double that acts like WSJTX. It connects to the wsjtx-go server as a client,
// but is actually a stub controlled by the integration test cases.
type WsjtxFake struct {
	t           *testing.T
	conn        *net.UDPConn
	ReceiveChan chan []byte
	stop        chan bool
}

// NewFake initializes a new fake WSJTX program on an OS-assigned port.
func NewFake(addr *net.UDPAddr, t *testing.T) (*WsjtxFake, error) {
	conn, err := net.DialUDP("udp", &net.UDPAddr{Port: 0}, addr)
	if err != nil {
		return &WsjtxFake{}, err
	}
	t.Logf("fake is connected to %v", conn.RemoteAddr())

	w := &WsjtxFake{t, conn, make(chan []byte, 5), make(chan bool, 1)}
	go w.handleReceive()
	return w, nil
}

// SendMessage immediately sends the given payload out from the WSJTX fake.
func (w *WsjtxFake) SendMessage(payload []byte) (int, error) {
	w.t.Log("sending message")
	return w.conn.Write(payload)
}

func (w *WsjtxFake) handleReceive() {
	b := make([]byte, 2048)
	w.t.Log("listening for receives")
	for {
		select {
		case <-w.stop:
			w.t.Log("stopping")
			close(w.ReceiveChan)
			_ = w.conn.Close()
			return
		default:
			_ = w.conn.SetDeadline(time.Now().Add(1 * time.Millisecond))
			n, err := w.conn.Read(b)
			if err != nil {
				if err != io.EOF && !errors.Is(err, os.ErrDeadlineExceeded) {
					w.t.Log(fmt.Errorf("got an error while reading UDP: %w", err))
				}
			}
			if n > 0 {
				tmp := make([]byte, n)
				copy(tmp, b[:n])
				w.t.Logf("received %d bytes, putting on channel", n)
				w.ReceiveChan <- tmp
			}
		}
	}
}

func (w *WsjtxFake) Stop() {
	w.stop <- true
}
