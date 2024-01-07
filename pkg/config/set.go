package config

import (
	"fmt"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
)

func (c *Config) Set(group, name, value string) error {
	if group == "settings" {
		switch name {
		case "chartType":
			fallthrough
		case "exportExcel":
			file.StringToAsciiFile(GetCacheFolder("")+name+".txt", value+"\n")
		default:
			return fmt.Errorf("unknown setting: %s", name)
		}
	}
	return nil
}
