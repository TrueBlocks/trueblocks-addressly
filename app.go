package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/rpc"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Export returns the number of transactions found at this address
func (a *App) Export(addressOrEns string) string {
	conn := rpc.NewConnection("mainnet", false, map[string]bool{})
	addr, ok := conn.GetEnsAddress(addressOrEns)
	if !ok {
		return fmt.Sprintln("GetEnsAddress returned not okay")
	}
	address := base.HexToAddress(addr)
	cmd := "chifra list --count " + address.Hex() + " --fmt csv --no_header --output file.fil"
	_ = utils.System(cmd)
	value := file.AsciiFileToString("file.fil")
	s := strings.Split(value, ",")
	n := address.Hex()
	if strings.ToLower(addressOrEns) != n {
		n = addressOrEns + " (" + address.Hex() + ")"
	}
	return fmt.Sprintf("%s has %s appearances", n, s[1]) //10)
}
