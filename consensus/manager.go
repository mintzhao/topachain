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
package consensus

import (
	"sync"

	"github.com/mintzhao/topachain/types"
	"github.com/pkg/errors"
)

var (
	// ErrNilTx indicate empty tx
	ErrNilTx = errors.New("nil tx")

	// ErrApplicationUnregistered means can not load application core
	ErrApplicationUnregistered = errors.New("application unregistered")
)

type handler struct {
	core   ConsensusCore
	stream types.Application_AppStreamServer
}

// Consensus Manager
type Manager struct {
	handlers sync.Map
}

// ReceiveTxSync receive tx from application synchronous.
func (m *Manager) ReceiveTxSync(application string, tx []byte) (*types.TxResponseSync, error) {
	if len(tx) == 0 {
		return nil, ErrNilTx
	}

	// get specific application consensus handler
	handlerInterface, ok := m.handlers.Load(application)
	if !ok {
		return nil, ErrApplicationUnregistered
	}
	handler := handlerInterface.(*handler)

	// valid tx
	handler.stream.Send()
}
