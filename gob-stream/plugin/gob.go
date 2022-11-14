package plugin

import (
	"bytes"
	"encoding/gob"

	"github.com/livexy/plugins/plugin/streamer"
)

type gobStream struct {}

func (g *gobStream) Marshal(val any) ([]byte, error) {
    buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(val)
	return buf.Bytes(), err
}

func (g *gobStream) Unmarshal(buf []byte, val any) error {
	dec := gob.NewDecoder(bytes.NewReader(buf))
	err := dec.Decode(val)
	return err
}

func NewStream(cfg streamer.StreamConfig) streamer.Streamer {
	return &gobStream{}
}
