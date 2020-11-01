package wsjtx

import (
	"encoding/hex"
	"net"
)

const Magic = 0xadbccbda
const BufLen = 1024

type Server struct {
	conn *net.UDPConn
}

// Create a UDP connection to communicate with WSJT-X.
func MakeServer() Server {
	// TODO: make address and port customizable?
	musticastAddr := "224.0.0.1"
	wsjtxPort := "2237"
	addr, err := net.ResolveUDPAddr("udp", musticastAddr+":"+wsjtxPort)
	check(err)
	conn, err := net.ListenMulticastUDP("udp", nil, addr)
	check(err)
	return Server{conn}
}

// Goroutine to listen for messages from WSJT-X. When heard, the messages are
// parsed and then placed in the given channel.
func (s *Server) ListenToWsjtx(c chan interface{}) {
	for {
		b := make([]byte, BufLen)
		length, _, err := s.conn.ReadFromUDP(b)
		check(err)
		message := parseMessage(b, length)
		if message != nil {
			c <- message
		}
	}
}

// Send a message to WSJT-X to clear the band activity window, the RX frequency
// window, or both.
func (s *Server) Clear(msg ClearMessage) error {
	// TODO: encode the given message
	msgBytes, _ := hex.DecodeString("adbccbda00000002000000030000000657534a542d5802")
	_, err := s.conn.Write(msgBytes)
	return err
}

// Send a message to WSJT-X to close the program.
func (s *Server) Close(msg CloseMessage) error {
	// TODO: encode the given message
	msgBytes, _ := hex.DecodeString("adbccbda00000002000000060000000657534a542d58")
	_, err := s.conn.Write(msgBytes)
	return err
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
