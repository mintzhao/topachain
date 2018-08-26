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
package application

import "github.com/mintzhao/topachain/types"

// Application
type Application interface {
	// Metadata define Application's metadata
	Metadata() (*types.AppMetadata, error)

	// Config
	Config() *types.AppConfig
}

// Run start the app, connecting with blockchain and hold
func Run(app Application) error {
	//cfg := app.Config()
	//
	//conn, err := net.Dial("tcp", cfg.MasterAddress)
	//if err != nil {
	//	return err
	//}

	return nil
}
