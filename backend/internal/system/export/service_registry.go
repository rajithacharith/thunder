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

package export

import (
	"github.com/asgardeo/thunder/internal/system/log"
)

// exportResourcesByType is a generic method that uses the registry to export any resource type.
// This eliminates the need for separate export methods for each resource type.
func (es *exportService) exportResourcesByType(
	resourceType string,
	resourceIDs []string,
	options *ExportOptions,
) ([]ExportFile, []ExportError) {
	logger := log.GetLogger().With(log.String("component", "ExportService"))
	exportFiles := make([]ExportFile, 0, len(resourceIDs))
	exportErrors := make([]ExportError, 0, len(resourceIDs))

	// Get the exporter from registry
	exporter, exists := es.registry.Get(resourceType)
	if !exists {
		logger.Error("No exporter registered for resource type", log.String("resourceType", resourceType))
		exportErrors = append(exportErrors, ExportError{
			ResourceType: resourceType,
			Error:        "Resource type not supported for export",
			Code:         "UNSUPPORTED_RESOURCE_TYPE",
		})
		return exportFiles, exportErrors
	}

	// Determine resource ID list (support wildcard)
	var resourceIDList []string
	if len(resourceIDs) == 1 && resourceIDs[0] == "*" {
		ids, err := exporter.GetAllResourceIDs()
		if err != nil {
			logger.Warn("Failed to get all resources",
				log.String("resourceType", resourceType),
				log.Any("error", err))
			exportErrors = append(exportErrors, ExportError{
				ResourceType: resourceType,
				Error:        err.Error,
				Code:         err.Code,
			})
			return exportFiles, exportErrors
		}
		resourceIDList = ids
	} else {
		resourceIDList = resourceIDs
	}

	// Export each resource
	for _, resourceID := range resourceIDList {
		// Get the resource
		resource, _, svcErr := exporter.GetResourceByID(resourceID)
		if svcErr != nil {
			logger.Warn("Failed to get resource for export",
				log.String("resourceType", resourceType),
				log.String("resourceID", resourceID),
				log.String("error", svcErr.Error))
			exportErrors = append(exportErrors, ExportError{
				ResourceType: resourceType,
				ResourceID:   resourceID,
				Error:        svcErr.Error,
				Code:         svcErr.Code,
			})
			continue
		}

		// Validate resource
		validatedName, exportErr := exporter.ValidateResource(resource, resourceID, logger)
		if exportErr != nil {
			exportErrors = append(exportErrors, *exportErr)
			continue
		}

		// Generate export content
		var content string
		if options.Format == formatJSON {
			logger.Warn("JSON format not yet implemented, falling back to YAML")
			options.Format = formatYAML
		}

		templateContent, err := es.generateTemplateFromStruct(
			resource,
			exporter.GetParameterizerType(),
			validatedName,
		)
		if err != nil {
			logger.Warn("Failed to generate template from struct",
				log.String("resourceType", resourceType),
				log.String("resourceID", resourceID),
				log.String("error", err.Error()))
			exportErrors = append(exportErrors, ExportError{
				ResourceType: resourceType,
				ResourceID:   resourceID,
				Error:        err.Error(),
				Code:         "TemplateGenerationError",
			})
			continue
		}
		content = templateContent

		// Generate file metadata
		fileName := es.generateFileName(validatedName, resourceType, resourceID, options)
		folderPath := es.generateFolderPath(resourceType, options)

		// Create export file
		exportFile := ExportFile{
			FileName:     fileName,
			Content:      content,
			FolderPath:   folderPath,
			ResourceType: resourceType,
			ResourceID:   resourceID,
		}
		exportFiles = append(exportFiles, exportFile)
	}

	return exportFiles, exportErrors
}

// Helper method to export resources with automatic registry lookup
func (es *exportService) exportResourcesFromRequest(
	resourceType string,
	resourceIDs []string,
	options *ExportOptions,
	exportFiles *[]ExportFile,
	exportErrors *[]ExportError,
	resourceCounts map[string]int,
) {
	if len(resourceIDs) > 0 {
		files, errors := es.exportResourcesByType(resourceType, resourceIDs, options)
		*exportFiles = append(*exportFiles, files...)
		*exportErrors = append(*exportErrors, errors...)
		resourceCounts[resourceType+"s"] = len(files)
	}
}
