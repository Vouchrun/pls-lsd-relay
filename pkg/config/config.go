// Copyright 2021 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Eth1Endpoint             string // url for eth1 rpc endpoint
	Eth2Endpoint             string // url for eth2 rpc endpoint
	Web3StorageApiToken      string //
	DirPath                  string
	LogFilePath              string
	Account                  string
	KeystorePath             string
	BlockstoreFilePath       string
	GasLimit                 string
	MaxGasPrice              string
	BatchRequestBlocksNumber uint64

	RunForEntrustedLsdNetwork bool

	Contracts Contracts
}

type Contracts struct {
	LsdTokenAddress   string
	LsdFactoryAddress string
}

func Load(configFilePath string) (*Config, error) {
	var cfg = Config{}
	if err := loadSysConfig(configFilePath, &cfg); err != nil {
		return nil, err
	}
	cfg.DirPath = strings.TrimSuffix(cfg.DirPath, "/")
	if cfg.DirPath == "" {
		return nil, fmt.Errorf("dirPath must be set")
	}
	cfg.LogFilePath = cfg.DirPath + "/log_data/"
	cfg.KeystorePath = cfg.DirPath + "/keystore/"
	cfg.BlockstoreFilePath = cfg.DirPath + "/blockstore"

	return &cfg, nil
}

func loadSysConfig(path string, config *Config) error {
	_, err := os.Open(path)
	if err != nil {
		return err
	}
	if _, err := toml.DecodeFile(path, config); err != nil {
		return err
	}
	fmt.Println("load config success")
	return nil
}
