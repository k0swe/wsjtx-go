package integration

import (
	"net"
	"strconv"
)

type WsjtxFake struct {
	port uint
	conn net.Conn
	resp []byte
}

// NewFake initializes a new fake WSJTX program on an OS-assigned port.
func NewFake(addr net.Addr) (*WsjtxFake, error) {
	conn, err := net.Dial("udp", addr.String())
	if err != nil {
		return &WsjtxFake{}, err
	}

	_, portStr, _ := net.SplitHostPort(conn.LocalAddr().String())
	port, _ := strconv.Atoi(portStr)
	w := &WsjtxFake{uint(port), conn, []byte("")}
	//go w.handleRequests()
	return w, nil
}

// SendMessage immediately sends the given payload out from the WSJTX fake.
func (w *WsjtxFake) SendMessage(payload []byte) (int, error) {
	return w.conn.Write(payload)
}

//func (w *WsjtxFake) SeedResponse(payload []byte) {
//	w.resp = payload
//}
//
//func (w *WsjtxFake) handleRequests() {
//	b := make([]byte, 2048)
//	for {
//		n, err := w.conn.Read(b)
//		if err != nil {
//			if err != io.EOF {
//				fmt.Println("read error:", err)
//			}
//			break
//		}
//		//fmt.Println("got", n, "bytes.")
//		buf = append(buf, tmp[:n]...)
//
//	}
//}

func (w *WsjtxFake) Stop() {
	_ = w.conn.Close()
}
