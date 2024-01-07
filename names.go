package main

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
)

func (a *App) ChangeName(newName, address string) {
	n := a.namesMap[a.dataFile.Chain][base.HexToAddress(address)]
	n.Name = newName
	a.namesMap[a.dataFile.Chain][base.HexToAddress(address)] = n

	dfKey := a.dataFile
	a.resolved[dfKey] = false
	a.months[dfKey] = make(map[string]int)
	a.years[dfKey] = make(map[string]int)
	a.asSender[dfKey] = make(map[string]int)
	a.asRecipient[dfKey] = make(map[string]int)
	a.Summarize(a.dataFile.Chain)
}
