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

package notification

import (
	"context"
	"errors"
	"strings"

	"github.com/asgardeo/thunder/internal/notification/common"
)

var (
	ErrNotificationSenderIsImmutable = errors.New("notification sender is immutable")
)

type compositeNotificationStore struct {
	fileStore notificationStoreInterface
	dbStore   notificationStoreInterface
}

func newCompositeNotificationStore(fileStore, dbStore notificationStoreInterface) notificationStoreInterface {
	return &compositeNotificationStore{
		fileStore: fileStore,
		dbStore:   dbStore,
	}
}

func (c *compositeNotificationStore) createSender(ctx context.Context, sender common.NotificationSenderDTO) error {
	return c.dbStore.createSender(ctx, sender)
}

func (c *compositeNotificationStore) listSenders(ctx context.Context) ([]common.NotificationSenderDTO, error) {
	dbSenders, err := c.dbStore.listSenders(ctx)
	if err != nil {
		return nil, err
	}

	fileSenders, err := c.fileStore.listSenders(ctx)
	if err != nil {
		return nil, err
	}

	seen := make(map[string]struct{})
	result := make([]common.NotificationSenderDTO, 0, len(dbSenders)+len(fileSenders))

	for _, sender := range dbSenders {
		if _, ok := seen[sender.ID]; ok {
			continue
		}
		seen[sender.ID] = struct{}{}
		result = append(result, sender)
	}

	for _, sender := range fileSenders {
		if _, ok := seen[sender.ID]; ok {
			continue
		}
		seen[sender.ID] = struct{}{}
		result = append(result, sender)
	}

	return result, nil
}

func (c *compositeNotificationStore) getSenderByID(ctx context.Context, id string) (*common.NotificationSenderDTO, error) {
	sender, err := c.dbStore.getSenderByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if sender != nil {
		return sender, nil
	}

	sender, err = c.fileStore.getSenderByID(ctx, id)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return nil, nil
		}
		return nil, err
	}

	return sender, nil
}

func (c *compositeNotificationStore) getSenderByName(
	ctx context.Context, name string,
) (*common.NotificationSenderDTO, error) {
	sender, err := c.dbStore.getSenderByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if sender != nil {
		return sender, nil
	}

	return c.fileStore.getSenderByName(ctx, name)
}

func (c *compositeNotificationStore) updateSender(
	ctx context.Context, id string, sender common.NotificationSenderDTO,
) error {
	immutable, err := c.isSenderDeclarative(ctx, id)
	if err != nil {
		return err
	}
	if immutable {
		return ErrNotificationSenderIsImmutable
	}

	return c.dbStore.updateSender(ctx, id, sender)
}

func (c *compositeNotificationStore) deleteSender(ctx context.Context, id string) error {
	immutable, err := c.isSenderDeclarative(ctx, id)
	if err != nil {
		return err
	}
	if immutable {
		return ErrNotificationSenderIsImmutable
	}

	return c.dbStore.deleteSender(ctx, id)
}

func (c *compositeNotificationStore) isSenderDeclarative(ctx context.Context, id string) (bool, error) {
	sender, err := c.fileStore.getSenderByID(ctx, id)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return false, nil
		}
		return false, err
	}

	return sender != nil, nil
}
