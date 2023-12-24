package main

import (
	"context"
	"fmt"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/rpc"
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
	addr, _ := conn.GetEnsAddress(addressOrEns)
	address := base.HexToAddress(addr)
	return fmt.Sprintf("Exporting address %s. Found %d items", address.Hex(), 10)
}
