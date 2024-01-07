package config

import (
	"encoding/json"
	"os"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
)

var defaultConfig = Config{
	WindowState: WindowState{
		Width:  1400,
		Height: 1000,
		X:      305,
		Y:      0,
	},
}

func (c *Config) Load() error {
	*c = defaultConfig

	fn := GetCacheFolder("") + "config.json"
	if file.FileExists(fn) {
		if contents, err := os.ReadFile(fn); err != nil {
			return err
		} else {
			if err = json.Unmarshal(contents, &c); err != nil {
				return err
			}
		}
	}

	return nil
}
