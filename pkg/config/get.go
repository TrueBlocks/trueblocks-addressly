package config

import (
	"fmt"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
)

func (c *Config) Get(group, name, def string) (string, error) {
	if group == "settings" {
		switch name {
		case "chartType":
			fallthrough
		case "exportExcel":
			fn := GetCacheFolder("") + name + ".txt"
			value := strings.Trim(file.AsciiFileToString(fn), "\n")
			if value == "" {
				value = def
			}
			logger.Info("Get", group, name, def, value)
			return value, nil
		default:
			return def, fmt.Errorf("unknown setting: %s", name)
		}
	}
	return def, nil

	// windowState = config.WindowState{
	// 	Width:  1024,
	// 	Height: 1000,
	// }

	// fn := GetCacheFolder("") + "windowstate.json"
	// if state, err := os.ReadFile(fn); err != nil {
	// 	// logger.Error("Error reading window state", err)
	// } else {
	// 	json.Unmarshal(state, &app.windowState)
	// }

	// func ReadConfigItem(group, name string) (string, error) {
	// fn := GetCacheFolder("") + "windowstate.json"
	// if state, err := os.ReadFile(fn); err != nil {
	// 	// logger.Error("Error reading window state", err)
	// } else {
	// 	json.Unmarshal(state, &app.windowState)
	// }

}
