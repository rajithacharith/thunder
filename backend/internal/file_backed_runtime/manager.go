package filebackedruntime

import (
	"os"
	"path/filepath"

	"github.com/asgardeo/thunder/internal/system/config"
	"github.com/asgardeo/thunder/internal/system/log"
)

var instance *ConfigManager

type ConfigManager struct {
	directoryPath string
	config        FileBackedRuntimeConf
}

func NewFileConfigManager(directoryPath string) *ConfigManager {
	instance = &ConfigManager{
		directoryPath: directoryPath,
		config:        FileBackedRuntimeConf{},
	}
	return instance
}

func (m *ConfigManager) LoadConfig() error {
	logger := log.GetLogger().With(log.String(log.LoggerKeyComponentName, "FileConfigManager"))
	notificationSendersPath := m.directoryPath + "/notification_senders/"
	files, err := os.ReadDir(notificationSendersPath)
	if err != nil {
		logger.Error("Failed to read notification senders directory", log.String("path", notificationSendersPath), log.Error(err))
		return err
	}
	var notificationSenders [][]byte
	for _, file := range files {
		if !file.IsDir() {
			fileContent, err := os.ReadFile(filepath.Join(notificationSendersPath, file.Name()))
			if err != nil {
				logger.Error("Failed to read file notification sender configuration", log.String("filePath", file.Name()), log.Error(err))
				continue
			}
			notificationSenders = append(notificationSenders, fileContent)
		}
	}
	m.config.NotificationSenders = notificationSenders

	idpsPath := m.directoryPath + "/identity_providers/"
	files, err = os.ReadDir(idpsPath)
	if err != nil {
		logger.Error("Failed to read identity providers directory", log.String("path", idpsPath), log.Error(err))
		return err
	}
	var idps [][]byte
	for _, file := range files {
		if !file.IsDir() {
			fileContent, err := os.ReadFile(filepath.Join(idpsPath, file.Name()))
			if err != nil {
				logger.Error("Failed to read file IDP configuration", log.String("filePath", file.Name()), log.Error(err))
				continue
			}
			idps = append(idps, fileContent)
		}
	}
	m.config.IDPs = idps

	logger.Info("Successfully loaded immutable gateway configurations", log.Int("notificationSendersCount", len(notificationSenders)), log.Int("idpsCount", len(idps)))
	return nil
}

func GetConfig() FileBackedRuntimeConf {
	if !config.GetThunderRuntime().Config.ImmutableGateway.Enabled {
		return FileBackedRuntimeConf{}
	}
	return instance.config
}
