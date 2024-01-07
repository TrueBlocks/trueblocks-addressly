package main

import (
	"context"
	"strings"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/names"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) startup(ctx context.Context) {
	logger.Info("Calling app.startup")
	a.ctx = ctx
	nameParts := names.Custom | names.Prefund | names.Regular | names.Baddress
	m, err := names.LoadNamesMap(a.dataFile.Chain, nameParts, nil)
	a.namesMap[a.dataFile.Chain] = m
	if err != nil {
		logger.Error("Failed to load names map:", err)
	}
}

func (a *App) domReady(ctx context.Context) {
	runtime.WindowSetPosition(a.ctx, a.config.WindowState.X, a.config.WindowState.Y)
	runtime.WindowShow(ctx)

	logger.Info("Calling app.domReady")
	chartType, err := a.config.Get("settings", "chartType", "month")
	if err != nil {
		logger.Error(err)
	}
	logger.Info("domReady", "chartType", chartType)

	exportExcel, err := a.config.Get("settings", "exportExcel", "false")
	if err != nil {
		logger.Error(err)
	}
	logger.Info("domReady", "exportExcel", exportExcel)

	a.Export(a.dataFile.Address.Hex(), "")
	a.updateState()

	if file.FileExists("addresses.txt") {
		lines := file.AsciiFileToLines("addresses.txt")
		m := make([]string, 0, len(lines))
		for _, line := range lines {
			if strings.HasPrefix(line, "#") {
				continue
			}
			m = append(m, line)
		}
		runtime.EventsEmit(a.ctx, "monitors", strings.Join(m, "\n"))
	}

	runtime.EventsEmit(a.ctx, "chartType", chartType)
	runtime.EventsEmit(a.ctx, "exportExcel", exportExcel)
	runtime.EventsEmit(a.ctx, "progress", "Loading names map...")
	time.Sleep(1 * time.Second)
	runtime.EventsEmit(a.ctx, "progress", "")
}

func (a *App) shutdown(ctx context.Context) {
	a.config.WindowState.Width, a.config.WindowState.Height = runtime.WindowGetSize(ctx)
	a.config.WindowState.X, a.config.WindowState.Y = runtime.WindowGetPosition(ctx)
	a.config.Save()
}
