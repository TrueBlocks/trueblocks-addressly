package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	rt "runtime"
	"sort"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/monitor"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/names"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/rpc"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
)

type DataFileKey struct {
	Chain   string
	Address base.Address
}

type App struct {
	ctx         context.Context
	conn        *rpc.Connection
	chain       string
	monitors    map[DataFileKey]*monitor.Monitor
	lines       map[DataFileKey][]string
	months      map[DataFileKey]map[string]int
	years       map[DataFileKey]map[string]int
	resolved    map[DataFileKey]bool
	info        map[DataFileKey]string
	asSender    map[DataFileKey]map[string]int
	asRecipient map[DataFileKey]map[string]int
	namesMap    map[string]map[base.Address]types.SimpleName
}

func (a *App) Clear(chain string) {
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
	a.chain = chain
	a.monitors = make(map[DataFileKey]*monitor.Monitor)
	a.lines = make(map[DataFileKey][]string)
	a.months = make(map[DataFileKey]map[string]int)
	a.years = make(map[DataFileKey]map[string]int)
	a.resolved = make(map[DataFileKey]bool)
	a.info = make(map[DataFileKey]string)
	a.asSender = make(map[DataFileKey]map[string]int)
	a.asRecipient = make(map[DataFileKey]map[string]int)
	a.namesMap = make(map[string]map[base.Address]types.SimpleName)
}

// NewApp creates a new App application struct
func NewApp(chain string) *App {
	var app App
	app.Clear(chain)
	nameParts := names.Custom | names.Prefund | names.Regular
	m, err := names.LoadNamesMap(chain, nameParts, nil)
	app.namesMap[chain] = m
	if err != nil {
		return nil
	}
	return &app
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	runtime.EventsOn(ctx, "openExcel", func(optionalData ...interface{}) {
		logger.Info("openExcel", optionalData)
	})
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

var initialized = false

func (a *App) domReady(ctx context.Context) {
	initialized = true
	a.updateState()
}

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
	address := base.HexToAddress(addrStr)
	dfKey := DataFileKey{a.chain, address}
	a.monitors[dfKey] = nil
	a.Export(addressOrEns, mode, openExcel)
}

func (a *App) Export(addressOrEns, mode string, openExcel bool) {
	defer func() {
		runtime.EventsEmit(a.ctx, "progress", "")
	}()

	if address, err := a.initExport(addressOrEns); err != nil {
		runtime.EventsEmit(a.ctx, "error", fmt.Sprintf("Error initializing export: %s", err))
		return

	} else {
		dfKey := DataFileKey{a.chain, address}

		logger.Info("Export")
		folder := "downloads/" + a.chain + "/"
		file.EstablishFolder(folder)
		fn := folder + address.Hex() + ".csv"

		if a.monitors[dfKey] != nil {
			a.Summarize(address, addressOrEns)
			return
		}

		if mon, err := monitor.NewMonitor(a.chain, address, true); err != nil {
			runtime.EventsEmit(a.ctx, "error", fmt.Sprintf("Error creating monitor: %s", err))
			return
		} else {
			if mon.Count() == 0 {
				a.getInfo(address)
			}
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
				msg := fmt.Sprintf("Exporting %6d of %6d for %s", prog.Count, prog.Total, address.Hex())
				runtime.EventsEmit(a.ctx, "error", "")
				runtime.EventsEmit(a.ctx, "progress", msg)
				logger.Info("Progress:", msg)
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
			Silent:     false,
			Chain:      a.chain,
		}

		logger.Info("Running command: ", cmd.String())
		_ = utils.System(cmd.String())
		logger.Info("Done...")

		a.lines[dfKey] = file.AsciiFileToLines(fn)
		if len(a.lines[dfKey]) == 0 {
			runtime.EventsEmit(a.ctx, "error", fmt.Sprintf("No transactions found for %s", address.Hex()))
			return
		} else {
			if openExcel {
				_ = utils.System("open " + fn)
			}
			a.Summarize(address, addressOrEns)
		}
	}
}

func (a *App) Summarize(address base.Address, addressOrEns string) {
	dfKey := DataFileKey{a.chain, address}

	if !a.resolved[dfKey] {
		a.info[dfKey] = strings.ToLower(addressOrEns) + "," + a.getInfo(address) + "," + a.getBalance(address)
	}
	runtime.EventsEmit(a.ctx, "info", a.info[dfKey])

	if !a.resolved[dfKey] {
		for _, line := range a.lines[dfKey] {
			parts := strings.Split(line, ",")
			if len(parts) > 3 {
				pp := strings.Split(parts[3], "-")
				if len(pp) > 1 {
					mKey := pp[0] + "_" + pp[1]
					if a.months[dfKey] == nil {
						a.months[dfKey] = make(map[string]int)
					}
					a.months[dfKey][mKey]++
					yKey := pp[0]
					if a.years[dfKey] == nil {
						a.years[dfKey] = make(map[string]int)
					}
					a.years[dfKey][yKey]++
				}
			}

			if len(parts) > 5 {
				from := base.HexToAddress(parts[4])
				to := base.HexToAddress(parts[5])
				if from == address {
					if a.asSender[dfKey] == nil {
						a.asSender[dfKey] = make(map[string]int)
					}
					a.asSender[dfKey][a.CounterKey(from, to)]++
				}
				if to == address {
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

func (a *App) SetChain(chain string, addressOrEns string) {
	logger.Info("Setting chain to: ", chain)
	a.chain = chain
	a.Export(addressOrEns, "", false)
}

func (a *App) CounterKey(from, to base.Address) string {
	f := from.Hex()
	fN := strings.Replace(strings.Replace(a.namesMap[a.chain][from].Name, "|", " ", -1), ",", " ", -1)
	t := to.Hex()
	tN := strings.Replace(strings.Replace(a.namesMap[a.chain][to].Name, "|", " ", -1), ",", " ", -1)
	return f + "|" + fN + "|" + t + "|" + tN
}
