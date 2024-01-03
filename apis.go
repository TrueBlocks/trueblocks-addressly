package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
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

func (a *App) updateState() {
	if !initialized {
		return
	}

	var price string
	var err error
	if price, err = getEthUsdPrice(); err != nil {
		logger.Info(colors.Red, "Error fetching price: ", err, colors.Off)
	}
	logger.Info(colors.Yellow, "Got price okay: ", price, colors.Off)
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
	_ = utils.System(cmd.String())
	contents := file.AsciiFileToString(fn)
	parts := strings.Split(contents, ",")
	state := "||" + price
	if len(parts) > 2 {
		state = parts[0] + "|" + parts[2] + "|" + price
	} else if len(parts) > 1 {
		state = parts[0] + "||" + price
	}
	logger.Info("Sending state: ", state)
	runtime.EventsEmit(a.ctx, "chainState", state)
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

func getEthUsdPrice() (string, error) {
	url := "https://api.coingecko.com/api/v3/simple/price?ids=ethereum&vs_currencies=usd"

	resp, err := http.Get(url)
	if err != nil {
		logger.Error("Error 1: ", err)
		return "", err
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("price server returned error %d %s", resp.StatusCode, resp.Status)
		logger.Error("Error 4: ", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Error 2: ", err)
		return "", err
	}

	var apiResponse ApiResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		logger.Error("Error 3: ", err)
		return "", err
	}

	logger.Info("fetched price: ", apiResponse.Ethereum.USD)
	return fmt.Sprintf("%-10.2f", apiResponse.Ethereum.USD), nil
}
