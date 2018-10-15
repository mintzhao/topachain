// Copyright © 2018 Zhao Ming <mint.zhao.chiu@gmail.com>.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mintzhao/topachain/common/genesis"
	"github.com/mintzhao/topachain/types"
	"github.com/spf13/cobra"
)

// genesisCmd represents the genesis command
var genesisCmd = &cobra.Command{
	Use:   "genesis",
	Short: "Generate the application genesis block",
	Long:  `Based on config flags to generate application the first block, aka genesis block`,
	Run: func(cmd *cobra.Command, args []string) {

		// step 1: generate genesis block
		logger.Info("generate genesis block")
		blk, err := genesis.GenesisBlock(appName, &types.AppConfig{
			BlockInterval: appBlockInterval,
			BlockTxCount:  appBlockTxCount,
			Hash:          appHash,
		})
		if err != nil {
			logger.Errorf("generate genesis block error: %s", err)
			os.Exit(-1)
		}

		// step 2: create output block file
		logger.Info("create genesis block output file")
		outfile := filepath.Join(genesisBlockOutputDir, fmt.Sprintf("%s.block", appName))
		f, err := os.Create(outfile)
		if err != nil {
			logger.Errorf("create output file error: %s", err)
			os.Exit(-1)
		}

		// step 3: write output block file
		logger.Info("write output genesis block")
		if err := genesis.WriteGenesisBlock(blk, f); err != nil {
			logger.Errorf("write genesis block file error: %s", err)
			os.Exit(-1)
		}

		logger.Info("generate genesis block done!")
	},
}

var (
	appName               string
	appBlockInterval      int64
	appBlockTxCount       int64
	appHash               string
	genesisBlockOutputDir string
)

func init() {
	rootCmd.AddCommand(genesisCmd)

	genesisCmd.Flags().StringVarP(&appName, "appName", "n", "", "the application name")
	genesisCmd.MarkFlagRequired("appName")
	genesisCmd.Flags().Int64VarP(&appBlockInterval, "blockInterval", "t", 2000, "max block interval, unit: millisecond")
	genesisCmd.Flags().Int64VarP(&appBlockTxCount, "blockTxCount", "c", 10, "max block tx count")
	genesisCmd.Flags().StringVarP(&appHash, "appHash", "", "SHA256", "application hash algorithm")
	genesisCmd.Flags().StringVarP(&genesisBlockOutputDir, "output", "o", "./", "genesis block output folder")
}
