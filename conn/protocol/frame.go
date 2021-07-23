package protocol

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"sync"
	"time"

	"github.com/eggz6/common/conn"
)

var framePool sync.Pool

func init() {
	framePool.New = func() interface{} {
		return &Frame{}
	}
}

type Frame struct {
	Header
	Data []byte
}

func (f *Frame)BodyLength()uint64{
	return f.bodyLength
}

const ver uint8 = 1

func EncodeFrame(protoc conn.Protoc,
	channel uint32,
	from uint32,
	start time.Time,
	data []byte) []byte {
	l := len(data)

	h := AllocHeader(
		WithVer(ver),
		WithDuration(uint32(time.Since(start).Milliseconds())),
		WithChannel(channel),
		WithFrom(from),
		SwitchProtoc(protoc),
		WithBodyLength(uint64(l)),
	)

	res := make([]byte,0)
	buf := bytes.NewBuffer(res)
	b := EncodeHeader(h)
	buf.Write(b)
	buf.Write(data)

	return buf.Bytes()
}

func DecodeFrame(r io.Reader) (*Frame, error) {
	h, err := DecodeHeader(r)
	if err != nil {
		return nil, err
	}
	data := make([]byte, h.bodyLength)
	if h.bodyLength == 0 {
		f := framePool.Get().(*Frame)
		f.Header = *h

		return f, err

	}

	read, err := r.Read(data)
	if err != nil {
		return nil, err
	}

	if uint64(read) != h.bodyLength {
		return nil, errors.New("read body unexpect length. ")
	}

	f := framePool.Get().(*Frame)
	f.Header = *h
	f.Data = data

	return f, nil
}

func (f *Frame) Unmarshal(val interface{}) error {
	switch conn.Protoc(f.Header.protocol) {
	case conn.JSONProtocol:
		return json.Unmarshal(f.Data, val)
	default:
		return json.Unmarshal(f.Data, val)
	}

	return nil
}
