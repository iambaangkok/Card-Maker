package renderer

import (
	"context"
	"os"
	"path/filepath"

	"github.com/chromedp/chromedp"
	"github.com/iambaangkok/Card-Maker/internal/config"
)

type ChromeRenderer interface {
	RenderElement()
}

type ChromeRendererImpl struct {
	Config config.Config
	ctx context.Context
	cancel context.CancelFunc
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