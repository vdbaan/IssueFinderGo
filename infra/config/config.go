package config

import (
	"encoding/json"
	"path/filepath"

	"github.com/shibukawa/configdir"
)

type Config struct {
	Filters []string
}

type Handler interface {
	GetConfig() *Config
	UpdateConfig(config *Config)
	SaveConfig()
}

type ConfigHandler struct {
	config *Config
}

func NewHandler() Handler {
	h := new(ConfigHandler)
	return h
}

func (h *ConfigHandler) GetConfig() *Config {

	if h.config == nil {
		DefaultConfig := Config{
			Filters: []string{
				"IP != 127.0.0.1",
				"SCANNER == 'nmap'",
				"SERVICE LIKE 'http'",
				"PORT LIKE 443",
				"!EXPLOITABLE",
				"(SERVICE LIKE 'SMB') && EXPLOITABLE",
			},
		}

		configDirs := configdir.New("vdbaan", "issuefinder")
		configDirs.LocalPath, _ = filepath.Abs(".")
		folder := configDirs.QueryFolderContainsFile("settings.json")
		if folder != nil {
			data, _ := folder.ReadFile("setting.json")
			json.Unmarshal(data, &h.config)
			if h.config == nil {
				h.config = &DefaultConfig
			}
		} else {
			h.config = &DefaultConfig
		}
	}
	return h.config
}

func (h *ConfigHandler) UpdateConfig(cfg *Config) {
	h.config = cfg
}

func (h *ConfigHandler) SaveConfig() {
	data, _ := json.Marshal(&h.config)
	configDirs := configdir.New("vdbaan", "issuefinder")
	configDirs.LocalPath, _ = filepath.Abs(".")
	folder := configDirs.QueryFolderContainsFile("settings.json")
	if folder == nil {
		folders := configDirs.QueryFolders(configdir.Global)
		folder = folders[0]
	}
	folder.WriteFile("issuefinder.json", data)
}
