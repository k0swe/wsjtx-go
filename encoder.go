package wsjtx

import (
	"bytes"
	"encoding/binary"
	"math"
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

type encoder struct {
	buf *bytes.Buffer
}

func newEncoder() encoder {
	e := encoder{bytes.NewBuffer(make([]byte, bufLen))}
	e.buf.Reset()
	e.encodeUint32(magic)
	e.encodeUint32(schema)
	return e
}

func (e encoder) encodeUint8(num uint8) {
	e.buf.WriteByte(num)
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

func (e encoder) finish() ([]byte, error) {
	ret := make([]byte, e.buf.Len())
	_, err := e.buf.Read(ret)
	return ret, err
}
