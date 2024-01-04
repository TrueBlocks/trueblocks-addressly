package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

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
	fn := "/tmp/latest"
	defer os.Remove(fn)
	cmd := Command{
		MaxRecords: int(maxRecords),
		Filename:   fn,
		Format:     "csv",
		Subcommand: "when",
		Rest:       "latest --no_header",
		Silent:     true,
		Chain:      a.dataFile.Chain,
	}
	_ = utils.System(cmd.String())
	contents := file.AsciiFileToString(fn)
	parts := strings.Split(contents, ",")
	state := "||" + strings.Trim(price, " ")
	if len(parts) > 2 {
		state = parts[0] + "|" + parts[2] + "|" + price
	} else if len(parts) > 1 {
		state = parts[0] + "||" + price
	}

	state += "|" + a.dataFile.Chain
	logger.Info("Sending state: ", state)
	runtime.EventsEmit(a.ctx, "chainState", state)
}

func (a *App) getBalance() string {
	fn := "/tmp/" + a.dataFile.Chain + "_" + a.dataFile.Address.Hex() + ".balance"
	defer os.Remove(fn)

	cmd := Command{
		MaxRecords: int(maxRecords),
		Address:    a.dataFile.Address,
		Filename:   fn,
		Format:     "csv",
		Subcommand: "state",
		Rest:       "--ether --no_header",
		Silent:     true,
		Chain:      a.dataFile.Chain,
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

func (a *App) getInfo(addressOrEns string) string {
	fn := "/tmp/" + a.dataFile.Chain + "_" + a.dataFile.Address.Hex() + ".info"
	defer os.Remove(fn)

	cmd := Command{
		MaxRecords: int(maxRecords),
		Address:    a.dataFile.Address,
		Filename:   fn,
		Format:     "csv",
		Subcommand: "list",
		Rest:       "--bounds --no_header",
		Silent:     true,
		Chain:      a.dataFile.Chain,
	}
	logger.Info("Running command: ", cmd.String())
	_ = utils.System(cmd.String())
	contents := strings.ToLower(addressOrEns) + "," + file.AsciiFileToString(fn) + "," + a.getBalance()
	parts := strings.Split(contents, ",")
	if len(parts) < 11 {
		return ""
	}
	parts[0] = a.namesMap[a.dataFile.Chain][a.dataFile.Address].Name
	parts[3] = parts[3] + " " + parts[5]
	parts[4] = parts[6] + " " + parts[8]
	parts[5] = parts[11]
	parts[6] = parts[9]
	parts[7] = parts[10]
	parts[8] = ""
	return strings.Join(parts[:9], ",")
	// return file.AsciiFileToString(fn)
	//trueblocks.eth,0xf503017d7baf7fbc0fff7492b751025c6a78179b,4158,8854723.61,1572639538,2019-11-01 20:18:58 UTC,18752751.88,1702174307,2023-12-10 02:11:47 UTC,9898028,2380 ,71.713067671079880299
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

	return fmt.Sprintf("%-10.2f", apiResponse.Ethereum.USD), nil
}
