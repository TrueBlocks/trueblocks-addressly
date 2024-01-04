package main

import (
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
	a.Export(addressOrEns, "", false)
	logger.Info("Setting chain to:", chain)
}

func (a *App) SetChartType(chartType string) {
	// a.chartType = chartType
	logger.Info("Chart type changed to:", chartType)
}

func (a *App) SetExportExcel(onOff bool) {
	// a.exportExcel = onOff
	logger.Info("ExportExcel changed to: ", onOff)
}
