package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	rt "runtime"
	"strings"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/monitor"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/names"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/rpc"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type DataFile struct {
	Chain       string
	Address     base.Address
	ChartType   string
	ExportExcel bool
}

type App struct {
	ctx         context.Context
	conn        *rpc.Connection
	namesMap    map[string]map[base.Address]types.SimpleName
	dataFile    DataFile
	monitors    map[DataFile]*monitor.Monitor
	lines       map[DataFile][]string
	months      map[DataFile]map[string]int
	years       map[DataFile]map[string]int
	resolved    map[DataFile]bool
	info        map[DataFile]string
	asSender    map[DataFile]map[string]int
	asRecipient map[DataFile]map[string]int
}

func (a *App) Clear(chain string, address base.Address) {
	a.conn = rpc.NewConnection(chain, false, map[string]bool{
		"blocks":       true,
		"receipts":     true,
		"transactions": true,
		"traces":       true,
		"logs":         true,
		"statements":   true,
		"state":        true,
		"tokens":       true,
		"results":      true,
	})
	a.dataFile.Chain = chain
	a.dataFile.Address = address
	a.monitors = make(map[DataFile]*monitor.Monitor)
	a.lines = make(map[DataFile][]string)
	a.months = make(map[DataFile]map[string]int)
	a.years = make(map[DataFile]map[string]int)
	a.resolved = make(map[DataFile]bool)
	a.info = make(map[DataFile]string)
	a.asSender = make(map[DataFile]map[string]int)
	a.asRecipient = make(map[DataFile]map[string]int)
	a.namesMap = make(map[string]map[base.Address]types.SimpleName)
}

// NewApp creates a new App application struct
func NewApp(chain string, address base.Address) *App {
	var app App
	app.Clear(chain, address)
	return &app
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	nameParts := names.Custom | names.Prefund | names.Regular
	m, err := names.LoadNamesMap(a.dataFile.Chain, nameParts, nil)
	a.namesMap[a.dataFile.Chain] = m
	if err != nil {
		logger.Error("Failed to load names map:", err)
	}
}

func (a *App) OpenUrl(url string) {
	var err error

	switch rt.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		log.Println("Failed to open URL:", err)
	}
}

func (a *App) domReady(ctx context.Context) {
	a.Export(a.dataFile.Address.Hex(), "")
	a.updateState()
	runtime.EventsEmit(a.ctx, "chartType", strings.Trim(file.AsciiFileToString("chartType.txt"), "\n"))
	runtime.EventsEmit(a.ctx, "exportExcel", strings.Trim(file.AsciiFileToString("exportExcel.txt"), "\n"))
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
	runtime.EventsEmit(a.ctx, "progress", "Loading names map...")
	time.Sleep(1 * time.Second)
	runtime.EventsEmit(a.ctx, "progress", "")
}
