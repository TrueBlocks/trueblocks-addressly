package main

import (
	"embed"
	"time"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()
	opts := options.App{
		Title:  "TrueBlocks Account Explorer",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 33, G: 37, B: 41, A: 1},
		OnStartup:        app.startup,
		OnDomReady:       app.domReady,
		Bind: []interface{}{
			app,
		},
		WindowStartState: options.Maximised,
	}

	done := false
	go func() {
		for {
			if done {
				return
			}
			app.updateState()
			time.Sleep(30 * time.Second)
		}
	}()
	defer func() {
		done = true
	}()

	if err := wails.Run(&opts); err != nil {
		println("Error:", err.Error())
	}
}
