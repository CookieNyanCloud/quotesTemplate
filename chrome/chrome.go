package chrome

import (
	"context"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func Screen(url string, imageBuf *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.ActionFunc(func(abc context.Context) (err error) {
			*imageBuf, err = page.CaptureScreenshot().WithQuality(100).Do(abc)
			return err
		}),
	}
}
