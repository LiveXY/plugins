package plugin

import (
	"github.com/livexy/plugin/jsoner"

	gojson "github.com/goccy/go-json"
)

type goJson struct {}

func (j *goJson) Marshal(val any) ([]byte, error) {
    return gojson.Marshal(val)
}

func (j *goJson) MarshalString(val any) (string, error) {
    v, err := j.Marshal(val)
    return string(v), err
}

func (j *goJson) Unmarshal(buf []byte, val any) error {
    return gojson.Unmarshal(buf, val)
}

func (j *goJson) UnmarshalString(buf string, val any) error {
    return j.Unmarshal([]byte(buf), val)
}

func NewGoJson(cfg jsoner.JsonConfig) jsoner.Jsoner {
	return &goJson{}
}
