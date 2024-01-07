package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/filter"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/monitor"
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
		// logger.Info("Sending state: ", state)
		runtime.EventsEmit(a.ctx, "chainState", state)
	}
}

func (a *App) getInfo(addressOrEns string) string {
	filter := filter.NewFilter(
		false,
		false,
		[]string{},
		base.BlockRange{First: 0, Last: utils.NOPOS},
		base.RecordRange{First: 0, Last: utils.NOPOS},
	)

	chain := a.dataFile.Chain
	addrs := []string{a.dataFile.Address.Hex()}
	monitorArray := make([]monitor.Monitor, 0, len(addrs))
	updater := monitor.NewUpdater(chain, false, true, addrs)
	if canceled, err := updater.FreshenMonitors(&monitorArray); err != nil || canceled {
		logger.Error("Error fetching monitor: ", err)
		return ""
	} else {
		mon := monitorArray[0]
		if apps, cnt, err := mon.ReadAndFilterAppearances(filter, true); err != nil {
			logger.Error("Error fetching monitor: ", err)
			return ""
		} else if cnt == 0 {
			logger.Error("no appearances found for", mon.Address.Hex())
			return ""
		} else {
			firstApp := apps[0]
			latestApp := apps[len(apps)-1]
			rng := latestApp.BlockNumber - firstApp.BlockNumber
			name := a.namesMap[a.dataFile.Chain][a.dataFile.Address].Name
			if name == "" {
				name = addressOrEns
			}
			ethBalance, usdBalance := a.getPrices()
			res := struct {
				name      string
				addr      string
				cnt       string
				firstApp  string
				latestApp string
				ethBal    string
				blkRng    string
				blkFreq   string
				usdBal    string
			}{
				name:      name,
				addr:      a.dataFile.Address.Hex(),
				cnt:       fmt.Sprintf("%d", cnt),
				firstApp:  fmt.Sprintf("%d.%d %s", firstApp.BlockNumber, firstApp.TransactionIndex, firstApp.Date()),
				latestApp: fmt.Sprintf("%d.%d %s", latestApp.BlockNumber, latestApp.TransactionIndex, latestApp.Date()),
				ethBal:    ethBalance,
				blkRng:    fmt.Sprintf("%d", rng),
				blkFreq:   fmt.Sprintf("%d", int64(float64(rng)/float64(cnt))),
				usdBal:    usdBalance,
			}

			return fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,%s",
				strings.Replace(res.name, ",", " ", -1),
				res.addr,
				res.cnt,
				res.firstApp,
				res.latestApp,
				res.ethBal,
				res.blkRng,
				res.blkFreq,
				res.usdBal,
			)
		}
	}
}

func (a *App) getPrices() (string, string) {
	ethBalance := ""
	usdBalance := ""

	if price, err := getEthUsdPrice(); err != nil {
		logger.Error(colors.Red, "error fetching price: ", err, colors.Off)

	} else {
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
	}

	return ethBalance, usdBalance
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

	body, err := io.ReadAll(resp.Body)
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
