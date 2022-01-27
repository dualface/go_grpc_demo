package config

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"time"
)

type (
	GlobalConfig struct {
		valid bool

		Server *struct {
			Listen      string        `toml:"listen"`
			AuthTimeout time.Duration `toml:"auth_timeout"`
		} `toml:"server"`

		Redis *struct {
			Enabled      bool          `toml:"enabled"`
			Addr         string        `toml:"addr"`
			Password     string        `toml:"password"`
			DB           int           `toml:"db"`
			DialTimeout  time.Duration `toml:"dial_timeout"`
			ReadTimeout  time.Duration `toml:"read_timeout"`
			WriteTimeout time.Duration `toml:"write_timeout"`
		} `toml:"redis"`

		MySQL *struct {
			Enabled      bool          `toml:"enabled"`
			Addr         string        `toml:"addr"`
			Username     string        `toml:"username"`
			Password     string        `toml:"password"`
			DB           string        `toml:"db"`
			Collation    string        `toml:"collation"`
			DialTimeout  time.Duration `toml:"dial_timeout"`
			ReadTimeout  time.Duration `toml:"read_timeout"`
			WriteTimeout time.Duration `toml:"write_timeout"`
		} `toml:"mysql"`
	}
)

var config = GlobalConfig{}

func LoadGlobalConfig(filename string, dump bool) (*GlobalConfig, error) {
	_, err := toml.DecodeFile(filename, &config)
	if err != nil {
		return nil, fmt.Errorf("load config file '%s' failed, %s", filename, err)
	}
	config.valid = true

	if dump {
		println("dump config:")
		b, _ := json.MarshalIndent(config, "", "  ")
		println(string(b))
		println()
	}

	return &config, nil
}
