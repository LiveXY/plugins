package plugin

import (
	"context"
	"io/ioutil"
	"path"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/livexy/plugin/pdfer"
)

func NewChromePdf() pdfer.PDFer {
	return &ChromePdf{}
}

type ChromePdf struct {}

func (o *ChromePdf) Html2Pdf(htmlpath, outpath string) error {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	var (
		buf []byte
		err error
	)
	err = chromedp.Run(ctx, chromedp.Tasks{
		chromedp.Navigate(path.Join("file:///", htmlpath)),
		chromedp.WaitReady("body"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			buf, _, err = page.PrintToPDF().
				Do(ctx)
			return err
		}),
	})
	if err != nil {
		return err
	}
	return ioutil.WriteFile(outpath, buf, 0644)
}
