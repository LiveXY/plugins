package plugin

import (
	"errors"
	"strconv"
	"strings"

	"github.com/livexy/plugin/exceler"

	"github.com/livexy/pkg/check"

	"github.com/xuri/excelize/v2"
)

type xlsxExcel struct {}

func (excel *xlsxExcel) Read(xlsxpath string, names ...string) ([]exceler.Field, [][]string, error) {
	var lines [][]string
	var fields []exceler.Field
	var name string
	xlsx, err := excelize.OpenFile(xlsxpath)
	if err != nil {
		return fields, lines, err
	}
	if len(names) > 0 {
		name = names[0]
		lines, err = xlsx.GetRows(name)
	} else {
		err = errors.New("no sheet")
	}
	if err != nil {
		name = getSheetName(xlsx)
		lines, err = xlsx.GetRows(name)
	}
	if err != nil {
		return fields, lines, err
	}
	fields, err = getStruct(xlsx, name)
	if len(lines) < 2 || len(lines[0]) < 1 {
		return fields, lines, err
	}
	return fields, lines, nil
}

func (excel *xlsxExcel) Export(xlsxpath, name string, data [][]string) bool {
	xlsx := excelize.NewFile()
	_ = xlsx.SetSheetName("Sheet1", name)
	ename, _ := excelize.ColumnNumberToName(len(data[0]))
	_ = xlsx.SetColWidth(name, "A", ename, 30)
	for i, l := range data {
		line := strconv.Itoa(i + 1)
		for k, v := range l {
			ename, _ := excelize.ColumnNumberToName(k + 1)
			_ = xlsx.SetCellValue(name, ename+line, v)
		}
	}
	// 冻结第一行
	panes := &excelize.Panes{
		Freeze: true, Split: false, XSplit: 0, YSplit: 1, TopLeftCell: "A2", ActivePane: "bottomLeft", Selection: []excelize.Selection{
			{SQRef: "A1:XFD1", ActiveCell: "A1", Pane: "bottomLeft"},
		},
	}
	_ = xlsx.SetPanes(name, panes)
	xlsx.SetActiveSheet(0)
	return xlsx.SaveAs(xlsxpath) == nil
}

func (excel *xlsxExcel) MultiExport(xlsxpath string, data ...exceler.ExportData) bool {
	xlsx := excelize.NewFile()
	for j, d := range data {
		_ = xlsx.SetSheetName("Sheet"+strconv.Itoa(j+1), d.Name)
		for i, w := range d.Width {
			ename, _ := excelize.ColumnNumberToName(i + 1)
			_ = xlsx.SetColWidth(d.Name, "A", ename, w)
		}
		for i, l := range d.Data {
			line := strconv.Itoa(i + 1)
			for k, v := range l {
				ename, _ := excelize.ColumnNumberToName(k + 1)
				_ = xlsx.SetCellValue(d.Name, ename+line, v)
			}
		}
		// 冻结第一行
		panes := &excelize.Panes{
			Freeze: true, Split: false, XSplit: 0, YSplit: 1, TopLeftCell: "A2", ActivePane: "bottomLeft", Selection: []excelize.Selection{
				{SQRef: "A1:XFD1", ActiveCell: "A1", Pane: "bottomLeft"},
			},
		}
		_ = xlsx.SetPanes(d.Name, panes)
	}
	xlsx.SetActiveSheet(0)
	return xlsx.SaveAs(xlsxpath) == nil
}

func (excel *xlsxExcel) AdvancedExport(xlsxpath, name string, header []string, names map[string]string, widths map[string]float64, data [][]string) bool {
	xlsx := excelize.NewFile()
	_ = xlsx.SetSheetName("Sheet1", name)
	for k, v := range header {
		nname := names[v]
		col, _ := excelize.ColumnNumberToName(k + 1)
		comment := excelize.Comment{Cell: col+"1", Author: "Field:", Text: v}
		_ = xlsx.AddComment(name, comment)
		_ = xlsx.SetCellValue(name, col+"1", nname)
		width := widths[v]
		if width < 1 {
			width = 10
		}
		_ = xlsx.SetColWidth(name, col, col, width)
	}
	// 冻结第一行
	panes := &excelize.Panes{
		Freeze: true, Split: false, XSplit: 0, YSplit: 1, TopLeftCell: "A2", ActivePane: "bottomLeft", Selection: []excelize.Selection{
			{SQRef: "A1:XFD1", ActiveCell: "A1", Pane: "bottomLeft"},
		},
	}
	_ = xlsx.SetPanes(name, panes)
	for i, l := range data {
		line := strconv.Itoa(i + 2)
		for k, v := range l {
			col, _ := excelize.ColumnNumberToName(k + 1)
			_ = xlsx.SetCellValue(name, col+line, v)
		}
	}
	xlsx.SetActiveSheet(0)
	return xlsx.SaveAs(xlsxpath) == nil
}

// 获取Excel结构
func (excel *xlsxExcel) Struct(xlsxpath string) ([]exceler.Field, error) {
	xlsx, err := excelize.OpenFile(xlsxpath)
	if err != nil {
		return []exceler.Field{}, err
	}
	name := getSheetName(xlsx)
	return getStruct(xlsx, name)
}

func getStruct(xlsx *excelize.File, name string) ([]exceler.Field, error) {
	var list []exceler.Field
	columns, istip := []string{}, false
	rows, err := xlsx.Rows(name)
	if rows.Next() {
		columns, _ = rows.Columns()
	}
	clen := len(columns)
	if clen == 0 {
		return list, err
	}
	comments, err := xlsx.GetComments(name)
	if err != nil {
		return list, err
	}
	fields := make(map[int]string)
	for i, v := range comments {
		field := v.Text
		if field == "Tip:Error" {
			istip = true
			continue
		}
		if strings.Index(field, "Field:") == 0 {
			field = field[6:]
		}
		if !check.IsUserName(field) {
			field = ""
		}
		fields[i] = field
	}
	if istip {
		columns = columns[1:]
	}
	for i, v := range columns {
		index := i
		if istip {
			index = i + 1
		}
		list = append(list, exceler.Field{Index: index, Key: fields[index], Name: v})
	}
	return list, nil
}

func getSheetName(xlsx *excelize.File) string {
	index := xlsx.GetActiveSheetIndex()
	name := xlsx.GetSheetName(index)
	return name
}

func NewQAXExcel(cfg exceler.ExcelConfig) exceler.Exceler {
	return &xlsxExcel{}
}
