package main

import (
	"fmt"
	"os"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/monitor"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
)

var maxRecords = utils.NOPOSI

type Progress struct {
	Count    int
	Total    int
	Finished bool
}

func (p *Progress) Finish() {
	p.Finished = true
}

func (a *App) initExport(addressOrEns string) (base.Address, error) {
	runtime.EventsEmit(a.ctx, "error", "")
	runtime.EventsEmit(a.ctx, "info", "Loading...")
	runtime.EventsEmit(a.ctx, "years", "")
	runtime.EventsEmit(a.ctx, "months", "")
	runtime.EventsEmit(a.ctx, "asSender", "")
	runtime.EventsEmit(a.ctx, "asRecipient", "")
	runtime.EventsEmit(a.ctx, "progress", "Scanning Unchained Index...")
	if !base.IsValidAddress(addressOrEns) {
		return base.Address{}, fmt.Errorf("invalid address or ENS name: %s", addressOrEns)
	}

	a.SetAddress(addressOrEns)
	if a.dataFile.Address.IsZero() {
		return base.Address{}, fmt.Errorf("no ENS address found: %s", addressOrEns)
	}

	if a.namesMap[a.dataFile.Chain][a.dataFile.Address].Name == "" {
		n := types.SimpleName{Name: addressOrEns, Address: a.dataFile.Address}
		a.namesMap[a.dataFile.Chain][a.dataFile.Address] = n
	}

	return a.dataFile.Address, nil
}

var exportToExcel = false

func (a *App) Export(addressOrEns, mode string) {
	defer func() {
		runtime.EventsEmit(a.ctx, "progress", "")
	}()

	var err error
	if a.dataFile.Address, err = a.initExport(addressOrEns); err != nil {
		runtime.EventsEmit(a.ctx, "error", fmt.Sprintf("Error initializing export: %s", err))
		return

	} else {
		dfKey := a.dataFile
		folder := "/Users/jrush/Development/trueblocks-addressly/downloads/" + a.dataFile.Chain + "/"
		file.EstablishFolder(folder)
		fn := folder + a.dataFile.Address.Hex() + ".csv"

		if a.monitors[dfKey] != nil {
			a.Summarize(addressOrEns)
			return
		}

		if mon, err := monitor.NewMonitor(a.dataFile.Chain, a.dataFile.Address, true); err != nil {
			runtime.EventsEmit(a.ctx, "error", fmt.Sprintf("Error creating monitor: %s", err))
			return
		} else {
			runtime.EventsEmit(a.ctx, "info", a.getInfo(addressOrEns))
			a.monitors[dfKey] = &mon
			os.Remove(fn)
			// logger.Info("Count:", a.monitors[dfKey].Count())
		}

		prog := Progress{Count: 0, Total: int(a.monitors[dfKey].Count())}
		defer prog.Finish()

		go func() {
			for {
				if prog.Finished {
					return
				}
				prog.Count, _ = file.WordCount(fn, true)
				msg := fmt.Sprintf("Exporting %6d of %6d for %s", prog.Count, prog.Total, a.dataFile.Address.Hex())
				runtime.EventsEmit(a.ctx, "error", "")
				runtime.EventsEmit(a.ctx, "progress", msg)
				time.Sleep(100 * time.Millisecond)
			}
		}()

		cmd := Command{
			MaxRecords: int(maxRecords),
			Address:    a.dataFile.Address,
			Filename:   fn,
			Format:     "csv",
			Subcommand: "export",
			Rest:       mode,
			Silent:     false,
			Chain:      a.dataFile.Chain,
		}
		_ = utils.System(cmd.String())

		a.lines[dfKey] = file.AsciiFileToLines(fn)
		if len(a.lines[dfKey]) == 0 {
			runtime.EventsEmit(a.ctx, "error", fmt.Sprintf("No transactions found for %s", a.dataFile.Address.Hex()))
			return
		} else {
			if exportToExcel {
				_ = utils.System("open " + fn)
			}
			a.Summarize(addressOrEns)
		}
	}
}
