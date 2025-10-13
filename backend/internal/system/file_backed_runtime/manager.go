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
	"path"
	"path/filepath"
	"sync"

	"github.com/asgardeo/thunder/internal/system/config"
	"github.com/asgardeo/thunder/internal/system/log"
)

// GetConfigs reads all configuration files from the specified directory within the immutable_configs directory
func GetConfigs(configDirectoryPath string) ([][]byte, error) {
	thunderHome := config.GetThunderRuntime().ThunderHome
	immutableConfigFilePath := path.Join(thunderHome, "repository/conf/immutable_configs/")
	absoluteDirectoryPath := immutableConfigFilePath + "/" + configDirectoryPath + "/"
	files, err := os.ReadDir(absoluteDirectoryPath)
	if err != nil {
		log.GetLogger().Error("Failed to read configuration directory",
			log.String("path", absoluteDirectoryPath), log.Error(err))
		return nil, err
	}

	var configs [][]byte

	// Count non-directory files
	var fileCount int
	for _, file := range files {
		if !file.IsDir() {
			fileCount++
		}
	}

	if fileCount == 0 {
		return configs, nil
	}

	// Use channels to collect results from goroutines
	configChan := make(chan []byte, fileCount)
	var wg sync.WaitGroup

	for _, file := range files {
		if !file.IsDir() {
			wg.Add(1)
			go func(fileName string) {
				defer wg.Done()
				filePath := filepath.Join(absoluteDirectoryPath, fileName)
				filePath = filepath.Clean(filePath)
				// #nosec G304 -- File path is controlled and within a trusted directory
				fileContent, err := os.ReadFile(filePath)
				if err != nil {
					log.GetLogger().Warn("Failed to read configuration file", log.String("filePath", fileName), log.Error(err))
					return
				}
				configChan <- fileContent
			}(file.Name())
		}
	}

	// Wait for all goroutines to complete and close the channel
	go func() {
		wg.Wait()
		close(configChan)
	}()

	// Collect results from the channel
	for content := range configChan {
		configs = append(configs, content)
	}
	return configs, nil
}
