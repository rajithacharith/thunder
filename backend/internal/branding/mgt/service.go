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

// Package brandingmgt provides handlers and services for managing branding configurations.
package brandingmgt

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/asgardeo/thunder/internal/branding/common"
	serverconst "github.com/asgardeo/thunder/internal/system/constants"
	"github.com/asgardeo/thunder/internal/system/error/serviceerror"
	"github.com/asgardeo/thunder/internal/system/log"
	"github.com/asgardeo/thunder/internal/system/utils"
)

const loggerComponentName = "BrandingMgtService"

// BrandingMgtServiceInterface defines the interface for the branding management service.
type BrandingMgtServiceInterface interface {
	GetBrandingList(limit, offset int) (*BrandingList, *serviceerror.ServiceError)
	CreateBranding(branding CreateBrandingRequest) (*common.Branding, *serviceerror.ServiceError)
	GetBranding(id string) (*common.Branding, *serviceerror.ServiceError)
	UpdateBranding(id string, branding UpdateBrandingRequest) (*common.Branding, *serviceerror.ServiceError)
	DeleteBranding(id string) *serviceerror.ServiceError
	IsBrandingExist(id string) (bool, *serviceerror.ServiceError)
}

// brandingMgtService is the default implementation of the BrandingMgtServiceInterface.
type brandingMgtService struct {
	brandingMgtStore brandingMgtStoreInterface
	logger           *log.Logger
}

// newBrandingMgtService creates a new instance of BrandingMgtService with injected dependencies.
func newBrandingMgtService(brandingMgtStore brandingMgtStoreInterface) BrandingMgtServiceInterface {
	logger := log.GetLogger().With(log.String(log.LoggerKeyComponentName, loggerComponentName))
	return &brandingMgtService{
		brandingMgtStore: brandingMgtStore,
		logger:           logger,
	}
}

// GetBrandingList retrieves a list of branding configurations.
func (bs *brandingMgtService) GetBrandingList(limit, offset int) (*BrandingList, *serviceerror.ServiceError) {
	if err := validatePaginationParams(limit, offset); err != nil {
		return nil, err
	}

	totalCount, err := bs.brandingMgtStore.GetBrandingListCount()
	if err != nil {
		bs.logger.Error("Failed to get branding count", log.Error(err))
		return nil, &serviceerror.InternalServerError
	}

	brandings, err := bs.brandingMgtStore.GetBrandingList(limit, offset)
	if err != nil {
		bs.logger.Error("Failed to list brandings", log.Error(err))
		return nil, &serviceerror.InternalServerError
	}

	response := &BrandingList{
		TotalResults: totalCount,
		Brandings:    brandings,
		StartIndex:   offset + 1,
		Count:        len(brandings),
		Links:        buildPaginationLinks(limit, offset, totalCount),
	}

	return response, nil
}

// CreateBranding creates a new branding configuration.
func (bs *brandingMgtService) CreateBranding(
	branding CreateBrandingRequest,
) (*common.Branding, *serviceerror.ServiceError) {
	bs.logger.Debug("Creating branding configuration")

	if branding.DisplayName == "" {
		return nil, &common.ErrorMissingDisplayName
	}

	if err := bs.validateBrandingPreferences(branding.Preferences); err != nil {
		return nil, err
	}

	id, err := utils.GenerateUUIDv7()
	if err != nil {
		bs.logger.Error("Failed to generate UUID", log.Error(err))
		return nil, &serviceerror.InternalServerError
	}

	if err := bs.brandingMgtStore.CreateBranding(id, branding); err != nil {
		bs.logger.Error("Failed to create branding", log.Error(err))
		return nil, &serviceerror.InternalServerError
	}

	createdBranding := &common.Branding{
		ID:          id,
		DisplayName: branding.DisplayName,
		Preferences: branding.Preferences,
	}

	bs.logger.Debug("Successfully created branding", log.String("id", id))
	return createdBranding, nil
}

// GetBranding retrieves a specific branding configuration by its id.
func (bs *brandingMgtService) GetBranding(id string) (*common.Branding, *serviceerror.ServiceError) {
	bs.logger.Debug("Retrieving branding", log.String("id", id))

	if id == "" {
		return nil, &common.ErrorMissingBrandingID
	}

	branding, err := bs.brandingMgtStore.GetBranding(id)
	if err != nil {
		if errors.Is(err, common.ErrBrandingNotFound) {
			bs.logger.Debug("Branding not found", log.String("id", id))
			return nil, &common.ErrorBrandingNotFound
		}
		bs.logger.Error("Failed to retrieve branding", log.String("id", id), log.Error(err))
		return nil, &serviceerror.InternalServerError
	}

	bs.logger.Debug("Successfully retrieved branding", log.String("id", branding.ID))
	return &branding, nil
}

// UpdateBranding updates an existing branding configuration.
func (bs *brandingMgtService) UpdateBranding(
	id string, branding UpdateBrandingRequest) (*common.Branding, *serviceerror.ServiceError) {
	bs.logger.Debug("Updating branding", log.String("id", id))

	if id == "" {
		return nil, &common.ErrorMissingBrandingID
	}

	if branding.DisplayName == "" {
		return nil, &common.ErrorMissingDisplayName
	}

	if err := bs.validateBrandingPreferences(branding.Preferences); err != nil {
		return nil, err
	}

	if err := bs.brandingMgtStore.UpdateBranding(id, branding); err != nil {
		if errors.Is(err, common.ErrBrandingNotFound) {
			bs.logger.Debug("Branding not found for update", log.String("id", id))
			return nil, &common.ErrorBrandingNotFound
		}
		bs.logger.Error("Failed to update branding", log.String("id", id), log.Error(err))
		return nil, &serviceerror.InternalServerError
	}

	updatedBranding := &common.Branding{
		ID:          id,
		DisplayName: branding.DisplayName,
		Preferences: branding.Preferences,
	}

	bs.logger.Debug("Successfully updated branding", log.String("id", id))
	return updatedBranding, nil
}

// DeleteBranding deletes the specified branding configuration by its id.
func (bs *brandingMgtService) DeleteBranding(id string) *serviceerror.ServiceError {
	bs.logger.Debug("Deleting branding", log.String("id", id))

	if id == "" {
		return &common.ErrorMissingBrandingID
	}

	exists, err := bs.brandingMgtStore.IsBrandingExist(id)
	if err != nil {
		bs.logger.Error("Failed to check branding existence", log.String("id", id), log.Error(err))
		return &serviceerror.InternalServerError
	}
	if !exists {
		bs.logger.Debug("Branding not found", log.String("id", id))
		return nil
	}

	// Check if branding has any applications before deleting
	appCount, err := bs.brandingMgtStore.GetApplicationsCountByBrandingID(id)
	if err != nil {
		bs.logger.Error("Failed to get applications count", log.String("id", id), log.Error(err))
		return &serviceerror.InternalServerError
	}
	if appCount > 0 {
		bs.logger.Debug("Cannot delete branding with active applications",
			log.String("id", id), log.Int("appCount", appCount))
		return &common.ErrorCannotDeleteBranding
	}

	if err := bs.brandingMgtStore.DeleteBranding(id); err != nil {
		bs.logger.Error("Failed to delete branding", log.String("id", id), log.Error(err))
		return &serviceerror.InternalServerError
	}

	bs.logger.Debug("Successfully deleted branding", log.String("id", id))
	return nil
}

// validateBrandingPreferences validates the branding preferences JSON.
func (bs *brandingMgtService) validateBrandingPreferences(preferences json.RawMessage) *serviceerror.ServiceError {
	if len(preferences) == 0 {
		return &common.ErrorMissingPreferences
	}

	// Validate that preferences is valid JSON
	var test interface{}
	if err := json.Unmarshal(preferences, &test); err != nil {
		return &common.ErrorInvalidPreferences
	}

	// Ensure it's an object, not an array or primitive
	if _, ok := test.(map[string]interface{}); !ok {
		return &common.ErrorInvalidPreferences
	}

	// TODO: Add additional validations against the JSON Schema
	return nil
}

// validatePaginationParams validates pagination parameters.
func validatePaginationParams(limit, offset int) *serviceerror.ServiceError {
	if limit < 1 || limit > serverconst.MaxPageSize {
		return &common.ErrorInvalidLimit
	}
	if offset < 0 {
		return &common.ErrorInvalidOffset
	}
	return nil
}

// buildPaginationLinks builds pagination links for the response.
func buildPaginationLinks(limit, offset, totalCount int) []Link {
	const brandingBasePath = "/branding"
	links := make([]Link, 0)

	if offset > 0 {
		links = append(links, Link{
			Href: fmt.Sprintf("%s?offset=0&limit=%d", brandingBasePath, limit),
			Rel:  "first",
		})

		prevOffset := offset - limit
		if prevOffset < 0 {
			prevOffset = 0
		}
		links = append(links, Link{
			Href: fmt.Sprintf("%s?offset=%d&limit=%d", brandingBasePath, prevOffset, limit),
			Rel:  "prev",
		})
	}

	if offset+limit < totalCount {
		nextOffset := offset + limit
		links = append(links, Link{
			Href: fmt.Sprintf("%s?offset=%d&limit=%d", brandingBasePath, nextOffset, limit),
			Rel:  "next",
		})
	}

	lastPageOffset := ((totalCount - 1) / limit) * limit
	if offset < lastPageOffset {
		links = append(links, Link{
			Href: fmt.Sprintf("%s?offset=%d&limit=%d", brandingBasePath, lastPageOffset, limit),
			Rel:  "last",
		})
	}

	return links
}

// IsBrandingExist checks if a branding configuration exists by its ID.
func (bs *brandingMgtService) IsBrandingExist(id string) (bool, *serviceerror.ServiceError) {
	if id == "" {
		return false, nil
	}

	exists, err := bs.brandingMgtStore.IsBrandingExist(id)
	if err != nil {
		bs.logger.Error("Failed to check branding existence", log.String("id", id), log.Error(err))
		return false, &serviceerror.InternalServerError
	}

	return exists, nil
}
