package wsjtx

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"reflect"
	"time"

	"github.com/leemcloughlin/jdn"
)

type parser struct {
	buffer []byte
	length int
	cursor int
}

//nolint:staticcheck // ST1012: exported name kept for API stability (published v4 module)
var ParseError = errors.New("parse error")

//nolint:staticcheck // ST1012: name kept alongside exported ParseError for consistency
var notEnoughBytes = fmt.Errorf("%w: fewer bytes than expected, maybe an older version of WSJTX", ParseError)

// Parse messages following the interface laid out in
// https://sourceforge.net/p/wsjt/wsjtx/ci/master/tree/Network/NetworkMessage.hpp. This only parses
// "Out" or "In/Out" message types and does not include "In" types because they will never be
// received by WSJT-X.
func parseMessage(buffer []byte, length int) (interface{}, error) {
	p := parser{buffer: buffer, length: length, cursor: 0}
	m, err := p.parseUint32()
	if err != nil {
		return nil, ParseError
	}
	if m != magic {
		return nil, fmt.Errorf("%w: packet is not speaking the WSJT-X protocol", ParseError)
	}
	sch, err := p.parseUint32()
	if err != nil {
		return nil, ParseError
	}
	if !supportedSchemas[sch] {
		return nil, fmt.Errorf("%w: got a schema version I wasn't expecting: %d", ParseError, sch)
	}

	messageType, _ := p.parseUint32()
	switch messageType {
	case heartbeatNum:
		heartbeat, err := p.parseHeartbeat()
		if err != nil {
			return heartbeat, err
		}
		err = p.checkParse(heartbeat)
		return heartbeat, err
	case statusNum:
		status, err := p.parseStatus()
		if err != nil {
			return status, err
		}
		err = p.checkParse(status)
		return status, err
	case decodeNum:
		decode, err := p.parseDecode()
		if err != nil {
			return decode, err
		}
		err = p.checkParse(decode)
		return decode, err
	case clearNum:
		clear, err := p.parseClear()
		if err != nil {
			return clear, err
		}
		err = p.checkParse(clear)
		return clear, err
	case qsoLoggedNum:
		qsoLogged, err := p.parseQsoLogged()
		if err != nil {
			return qsoLogged, err
		}
		err = p.checkParse(qsoLogged)
		return qsoLogged, err
	case closeNum:
		closeMsg, err := p.parseClose()
		if err != nil {
			return closeMsg, err
		}
		err = p.checkParse(closeMsg)
		return closeMsg, err
	case wsprDecodeNum:
		wspr, err := p.parseWsprDecode()
		if err != nil {
			return wspr, err
		}
		err = p.checkParse(wspr)
		return wspr, err
	case loggedAdifNum:
		loggedAdif, err := p.parseLoggedAdif()
		if err != nil {
			return loggedAdif, err
		}
		err = p.checkParse(loggedAdif)
		return loggedAdif, err
	}
	return nil, fmt.Errorf("%w: unknown message type %d", ParseError, messageType)
}

// Quick sanity check that we parsed all of the message bytes
func (p *parser) checkParse(message interface{}) error {
	if p.cursor != p.length {
		return fmt.Errorf("%w %s: there were %d bytes left over",
			ParseError, reflect.TypeOf(message).Name(), p.length-p.cursor)
	}
	return nil
}

func (p *parser) parseHeartbeat() (HeartbeatMessage, error) {
	var err error
	heartbeatMessage := HeartbeatMessage{}
	heartbeatMessage.Id, _ = p.parseUtf8()
	heartbeatMessage.MaxSchema, _ = p.parseUint32()
	heartbeatMessage.Version, err = p.parseUtf8()

	// JTDX Packet
	if !p.isDataAvailable() {
		return heartbeatMessage, err
	}
	heartbeatMessage.Revision, err = p.parseUtf8()
	return heartbeatMessage, err
}

func (p *parser) parseStatus() (StatusMessage, error) {
	var err error
	statusMessage := StatusMessage{}
	statusMessage.Id, _ = p.parseUtf8()
	statusMessage.DialFrequency, _ = p.parseUint64()
	statusMessage.Mode, _ = p.parseUtf8()
	statusMessage.DxCall, _ = p.parseUtf8()
	statusMessage.Report, _ = p.parseUtf8()
	statusMessage.TxMode, _ = p.parseUtf8()
	statusMessage.TxEnabled, _ = p.parseBool()
	statusMessage.Transmitting, _ = p.parseBool()
	statusMessage.Decoding, _ = p.parseBool()
	statusMessage.RxDF, _ = p.parseUint32()
	statusMessage.TxDF, _ = p.parseUint32()
	statusMessage.DeCall, _ = p.parseUtf8()
	statusMessage.DeGrid, _ = p.parseUtf8()
	statusMessage.DxGrid, _ = p.parseUtf8()
	statusMessage.TxWatchdog, _ = p.parseBool()
	statusMessage.SubMode, _ = p.parseUtf8()
	statusMessage.FastMode, err = p.parseBool()

	// A JTDX Format packet (Tx first, bool)
	if (p.length - p.cursor) == 1 {
		_, _ = p.parseBool()
		return statusMessage, err
	}
	statusMessage.SpecialOperationMode, _ = p.parseUint8()
	statusMessage.FrequencyTolerance, _ = p.parseUint32()
	statusMessage.TRPeriod, _ = p.parseUint32()
	statusMessage.ConfigurationName, _ = p.parseUtf8()
	statusMessage.TxMessage, err = p.parseUtf8()
	return statusMessage, err
}

func (p *parser) parseDecode() (DecodeMessage, error) {
	var err error
	decodeMessage := DecodeMessage{}
	decodeMessage.Id, _ = p.parseUtf8()
	decodeMessage.New, _ = p.parseBool()
	decodeMessage.Time, _ = p.parseUint32()
	decodeMessage.Snr, _ = p.parseInt32()
	decodeMessage.DeltaTimeSec, _ = p.parseFloat64()
	decodeMessage.DeltaFrequencyHz, _ = p.parseUint32()
	decodeMessage.Mode, _ = p.parseUtf8()
	decodeMessage.Message, _ = p.parseUtf8()
	decodeMessage.LowConfidence, _ = p.parseBool()
	decodeMessage.OffAir, err = p.parseBool()
	return decodeMessage, err
}

func (p *parser) parseClear() (ClearMessage, error) {
	var err error
	clearMessage := ClearMessage{}
	clearMessage.Id, err = p.parseUtf8()
	return clearMessage, err
}

func (p *parser) parseQsoLogged() (QsoLoggedMessage, error) {
	var err error
	qsoLoggedMessage := QsoLoggedMessage{}
	qsoLoggedMessage.Id, _ = p.parseUtf8()
	qsoLoggedMessage.DateTimeOff, _ = p.parseQDateTime()
	qsoLoggedMessage.DxCall, _ = p.parseUtf8()
	qsoLoggedMessage.DxGrid, _ = p.parseUtf8()
	qsoLoggedMessage.TxFrequency, _ = p.parseUint64()
	qsoLoggedMessage.Mode, _ = p.parseUtf8()
	qsoLoggedMessage.ReportSent, _ = p.parseUtf8()
	qsoLoggedMessage.ReportReceived, _ = p.parseUtf8()
	qsoLoggedMessage.TxPower, _ = p.parseUtf8()
	qsoLoggedMessage.Comments, _ = p.parseUtf8()
	qsoLoggedMessage.Name, _ = p.parseUtf8()
	qsoLoggedMessage.DateTimeOn, _ = p.parseQDateTime()
	qsoLoggedMessage.OperatorCall, _ = p.parseUtf8()
	qsoLoggedMessage.MyCall, _ = p.parseUtf8()
	qsoLoggedMessage.MyGrid, err = p.parseUtf8()

	// JTDX Packet
	if !p.isDataAvailable() {
		return qsoLoggedMessage, err
	}
	qsoLoggedMessage.ExchangeSent, _ = p.parseUtf8()
	qsoLoggedMessage.ExchangeReceived, err = p.parseUtf8()

	// Older WSJT-X packet doesn't have Propagation Mode
	if p.isDataAvailable() {
		qsoLoggedMessage.ADIFPropagationMode, err = p.parseUtf8()
	}
	return qsoLoggedMessage, err
}

func (p *parser) parseClose() (CloseMessage, error) {
	var err error
	closeMessage := CloseMessage{}
	closeMessage.Id, err = p.parseUtf8()
	return closeMessage, err
}

func (p *parser) parseWsprDecode() (WSPRDecodeMessage, error) {
	var err error
	wsprDecodeMessage := WSPRDecodeMessage{}
	wsprDecodeMessage.Id, _ = p.parseUtf8()
	wsprDecodeMessage.New, _ = p.parseBool()
	wsprDecodeMessage.Time, _ = p.parseUint32()
	wsprDecodeMessage.Snr, _ = p.parseInt32()
	wsprDecodeMessage.DeltaTime, _ = p.parseFloat64()
	wsprDecodeMessage.Frequency, _ = p.parseUint64()
	wsprDecodeMessage.Drift, _ = p.parseInt32()
	wsprDecodeMessage.Callsign, _ = p.parseUtf8()
	wsprDecodeMessage.Grid, _ = p.parseUtf8()
	wsprDecodeMessage.Power, err = p.parseInt32()

	// JTDX Packet
	if !p.isDataAvailable() {
		return wsprDecodeMessage, err
	}
	wsprDecodeMessage.OffAir, err = p.parseBool()
	return wsprDecodeMessage, err
}

func (p *parser) parseLoggedAdif() (LoggedAdifMessage, error) {
	var err error
	loggedAdifMessage := LoggedAdifMessage{}
	loggedAdifMessage.Id, _ = p.parseUtf8()
	loggedAdifMessage.Adif, err = p.parseUtf8()
	return loggedAdifMessage, err
}

func (p *parser) parseUint8() (uint8, error) {
	if len(p.buffer) < p.cursor {
		return 0, notEnoughBytes
	}
	value := p.buffer[p.cursor]
	p.cursor += 1
	return value, nil
}

func (p *parser) parseUtf8() (string, error) {
	strlen, err := p.parseUint32()
	if err != nil {
		return "", err
	}
	if strlen == uint32(qDataStreamNull) {
		// this is a sentinel value meaning "null" in QDataStream, but Golang can't have nil strings
		strlen = 0
	}
	end := p.cursor + int(strlen)
	value := string(p.buffer[p.cursor:end])
	p.cursor += int(strlen)
	return value, nil
}

func (p *parser) parseUint32() (uint32, error) {
	end := p.cursor + 4
	if len(p.buffer) < end {
		return 0, notEnoughBytes
	}
	value := binary.BigEndian.Uint32(p.buffer[p.cursor:end])
	p.cursor += 4
	return value, nil
}

func (p *parser) parseInt32() (int32, error) {
	end := p.cursor + 4
	if len(p.buffer) < end {
		return 0, notEnoughBytes
	}
	value := int32(binary.BigEndian.Uint32(p.buffer[p.cursor:end]))
	p.cursor += 4
	return value, nil
}

func (p *parser) parseUint64() (uint64, error) {
	end := p.cursor + 8
	if len(p.buffer) < end {
		return 0, notEnoughBytes
	}
	value := binary.BigEndian.Uint64(p.buffer[p.cursor:end])
	p.cursor += 8
	return value, nil
}

func (p *parser) parseFloat64() (float64, error) {
	end := p.cursor + 8
	if len(p.buffer) < end {
		return 0, notEnoughBytes
	}
	bits := binary.BigEndian.Uint64(p.buffer[p.cursor:end])
	value := math.Float64frombits(bits)
	p.cursor += 8
	return value, nil
}

func (p *parser) parseBool() (bool, error) {
	if len(p.buffer) < p.cursor {
		return false, notEnoughBytes
	}
	value := p.buffer[p.cursor] != 0
	p.cursor += 1
	return value, nil
}

func (p *parser) parseQDateTime() (time.Time, error) {
	julianDay, _ := p.parseUint64()
	year, month, day := jdn.FromNumber(int(julianDay))
	msMid, _ := p.parseUint32()
	msSinceMidnight := int(msMid)
	hour := msSinceMidnight / 3600000
	msSinceMidnight -= hour * 3600000
	minute := msSinceMidnight / 60000
	msSinceMidnight -= minute * 60000
	second := msSinceMidnight / 1000
	timespec, err := p.parseUint8()
	var value time.Time
	switch timespec {
	case 0:
		// local
		value = time.Date(year, month, day, hour, minute, second, 0, time.Local)
	case 1:
		// UTC
		value = time.Date(year, month, day, hour, minute, second, 0, time.UTC)
	default:
		return value, fmt.Errorf("got a timespec I wasn't expecting: %d", timespec)
	}
	return value, err
}

func (p *parser) isDataAvailable() bool {
	return p.cursor < p.length
}
