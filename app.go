package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/monitor"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/rpc"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
)

// App struct
type App struct {
	ctx      context.Context
	conn     *rpc.Connection
	monitors map[base.Address]*monitor.Monitor
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		conn:     rpc.NewConnection("mainnet", false, map[string]bool{}),
		monitors: make(map[base.Address]*monitor.Monitor),
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	runtime.EventsOn(ctx, "openExcel", func(optionalData ...interface{}) {
		logger.Info("openExcel", optionalData)
	})
}

var initialized = false

func (a *App) domReady(ctx context.Context) {
	initialized = true
	a.updateState()
}

var maxRecords = utils.NOPOSI

type Progress struct {
	Count int
	Total int
}

func (p *Progress) Finish() {
	p.Count = p.Total
}

func (a *App) Reload(addressOrEns, mode string, openExcel bool) string {
	logger.Info("Reload")
	addrStr, _ := a.conn.GetEnsAddress(addressOrEns)
	address := base.HexToAddress(addrStr)
	a.monitors[address] = nil
	return a.Export(addressOrEns, mode, openExcel)
}

func (a *App) Export(addressOrEns, mode string, openExcel bool) string {
	logger.Info("Export")
	if !base.IsValidAddress(addressOrEns) {
		return fmt.Sprintf("Invalid address or ENS name: %s", addressOrEns)
	}

	addrStr, _ := a.conn.GetEnsAddress(addressOrEns)
	address := base.HexToAddress(addrStr)
	fn := "downloads/" + address.Hex() + ".csv"

	if a.monitors[address] != nil {
		return fmt.Sprintf("Exported %d transactions to %s", a.monitors[address].Count(), fn)
	}

	if mon, err := monitor.NewMonitor("mainnet", address, true); err != nil {
		return fmt.Sprintf("Error creating monitor: %s", err)
	} else {
		a.monitors[address] = &mon
		os.Remove(fn)
		logger.Info("Count:", a.monitors[address].Count())
	}

	prog := Progress{Count: 0, Total: int(a.monitors[address].Count())}
	defer prog.Finish()

	go func() {
		for {
			if prog.Total == prog.Count {
				return
			}
			prog.Count, _ = file.WordCount(fn, true)
			msg := fmt.Sprintf("Exporting %6d of %6d for %s", prog.Count, prog.Total, address.Hex())
			runtime.EventsEmit(a.ctx, "progress", msg)
			time.Sleep(100 * time.Millisecond)
		}
	}()

	cmd := Command{
		MaxRecords: int(maxRecords),
		Address:    address,
		Filename:   fn,
		Format:     "csv",
		Subcommand: "export",
		Rest:       mode,
	}

	logger.Info("Running command: ", cmd.String())
	_ = utils.System(cmd.String())
	logger.Info("Done...")

	lines := file.AsciiFileToLines(fn)

	if len(lines) == 0 {
		return fmt.Sprintf("No transactions found for %s", address.Hex())
	} else {
		if openExcel {
			_ = utils.System("open " + fn)
		}
		return fmt.Sprintf("Exported %d transactions to %s", len(lines)-1, fn)
	}
}
