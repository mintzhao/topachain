// Copyright © 2018 Zhao Ming <mint.zhao.chiu@gmail.com>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
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
	"runtime"

	"github.com/mitchellh/go-homedir"
	"github.com/op/go-logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "topa",
	Short: "Another Enterprise-Level Consortium Blackchain",
	Long:  `Topa is Another Enterprise-Level Consortium Blackchain.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
	Version: func() string {
		v := fmt.Sprintf("%s, powered by %s_%s_%s", Version, runtime.Version(), runtime.GOOS, runtime.GOARCH)

		return v
	}(),
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var (
	cfgFile string
	Version string
)

func init() {
	cobra.OnInitialize(initLog)
	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.topa.yaml)")
}

func initConfig() {
	// Don't forget to read config either from cfgFile or from home directory!
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".topa" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".topa")
	}

	viper.SetEnvPrefix("TOPA")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}

	viper.Debug()
}

var (
	logger = logging.MustGetLogger("cmd")
)

func initLog() {
	loggings := viper.GetStringMap("logging")
	for module, levelInterface := range loggings {
		levelstring, ok := levelInterface.(string)
		if !ok {
			continue
		}

		level, err := logging.LogLevel(levelstring)
		if err != nil {
			logger.Warningf("invalid log level: %s", levelstring)
			continue
		}

		logging.SetLevel(level, module)
	}

	logger.Info("inited log")
}
