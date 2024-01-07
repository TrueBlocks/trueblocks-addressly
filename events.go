package main

import (
	"fmt"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
)

func (a *App) SetAddress(addressOrEns string) {
	addrStr, _ := a.conn.GetEnsAddress(addressOrEns)
	a.dataFile.Address = base.HexToAddress(addrStr)
	logger.Info("Setting query to:", addressOrEns)
	logger.Info("Setting address to:", addrStr)
}

func (a *App) SetChain(chain string, addressOrEns string) {
	a.dataFile.Chain = chain
	a.Export(addressOrEns, "")
	logger.Info("Setting chain to:", chain)
}

func (a *App) SetChartType(chartType string) {
	a.config.Set("settings", "chartType", chartType)
}

func (a *App) SetExportExcel(exportExcel bool) {
	a.config.Set("settings", "exportExcel", fmt.Sprintf("%t", exportExcel))
}
