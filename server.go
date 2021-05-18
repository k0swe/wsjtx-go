package wsjtx

import (
	"net"
)

const magic = 0xadbccbda
const schema = 2
const qDataStreamNull = 0xffffffff
const bufLen = 1024
const multicastAddr = "224.0.0.1"
const wsjtxPort = "2237"

type Server struct {
	conn       *net.UDPConn
	remoteAddr *net.UDPAddr
}

// Create a UDP connection to communicate with WSJT-X.
func MakeServer() (Server, error) {
	return MakeMulticastServer(multicastAddr, wsjtxPort)
}

func MakeMulticastServer(addrStr string, portStr string) (Server, error) {
	var empty Server
	addr, err := net.ResolveUDPAddr("udp", addrStr+":"+portStr)
	if err != nil {
		return empty, err
	}
	conn, err := net.ListenMulticastUDP(addr.Network(), nil, addr)
	if err != nil {
		return empty, err
	}
	return Server{conn, nil}, nil
}

func MakeUnicastServer(addrStr string, portStr string) (Server, error) {
	var empty Server
	addr, err := net.ResolveUDPAddr("udp", addrStr+":"+portStr)
	if err != nil {
		return empty, err
	}
	conn, err := net.ListenUDP(addr.Network(), addr)
	if err != nil {
		return empty, err
	}
	return Server{conn, nil}, nil
}

// Goroutine to listen for messages from WSJT-X. When heard, the messages are
// parsed and then placed in the given channel.
func (s *Server) ListenToWsjtx(c chan interface{}) {
	for {
		b := make([]byte, bufLen)
		length, rAddr, err := s.conn.ReadFromUDP(b)
		c <- err
		s.remoteAddr = rAddr
		message, err := parseMessage(b, length)
		if message != nil {
			c <- message
		}
	}
}

// Send a heartbeat message to WSJT-X.
func (s *Server) Heartbeat(msg HeartbeatMessage) error {
	msgBytes, _ := encodeHeartbeat(msg)
	_, err := s.conn.WriteTo(msgBytes, s.remoteAddr)
	return err
}

// Send a message to WSJT-X to clear the band activity window, the RX frequency
// window, or both.
func (s *Server) Clear(msg ClearMessage) error {
	msgBytes, _ := encodeClear(msg)
	_, err := s.conn.WriteTo(msgBytes, s.remoteAddr)
	return err
}

// Initiate a reply to an earlier decode. The decode message must have started
// with CQ or QRZ.
func (s *Server) Reply(msg ReplyMessage) error {
	msgBytes, _ := encodeReply(msg)
	_, err := s.conn.WriteTo(msgBytes, s.remoteAddr)
	return err
}

// Send a message to WSJT-X to close the program.
func (s *Server) Close(msg CloseMessage) error {
	msgBytes, _ := encodeClose(msg)
	_, err := s.conn.WriteTo(msgBytes, s.remoteAddr)
	return err
}

// Send a message to WSJT-X to replay QSOs in the Band Activity window.
func (s *Server) Replay(msg ReplayMessage) error {
	msgBytes, _ := encodeReplay(msg)
	_, err := s.conn.WriteTo(msgBytes, s.remoteAddr)
	return err
}

// Send a message to WSJT-X to halt transmission.
func (s *Server) HaltTx(msg HaltTxMessage) error {
	msgBytes, _ := encodeHaltTx(msg)
	_, err := s.conn.WriteTo(msgBytes, s.remoteAddr)
	return err
}

// Send a message to WSJT-X to set the free text.
func (s *Server) FreeText(msg FreeTextMessage) error {
	msgBytes, _ := encodeFreeText(msg)
	_, err := s.conn.WriteTo(msgBytes, s.remoteAddr)
	return err
}

// Send a message to WSJT-X to set this station's Maidenhead grid.
func (s *Server) Location(msg LocationMessage) error {
	msgBytes, _ := encodeLocation(msg)
	_, err := s.conn.WriteTo(msgBytes, s.remoteAddr)
	return err
}

// Send a message to WSJT-X to set callsign highlighting.
func (s *Server) HighlightCallsign(msg HighlightCallsignMessage) error {
	msgBytes, _ := encodeHighlightCallsign(msg)
	_, err := s.conn.WriteTo(msgBytes, s.remoteAddr)
	return err
}

// Send a message to WSJT-X to switch to a different pre-defined configuration.
func (s *Server) SwitchConfiguration(msg SwitchConfigurationMessage) error {
	msgBytes, _ := encodeSwitchConfiguration(msg)
	_, err := s.conn.WriteTo(msgBytes, s.remoteAddr)
	return err
}

// Send a message to WSJT-X to change various configuration options.
func (s *Server) Configure(msg ConfigureMessage) error {
	msgBytes, _ := encodeConfigure(msg)
	_, err := s.conn.WriteTo(msgBytes, s.remoteAddr)
	return err
}
