package worder

type WordConfig struct {
	Path   string `yaml:"path"`   // 路径
	Driver string `yaml:"driver"` // 类型
}

// WORD接口
type Worder interface {
	Read(docxPath string, prefix ...string) ([][]string, error)
	Write(docxPath string, title string, data []Question) error
}

type Question struct {
	Title        string
	Difficulty   string
	Analysis     string
	Score        string
	Options      []string
	Answers      []int8
	QuestionType int8
	Required     bool
}
