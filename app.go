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
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/rpc"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
)

type PerChain struct {
	Chain   string
	Address base.Address
}

type App struct {
	ctx        context.Context
	conn       *rpc.Connection
	chain      string
	monitors   map[PerChain]*monitor.Monitor
	lines      map[PerChain][]string
	months     map[PerChain]map[string]int
	years      map[PerChain]map[string]int
	resolved   map[PerChain]bool
	info       map[PerChain]string
	to         map[PerChain]map[string]int
	from       map[PerChain]map[string]int
	toCounts   map[int]int
	fromCounts map[int]int
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
	a.monitors = make(map[PerChain]*monitor.Monitor)
	a.lines = make(map[PerChain][]string)
	a.months = make(map[PerChain]map[string]int)
	a.years = make(map[PerChain]map[string]int)
	a.resolved = make(map[PerChain]bool)
	a.info = make(map[PerChain]string)
	a.to = make(map[PerChain]map[string]int)
	a.from = make(map[PerChain]map[string]int)
	a.toCounts = make(map[int]int)
	a.fromCounts = make(map[int]int)
}

// NewApp creates a new App application struct
func NewApp(chain string) *App {
	var app App
	app.Clear(chain)
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
	runtime.EventsEmit(a.ctx, "toCount", "")
	runtime.EventsEmit(a.ctx, "fromCount", "")
	runtime.EventsEmit(a.ctx, "toTopTen", "")
	runtime.EventsEmit(a.ctx, "fromTopTen", "")
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
	a.monitors[PerChain{a.chain, address}] = nil
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
		logger.Info("Export")
		folder := "downloads/" + a.chain + "/"
		file.EstablishFolder(folder)
		fn := folder + address.Hex() + ".csv"

		if a.monitors[PerChain{a.chain, address}] != nil {
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
			a.monitors[PerChain{a.chain, address}] = &mon
			os.Remove(fn)
			logger.Info("Count:", a.monitors[PerChain{a.chain, address}].Count())
		}

		prog := Progress{Count: 0, Total: int(a.monitors[PerChain{a.chain, address}].Count())}
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

		a.lines[PerChain{a.chain, address}] = file.AsciiFileToLines(fn)
		if len(a.lines[PerChain{a.chain, address}]) == 0 {
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
	if !a.resolved[PerChain{a.chain, address}] {
		for _, line := range a.lines[PerChain{a.chain, address}] {
			parts := strings.Split(line, ",")
			if len(parts) > 3 {
				pp := strings.Split(parts[3], "-")
				if len(pp) > 1 {
					mKey := pp[0] + "_" + pp[1]
					if a.months[PerChain{a.chain, address}] == nil {
						a.months[PerChain{a.chain, address}] = make(map[string]int)
					}
					a.months[PerChain{a.chain, address}][mKey]++
					yKey := pp[0]
					if a.years[PerChain{a.chain, address}] == nil {
						a.years[PerChain{a.chain, address}] = make(map[string]int)
					}
					a.years[PerChain{a.chain, address}][yKey]++
				}
			}
			if len(parts) > 5 {
				from := base.HexToAddress(parts[4])
				to := base.HexToAddress(parts[5])
				if from == address {
					if a.to[PerChain{a.chain, address}] == nil {
						a.to[PerChain{a.chain, address}] = make(map[string]int)
					}
					a.to[PerChain{a.chain, address}][from.Hex()+"|"+to.Hex()]++
				}
				if to == address {
					if a.from[PerChain{a.chain, address}] == nil {
						a.from[PerChain{a.chain, address}] = make(map[string]int)
					}
					a.from[PerChain{a.chain, address}][from.Hex()+"|"+to.Hex()]++
				}
			}
		}
		a.info[PerChain{a.chain, address}] = strings.ToLower(addressOrEns) + "," + a.getInfo(address) + "," + a.getBalance(address)
		a.resolved[PerChain{a.chain, address}] = true
	}

	a.toCounts = make(map[int]int)
	for _, value := range a.to[PerChain{a.chain, address}] {
		a.toCounts[value]++
	}
	a.fromCounts = make(map[int]int)
	for _, value := range a.from[PerChain{a.chain, address}] {
		a.fromCounts[value]++
	}

	runtime.EventsEmit(a.ctx, "info", a.info[PerChain{a.chain, address}])
	runtime.EventsEmit(a.ctx, "years", asSortedArray1(a.years[PerChain{a.chain, address}], len(a.years[PerChain{a.chain, address}]), true, false))
	runtime.EventsEmit(a.ctx, "months", asSortedArray1(a.months[PerChain{a.chain, address}], len(a.months[PerChain{a.chain, address}]), true, false))
	runtime.EventsEmit(a.ctx, "fromTopTen", asSortedArray1(a.from[PerChain{a.chain, address}], len(a.from[PerChain{a.chain, address}]), false, false))
	runtime.EventsEmit(a.ctx, "toTopTen", asSortedArray1(a.to[PerChain{a.chain, address}], len(a.to[PerChain{a.chain, address}]), false, false))
	runtime.EventsEmit(a.ctx, "fromCount", asSortedArray2(a.fromCounts, len(a.fromCounts), true, true))
	runtime.EventsEmit(a.ctx, "toCount", asSortedArray2(a.toCounts, len(a.toCounts), true, true))
}

type kv1 struct {
	Key   string
	Value int
}

type Sortable interface {
	int | string
}

func oneSort[T Sortable](iKey, jKey T, iVal, jVal int, bykey, reversed bool) bool {
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

func asSortedArray1(m map[string]int, limit int, bykey, reversed bool) string {
	var ss []kv1
	for k, v := range m {
		ss = append(ss, kv1{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return oneSort(ss[i].Key, ss[j].Key, ss[i].Value, ss[j].Value, bykey, reversed)
	})

	if len(ss) > limit {
		rest := ss[limit:]
		ss = ss[:limit]
		ss = append(ss, kv1{"other", 0})
		for _, v := range rest {
			ss[len(ss)-1].Value += v.Value
		}
	}

	ret := ""

	for _, kv := range ss {
		ret += fmt.Sprintf("%s-%d,", kv.Key, kv.Value)
	}

	return ret
}

type kv2 struct {
	Key   int
	Value int
}

func asSortedArray2(m map[int]int, limit int, bykey, reversed bool) string {
	var ss []kv2
	for k, v := range m {
		ss = append(ss, kv2{k, v})
	}
	sort.Slice(ss, func(i, j int) bool {
		return oneSort(ss[i].Key, ss[j].Key, ss[i].Value, ss[j].Value, bykey, reversed)
	})
	if len(ss) > limit {
		ss = ss[:limit]
	}
	ret := ""
	for _, kv := range ss {
		ret += fmt.Sprintf("%d-%d,", kv.Key, kv.Value)
	}
	return ret
}

func (a *App) SetChain(chain string, addressOrEns string) {
	logger.Info("Setting chain to: ", chain)
	a.chain = chain
	a.Export(addressOrEns, "", false)
}
