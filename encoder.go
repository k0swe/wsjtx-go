package wsjtx

import (
	"bytes"
	"encoding/binary"
	"math"

	"github.com/mazznoer/csscolorparser"
)

func encodeHeartbeat(msg HeartbeatMessage) ([]byte, error) {
	e := newEncoder()
	e.encodeUint32(heartbeatNum)
	e.encodeUtf8(msg.Id)
	e.encodeUint32(msg.MaxSchema)
	e.encodeUtf8(msg.Version)
	e.encodeUtf8(msg.Revision)
	return e.finish()
}

func encodeClear(msg ClearMessage) ([]byte, error) {
	e := newEncoder()
	e.encodeUint32(clearNum)
	e.encodeUtf8(msg.Id)
	e.encodeUint8(msg.Window)
	return e.finish()
}

func encodeReply(msg ReplyMessage) ([]byte, error) {
	e := newEncoder()
	e.encodeUint32(replyNum)
	e.encodeUtf8(msg.Id)
	e.encodeUint32(msg.Time)
	e.encodeInt32(msg.Snr)
	e.encodeFloat64(msg.DeltaTimeSec)
	e.encodeUint32(msg.DeltaFrequencyHz)
	e.encodeUtf8(msg.Mode)
	e.encodeUtf8(msg.Message)
	e.encodeBool(msg.LowConfidence)
	e.encodeUint8(msg.Modifiers)
	return e.finish()
}

func encodeClose(msg CloseMessage) ([]byte, error) {
	e := newEncoder()
	e.encodeUint32(closeNum)
	e.encodeUtf8(msg.Id)
	return e.finish()
}

func encodeReplay(msg ReplayMessage) ([]byte, error) {
	e := newEncoder()
	e.encodeUint32(replayNum)
	e.encodeUtf8(msg.Id)
	return e.finish()
}

func encodeHaltTx(msg HaltTxMessage) ([]byte, error) {
	e := newEncoder()
	e.encodeUint32(haltTxNum)
	e.encodeUtf8(msg.Id)
	e.encodeBool(msg.AutoTxOnly)
	return e.finish()
}

func encodeFreeText(msg FreeTextMessage) ([]byte, error) {
	e := newEncoder()
	e.encodeUint32(freeTextNum)
	e.encodeUtf8(msg.Id)
	e.encodeUtf8(msg.Text)
	e.encodeBool(msg.Send)
	return e.finish()
}

func encodeLocation(msg LocationMessage) ([]byte, error) {
	e := newEncoder()
	e.encodeUint32(locationNum)
	e.encodeUtf8(msg.Id)
	e.encodeUtf8(msg.Location)
	return e.finish()
}

func encodeHighlightCallsign(msg HighlightCallsignMessage) ([]byte, error) {
	e := newEncoder()
	e.encodeUint32(highlightCallsignNum)
	e.encodeUtf8(msg.Id)
	e.encodeUtf8(msg.Callsign)
	if err := e.encodeColor(msg.BackgroundColor, msg.Reset); err != nil {
		return []byte{}, err
	}
	if err := e.encodeColor(msg.ForegroundColor, msg.Reset); err != nil {
		return []byte{}, err
	}
	e.encodeBool(msg.HighlightLast)
	return e.finish()
}

func encodeSwitchConfiguration(msg SwitchConfigurationMessage) ([]byte, error) {
	e := newEncoder()
	e.encodeUint32(switchConfigurationNum)
	e.encodeUtf8(msg.Id)
	e.encodeUtf8(msg.ConfigurationName)
	return e.finish()
}

func encodeConfigure(msg ConfigureMessage) ([]byte, error) {
	e := newEncoder()
	e.encodeUint32(configureNum)
	e.encodeUtf8(msg.Id)
	e.encodeUtf8(msg.Mode)
	e.encodeUint32(msg.FrequencyTolerance)
	e.encodeUtf8(msg.Submode)
	e.encodeBool(msg.FastMode)
	e.encodeUint32(msg.TRPeriod)
	e.encodeUint32(msg.RxDF)
	e.encodeUtf8(msg.DXCall)
	e.encodeUtf8(msg.DXGrid)
	e.encodeBool(msg.GenerateMessages)
	return e.finish()
}

type encoder struct {
	buf *bytes.Buffer
}

func newEncoder() encoder {
	e := encoder{bytes.NewBuffer(make([]byte, bufLen))}
	e.buf.Reset()
	e.encodeUint32(magic)
	e.encodeUint32(qDataStream5_2)
	return e
}

func (e encoder) encodeUint8(num uint8) {
	e.buf.WriteByte(num)
}

func (e encoder) encodeUint16(num uint16) {
	bin := make([]byte, 2)
	binary.BigEndian.PutUint16(bin, num)
	e.buf.Write(bin)
}

func (e encoder) encodeUint32(num uint32) {
	bin := make([]byte, 4)
	binary.BigEndian.PutUint32(bin, num)
	e.buf.Write(bin)
}

func (e encoder) encodeUint64(num uint64) {
	bin := make([]byte, 8)
	binary.BigEndian.PutUint64(bin, num)
	e.buf.Write(bin)
}

func (e encoder) encodeBool(b bool) {
	if b {
		e.encodeUint8(1)
	} else {
		e.encodeUint8(0)
	}

}

func (e encoder) encodeInt32(num int32) {
	e.encodeUint32(uint32(num))
}

func (e encoder) encodeFloat64(num float64) {
	e.encodeUint64(math.Float64bits(num))
}

func (e encoder) encodeUtf8(str string) {
	strlen := uint32(len(str))
	if strlen == 0 {
		e.encodeUint32(qDataStreamNull)
		return
	}
	e.encodeUint32(strlen)
	e.buf.WriteString(str)
}

func (e encoder) encodeColor(color string, invalid bool) error {
	// Spec enum: https://github.com/radekp/qt/blob/b881d8fb/src/gui/painting/qcolor.h#L70
	const invalidSpec = uint8(0)
	const rgbSpec = uint8(1)
	const pad = uint16(0)

	spec := rgbSpec
	if invalid {
		spec = invalidSpec
	}

	// pre-multiplied to range 0x0 to 0xffff
	c, err := csscolorparser.Parse(color)
	if err != nil {
		return err
	}
	r, g, b, a := c.RGBA()

	// Field type and order: https://github.com/radekp/qt/blob/b881d8fb/src/gui/painting/qcolor.cpp#L2506
	e.encodeUint8(spec)
	e.encodeUint16(uint16(a))
	e.encodeUint16(uint16(r))
	e.encodeUint16(uint16(g))
	e.encodeUint16(uint16(b))
	e.encodeUint16(pad)
	return nil
}

func (e encoder) finish() ([]byte, error) {
	ret := make([]byte, e.buf.Len())
	_, err := e.buf.Read(ret)
	return ret, err
}
