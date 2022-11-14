package osser

type OSSConfig struct {
	Path            string `yaml:"path"`     // 路径
	Driver          string `yaml:"driver"`   // 对象存储类型
	Endpoint        string `yaml:"endpoint"` // Endpoint 访问域名
	AccessKeyID     string `yaml:"key"`      // AccessKeyID
	AccessKeySecret string `yaml:"secret"`   // AccessKeySecret
	BucketName      string `yaml:"bucket"`   // 桶名称
	Region          string `yaml:"region"`   // 区域
}

type OSSer interface {
	Upload(objname, localfile string) error
	Delete(objname string) error
}
