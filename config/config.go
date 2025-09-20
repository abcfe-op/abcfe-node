package config

import (
	"fmt"
	"os"
	"path"

	"github.com/abcfe-op/abcfe-node/common/utils"
	"github.com/naoina/toml"
)

type Common struct {
	Mode        string
	ServiceName string
}

type LogInfos struct {
	Fpath      string
	MaxAgeHour int
	RotateHour int
	ProdTelKey string
	ProdChatId int64
	DevTelKey  string
	DevChatId  int64
}

type Config struct {
	Common  Common
	LogInfo LogInfos
}

func NewConfig(filepath string) *Config {
	if filepath == "" {
		fmt.Println(os.Getwd())
		filepath = "./config/config.toml"
	}
	// fmt.Println(os.Getwd())
	if file, err := os.Open(filepath); err != nil {
		return nil
	} else {
		defer file.Close()

		c := new(Config)
		if err := toml.NewDecoder(file).Decode(c); err != nil {
			return nil
		} else {
			c.sanitize()
			return c
		}
	}
}

func (p *Config) sanitize() {
	if p.LogInfo.Fpath[0] == byte('~') {
		p.LogInfo.Fpath = path.Join(utils.HomeDir(), p.LogInfo.Fpath[1:])
	}
}

func (p *Config) GetConfig() *Config {
	return p
}

func (p *Config) GetLogInfoConfig() *LogInfos {
	return &p.LogInfo
}
