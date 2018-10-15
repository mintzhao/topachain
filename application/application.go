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

import (
	"context"
	"time"

	"github.com/mintzhao/topachain/common/comm"
	"github.com/mintzhao/topachain/types"
)

// Application
type Application interface {
	// Metadata define Application's metadata
	Metadata() (*types.AppMetadata, error)

	// Config
	Config() *types.AppConfig
}

// Run start the app, connecting with blockchain and hold
func Run(app Application) error {
	cfg := app.Config()

	conn, err := comm.NewgRPCClient(cfg.MasterAddress)
	if err != nil {
		return err
	}
	appCli := types.NewApplicationClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if _, err := appCli.Ping(ctx, &types.Empty{}); err != nil {
		return err
	}

	stream, err := appCli.Stream(context.Background())
	if err != nil {
		return err
	}

	for {
		stream.Recv()

	}

	return nil
}
