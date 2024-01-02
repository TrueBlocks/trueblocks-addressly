package main

import (
	"context"
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
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/rpc"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
)

// App struct
type App struct {
	ctx        context.Context
	conn       *rpc.Connection
	monitors   map[base.Address]*monitor.Monitor
	lines      map[base.Address][]string
	months     map[base.Address]map[string]int
	years      map[base.Address]map[string]int
	resolved   map[base.Address]bool
	info       map[base.Address]string
	to         map[base.Address]map[string]int
	from       map[base.Address]map[string]int
	toCounts   map[int]int
	fromCounts map[int]int
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		conn: rpc.NewConnection("mainnet", false, map[string]bool{
			"blocks":       true,
			"receipts":     true,
			"transactions": true,
			"traces":       true,
			"logs":         true,
			"statements":   true,
			"state":        true,
			"tokens":       true,
			"results":      true,
		}),
		monitors:   make(map[base.Address]*monitor.Monitor),
		lines:      make(map[base.Address][]string),
		months:     make(map[base.Address]map[string]int),
		years:      make(map[base.Address]map[string]int),
		resolved:   make(map[base.Address]bool),
		info:       make(map[base.Address]string),
		to:         make(map[base.Address]map[string]int),
		from:       make(map[base.Address]map[string]int),
		toCounts:   make(map[int]int),
		fromCounts: make(map[int]int),
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
	Count    int
	Total    int
	Finished bool
}

func (p *Progress) Finish() {
	p.Finished = true
}

func (a *App) Reload(addressOrEns, mode string, openExcel bool) {
	logger.Info("Reload")
	addrStr, _ := a.conn.GetEnsAddress(addressOrEns)
	address := base.HexToAddress(addrStr)
	a.monitors[address] = nil
	a.Export(addressOrEns, mode, openExcel)
}

func (a *App) Export(addressOrEns, mode string, openExcel bool) {
	runtime.EventsEmit(a.ctx, "error", "")
	runtime.EventsEmit(a.ctx, "progress", "")
	runtime.EventsEmit(a.ctx, "info", "Loading...")
	runtime.EventsEmit(a.ctx, "years", "")
	runtime.EventsEmit(a.ctx, "months", "")
	runtime.EventsEmit(a.ctx, "toCount", "")
	runtime.EventsEmit(a.ctx, "fromCount", "")
	runtime.EventsEmit(a.ctx, "toTopTen", "")
	runtime.EventsEmit(a.ctx, "fromTopTen", "")

	logger.Info("Export")
	if !base.IsValidAddress(addressOrEns) {
		runtime.EventsEmit(a.ctx, "error", fmt.Sprintf("Invalid address or ENS name: %s", addressOrEns))
		return
	}

	addrStr, _ := a.conn.GetEnsAddress(addressOrEns)
	address := base.HexToAddress(addrStr)
	if address.IsZero() {
		runtime.EventsEmit(a.ctx, "error", fmt.Sprintf("No ENS address found: %s", addressOrEns))
		return
	}
	fn := "downloads/" + address.Hex() + ".csv"

	if a.monitors[address] != nil {
		a.Summarize(address, addressOrEns)
		return
	}

	if mon, err := monitor.NewMonitor("mainnet", address, true); err != nil {
		runtime.EventsEmit(a.ctx, "error", fmt.Sprintf("Error creating monitor: %s", err))
		return
	} else {
		if mon.Count() == 0 {
			getInfo(address)
		}
		a.monitors[address] = &mon
		os.Remove(fn)
		logger.Info("Count:", a.monitors[address].Count())
	}

	prog := Progress{Count: 0, Total: int(a.monitors[address].Count())}
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
	}

	logger.Info("Running command: ", cmd.String())
	_ = utils.System(cmd.String())
	logger.Info("Done...")

	a.lines[address] = file.AsciiFileToLines(fn)
	if len(a.lines[address]) == 0 {
		runtime.EventsEmit(a.ctx, "error", fmt.Sprintf("No transactions found for %s", address.Hex()))
		return
	} else {
		if openExcel {
			_ = utils.System("open " + fn)
		}
		a.Summarize(address, addressOrEns)
	}
}

func (a *App) Summarize(address base.Address, addressOrEns string) {
	if !a.resolved[address] {
		for _, line := range a.lines[address] {
			parts := strings.Split(line, ",")
			if len(parts) > 3 {
				pp := strings.Split(parts[3], "-")
				if len(pp) > 1 {
					mKey := pp[0] + "_" + pp[1]
					if a.months[address] == nil {
						a.months[address] = make(map[string]int)
					}
					a.months[address][mKey]++
					yKey := pp[0]
					if a.years[address] == nil {
						a.years[address] = make(map[string]int)
					}
					a.years[address][yKey]++
				}
			}
			if len(parts) > 5 {
				to := base.HexToAddress(parts[5])
				if a.to[address] == nil {
					a.to[address] = make(map[string]int)
				}
				if len(to.Hex()) > 10 {
					a.to[address]["to:"+to.Hex()[0:10]]++
				}
			}
			if len(parts) > 4 {
				from := base.HexToAddress(parts[4])
				if a.from[address] == nil {
					a.from[address] = make(map[string]int)
				}
				if len(from.Hex()) > 10 {
					a.from[address]["from:"+from.Hex()[0:10]]++
				}
			}
		}
		a.info[address] = strings.ToLower(addressOrEns) + "," + getInfo(address) + "," + getBalance(address)
		a.resolved[address] = true
	}

	a.toCounts = make(map[int]int)
	for _, value := range a.to[address] {
		a.toCounts[value]++
	}
	a.fromCounts = make(map[int]int)
	for _, value := range a.from[address] {
		a.fromCounts[value]++
	}

	runtime.EventsEmit(a.ctx, "info", a.info[address])
	runtime.EventsEmit(a.ctx, "years", asSortedArray1(a.years[address], len(a.years[address]), true, false))
	runtime.EventsEmit(a.ctx, "months", asSortedArray1(a.months[address], len(a.months[address]), true, false))
	runtime.EventsEmit(a.ctx, "fromTopTen", asSortedArray1(a.from[address], 25, false, false))
	runtime.EventsEmit(a.ctx, "toTopTen", asSortedArray1(a.to[address], 25, false, false))
	runtime.EventsEmit(a.ctx, "fromCount", "fromCount:"+asSortedArray2(a.fromCounts, len(a.fromCounts), true, true))
	runtime.EventsEmit(a.ctx, "toCount", "toCount:"+asSortedArray2(a.toCounts, len(a.toCounts), true, true))
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
		ss = ss[:limit]
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
		ret += fmt.Sprintf("n:%d-count:%d,", kv.Key, kv.Value)
	}
	return ret
}
