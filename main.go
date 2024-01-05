package main

import (
	"embed"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp("mainnet", base.HexToAddress("0xf503017d7baf7fbc0fff7492b751025c6a78179b"))
	opts := options.App{
		Title:  "TrueBlocks Account Explorer",
		Width:  1024,
		Height: 1000,
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
			time.Sleep(30 * time.Second)
			app.updateState()
		}
	}()
	defer func() {
		done = true
	}()

	if err := wails.Run(&opts); err != nil {
		println("Error:", err.Error())
	}
}
