package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type ApiResponse struct {
	Ethereum struct {
		USD float64 `json:"usd"`
	} `json:"ethereum"`
}

func (a *App) updateState() {
	var price float64
	var err error
	if price, err = getEthUsdPrice(); err != nil {
		logger.Info(colors.Red, "Error fetching price: ", err, colors.Off)
	}

	bn := a.conn.GetLatestBlockNumber()
	if block, err := a.conn.GetBlockHeaderByNumber(bn); err != nil {
		logger.Info(colors.Red, "Error fetching price: ", err, colors.Off)
	} else {
		s := types.SimpleNamedBlock{
			BlockNumber: bn,
			Timestamp:   block.Timestamp,
		}
		state := fmt.Sprintf("%d|%s|%.2f", s.BlockNumber, s.Date(), price)
		logger.Info("Sending state: ", state)
		runtime.EventsEmit(a.ctx, "chainState", state)
	}
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
	// logger.Info("Running command: ", cmd.String())
	_ = utils.System(cmd.String())

	var price float64
	var err error
	if price, err = getEthUsdPrice(); err != nil {
		logger.Info(colors.Red, "Error fetching price: ", err, colors.Off)
	}

	ethBalance := ""
	usdBalance := ""
	if bal, err := a.conn.GetBalanceAt(a.dataFile.Address, a.conn.GetLatestBlockNumber()); err != nil {
		ethBalance = "0.000000000000000000"
		usdBalance = "0.00"
	} else {
		ethBalance = utils.FormattedValue(*bal, true, 18)
		eb := big.Float{}
		eb.SetString(ethBalance)
		p := big.Float{}
		p.SetFloat64(price)
		eb.Mul(&eb, &p)
		usdBalance = eb.Text('f', 2)
	}

	contents := strings.ToLower(addressOrEns) + "," + file.AsciiFileToString(fn) + "," + ethBalance + "," + usdBalance
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
	parts[8] = parts[12]
	return strings.Join(parts[:9], ",")
	// return file.AsciiFileToString(fn)
	//trueblocks.eth,0xf503017d7baf7fbc0fff7492b751025c6a78179b,4158,8854723.61,1572639538,2019-11-01 20:18:58 UTC,18752751.88,1702174307,2023-12-10 02:11:47 UTC,9898028,2380 ,71.713067671079880299
}

var lastPrice float64
var lastFetched int64
var m sync.Mutex

func getEthUsdPrice() (float64, error) {
	now := time.Now().Unix()
	if lastPrice != 0.0 {
		if now-lastFetched < (60 * 5) {
			return lastPrice, nil
		}
	}

	m.Lock()
	defer m.Unlock()

	lastFetched = now

	url := "https://api.coingecko.com/api/v3/simple/price?ids=ethereum&vs_currencies=usd"

	resp, err := http.Get(url)
	if err != nil {
		logger.Error("Error 1: ", err)
		return lastPrice, err
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("price server returned error %d %s", resp.StatusCode, resp.Status)
		logger.Error("Error 4: ", err)
		return lastPrice, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Error 2: ", err)
		return lastPrice, err
	}

	var apiResponse ApiResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		logger.Error("Error 3: ", err)
		return lastPrice, err
	}

	lastPrice = apiResponse.Ethereum.USD
	return apiResponse.Ethereum.USD, nil
}

func getCacheDir() string {
	dirName, _ := os.UserCacheDir()
	dirName += "/TrueBlocks/browse/"
	if !file.FolderExists(dirName) {
		file.EstablishFolder(dirName)
	}
	logger.Info("Cache dir: ", dirName, file.FolderExists(dirName))
	return dirName
}

func getConfigDir() string {
	dirName, _ := os.UserConfigDir()
	dirName += "/TrueBlocks/browse/"
	if !file.FolderExists(dirName) {
		file.EstablishFolder(dirName)
	}
	logger.Info("Config dir: ", dirName)
	return dirName
}
