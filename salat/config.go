package salat

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/alfajrimutawadhi/salat/common"
)

type Config struct {
	Location Location `json:"location"`
	TimeMode int8     `json:"time_mode"` // 1 = 12 Hours; 2 = 24 Hours
}

type Location struct {
	Country string `json:"country"`
	City    string `json:"city"`
}

func ReadConfig(path string) *Config {
	var c Config
	cf, err := os.ReadFile(fmt.Sprintf("%s/.salat/config.json", path))
	if err != nil {
		common.HandleError(err)
	}
	json.Unmarshal(cf, &c)
	return &c
}
