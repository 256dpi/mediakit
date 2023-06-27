package chromium

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/256dpi/xo"
	"github.com/chromedp/cdproto/log"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

const scrollThrough = `
new Promise((resolve) => {
	const step = window.innerHeight / 2; // pixel
	const frequency = 8; // 120 Hz
	
	let scrolls = 0;
	function scroll() {
		window.setTimeout(() => {
			window.scrollTo(0, scrolls * step);

			const total = document.body.scrollHeight / step;
			console.log(step, frequency, total, scrolls);

			if (scrolls < total) {
				scrolls += 1;
				scroll();
			}

			resolve(true);
		}, frequency);
	}
	
	scroll();
});
`

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

// ScreenshotOptions are the options used for taking a screenshot.
type ScreenshotOptions struct {
	Width    int64
	Height   int64
	Full     bool
	Scale    float64
	Pedantic bool
	Wait     time.Duration
}

// Screenshot will capture a screenshot of the given URL. A browser context may
// be provided using Allocate, otherwise a new one will be allocated.
func Screenshot(ctx context.Context, url string, opts ScreenshotOptions) ([]byte, error) {
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
		withTimeout(30*time.Second, "navigation failed", chromedp.Navigate(url)),
		withTimeout(30*time.Second, "awaiting body failed", chromedp.WaitReady("body")),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// scroll through page once
			if opts.Full {
				err := chromedp.Evaluate(scrollThrough, nil, func(params *runtime.EvaluateParams) *runtime.EvaluateParams {
					return params.WithAwaitPromise(true)
				}).Do(ctx)
				if err != nil {
					return err
				}
				err = chromedp.Evaluate(`window.scroll({top: 0})`, nil).Do(ctx)
				if err != nil {
					return err
				}
			}

			// wait some time
			if opts.Wait > 0 {
				time.Sleep(opts.Wait)
			}

			// get height
			var height int64
			err := chromedp.Evaluate(`document.body.scrollHeight`, &height).Do(ctx)
			if err != nil {
				return err
			}

			// capture screenshot
			buf, err = page.CaptureScreenshot().
				WithFormat(page.CaptureScreenshotFormatPng).
				WithCaptureBeyondViewport(opts.Full && height > opts.Height).
				Do(ctx)
			if err != nil {
				return err
			}

			return nil
		}),
	)
	if err != nil {
		return nil, err
	}

	// handle log errors
	if opts.Pedantic && len(logErrors) > 0 {
		return nil, fmt.Errorf("log errors: %s", logErrors)
	}

	return buf, nil
}

func withTimeout(timeout time.Duration, msg string, tasks ...chromedp.Action) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		// prepare context
		timeoutContext, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		// tun tasks
		err := chromedp.Tasks(tasks).Do(timeoutContext)

		// handle timeout
		if timeoutContext.Err() != nil && errors.Is(err, context.DeadlineExceeded) {
			return xo.F(msg)
		}

		return err
	}
}
