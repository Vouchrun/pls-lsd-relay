// Copyright 2021 stafiprotocol
// SPDX-License-Identifier: LGPL-3.0-only

package cmd

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
	"github.com/stafiprotocol/chainbridge/utils/keystore"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/config"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/log"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/utils"
	"github.com/stafiprotocol/eth-lsd-relay/service"
)

func startRelayCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start lsd relay",
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath, err := cmd.Flags().GetString(flagConfigPath)
			if err != nil {
				return err
			}
			fmt.Printf("config path: %s\n", configPath)

			logLevelStr, err := cmd.Flags().GetString(flagLogLevel)
			if err != nil {
				return err
			}
			logLevel, err := logrus.ParseLevel(logLevelStr)
			if err != nil {
				return err
			}
			logrus.SetLevel(logLevel)

			cfg, err := config.Load(configPath)
			if err != nil {
				return err
			}
			logrus.Infof(
				`config info:
  logFilePath: %s
  logLevel: %s
  eth1Endpoint: %s
  eth2Endpoint: %s
  account: %s
  lsdTokenAddress: %s
  factoryAddress: %s
  batchRequestBlocksNumber: %d`,
				cfg.LogFilePath, logLevelStr, cfg.Eth1Endpoint, cfg.Eth2Endpoint, cfg.Account,
				cfg.Contracts.LsdTokenAddress, cfg.Contracts.LsdFactoryAddress, cfg.BatchRequestBlocksNumber)

			err = log.InitLogFile(cfg.LogFilePath + "/relay")
			if err != nil {
				return fmt.Errorf("InitLogFile failed: %w", err)
			}

			//interrupt signal
			ctx := utils.ShutdownListener()

			// load trust node account
			kpI, err := keystore.KeypairFromAddress(cfg.Account, keystore.EthChain, cfg.KeystorePath, false)
			if err != nil {
				return err
			}
			kp, ok := kpI.(*secp256k1.Keypair)
			if !ok {
				return fmt.Errorf(" keypair err")
			}

			t, err := service.NewService(cfg, kp)
			if err != nil {
				return fmt.Errorf("NewService err: %w", err)
			}

			err = t.Start()
			if err != nil {
				logrus.Errorf("start err: %s", err)
				return err
			}
			defer func() {
				logrus.Infof("shutting down task ...")
				t.Stop()
			}()

			<-ctx.Done()
			return nil
		},
	}

	cmd.Flags().String(flagConfigPath, defaultConfigPath, "Config file path")
	cmd.Flags().String(flagLogLevel, logrus.InfoLevel.String(), "The logging level (trace|debug|info|warn|error|fatal|panic)")

	return cmd
}
