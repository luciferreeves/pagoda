package services

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"
	"shrine/repositories"
	"shrine/utils/logger"
	"shrine/utils/storage"
	"time"

	"github.com/chromedp/chromedp"
	"golang.org/x/image/draw"
)

const (
	thumbnailWidth  = 320
	thumbnailHeight = 200
	captureWidth    = 1280
	captureHeight   = 800
	captureTimeout  = 30 * time.Second
)

func GenerateThumbnails() {
	sites := repositories.ListApprovedSitesForThumbnail()
	if len(sites) == 0 {
		return
	}

	logger.Infof("Thumbnails", "Generating thumbnails for %d sites", len(sites))

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.WindowSize(captureWidth, captureHeight),
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
	)

	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer allocCancel()

	for _, site := range sites {
		var screenshot []byte

		ctx, cancel := chromedp.NewContext(allocCtx)
		timeoutCtx, timeoutCancel := context.WithTimeout(ctx, captureTimeout)

		err := chromedp.Run(timeoutCtx,
			chromedp.Navigate(site.URL),
			chromedp.Sleep(2*time.Second),
			chromedp.FullScreenshot(&screenshot, 90),
		)

		timeoutCancel()
		cancel()

		if err != nil {
			logger.Errorf("Thumbnails", "Failed to capture %s: %v", site.URL, err)
			continue
		}

		resized, err := resizeScreenshot(screenshot)
		if err != nil {
			logger.Errorf("Thumbnails", "Failed to resize %s: %v", site.URL, err)
			continue
		}

		path := fmt.Sprintf("districts/thumbnails/%s.png", site.Ref)
		err = storage.Upload(path, bytes.NewReader(resized), int64(len(resized)), "image/png")
		if err != nil {
			logger.Errorf("Thumbnails", "Failed to upload thumbnail for %s: %v", site.URL, err)
			continue
		}

		site.ThumbnailURL = path
		repositories.UpdateDistrictSite(&site)

		logger.Infof("Thumbnails", "Generated thumbnail for %s", site.URL)
	}

	logger.Successf("Thumbnails", "Thumbnail generation complete")
}

func GenerateSiteThumbnail(ref string, siteURL string) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.WindowSize(captureWidth, captureHeight),
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
	)

	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer allocCancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()
	timeoutCtx, timeoutCancel := context.WithTimeout(ctx, captureTimeout)
	defer timeoutCancel()

	var screenshot []byte
	err := chromedp.Run(timeoutCtx,
		chromedp.Navigate(siteURL),
		chromedp.Sleep(2*time.Second),
		chromedp.FullScreenshot(&screenshot, 90),
	)
	if err != nil {
		logger.Errorf("Thumbnails", "Failed to capture %s: %v", siteURL, err)
		return
	}

	resized, err := resizeScreenshot(screenshot)
	if err != nil {
		logger.Errorf("Thumbnails", "Failed to resize %s: %v", siteURL, err)
		return
	}

	path := fmt.Sprintf("districts/thumbnails/%s.png", ref)
	err = storage.Upload(path, bytes.NewReader(resized), int64(len(resized)), "image/png")
	if err != nil {
		logger.Errorf("Thumbnails", "Failed to upload thumbnail for %s: %v", siteURL, err)
		return
	}

	site, findErr := repositories.FindDistrictSiteByRef(ref)
	if findErr != nil {
		return
	}
	site.ThumbnailURL = path
	repositories.UpdateDistrictSite(site)

	logger.Infof("Thumbnails", "Generated thumbnail for %s", siteURL)
}

func resizeScreenshot(data []byte) ([]byte, error) {
	src, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	dst := image.NewRGBA(image.Rect(0, 0, thumbnailWidth, thumbnailHeight))
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)

	var buf bytes.Buffer
	if err := png.Encode(&buf, dst); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}