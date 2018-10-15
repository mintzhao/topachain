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
package genesis

import (
	"io"

	"github.com/gogo/protobuf/proto"
	"github.com/mintzhao/topachain/types"
)

// generate genesis block
func GenesisBlock(name string, config *types.AppConfig) (*types.Block, error) {
	gtxp := &types.GenesisTxProposal{
		Name:   name,
		Config: config,
	}

	gtxpBytes, err := proto.Marshal(gtxp)
	if err != nil {
		return nil, err
	}

	gtxs := &types.BlockTxs{
		Txs: []*types.Transaction{
			&types.Transaction{
				Payload: gtxpBytes,
			},
		},
	}

	txroot, err := gtxs.Hash(config.GetHash())
	if err != nil {
		return nil, err
	}

	blk := &types.Block{
		Header: &types.BlockHeader{
			BlockHeight:   0,
			PreviousBlock: nil,
			Txroot:        txroot,
		},
		Txs: gtxs,
	}

	return blk, nil
}

func WriteGenesisBlock(blk *types.Block, writer io.Writer) error {
	blkBytes, err := proto.Marshal(blk)
	if err != nil {
		return err
	}

	_, err = writer.Write(blkBytes)
	return err
}
