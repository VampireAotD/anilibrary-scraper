package client

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/corpix/uarand"
)

type ChromeDp struct{}

func NewChromeDp() ChromeDp {
	return ChromeDp{}
}

// resolveAllocatorOptions method overrides some default settings for chromedp allocator
// Returns slice with options
func (c ChromeDp) resolveAllocatorOptions() []chromedp.ExecAllocatorOption {
	return append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Headless,
		chromedp.DisableGPU,
		chromedp.NoSandbox,
		chromedp.WindowSize(1920, 1080),
		chromedp.UserAgent(uarand.GetRandom()),
	)
}

// FetchDocument method sends request to a given url with a defined timeout.
// Returns goquery.Document or error if any.
func (c ChromeDp) FetchDocument(timeout time.Duration, url string) (*goquery.Document, error) {
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), c.resolveAllocatorOptions()...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

	var html string
	err := chromedp.Run(
		ctx,
		chromedp.Navigate(url),
		chromedp.OuterHTML(`body`, &html, chromedp.ByQuery),
	)
	if err != nil {
		return nil, fmt.Errorf("request for qoquery document: %w", err)
	}

	return goquery.NewDocumentFromReader(strings.NewReader(html))
}

// FetchResponseBody method sends request to a given url with a defined timeout.
// Returns slice of bytes or error if any.
func (c ChromeDp) FetchResponseBody(timeout time.Duration, url string) ([]byte, error) {
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), c.resolveAllocatorOptions()...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

	done := make(chan struct{})

	var requestID network.RequestID

	chromedp.ListenTarget(ctx, func(v any) {
		switch ev := v.(type) {
		case *network.EventRequestWillBeSent:
			if ev.Request.URL == url {
				requestID = ev.RequestID
			}
		case *network.EventLoadingFinished:
			if ev.RequestID == requestID {
				close(done)
			}
		}
	})

	if err := chromedp.Run(ctx, chromedp.Navigate(url)); err != nil {
		return nil, fmt.Errorf("request for response body: %w", err)
	}

	<-done

	var buf []byte
	if err := chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		var err error
		buf, err = network.GetResponseBody(requestID).Do(ctx)
		return err
	})); err != nil {
		return nil, fmt.Errorf("fetching response body: %w", err)
	}

	return buf, nil
}
