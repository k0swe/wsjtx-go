package wsjtx

import (
	"fmt"
	"net"
)

const Magic = 0xadbccbda
const BufLen = 1024

type Server struct {
	conn *net.UDPConn
	id   string
}

func MakeServer() Server {
	// TODO: make address and port customizable?
	musticastAddr := "224.0.0.1"
	wsjtxPort := "2237"
	addr, err := net.ResolveUDPAddr("udp", musticastAddr+":"+wsjtxPort)
	check(err)
	conn, err := net.ListenMulticastUDP("udp", nil, addr)
	check(err)
	return Server{conn, "WSJT-X"}
}

func (s *Server) Clear(bandActivity bool, rxFrequency bool) error {
	if !bandActivity && !rxFrequency {
		return nil
	}
	var window uint8
	if bandActivity && rxFrequency {
		window = 2
	} else if bandActivity {
		window = 0
	} else if rxFrequency {
		window = 1
	}

	msg := ClearMessage{
		Id:     s.id,
		Window: window,
	}
	fmt.Printf("Pretend I'm sending Clear:%v", msg)
	// TODO
	//s.conn.Write();
	return nil
}

// Goroutine which will listen on a UDP port for messages from WSJT-X. When heard, the messages are
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

func check(err error) {
	if err != nil {
		panic(err)
	}
}
