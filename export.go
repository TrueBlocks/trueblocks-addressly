package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/monitor"
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

	addrStr, _ := a.conn.GetEnsAddress(addressOrEns)
	address := base.HexToAddress(addrStr)
	if address.IsZero() {
		return base.Address{}, fmt.Errorf("no ENS address found: %s", addressOrEns)
	}

	return address, nil
}

func (a *App) Reload(addressOrEns, mode string, openExcel bool) {
	logger.Info("Reload")
	addrStr, _ := a.conn.GetEnsAddress(addressOrEns)
	a.dataFile.Address = base.HexToAddress(addrStr)
	dfKey := a.dataFile
	a.monitors[dfKey] = nil
	a.Export(addressOrEns, mode, openExcel)
}

func (a *App) Export(addressOrEns, mode string, openExcel bool) {
	defer func() {
		runtime.EventsEmit(a.ctx, "progress", "")
	}()

	var err error
	if a.dataFile.Address, err = a.initExport(addressOrEns); err != nil {
		runtime.EventsEmit(a.ctx, "error", fmt.Sprintf("Error initializing export: %s", err))
		return

	} else {
		dfKey := a.dataFile
		logger.Info("Export")
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
			logger.Info("Count:", a.monitors[dfKey].Count())
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
				logger.Info("Progress:", msg)
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

		logger.Info("Running command: ", cmd.String())
		_ = utils.System(cmd.String())
		logger.Info("Done...")

		a.lines[dfKey] = file.AsciiFileToLines(fn)
		if len(a.lines[dfKey]) == 0 {
			runtime.EventsEmit(a.ctx, "error", fmt.Sprintf("No transactions found for %s", a.dataFile.Address.Hex()))
			return
		} else {
			if openExcel {
				_ = utils.System("open " + fn)
			}
			a.Summarize(addressOrEns)
		}
	}
}

func (a *App) Summarize(addressOrEns string) {
	dfKey := a.dataFile

	if !a.resolved[dfKey] {
		a.info[dfKey] = a.getInfo(addressOrEns)
	}
	runtime.EventsEmit(a.ctx, "info", a.info[dfKey])

	if !a.resolved[dfKey] {
		for _, line := range a.lines[dfKey] {
			parts := strings.Split(line, ",")
			if len(parts) > 3 {
				// logger.Info("Here", line)
				pp := strings.Split(parts[3], "-")
				if len(pp) > 1 {
					// logger.Info("Here", pp)
					mKey := pp[0] + "_" + pp[1]
					if a.months[dfKey] == nil {
						a.months[dfKey] = make(map[string]int)
					}
					a.months[dfKey][mKey]++
					logger.Info("Here", mKey, a.months[dfKey][mKey])
					yKey := pp[0]
					if a.years[dfKey] == nil {
						a.years[dfKey] = make(map[string]int)
					}
					a.years[dfKey][yKey]++
					logger.Info("Here", yKey, a.months[dfKey][yKey])
				}
			}

			if len(parts) > 5 {
				from := base.HexToAddress(parts[4])
				to := base.HexToAddress(parts[5])
				if from == a.dataFile.Address {
					if a.asSender[dfKey] == nil {
						a.asSender[dfKey] = make(map[string]int)
					}
					a.asSender[dfKey][a.CounterKey(from, to)]++
				}
				if to == a.dataFile.Address {
					if a.asRecipient[dfKey] == nil {
						a.asRecipient[dfKey] = make(map[string]int)
					}
					a.asRecipient[dfKey][a.CounterKey(from, to)]++
				}
			}
		}
		a.resolved[dfKey] = true
	}

	runtime.EventsEmit(a.ctx, "years", asSortedArray(a.years[dfKey], len(a.years[dfKey]), true, false))
	runtime.EventsEmit(a.ctx, "months", asSortedArray(a.months[dfKey], len(a.months[dfKey]), true, false))
	runtime.EventsEmit(a.ctx, "asRecipient", asSortedArray(a.asRecipient[dfKey], len(a.asRecipient[dfKey]), false, false))
	runtime.EventsEmit(a.ctx, "asSender", asSortedArray(a.asSender[dfKey], len(a.asSender[dfKey]), false, false))
}

type kv struct {
	Key   string
	Value int
}

func oneSort(iKey, jKey string, iVal, jVal int, bykey, reversed bool) bool {
	if reversed {
		iVal, jVal = jVal, iVal
		iKey, jKey = jKey, iKey
	}
	if bykey {
		if iKey == jKey {
			return iVal > jVal
		}
		return iKey < jKey
	} else {
		if iVal == jVal {
			return iKey < jKey
		}
		return iVal > jVal
	}
}

func asSortedArray(m map[string]int, limit int, bykey, reversed bool) string {
	var ss []kv
	for k, v := range m {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return oneSort(ss[i].Key, ss[j].Key, ss[i].Value, ss[j].Value, bykey, reversed)
	})

	if len(ss) > limit {
		rest := ss[limit:]
		ss = ss[:limit]
		ss = append(ss, kv{"other", 0})
		for _, v := range rest {
			ss[len(ss)-1].Value += v.Value
		}
	}

	ret := ""

	for _, kv := range ss {
		ret += fmt.Sprintf("%s|%d,", kv.Key, kv.Value)
	}

	return strings.Trim(ret, ", ")
}

func (a *App) CounterKey(from, to base.Address) string {
	f := from.Hex()
	fN := strings.Replace(strings.Replace(a.namesMap[a.dataFile.Chain][from].Name, "|", " ", -1), ",", " ", -1)
	t := to.Hex()
	tN := strings.Replace(strings.Replace(a.namesMap[a.dataFile.Chain][to].Name, "|", " ", -1), ",", " ", -1)
	return f + "|" + fN + "|" + t + "|" + tN
}
