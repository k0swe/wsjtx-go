package wsjtx

import (
	"bytes"
	"encoding/binary"
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

func encodeClose(msg CloseMessage) ([]byte, error) {
	e := newEncoder()
	e.encodeUint32(closeNum)
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

func (e encoder) encodeUint32(num uint32) {
	bin := make([]byte, 4)
	binary.BigEndian.PutUint32(bin, num)
	e.buf.Write(bin)
}

func (e encoder) encodeUint8(num uint8) {
	e.buf.WriteByte(num)
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
