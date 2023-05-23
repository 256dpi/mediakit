package chromium

import (
	"context"
	"fmt"

	"github.com/chromedp/cdproto/log"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// Allocate will allocate a new browser instance and return an associated
// context and cancel function.
func Allocate() (context.Context, context.CancelFunc, error) {
	// prepare context
	ctx, cancel := chromedp.NewContext(context.Background())

	// allocate browser
	err := chromedp.Run(ctx)
	if err != nil {
		cancel()
		return nil, nil, err
	}

	return ctx, cancel, nil
}

// Options are the options for the screenshot capture.
type Options struct {
	Width  int64
	Height int64
	Full   bool
	Scale  float64
}

// CaptureScreenshot will capture a screenshot of the given URL. A browser
// context may be provided using Allocate, otherwise a new one will be allocated.
func CaptureScreenshot(ctx context.Context, url string, opts Options) ([]byte, error) {
	// ensure context
	if ctx == nil {
		ctx = context.Background()
	}

	// wrap context
	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	// collect errors
	var logErrors []string
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		// check for errors
		if ev, ok := ev.(*log.EventEntryAdded); ok {
			if ev.Entry.Level == log.LevelError {
				logErrors = append(logErrors, fmt.Sprintf("%s (%s)", ev.Entry.Text, ev.Entry.URL))
			}
		}
	})

	// capture screenshot
	var buf []byte
	err := chromedp.Run(ctx,
		chromedp.EmulateViewport(opts.Width, opts.Height, chromedp.EmulateScale(opts.Scale)),
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			buf, err = page.CaptureScreenshot().
				WithFormat(page.CaptureScreenshotFormatPng).
				WithCaptureBeyondViewport(opts.Full).
				Do(ctx)
			return err
		}),
	)
	if err != nil {
		return nil, err
	}

	// handle log errors
	if len(logErrors) > 0 {
		return nil, fmt.Errorf("log errors: %s", logErrors)
	}

	return buf, nil
}
