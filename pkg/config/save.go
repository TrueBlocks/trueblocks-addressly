package config

import (
	"encoding/json"
	"os"
)

type WindowState struct {
	Width  int `json:"width"`
	Height int `json:"height"`
	X      int `json:"x"`
	Y      int `json:"y"`
}

func (ws WindowState) String() string {
	bytes, _ := json.MarshalIndent(ws, "", "  ")
	return string(bytes)
}

type Config struct {
	WindowState WindowState `json:"windowState"`
}

func (c Config) String() string {
	bytes, _ := json.Marshal(c)
	return string(bytes)
}

func (c *Config) Save() error {
	if state, err := json.Marshal(c); err != nil {
		return err
	} else {
		fn := GetCacheFolder("") + "config.json"
		return os.WriteFile(fn, append(state, '\n'), 0644)
	}
}
