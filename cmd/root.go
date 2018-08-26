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
	Short: "ToPa IS a Blockchain framework for developers to Learn Blockchain",
	Long: `A Fast and Flexible Static Site Generator built with
                love by spf13 and friends in Go.
                Complete documentation is available at http://hugo.spf13.com`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
	Version: func() string {
		v := fmt.Sprintf(": %s\n", Version)
		v += fmt.Sprintf("go   version : %s_%s_%s", runtime.Version(), runtime.GOOS, runtime.GOARCH)

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
	cobra.OnInitialize(initConfig, initLog)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.topa.yaml)")
	//rootCmd.PersistentFlags().StringVarP(&projectBase, "projectbase", "b", "", "base project directory eg. github.com/spf13/")
	//rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "Author name for copyright attribution")
	//rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "Name of license for the project (can provide `licensetext` in config)")
	//rootCmd.PersistentFlags().Bool("viper", true, "Use Viper for configuration")
	//viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	//viper.BindPFlag("projectbase", rootCmd.PersistentFlags().Lookup("projectbase"))
	//viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	//viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
	//viper.SetDefault("license", "apache")
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

		// Search config in home directory with name ".cobra" (without extension).
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
	logger    = logging.MustGetLogger("main")
	backend   = logging.NewLogBackend(os.Stderr, "", 0)
	formatter = logging.MustStringFormatter(`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`)
)

func initLog() {
	logging.SetFormatter(formatter)
	logging.SetBackend(backend)

	// set modules logging level
	logging.SetLevel(logging.INFO, "")
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
