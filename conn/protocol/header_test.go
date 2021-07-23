package protocol

import (
	"bytes"
	"testing"
)

func TestHeader(t *testing.T) {
	h := AllocHeader(WithVer(1),
		WithDuration(10),
		WithChannel(100),
		WithFrom(1000),
		WithJSON(),
		WithBodyLength(10000),
		WithExtra([]byte("extra")),
	)
	b := EncodeHeader(h)
	t.Logf("l=%d", len(b))

	res, err := DecodeHeader(bytes.NewBuffer(b))
	if err != nil {
		t.Fatalf("decode failed. err=%v", err)
	}

	t.Logf("res=%+v", res)

}
