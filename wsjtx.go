package wsjtx

import (
	"encoding/binary"
	"github.com/leemcloughlin/jdn"
	"log"
	"math"
	"net"
	"reflect"
	"time"
)

/*
The heartbeat  message shall be  sent on a periodic  basis every
15   seconds.  This
message is intended to be used by servers to detect the presence
of a  client and also  the unexpected disappearance of  a client
and  by clients  to learn  the schema  negotiated by  the server
after it receives  the initial heartbeat message  from a client.

Out/In.

https://sourceforge.net/p/wsjt/wsjtx/ci/8f99fcce/tree/Network/NetworkMessage.hpp#l110
*/
type HeartbeatMessage struct {
	Id        string `json:"id"`
	MaxSchema uint32 `json:"maxSchemaVersion"`
	Version   string `json:"version"`
	Revision  string `json:"revision"`
}

/*
WSJT-X  sends this  status message  when various  internal state
changes to allow the server to  track the relevant state of each
client without the need for  polling commands.

Out only.

https://sourceforge.net/p/wsjt/wsjtx/ci/8f99fcce/tree/Network/NetworkMessage.hpp#l141
*/
type StatusMessage struct {
	Id                   string `json:"id"`
	DialFrequency        uint64 `json:"dialFrequency"`
	Mode                 string `json:"mode"`
	DxCall               string `json:"dxCall"`
	Report               string `json:"report"`
	TxMode               string `json:"txMode"`
	TxEnabled            bool   `json:"txEnabled"`
	Transmitting         bool   `json:"transmitting"`
	Decoding             bool   `json:"decoding"`
	RxDF                 uint32 `json:"rxDeltaFreq"`
	TxDF                 uint32 `json:"txDeltaFreq"`
	DeCall               string `json:"deCall"`
	DeGrid               string `json:"deGrid"`
	DxGrid               string `json:"dxGrid"`
	TxWatchdog           bool   `json:"txWatchdog"`
	SubMode              string `json:"submode"`
	FastMode             bool   `json:"fastMode"`
	SpecialOperationMode uint8  `json:"specialMode"`
	FrequencyTolerance   uint32 `json:"frequencyTolerance"`
	TRPeriod             uint32 `json:"txRxPeriod"`
	ConfigurationName    string `json:"configName"`
}

/*
The decode message is sent when  a new decode is completed, in
this case the 'New' field is true. It is also used in response
to  a "Replay"  message where  each  old decode  in the  "Band
activity" window, that  has not been erased, is  sent in order
as a one of these messages  with the 'New' field set to false.

Out only.

https://sourceforge.net/p/wsjt/wsjtx/ci/8f99fcce/tree/Network/NetworkMessage.hpp#l206
*/
type DecodeMessage struct {
	Id               string  `json:"id"`
	New              bool    `json:"new"`
	Time             uint32  `json:"time"`
	Snr              int32   `json:"snr"`
	DeltaTimeSec     float64 `json:"deltaTime"`
	DeltaFrequencyHz uint32  `json:"deltaFrequency"`
	Mode             string  `json:"mode"`
	Message          string  `json:"message"`
	LowConfidence    bool    `json:"lowConfidence"`
	OffAir           bool    `json:"offAir"`
}

/*
This message is  send when all prior "Decode"  messages in the
"Band Activity"  window have been discarded  and therefore are
no long available for actioning  with a "Reply" message.

The Window  argument  can be  one  of the  following values:

	0  - clear the "Band Activity" window (default)
	1  - clear the "Rx Frequency" window
	2  - clear both "Band Activity" and "Rx Frequency" windows

Out/In.

https://sourceforge.net/p/wsjt/wsjtx/ci/8f99fcce/tree/Network/NetworkMessage.hpp#l232
*/
type ClearMessage struct {
	Id     string `json:"id"`
	Window uint8  `json:"window"` // In only
}

/*
The QSO logged message is sent when the WSJT-X user accepts the "Log  QSO" dialog by clicking
the "OK" button.

Out only.

https://sourceforge.net/p/wsjt/wsjtx/ci/8f99fcce/tree/Network/NetworkMessage.hpp#l293
*/
type QsoLoggedMessage struct {
	Id               string    `json:"id"`
	DateTimeOff      time.Time `json:"dateTimeOff"`
	DxCall           string    `json:"dxCall"`
	DxGrid           string    `json:"dxGrid"`
	TxFrequency      uint64    `json:"txFrequency"`
	Mode             string    `json:"mode"`
	ReportSent       string    `json:"reportSent"`
	ReportReceived   string    `json:"reportReceived"`
	TxPower          string    `json:"txPower"`
	Comments         string    `json:"comments"`
	Name             string    `json:"name"`
	DateTimeOn       time.Time `json:"dateTimeOn"`
	OperatorCall     string    `json:"operatorCall"`
	MyCall           string    `json:"myCall"`
	MyGrid           string    `json:"myGrid"`
	ExchangeSent     string    `json:"exchangeSent"`
	ExchangeReceived string    `json:"exchangeReceived"`
}

/*
Close is  sent by  a client immediately  prior to  it shutting
down gracefully.

Out/In.

https://sourceforge.net/p/wsjt/wsjtx/ci/8f99fcce/tree/Network/NetworkMessage.hpp#l318
*/
type CloseMessage struct {
	Id string `json:"id"`
}

/*
The decode message is sent when  a new decode is completed, in
this case the 'New' field is true.

Out only.

https://sourceforge.net/p/wsjt/wsjtx/ci/8f99fcce/tree/Network/NetworkMessage.hpp#l381).
*/
type WSPRDecodeMessage struct {
	Id        string  `json:"id"`
	New       bool    `json:"new"`
	Time      uint32  `json:"time"`
	Snr       int32   `json:"snr"`
	DeltaTime float64 `json:"deltaTime"`
	Frequency uint64  `json:"frequency"`
	Drift     int32   `json:"drift"`
	Callsign  string  `json:"callsign"`
	Grid      string  `json:"grid"`
	Power     int32   `json:"power"`
	OffAir    bool    `json:"offAir"`
}

/*
The  logged ADIF  message is  sent to  the server(s)  when the
WSJT-X user accepts the "Log  QSO" dialog by clicking the "OK"
button.

Out only.

https://sourceforge.net/p/wsjt/wsjtx/ci/8f99fcce/tree/Network/NetworkMessage.hpp#l421
*/
type LoggedAdifMessage struct {
	Id   string `json:"id"`
	Adif string `json:"adif"`
}

const Magic = 0xadbccbda
const BufLen = 1024

type Server struct {
	conn *net.UDPConn
}

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

func (s *Server) Clear() {
	// TODO
	//s.conn.Write();
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

type parser struct {
	buffer []byte
	length int
	cursor int
}

// Parse messages following the interface laid out in
// https://sourceforge.net/p/wsjt/wsjtx/ci/master/tree/Network/NetworkMessage.hpp. This only parses
// "Out" or "In/Out" message types and does not include "In" types because they will never be
// received by WSJT-X.
func parseMessage(buffer []byte, length int) interface{} {
	p := parser{buffer: buffer, length: length, cursor: 0}
	magic := p.parseUint32()
	if magic != Magic {
		// Packet is not speaking the WSJT-X protocol
		return nil
	}
	schema := p.parseUint32()
	if schema != 2 {
		log.Println("Got a schema version I wasn't expecting:", schema)
	}

	messageType := p.parseUint32()
	switch messageType {
	case 0:
		heartbeat := p.parseHeartbeat()
		p.checkParse(heartbeat)
		return heartbeat
	case 1:
		status := p.parseStatus()
		p.checkParse(status)
		return status
	case 2:
		decode := p.parseDecode()
		p.checkParse(decode)
		return decode
	case 3:
		clear := p.parseClear()
		p.checkParse(clear)
		return clear
	case 5:
		qsoLogged := p.parseQsoLogged()
		p.checkParse(qsoLogged)
		return qsoLogged
	case 6:
		closeMsg := p.parseClose()
		p.checkParse(closeMsg)
		return closeMsg
	case 10:
		wspr := p.parseWsprDecode()
		p.checkParse(wspr)
		return wspr
	case 12:
		loggedAdif := p.parseLoggedAdif()
		p.checkParse(loggedAdif)
		return loggedAdif
	}
	return nil
}

// Quick sanity check that we parsed all of the message bytes
func (p *parser) checkParse(message interface{}) {
	if p.cursor != p.length {
		log.Fatalf("Parsing WSJT-X %s: There were %d bytes left over\n",
			reflect.TypeOf(message).Name(), p.length-p.cursor)
	}
}

func (p *parser) parseHeartbeat() HeartbeatMessage {
	id := p.parseUtf8()
	maxSchema := p.parseUint32()
	version := p.parseUtf8()
	revision := p.parseUtf8()
	return HeartbeatMessage{
		Id:        id,
		MaxSchema: maxSchema,
		Version:   version,
		Revision:  revision,
	}
}

func (p *parser) parseStatus() StatusMessage {
	id := p.parseUtf8()
	dialFreq := p.parseUint64()
	mode := p.parseUtf8()
	dxCall := p.parseUtf8()
	report := p.parseUtf8()
	txMode := p.parseUtf8()
	txEnabled := p.parseBool()
	transmitting := p.parseBool()
	decoding := p.parseBool()
	rxDf := p.parseUint32()
	txDf := p.parseUint32()
	deCall := p.parseUtf8()
	deGrid := p.parseUtf8()
	dxGrid := p.parseUtf8()
	txWatchdog := p.parseBool()
	subMode := p.parseUtf8()
	fastMode := p.parseBool()
	specialMode := p.parseUint8()
	freqTolerance := p.parseUint32()
	trPeriod := p.parseUint32()
	configName := p.parseUtf8()
	return StatusMessage{
		Id:                   id,
		DialFrequency:        dialFreq,
		Mode:                 mode,
		DxCall:               dxCall,
		Report:               report,
		TxMode:               txMode,
		TxEnabled:            txEnabled,
		Transmitting:         transmitting,
		Decoding:             decoding,
		RxDF:                 rxDf,
		TxDF:                 txDf,
		DeCall:               deCall,
		DeGrid:               deGrid,
		DxGrid:               dxGrid,
		TxWatchdog:           txWatchdog,
		SubMode:              subMode,
		FastMode:             fastMode,
		SpecialOperationMode: specialMode,
		FrequencyTolerance:   freqTolerance,
		TRPeriod:             trPeriod,
		ConfigurationName:    configName,
	}
}

func (p *parser) parseDecode() DecodeMessage {
	id := p.parseUtf8()
	newDecode := p.parseBool()
	decodeTime := p.parseUint32()
	snr := p.parseInt32()
	deltaTime := p.parseFloat64()
	deltaFreq := p.parseUint32()
	mode := p.parseUtf8()
	message := p.parseUtf8()
	lowConfidence := p.parseBool()
	offAir := p.parseBool()
	return DecodeMessage{
		Id:               id,
		New:              newDecode,
		Time:             decodeTime,
		Snr:              snr,
		DeltaTimeSec:     deltaTime,
		DeltaFrequencyHz: deltaFreq,
		Mode:             mode,
		Message:          message,
		LowConfidence:    lowConfidence,
		OffAir:           offAir,
	}
}

func (p *parser) parseClear() ClearMessage {
	id := p.parseUtf8()
	return ClearMessage{
		Id: id,
	}
}

func (p *parser) parseQsoLogged() QsoLoggedMessage {
	id := p.parseUtf8()
	timeOff := p.parseQDateTime()
	dxCall := p.parseUtf8()
	dxGrid := p.parseUtf8()
	txFrequency := p.parseUint64()
	mode := p.parseUtf8()
	reportSent := p.parseUtf8()
	reportReceived := p.parseUtf8()
	txPower := p.parseUtf8()
	comments := p.parseUtf8()
	name := p.parseUtf8()
	timeOn := p.parseQDateTime()
	operatorCall := p.parseUtf8()
	myCall := p.parseUtf8()
	myGrid := p.parseUtf8()
	exchangeSent := p.parseUtf8()
	exchangeReceived := p.parseUtf8()
	return QsoLoggedMessage{
		Id:               id,
		DateTimeOff:      timeOff,
		DxCall:           dxCall,
		DxGrid:           dxGrid,
		TxFrequency:      txFrequency,
		Mode:             mode,
		ReportSent:       reportSent,
		ReportReceived:   reportReceived,
		TxPower:          txPower,
		Comments:         comments,
		Name:             name,
		DateTimeOn:       timeOn,
		OperatorCall:     operatorCall,
		MyCall:           myCall,
		MyGrid:           myGrid,
		ExchangeSent:     exchangeSent,
		ExchangeReceived: exchangeReceived,
	}
}

func (p *parser) parseClose() interface{} {
	id := p.parseUtf8()
	return CloseMessage{
		Id: id,
	}
}

func (p *parser) parseWsprDecode() interface{} {
	id := p.parseUtf8()
	newDecode := p.parseBool()
	decodeTime := p.parseUint32()
	snr := p.parseInt32()
	deltaTime := p.parseFloat64()
	frequency := p.parseUint64()
	drift := p.parseInt32()
	callsign := p.parseUtf8()
	grid := p.parseUtf8()
	power := p.parseInt32()
	offAir := p.parseBool()
	return WSPRDecodeMessage{
		Id:        id,
		New:       newDecode,
		Time:      decodeTime,
		Snr:       snr,
		DeltaTime: deltaTime,
		Frequency: frequency,
		Drift:     drift,
		Callsign:  callsign,
		Grid:      grid,
		Power:     power,
		OffAir:    offAir,
	}
}

func (p *parser) parseLoggedAdif() interface{} {
	id := p.parseUtf8()
	adif := p.parseUtf8()
	return LoggedAdifMessage{
		Id:   id,
		Adif: adif,
	}
}

func (p *parser) parseUint8() uint8 {
	value := p.buffer[p.cursor]
	p.cursor += 1
	return value
}

func (p *parser) parseUtf8() string {
	strlen := p.parseUint32()
	if strlen == uint32(0xffffffff) {
		// this is a sentinel value meaning "null" in QDataStream, but Golang can't have nil strings
		strlen = 0
	}
	value := string(p.buffer[p.cursor:(p.cursor + int(strlen))])
	p.cursor += int(strlen)
	return value
}

func (p *parser) parseUint32() uint32 {
	value := binary.BigEndian.Uint32(p.buffer[p.cursor : p.cursor+4])
	p.cursor += 4
	return value
}

func (p *parser) parseInt32() int32 {
	value := int32(binary.BigEndian.Uint32(p.buffer[p.cursor : p.cursor+4]))
	p.cursor += 4
	return value
}

func (p *parser) parseUint64() uint64 {
	value := binary.BigEndian.Uint64(p.buffer[p.cursor : p.cursor+8])
	p.cursor += 8
	return value
}

func (p *parser) parseFloat64() float64 {
	bits := binary.BigEndian.Uint64(p.buffer[p.cursor : p.cursor+8])
	value := math.Float64frombits(bits)
	p.cursor += 8
	return value
}

func (p *parser) parseBool() bool {
	value := p.buffer[p.cursor] != 0
	p.cursor += 1
	return value
}

func (p *parser) parseQDateTime() time.Time {
	julianDay := p.parseUint64()
	year, month, day := jdn.FromNumber(int(julianDay))
	msSinceMidnight := int(p.parseUint32())
	hour := msSinceMidnight / 3600000
	msSinceMidnight -= hour * 3600000
	minute := msSinceMidnight / 60000
	msSinceMidnight -= minute * 60000
	second := msSinceMidnight / 1000
	timespec := p.parseUint8()
	var value time.Time
	switch timespec {
	case 0:
		// local
		value = time.Date(year, month, day, hour, minute, second, 0, time.Local)
	case 1:
		// UTC
		value = time.Date(year, month, day, hour, minute, second, 0, time.UTC)
	default:
		log.Fatalln("WSJT-X parser: Got a timespec I wasn't expecting,", timespec)
	}
	return value
}
