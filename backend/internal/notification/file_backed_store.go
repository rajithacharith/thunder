package notification

import (
	"errors"

	filebackedruntime "github.com/asgardeo/thunder/internal/file_backed_runtime"
	"github.com/asgardeo/thunder/internal/notification/common"
	"github.com/asgardeo/thunder/internal/system/log"
	"github.com/stretchr/testify/assert/yaml"
)

type inMemoryStore struct {
	store map[string]common.NotificationSenderDTO
}

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
	var senders []common.NotificationSenderDTO
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
	for _, fileConfig := range fileConfigs {
		sender, err := convertFileSenderConfigToDTO(fileConfig)
		if err != nil {
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
		logger := log.GetLogger().With(log.String(log.LoggerKeyComponentName, "FileBasedNotificationStore"))
		logger.Error("Failed to unmarshal file notification sender configuration", log.Error(err))
		return common.NotificationSenderDTO{}, errors.New("failed to unmarshal file notification sender configuration: " + err.Error())
	}
	return sender, nil
}
