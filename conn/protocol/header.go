package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"sync"

	"github.com/eggz6/common/conn"
)

var headerPool sync.Pool

func init() {
	headerPool.New = func() interface{} {
		return &Header{}
	}
}

// 1 + 4 + 4 + 8 + 4 + 1 + 8
const headStableLen uint32 = 26

type Header struct {
	ver          uint8  // 1
	headerLength uint32 // 5
	duration     uint32 // 自建联以来的持续时间 9, 单位毫秒
	channel      uint32 // 17
	from         uint32 // 21
	protocol     uint8  // 22 上层协议类型
	bodyLength   uint64 // 30
	extra        []byte
}

func (h *Header) Duration() uint32 {
	return h.duration
}

func (h *Header) Extra() []byte {
	return h.extra
}

func (h *Header) Channel() uint32 {
	return h.channel
}

func (h *Header) From() uint32 {
	return h.from
}

func EncodeHeader(h *Header) []byte {
	res := []byte{}

	buf := bytes.NewBuffer(res)
	binary.Write(buf, binary.BigEndian, h.ver)
	binary.Write(buf, binary.BigEndian, h.headerLength)
	binary.Write(buf, binary.BigEndian, h.duration)
	binary.Write(buf, binary.BigEndian, h.channel)
	binary.Write(buf, binary.BigEndian, h.from)
	binary.Write(buf, binary.BigEndian, h.protocol)
	binary.Write(buf, binary.BigEndian, h.bodyLength)
	buf.Write(h.extra)

	return buf.Bytes()
}

type HeaderOption func(h *Header)

func WithVer(v uint8) HeaderOption {
	return func(h *Header) {
		h.ver = v
	}
}

func WithDuration(d uint32) HeaderOption {
	return func(h *Header) {
		h.duration = d
	}
}

func WithChannel(c uint32) HeaderOption {
	return func(h *Header) {
		h.channel = c
	}
}

func WithFrom(f uint32) HeaderOption {
	return func(h *Header) {
		h.from = f
	}
}

func WithJSON() HeaderOption {
	return func(h *Header) {
		h.protocol = 1
	}
}

func WithBodyLength(bodyLength uint64) HeaderOption {
	return func(h *Header) {
		h.bodyLength = bodyLength
	}
}

func WithExtra(extra []byte) HeaderOption {
	return func(h *Header) {
		h.extra = extra
	}
}

func AllocHeader(options ...HeaderOption) *Header {
	res := headerPool.Get().(*Header)
	for _, opt := range options {
		opt(res)
	}

	res.headerLength = headStableLen + uint32(len(res.extra))

	return res
}

type headerDecodeHandle func(r io.Reader, h *Header) error

func DecodeHeader(r io.Reader) (*Header, error) {
	res := headerPool.Get().(*Header)

	chain := []headerDecodeHandle{
		readVer(), readHeaderLength(), readDuration(), readChannel(), readFrom(), readProtocol(), readbodyLength(), readExtra(),
	}

	var err error
	for _, c := range chain {
		err = c(r, res)
		if err != nil {
			break
		}
	}

	return res, err
}

func readVer() headerDecodeHandle {
	return func(r io.Reader, h *Header) error {
		return binary.Read(r, binary.BigEndian, &h.ver)
	}
}

func readHeaderLength() headerDecodeHandle {
	return func(r io.Reader, h *Header) error {
		return binary.Read(r, binary.BigEndian, &h.headerLength)
	}
}

func readDuration() headerDecodeHandle {
	return func(r io.Reader, h *Header) error {
		return binary.Read(r, binary.BigEndian, &h.duration)
	}
}

func readChannel() headerDecodeHandle {
	return func(r io.Reader, h *Header) error {
		return binary.Read(r, binary.BigEndian, &h.channel)
	}
}

func readFrom() headerDecodeHandle {
	return func(r io.Reader, h *Header) error {
		return binary.Read(r, binary.BigEndian, &h.from)
	}
}

func readProtocol() headerDecodeHandle {
	return func(r io.Reader, h *Header) error {
		return binary.Read(r, binary.BigEndian, &h.protocol)
	}
}

func readbodyLength() headerDecodeHandle {
	return func(r io.Reader, h *Header) error {
		return binary.Read(r, binary.BigEndian, &h.bodyLength)
	}
}

func readExtra() headerDecodeHandle {
	return func(r io.Reader, h *Header) error {
		l := h.headerLength - headStableLen
		if l == 0 {
			return nil
		}

		res := make([]byte, l)
		read, err := r.Read(res)

		if err != nil {
			return err
		}

		if uint32(read) != l {
			return errors.New("header extra read unexpect length")

		}

		return nil
	}
}

func Recycle(h *Header) {
	h.ver = 0
	h.headerLength = 0
	h.duration = 0
	h.channel = 0
	h.from = 0
	h.protocol = 0
	h.bodyLength = 0
	h.extra = nil

	headerPool.Put(h)
}

func (h *Header) Dispose() {
	Recycle(h)
}

func SwitchProtoc(p conn.Protoc) HeaderOption {
	return func(h *Header) {
		h.protocol = uint8(p)
	}
}
