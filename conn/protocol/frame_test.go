package protocol

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/eggz6/common/conn"
)

func TestFrame(t *testing.T) {
	frame := EncodeFrame(conn.JSONProtocol, 1000, 10000, time.Now().Add(-10*time.Second), []byte(`{"msg":"hello"}`))

	f, err := DecodeFrame(bytes.NewBuffer(frame))
	if err != nil && err != io.EOF {
		t.Fatalf("decode failed. err=%v", err)
	}

	t.Logf("frame=%+v", f)
	if f.duration != uint32((10*time.Second).Milliseconds()){
		t.Fatal("duration not equal")
	}

	if f.channel != 1000{
		t.Fatal("channel not equal")
	}

	if string(f.Data) != `{"msg":"hello"}`{
		t.Fatal("data not equal")
	}

	if conn.Protoc(f.protocol) != conn.JSONProtocol{
		t.Fatal("data not equal")
	}
}
