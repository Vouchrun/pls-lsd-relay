// Copyright 2021 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

type Endpoint struct {
	Eth1 string
	Eth2 string
}

type Config struct {
	LogFilePath                string
	Account                    string
	KeystorePath               string
	BlockstoreFilePath         string
	GasLimit                   string
	MaxGasPrice                string // Gwei
	GasPriceMultiplier         float64
	BatchRequestBlocksNumber   uint64
	EventFilterMaxSpanBlocks   uint64
	MaxEjectedValPerCycle      int
	TrustNodeDepositAmount     uint64 // ether
	Eth2EffectiveBalance       uint64 // ether
	MaxPartialWithdrawalAmount uint64 // ether

	RunForEntrustedLsdNetwork bool
	TransferFeeAddresses      []string

	TotalUndistributedTransferFee      uint64 // unit ether
	DistributeBlockedTransferFeePerEra uint64 // unit ether

	Contracts   Contracts
	Endpoints   []Endpoint
	Web3Storage Web3Storage
	Pinata      Pinata
}

type Web3Storage struct {
	PrivateKey string
	SpaceDid   string
	ProofFile  string
}

type Pinata struct {
	Apikey   string
	Endpoint string
	PinDays  uint
}

type Contracts struct {
	LsdTokenAddress   string
	LsdFactoryAddress string
}

func Load(basePath string) (*Config, error) {
	basePath = strings.TrimSuffix(basePath, "/")
	configFilePath := basePath + "/config.toml"
	fmt.Printf("config path: %s\n", configFilePath)

	var cfg = Config{}
	if err := loadSysConfig(configFilePath, &cfg); err != nil {
		return nil, err
	}
	cfg.LogFilePath = basePath + "/log_data"
	cfg.KeystorePath = KeyStoreFilePath(basePath)
	cfg.BlockstoreFilePath = basePath + "/blockstore"

	// add default values
	if cfg.TrustNodeDepositAmount == 0 {
		cfg.TrustNodeDepositAmount = 1
	}
	if cfg.Eth2EffectiveBalance == 0 {
		cfg.Eth2EffectiveBalance = 32
	}
	if cfg.MaxPartialWithdrawalAmount == 0 {
		cfg.MaxPartialWithdrawalAmount = 8
	}
	if cfg.GasLimit == "" {
		cfg.GasLimit = "3000000"
	}
	if cfg.MaxGasPrice == "" {
		cfg.MaxGasPrice = "600"
	}
	if cfg.GasPriceMultiplier == 0 {
		cfg.GasPriceMultiplier = 1
	}
	if cfg.EventFilterMaxSpanBlocks == 0 {
		cfg.EventFilterMaxSpanBlocks = 3000
	}

	// handle invalid parameters
	if cfg.GasPriceMultiplier < 1 {
		return nil, fmt.Errorf("gas price multiplier can not be less than 1")
	}
	if cfg.GasPriceMultiplier > 10 {
		return nil, fmt.Errorf("gas price multiplier can not be greater than 10")
	}
	if cfg.BatchRequestBlocksNumber > 32 {
		return nil, fmt.Errorf("batchRequestBlocksNumber can not be greater than 32")
	}

	if cfg.TotalUndistributedTransferFee == 0 {
		cfg.TotalUndistributedTransferFee = 110_000_000
	}
	if cfg.DistributeBlockedTransferFeePerEra == 0 {
		cfg.DistributeBlockedTransferFeePerEra = // todo
	}

	return &cfg, nil
}

func KeyStoreFilePath(basePath string) string {
	basePath = strings.TrimSuffix(basePath, "/")
	return basePath + "/keystore"
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
