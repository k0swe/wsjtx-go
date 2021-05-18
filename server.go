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
	listening  bool
}

// MakeServer creates a multicast UDP connection to communicate with WSJT-X on the default address
// and port.
func MakeServer() (Server, error) {
	return MakeMulticastServer(multicastAddr, wsjtxPort)
}

// MakeMulticastServer creates a multicast UDP connection to communicate with WSJT-X on the given
// address and port.
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
	return Server{conn, nil, false}, nil
}

// MakeUnicastServer creates a unicast UDP connection to communicate with WSJT-X on the given
// address and port.
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
	return Server{conn, nil, false}, nil
}

// ListenToWsjtx listens for messages from WSJT-X. When heard, the messages are parsed and then
// placed in the given message channel. If parsing errors occur, those are reported on the errors
// channel. If a fatal error happens, e.g. the network connection gets closed, the channels are
// closed and the goroutine ends.
func (s *Server) ListenToWsjtx(c chan interface{}, e chan error) {
	s.listening = true
	defer close(c)
	defer close(e)

	for {
		b := make([]byte, bufLen)
		length, rAddr, err := s.conn.ReadFromUDP(b)
		if err != nil {
			e <- err
			s.listening = false
			return
		}
		s.remoteAddr = rAddr
		message, err := parseMessage(b, length)
		if err != nil {
			e <- err
		}
		if message != nil {
			c <- message
		}
	}
}

// Listening returns whether the ListenToWsjtx goroutine is currently running.
func (s *Server) Listening() bool {
	return s.listening
}

// Heartbeat sends a heartbeat message to WSJT-X.
func (s *Server) Heartbeat(msg HeartbeatMessage) error {
	msgBytes, _ := encodeHeartbeat(msg)
	_, err := s.conn.WriteTo(msgBytes, s.remoteAddr)
	return err
}

// Clear sends a message to WSJT-X to clear the band activity window, the RX frequency window, or
// both.
func (s *Server) Clear(msg ClearMessage) error {
	msgBytes, _ := encodeClear(msg)
	_, err := s.conn.WriteTo(msgBytes, s.remoteAddr)
	return err
}

// Reply initiates a reply to an earlier decode. The decode message must have started with CQ or
// QRZ.
func (s *Server) Reply(msg ReplyMessage) error {
	msgBytes, _ := encodeReply(msg)
	_, err := s.conn.WriteTo(msgBytes, s.remoteAddr)
	return err
}

// Close sends a message to WSJT-X to close the program.
func (s *Server) Close(msg CloseMessage) error {
	msgBytes, _ := encodeClose(msg)
	_, err := s.conn.WriteTo(msgBytes, s.remoteAddr)
	return err
}

// Replay sends a message to WSJT-X to replay QSOs in the Band Activity window.
func (s *Server) Replay(msg ReplayMessage) error {
	msgBytes, _ := encodeReplay(msg)
	_, err := s.conn.WriteTo(msgBytes, s.remoteAddr)
	return err
}

// HaltTx sends a message to WSJT-X to halt transmission.
func (s *Server) HaltTx(msg HaltTxMessage) error {
	msgBytes, _ := encodeHaltTx(msg)
	_, err := s.conn.WriteTo(msgBytes, s.remoteAddr)
	return err
}

// FreeText sends a message to WSJT-X to set the free text of the TX message.
func (s *Server) FreeText(msg FreeTextMessage) error {
	msgBytes, _ := encodeFreeText(msg)
	_, err := s.conn.WriteTo(msgBytes, s.remoteAddr)
	return err
}

// Location sends a message to WSJT-X to set this station's Maidenhead grid.
func (s *Server) Location(msg LocationMessage) error {
	msgBytes, _ := encodeLocation(msg)
	_, err := s.conn.WriteTo(msgBytes, s.remoteAddr)
	return err
}

// HighlightCallsign sends a message to WSJT-X to set callsign highlighting.
func (s *Server) HighlightCallsign(msg HighlightCallsignMessage) error {
	msgBytes, _ := encodeHighlightCallsign(msg)
	_, err := s.conn.WriteTo(msgBytes, s.remoteAddr)
	return err
}

// SwitchConfiguration sends a message to WSJT-X to switch to a different pre-defined configuration.
func (s *Server) SwitchConfiguration(msg SwitchConfigurationMessage) error {
	msgBytes, _ := encodeSwitchConfiguration(msg)
	_, err := s.conn.WriteTo(msgBytes, s.remoteAddr)
	return err
}

// Configure sends a message to WSJT-X to change various configuration options.
func (s *Server) Configure(msg ConfigureMessage) error {
	msgBytes, _ := encodeConfigure(msg)
	_, err := s.conn.WriteTo(msgBytes, s.remoteAddr)
	return err
}
