/*
 * Copyright (c) 2025, WSO2 LLC. (https://www.wso2.com).
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

package idp

import (
	"errors"

	filebackedruntime "github.com/asgardeo/thunder/internal/system/file_backed_runtime"
	"github.com/asgardeo/thunder/internal/system/log"

	"gopkg.in/yaml.v3"
)

type fileBasedStore struct {
	store map[string]IDPDTO
}

var logger = log.GetLogger().With(log.String(log.LoggerKeyComponentName, "FileBasedIDPStore"))

// CreateIdentityProvider implements idpStoreInterface.
func (f *fileBasedStore) CreateIdentityProvider(idp IDPDTO) error {
	return errors.New("create operation is not supported in file-based store")
}

// DeleteIdentityProvider implements idpStoreInterface.
func (f *fileBasedStore) DeleteIdentityProvider(idpID string) error {
	return errors.New("delete operation is not supported in file-based store")
}

// GetIdentityProvider implements idpStoreInterface.
func (f *fileBasedStore) GetIdentityProvider(idpID string) (*IDPDTO, error) {
	idp, exists := f.store[idpID]
	if !exists {
		return nil, nil
	}
	return &idp, nil
}

// GetIdentityProviderByName implements idpStoreInterface.
func (f *fileBasedStore) GetIdentityProviderByName(idpName string) (*IDPDTO, error) {
	for _, idp := range f.store {
		if idp.Name == idpName {
			return &idp, nil
		}
	}
	return nil, nil
}

// GetIdentityProviderList implements idpStoreInterface.
func (f *fileBasedStore) GetIdentityProviderList() ([]BasicIDPDTO, error) {
	idpList := make([]BasicIDPDTO, 0, len(f.store))
	for _, idp := range f.store {
		idpList = append(idpList, BasicIDPDTO{
			ID:   idp.ID,
			Name: idp.Name,
		})
	}
	return idpList, nil
}

// UpdateIdentityProvider implements idpStoreInterface.
func (f *fileBasedStore) UpdateIdentityProvider(idp *IDPDTO) error {
	return errors.New("update operation is not supported in file-based store")
}

var _ idpStoreInterface = (*fileBasedStore)(nil)

func newFileBasedStore() *fileBasedStore {
	fileConfigs, err := filebackedruntime.GetConfigs("identity_providers")
	if err != nil {
		logger.Fatal("Failed to load identity provider configurations from files", log.Error(err))
	}
	store := make(map[string]IDPDTO)
	for idx, fileConfig := range fileConfigs {
		idp, err := convertFileIDPConfigToDTO(fileConfig)
		if err != nil {
			logger.Warn("Skipping invalid IDP configuration", log.Int("index", idx), log.Error(err))
			continue
		}
		store[idp.ID] = idp
	}
	return &fileBasedStore{
		store: store,
	}
}

func convertFileIDPConfigToDTO(fileContent []byte) (IDPDTO, error) {
	var idp IDPDTO
	err := yaml.Unmarshal(fileContent, &idp)
	if err != nil {
		return IDPDTO{}, errors.New("failed to unmarshal file IDP configuration: " + err.Error())
	}
	if idp.ID == "" || idp.Name == "" {
		return IDPDTO{}, errors.New("invalid file IDP configuration, missing required fields")
	}
	return idp, nil
}
