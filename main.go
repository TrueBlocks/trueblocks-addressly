package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()
	opts := options.App{
		Title: "TrueBlocks Account Explorer",
		// Width:  1024,
		// Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 33, G: 37, B: 41, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		WindowStartState: options.Maximised,
	}

	if err := wails.Run(&opts); err != nil {
		println("Error:", err.Error())
	}
}
