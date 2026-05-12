/*
 * Copyright (c) 2026, WSO2 LLC. (https://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package declarativeresource

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

// JSONRawField decodes a JSON blob from YAML. It handles both a literal block
// scalar (string) produced by the exporter and an inline YAML mapping/sequence,
// normalising both to raw JSON bytes at decode time.
type JSONRawField []byte

// UnmarshalYAML implements yaml.Unmarshaler.
func (f *JSONRawField) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.ScalarNode {
		*f = JSONRawField(value.Value)
		return nil
	}
	var v interface{}
	if err := value.Decode(&v); err != nil {
		return err
	}
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	*f = JSONRawField(b)
	return nil
}
