package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	rt "runtime"

	"github.com/TrueBlocks/trueblocks-addressly/pkg/config"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/monitor"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/rpc"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
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
	config      config.Config
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
	_ = app.config.Load()
	return &app
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
