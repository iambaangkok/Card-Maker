package renderer

import (
	"context"
	"os"
	"path/filepath"
	"sync"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/iambaangkok/Card-Maker/internal/config"
)

// Card viewport for PNG capture: standard poker playing card at CSS 96 DPI —
// 2.5 in × 3.5 in (63.5 mm × 88.9 mm) → 240 × 336 px.
const CardViewportWidth = 240
const CardViewportHeight = 336

type ChromeRenderer interface {
	RenderElement()
}

type ChromeRendererImpl struct {
	Config config.Config
}

// RenderHTMLToPNG renders HTML to a PNG clipped to viewportWidth×viewportHeight CSS pixels (Scale 2 device pixels).
func (c ChromeRendererImpl) RenderHTMLToPNG(html string, outputFileName string, viewportWidth, viewportHeight float64) error {
	if viewportWidth <= 0 {
		viewportWidth = CardViewportWidth
	}
	if viewportHeight <= 0 {
		viewportHeight = CardViewportHeight
	}
	outputPath := outputFileName
	if !filepath.IsAbs(outputFileName) && filepath.Dir(outputFileName) == "." {
		// No directory component provided; write into the configured output directory.
		outputPath = filepath.Join(c.Config.Renderer.OutputDir, outputFileName)
	}

	// Fix CORS when loading images from localhost FileServer
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		// This disables the Private Network Access checks entirely
		chromedp.Flag("disable-features", "BlockInsecurePrivateNetworkRequests"),
		// Common flags for dev environments
		chromedp.Flag("disable-web-security", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)

	return chromedp.Run(ctx,
		chromedp.Navigate("about:blank"),
		// setup the listener to listen for the page.EventLoadEventFired
		chromedp.ActionFunc(func(ctx context.Context) error {
			lctx, cancel := context.WithCancel(ctx)
			chromedp.ListenTarget(lctx, func(ev interface{}) {
				if _, ok := ev.(*page.EventLoadEventFired); ok {
					wg.Done()
					// remove the event listener
					cancel()
				}
			})
			return nil
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}

			return page.SetDocumentContent(frameTree.Frame.ID, html).Do(ctx)
		}),
		// wait for page.EventLoadEventFired
		chromedp.ActionFunc(func(ctx context.Context) error {
			wg.Wait()
			return nil
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, err := page.CaptureScreenshot().
				WithFormat(page.CaptureScreenshotFormatPng).
				WithQuality(100).
				WithClip(&page.Viewport{
					X:      0,
					Y:      0,
					Width:  viewportWidth,
					Height: viewportHeight,
					Scale:  2,
				}).
				Do(ctx)
			if err != nil {
				return err
			}
			return os.WriteFile(outputPath, buf, 0644)
		}),
	)
}


func (c ChromeRendererImpl) RenderHTMLToPDF(html string, outputFileName string) error {
	outputPath := outputFileName
	if !filepath.IsAbs(outputFileName) && filepath.Dir(outputFileName) == "." {
		outputPath = filepath.Join(c.Config.Renderer.OutputDir, outputFileName)
	}

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	return chromedp.Run(ctx,
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}

			return page.SetDocumentContent(frameTree.Frame.ID, html).Do(ctx)
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().WithPrintBackground(false).Do(ctx)
			if err != nil {
				return err
			}
			return os.WriteFile(outputPath, buf, 0644)
		}),
	)
}

func (c ChromeRendererImpl) RenderElement(sel string, outputFileName string) error {
	ctx, cancel := chromedp.NewContext(
		context.Background(),
	)
	defer cancel()

	url := c.Config.Renderer.URL
	outputDir := c.Config.Renderer.OutputDir

	var buf []byte
	if err := chromedp.Run(ctx, elementScreenshot(url, sel, &buf)); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(outputDir, outputFileName) , buf, 0o644); err != nil {
		return err
	}
	return nil
}


// elementScreenshot takes a screenshot of a specific element.
func elementScreenshot(urlstr, sel string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.Screenshot(sel, res, chromedp.NodeVisible),
	}
}

// fullScreenshot takes a screenshot of the entire browser viewport.
//
// Note: chromedp.FullScreenshot overrides the device's emulation settings. Use
// device.Reset to reset the emulation and viewport settings.
func fullScreenshot(urlstr string, quality int, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.FullScreenshot(res, quality),
	}
}