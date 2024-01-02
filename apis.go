package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type ApiResponse struct {
	Ethereum struct {
		USD float64 `json:"usd"`
	} `json:"ethereum"`
}

var price = 0.0
var latest = 0

func (a *App) updateState() {
	if !initialized {
		return
	}

	var err error
	prev := price
	if price, err = getEthUsdPrice(); err != nil {
		runtime.EventsEmit(a.ctx, "price", prev)
	} else {
		runtime.EventsEmit(a.ctx, "price", price)
	}
	blk := latest
	if latest, err = getLatestBlock(); err != nil {
		runtime.EventsEmit(a.ctx, "status", err)
		runtime.EventsEmit(a.ctx, "latest", blk)
	} else {
		runtime.EventsEmit(a.ctx, "latest", latest)
	}
}

func getBalance(address base.Address) string {
	fn := "/tmp/" + address.Hex() + ".balance"
	defer os.Remove(fn)

	cmd := Command{
		MaxRecords: int(maxRecords),
		Address:    address,
		Filename:   fn,
		Format:     "csv",
		Subcommand: "state",
		Rest:       "--ether --no_header",
		Silent:     true,
	}
	logger.Info("Running command: ", cmd.String())
	_ = utils.System(cmd.String())
	res := file.AsciiFileToString(fn)
	parts := strings.Split(res, ",")
	if len(parts) < 3 {
		return "0"
	}
	return parts[2]
}

func getInfo(address base.Address) string {
	fn := "/tmp/" + address.Hex() + ".info"
	defer os.Remove(fn)

	cmd := Command{
		MaxRecords: int(maxRecords),
		Address:    address,
		Filename:   fn,
		Format:     "csv",
		Subcommand: "list",
		Rest:       "--bounds --no_header",
		Silent:     true,
	}
	logger.Info("Running command: ", cmd.String())
	_ = utils.System(cmd.String())
	return file.AsciiFileToString(fn)
}

func getLatestBlock() (int, error) {
	fn := "/tmp/latest"
	defer os.Remove(fn)

	cmd := Command{
		MaxRecords: int(maxRecords),
		Filename:   fn,
		Format:     "csv",
		Subcommand: "when",
		Rest:       "latest --no_header",
		Silent:     true,
	}

	// logger.Info("Running command: ", cmd.String())
	_ = utils.System(cmd.String())
	contents := file.AsciiFileToString(fn)
	parts := strings.Split(contents, ",")
	// logger.Info(parts)
	// logger.Info(utils.MustParseInt(parts[0]))
	return int(utils.MustParseInt(parts[0])), nil
}

func getEthUsdPrice() (float64, error) {
	prev := price
	url := "https://api.coingecko.com/api/v3/simple/price?ids=ethereum&vs_currencies=usd"

	resp, err := http.Get(url)
	if err != nil {
		return prev, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return prev, err
	}

	var apiResponse ApiResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return prev, err
	}

	return apiResponse.Ethereum.USD, err
}
