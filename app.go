package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/rpc"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
)

// App struct
type App struct {
	ctx  context.Context
	conn *rpc.Connection
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		conn: rpc.NewConnection("mainnet", false, map[string]bool{}),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

var maxRecords = utils.NOPOSI

// Export returns the number of transactions found at this address
func (a *App) Export(addressOrEns string) string {
	if len(addressOrEns) == 0 {
		addressOrEns = "trueblocks.eth"
	}
	if !base.IsValidAddress(addressOrEns) {
		return fmt.Sprintf("Invalid address or ENS name: %s", addressOrEns)
	}
	addrStr, _ := a.conn.GetEnsAddress(addressOrEns)
	address := base.HexToAddress(addrStr)
	if strings.HasSuffix(addressOrEns, ".eth") {
		addrStr = addressOrEns + " (" + address.Hex() + ")"
	}

	fn := "downloads/" + address.Hex() + ".csv"
	cmd := Command{
		MaxRecords: int(maxRecords),
		Address:    address,
		// Filename:   fn,
		Format: "csv",
		Rest:   "--count --no_header | cut -f2 -d, > " + fn}

	logger.Info("Running command: ", cmd.String())
	_ = utils.System(cmd.String())

	// lines := file.AsciiFileToLines(fn)
	// cnt := len(lines)
	contents := strings.Trim(file.AsciiFileToString(fn), "\n")
	cnt := utils.Min(maxRecords, utils.MustParseInt(contents))

	return fmt.Sprintf("%s has\n%d appearances", addrStr, cnt)
}
