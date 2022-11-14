package plugin

import (
	"github.com/livexy/plugins/plugin/jsoner"

	"github.com/bytedance/sonic"
)

type sonicJson struct {}

func (j *sonicJson) Marshal(val any) ([]byte, error) {
    return sonic.Marshal(val)
}

func (j *sonicJson) MarshalString(val any) (string, error) {
    return sonic.MarshalString(val)
}

func (j *sonicJson) Unmarshal(buf []byte, val any) error {
    return sonic.Unmarshal(buf, val)
}

func (j *sonicJson) UnmarshalString(buf string, val any) error {
    return sonic.UnmarshalString(buf, val)
}

func NewSonicJson(cfg jsoner.JsonConfig) jsoner.Jsoner {
	return &sonicJson{}
}
