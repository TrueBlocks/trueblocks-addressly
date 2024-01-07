package main

import "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"

func (a *App) KeyPress(key string) {
	if key == "Escape" {
		logger.Info("Escape pressed")
		escPressed = true
	}
}
