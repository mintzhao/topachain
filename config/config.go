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
package config

import (
	"strings"

	"github.com/op/go-logging"
	"github.com/spf13/viper"
)

// Config
type Config struct {
	Logging map[string]string
	Common  *Common
}

// Common
type Common struct {
	Crypto   *Crypto
	Database *Database
}

// Crypto
type Crypto struct {
	Hash string
	Sign string
}

// Database
type Database struct {
	Type   string
	Badger *Badger
}

// Badger
type Badger struct {
	Dir string
}

var (
	// logger
	logger = logging.MustGetLogger("config")

	defaults = &Config{}
)

func Load(configfile string) *Config {
	if configfile == "" {
		return defaults
	}

	v := viper.New()
	v.SetEnvPrefix("TOPA")
	v.AutomaticEnv()
	v.SetConfigFile(configfile)
	replacer := strings.NewReplacer(".", "_")
	v.SetEnvKeyReplacer(replacer)

	if err := v.ReadInConfig(); err != nil {
		logger.Panicf("reading config error: %s", err)
	}

	conf := new(Config)
	if err := v.Unmarshal(conf); err != nil {
		logger.Panicf("unmarshal config error: %s", err)
	}

	return conf
}
