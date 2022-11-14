package exceler

type ExcelConfig struct {
	Path   string `yaml:"path"`   // 路径
	Driver string `yaml:"driver"` // 类型
}

// Excel接口
type Exceler interface {
	Read(xlsxpath string, name ...string) ([]Field, [][]string, error)
	Struct(xlsxpath string) ([]Field, error)
	Export(xlsxpath string, name string, data [][]string) bool
	MultiExport(xlsxpath string, data ...ExportData) bool
	AdvancedExport(xlsxpath, name string, header []string, names map[string]string, widths map[string]float64, data [][]string) bool
}

type ExportData struct {
	Name  string     `json:"name"`  // 名称
	Data  [][]string `json:"data"`  // 数据
	Width []float64  `json:"width"` // 宽
}

type Field struct {
	Key   string `json:"key"`
	Name  string `json:"name"`
	Index int    `json:"index"`
}
