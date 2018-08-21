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
package merkle

// MerkleTree abstract funcs
type MerkleTree interface {

	// Root returns merkletree's root hash
	Root() []byte

	// Reconstruct construct merkletree again,
	// if vals not empty, using vals as new data source, otherwise using old data
	Reconstruct(vals ...interface{}) error

	// VerifyRoot verify tree root hash whether matches calculate root hash,
	// if match, return true, otherwise return false.
	VerifyRoot() bool

	// VerifyContent indicates whether a given content is in the tree and the hashes are valid for that content.
	// Returns true if the expected Root equals to the calculated root return true, otherwise false.
	VerifyValue(exceptedRoot []byte, val interface{}) bool
}
