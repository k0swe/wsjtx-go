package wsjtx

import (
	"encoding/binary"
	"github.com/leemcloughlin/jdn"
	"log"
	"math"
	"reflect"
	"time"
)

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
	m := p.parseUint32()
	if m != magic {
		// Packet is not speaking the WSJT-X protocol
		return nil
	}
	sch := p.parseUint32()
	if sch != schema {
		log.Println("Got a schema version I wasn't expecting:", sch)
	}

	messageType := p.parseUint32()
	switch messageType {
	case heartbeatNum:
		heartbeat := p.parseHeartbeat()
		p.checkParse(heartbeat)
		return heartbeat
	case statusNum:
		status := p.parseStatus()
		p.checkParse(status)
		return status
	case decodeNum:
		decode := p.parseDecode()
		p.checkParse(decode)
		return decode
	case clearNum:
		clear := p.parseClear()
		p.checkParse(clear)
		return clear
	case qsoLoggedNum:
		qsoLogged := p.parseQsoLogged()
		p.checkParse(qsoLogged)
		return qsoLogged
	case closeNum:
		closeMsg := p.parseClose()
		p.checkParse(closeMsg)
		return closeMsg
	case wsprDecodeNum:
		wspr := p.parseWsprDecode()
		p.checkParse(wspr)
		return wspr
	case loggedAdifNum:
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
	if strlen == uint32(qDataStreamNull) {
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
