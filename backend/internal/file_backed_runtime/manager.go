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

// Package filebackedruntime provides functionalities to manage file-backed runtime configurations.
package filebackedruntime

import (
	"os"
	"path/filepath"

	"github.com/asgardeo/thunder/internal/system/log"
)

var instance *ConfigManager

// ConfigManager manages the file-backed runtime configurations.
type ConfigManager struct {
	directoryPath string
	config        FileBackedRuntimeConf
}

// NewFileConfigManager creates a new instance of ConfigManager and loads the configurations
// from the specified directory.
func NewFileConfigManager(directoryPath string) *ConfigManager {
	instance = &ConfigManager{
		directoryPath: directoryPath,
		config:        FileBackedRuntimeConf{},
	}
	if err := instance.loadConfig(); err != nil {
		log.GetLogger().Fatal("Failed to load file-backed runtime configuration", log.Error(err))
	}
	return instance
}

func (m *ConfigManager) loadConfig() error {
	logger := log.GetLogger().With(log.String(log.LoggerKeyComponentName, "FileConfigManager"))
	notificationSendersPath := m.directoryPath + "/notification_senders/"
	files, err := os.ReadDir(notificationSendersPath)
	if err != nil {
		logger.Error("Failed to read notification senders directory",
			log.String("path", notificationSendersPath), log.Error(err))
		return err
	}
	var notificationSenders [][]byte
	for _, file := range files {
		if !file.IsDir() {
			// #nosec G304 -- File path is controlled and within a trusted directory
			fileContent, err := os.ReadFile(filepath.Join(notificationSendersPath, file.Name()))
			if err != nil {
				logger.Error("Failed to read file notification sender configuration",
					log.String("filePath", file.Name()), log.Error(err))
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
			// #nosec G304 -- File path is controlled and within a trusted directory
			fileContent, err := os.ReadFile(filepath.Join(idpsPath, file.Name()))
			if err != nil {
				logger.Error("Failed to read file IDP configuration", log.String("filePath", file.Name()), log.Error(err))
				continue
			}
			idps = append(idps, fileContent)
		}
	}
	m.config.IDPs = idps
	return nil
}

// GetConfig returns the loaded file-backed runtime configuration.
func GetConfig() FileBackedRuntimeConf {
	if instance == nil {
		log.GetLogger().Warn("FileConfigManager is not initialized. Returning empty configuration.")
		return FileBackedRuntimeConf{}
	}
	return instance.config
}
