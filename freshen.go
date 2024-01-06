package main

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/names"
)

func (a *App) Freshen(addressOrEns, mode string) {
	addrStr, _ := a.conn.GetEnsAddress(addressOrEns)
	a.dataFile.Address = base.HexToAddress(addrStr)
	dfKey := a.dataFile
	a.monitors[dfKey] = nil
	nameParts := names.Custom | names.Prefund | names.Regular | names.Baddress
	m, err := names.LoadNamesMap(a.dataFile.Chain, nameParts, nil)
	a.namesMap[a.dataFile.Chain] = m
	if err != nil {
		logger.Error("Failed to load names map:", err)
	}
	exportToExcel = file.AsciiFileToString("exportExcel.txt") == "true"
	a.Export(addressOrEns, mode)
	exportToExcel = false
}
