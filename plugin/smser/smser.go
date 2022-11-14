package smser

type SMSConfig struct {
	Path         string `yaml:"path"`     // 路径
	Driver       string `yaml:"driver"`   // 发短信类型
	SignName     string `yaml:"signName"` // 名称
	AccessID     string `yaml:"accessID"`
	AccessSecret string `yaml:"accessSecret"`
	ExtendData   string `yaml:"extendData"`
	Template     struct {
		VerificationCode string `yaml:"verificationCode"`
	} `yaml:"template"`
}

type SMSer interface {
	Send(template, mobile string, data map[string]any) error
}
