package wsjtx

import (
	"errors"
	"fmt"
	"net"
	"runtime"
)

const magic = 0xadbccbda
const qDataStreamNull = 0xffffffff
const bufLen = 1024
const localhostAddr = "127.0.0.1"
const multicastAddr = "224.0.0.1"
const wsjtxPort = 2237

const qDataStream5_2 uint32 = 2
const qDataStream5_4 uint32 = 3

var supportedSchemas = map[uint32]bool{
	qDataStream5_2: true,
	qDataStream5_4: true,
}

type Server struct {
	ServingAddr net.Addr
	conn        *net.UDPConn
	remoteAddr  *net.UDPAddr
	listening   bool
}

var NotConnectedError = fmt.Errorf("haven't heard from wsjtx yet, don't know where to send commands")

// MakeServer creates a multicast UDP connection to communicate with WSJT-X on the default address
// and port.
func MakeServer() (Server, error) {
	var defaultWsjtxAddr net.IP
	switch runtime.GOOS {
	case "windows":
		defaultWsjtxAddr = net.ParseIP(localhostAddr)
	default:
		defaultWsjtxAddr = net.ParseIP(multicastAddr)
	}
	return MakeServerGiven(defaultWsjtxAddr, wsjtxPort)
}

// MakeServerGiven creates a UDP connection to communicate with WSJT-X on the given address and
// port. Port 0 is allowed, and will cause the OS to assign a port number.
func MakeServerGiven(ipAddr net.IP, port uint) (Server, error) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%v:%d", ipAddr, port))
	if err != nil {
		return Server{}, err
	}
	var conn *net.UDPConn
	if ipAddr.IsMulticast() {
		conn, err = net.ListenMulticastUDP(addr.Network(), nil, addr)
	} else {
		conn, err = net.ListenUDP(addr.Network(), addr)
	}
	if err != nil {
		return Server{}, err
	}
	if conn == nil {
		return Server{}, errors.New("wsjtx udp connection not opened")
	}
	return Server{conn.LocalAddr(), conn, nil, false}, nil
}

func (s *Server) LocalAddr() net.Addr {
	return s.conn.LocalAddr()
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
		if s.conn == nil {
			e <- errors.New("wsjtx connection is nil")
			s.listening = false
			return
		}
		length, rAddr, err := s.conn.ReadFromUDP(b)
		if err != nil {
			e <- fmt.Errorf("problem reading from wsjtx: %w", err)
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
	return s.tryWrite(msgBytes)
}

// Clear sends a message to WSJT-X to clear the band activity window, the RX frequency window, or
// both.
func (s *Server) Clear(msg ClearMessage) error {
	msgBytes, _ := encodeClear(msg)
	return s.tryWrite(msgBytes)
}

// Reply initiates a reply to an earlier decode. The decode message must have started with CQ or
// QRZ.
func (s *Server) Reply(msg ReplyMessage) error {
	msgBytes, _ := encodeReply(msg)
	return s.tryWrite(msgBytes)
}

// Close sends a message to WSJT-X to close the program.
func (s *Server) Close(msg CloseMessage) error {
	msgBytes, _ := encodeClose(msg)
	return s.tryWrite(msgBytes)
}

// Replay sends a message to WSJT-X to replay QSOs in the Band Activity window.
func (s *Server) Replay(msg ReplayMessage) error {
	msgBytes, _ := encodeReplay(msg)
	return s.tryWrite(msgBytes)
}

// HaltTx sends a message to WSJT-X to halt transmission.
func (s *Server) HaltTx(msg HaltTxMessage) error {
	msgBytes, _ := encodeHaltTx(msg)
	return s.tryWrite(msgBytes)
}

// FreeText sends a message to WSJT-X to set the free text of the TX message.
func (s *Server) FreeText(msg FreeTextMessage) error {
	msgBytes, _ := encodeFreeText(msg)
	return s.tryWrite(msgBytes)
}

// Location sends a message to WSJT-X to set this station's Maidenhead grid.
func (s *Server) Location(msg LocationMessage) error {
	msgBytes, _ := encodeLocation(msg)
	return s.tryWrite(msgBytes)
}

// HighlightCallsign sends a message to WSJT-X to set callsign highlighting.
func (s *Server) HighlightCallsign(msg HighlightCallsignMessage) error {
	msgBytes, _ := encodeHighlightCallsign(msg)
	return s.tryWrite(msgBytes)
}

// SwitchConfiguration sends a message to WSJT-X to switch to a different pre-defined configuration.
func (s *Server) SwitchConfiguration(msg SwitchConfigurationMessage) error {
	msgBytes, _ := encodeSwitchConfiguration(msg)
	return s.tryWrite(msgBytes)
}

// Configure sends a message to WSJT-X to change various configuration options.
func (s *Server) Configure(msg ConfigureMessage) error {
	msgBytes, _ := encodeConfigure(msg)
	return s.tryWrite(msgBytes)
}

func (s *Server) tryWrite(msgBytes []byte) error {
	if s.remoteAddr == nil {
		return NotConnectedError
	}
	_, err := s.conn.WriteTo(msgBytes, s.remoteAddr)
	return err
}
