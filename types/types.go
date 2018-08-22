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
package types

import (
	"bytes"
	"encoding/gob"
)

// Marshaler is the interface representing objects that can marshal into bytes.
type Marshaler interface {
	Marshal() ([]byte, error)
}

// Unmarshaler is the interface representing objects that can unmarshal from bytes.
type Unmarshaler interface {
	Unmarshal([]byte) error
}

func Marshal(v interface{}) ([]byte, error) {
	if m, ok := v.(Marshaler); ok {
		return m.Marshal()
	}

	// Using gob as default data serialization format
	// Create an gob encoder
	var vBuffer bytes.Buffer
	enc := gob.NewEncoder(&vBuffer)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}

	return vBuffer.Bytes(), nil
}

func Unmarshal(data []byte, v interface{}) error {
	if m, ok := v.(Unmarshaler); ok {
		return m.Unmarshal(data)
	}

	// Using gob as default data serialization format
	// Create an gob decoder
	dec := gob.NewDecoder(bytes.NewBuffer(data))
	if err := dec.Decode(v); err != nil {
		return err
	}

	return nil
}
