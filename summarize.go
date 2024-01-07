package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
)

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
					// logger.Info("Here", mKey, a.months[dfKey][mKey])
					yKey := pp[0]
					if a.years[dfKey] == nil {
						a.years[dfKey] = make(map[string]int)
					}
					a.years[dfKey][yKey]++
					// logger.Info("Here", yKey, a.months[dfKey][yKey])
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
