package kiosk

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

// GrafanaKioskAnonymous creates a chrome-based kiosk using a local grafana-server account.
func GrafanaKioskAnonymous(cfg *Config) {
	dir, err := os.MkdirTemp(os.TempDir(), "chromedp-kiosk")
	if err != nil {
		panic(err)
	}

	log.Println("Using temp dir:", dir)
	defer os.RemoveAll(dir)

	opts := generateExecutorOptions(dir, cfg.General.WindowPosition, cfg.Target.IgnoreCertificateErrors)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// also set up a custom logger
	taskCtx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	listenChromeEvents(taskCtx, consoleAPICall|targetCrashed)

	// ensure that the browser process is started
	if err := chromedp.Run(taskCtx); err != nil {
		panic(err)
	}

	// Give browser time to load next page (this can be prone to failure, explore different options vs sleeping)
	time.Sleep(2000 * time.Millisecond)

	var generatedURL = GenerateURL(cfg.Target.URL, cfg.General.Mode, cfg.General.AutoFit, cfg.Target.IsPlayList)

	log.Println("Navigating to ", generatedURL)
	/*
		Launch chrome and look for main-view element
	*/
	if err := chromedp.Run(taskCtx,
		chromedp.Navigate(generatedURL),
		chromedp.WaitVisible(`//div[@class="main-view"]`, chromedp.BySearch),
		// wait forever (for now)
		chromedp.WaitVisible("notinputPassword", chromedp.ByID),
	); err != nil {
		panic(err)
	}

	log.Println("Sleep before exit...")
	// wait here for the process to exit
	time.Sleep(2000 * time.Millisecond)
	log.Println("Exit...")
}
