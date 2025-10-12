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

package notification

import (
	"errors"

	filebackedruntime "github.com/asgardeo/thunder/internal/file_backed_runtime"
	"github.com/asgardeo/thunder/internal/notification/common"
	"github.com/asgardeo/thunder/internal/system/log"

	"gopkg.in/yaml.v3"
)

type inMemoryStore struct {
	store map[string]common.NotificationSenderDTO
}

var logger = log.GetLogger().With(log.String(log.LoggerKeyComponentName, "InMemoryNotificationStore"))

// createSender implements notificationStoreInterface.
func (i *inMemoryStore) createSender(sender common.NotificationSenderDTO) error {
	return errors.New("create operation is not supported in in-memory store")
}

// deleteSender implements notificationStoreInterface.
func (i *inMemoryStore) deleteSender(id string) error {
	return errors.New("delete operation is not supported in in-memory store")
}

// getSenderByID implements notificationStoreInterface.
func (i *inMemoryStore) getSenderByID(id string) (*common.NotificationSenderDTO, error) {
	sender, exists := i.store[id]
	if !exists {
		return nil, nil
	}
	return &sender, nil
}

// getSenderByName implements notificationStoreInterface.
func (i *inMemoryStore) getSenderByName(name string) (*common.NotificationSenderDTO, error) {
	for _, sender := range i.store {
		if sender.Name == name {
			return &sender, nil
		}
	}
	return nil, nil
}

// listSenders implements notificationStoreInterface.
func (i *inMemoryStore) listSenders() ([]common.NotificationSenderDTO, error) {
	senders := make([]common.NotificationSenderDTO, 0, len(i.store))
	for _, sender := range i.store {
		senders = append(senders, sender)
	}
	return senders, nil
}

// updateSender implements notificationStoreInterface.
func (i *inMemoryStore) updateSender(id string, sender common.NotificationSenderDTO) error {
	return errors.New("update operation is not supported in in-memory store")
}

var _ notificationStoreInterface = (*inMemoryStore)(nil)

func newInMemoryStore() *inMemoryStore {
	fileConfigs := filebackedruntime.GetConfig().NotificationSenders
	store := make(map[string]common.NotificationSenderDTO)
	for idx, fileConfig := range fileConfigs {
		sender, err := convertFileSenderConfigToDTO(fileConfig)
		if err != nil {
			logger.Warn("Skipping invalid notification sender configuration", log.Int("index", idx), log.Error(err))
			continue
		}
		store[sender.ID] = sender
	}
	return &inMemoryStore{
		store: store,
	}
}

func convertFileSenderConfigToDTO(fileConfig []byte) (common.NotificationSenderDTO, error) {
	var sender common.NotificationSenderDTO
	err := yaml.Unmarshal(fileConfig, &sender)
	if err != nil {
		return common.NotificationSenderDTO{},
			errors.New("failed to unmarshal file notification sender configuration: " + err.Error())
	}
	return sender, nil
}
