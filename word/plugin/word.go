package plugin

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/livexy/plugin/worder"

	"github.com/gonfva/docxlib"
)

type docxWord struct {}
const fontSize = 11

func (word *docxWord) Read(docxPath string, prefix ...string) ([][]string, error) {
	var lines [][]string
	file, err := os.Open(filepath.Clean(docxPath))
	if err != nil {
		return lines, err
	}
	info, err := file.Stat()
	if err != nil {
		return lines, err
	}
	size := info.Size()
	doc, err := docxlib.Parse(file, int64(size))
	if err != nil {
		return lines, err
	}
    nums, lasts, isanalysis := 0, []string{}, false
	for _, para := range doc.Paragraphs() {
		if len(para.Children()) == 0 {
			if nums > 1 || isanalysis {
				lines = append(lines, lasts)
			}
			nums, lasts, isanalysis = 0, []string{}, false
			continue
		}
		nums++
		qs := []string{}
		for _, child := range para.Children() {
			if child.Run == nil {
				continue
			}
			text := strings.TrimSpace(child.Run.Text.Text)
			qs = append(qs, text)
			for _, preval := range prefix {
				if strings.Contains(child.Run.Text.Text, preval) {
					isanalysis = true
				}
			}
		}
		if len(qs) == 0 || strings.Join(qs, "") == "" {
			if nums > 1 || isanalysis {
				lines = append(lines, lasts)
			}
			nums, lasts, isanalysis = 0, []string{}, false
			continue
		}
		lasts = append(lasts, strings.Join(qs, ""))
	}
	if len(lasts) > 0 {
		lines = append(lines, lasts)
	}
	return lines, nil
}

func (word *docxWord) Write(docxPath string, title string, data []worder.Question) error {
	docx := docxlib.New()
	addTitle(docx, title)
	index := 0
	for _, d := range data {
		if addData(docx, d, index) {
			index++
		}
	}
	f, err := os.Create(filepath.Clean(docxPath))
	if err != nil {
		return err
	}
	return docx.Write(f)
}

func addData(docx *docxlib.DocxLib, v worder.Question, index int) bool {
	p := docx.AddParagraph()
	p.AddText(strconv.Itoa(index + 1)).Size(fontSize)
	p.AddText(".").Size(fontSize)
	p.AddText(v.Title).Size(fontSize)
	if len(v.Difficulty) > 0 {
		p.AddText("(" + v.Difficulty + ")").Size(fontSize - 1)
	}
	if len(v.Score) > 0 {
		p.AddText("(" + v.Score + "分)").Size(fontSize - 1)
	}
	for i, v2 := range v.Options {
		p = docx.AddParagraph()
		p.AddText(string(rune(65 + i))).Size(fontSize)
		p.AddText(".").Size(fontSize)
		p.AddText(v2).Size(fontSize)
	}
	p = docx.AddParagraph()
	p.AddText("解析：").Size(fontSize - 1)
	p.AddText(v.Analysis).Size(fontSize - 1)
	docx.AddParagraph()
	return true
}
func addTitle(docx *docxlib.DocxLib, title string) {
	p := docx.AddParagraph()
	p.AddText("\t\t\t\t" + title).Size(fontSize + 1) //.TextAlign("center")
	docx.AddParagraph()
}

func NewDocxWord(cfg worder.WordConfig) worder.Worder {
	return &docxWord{}
}
