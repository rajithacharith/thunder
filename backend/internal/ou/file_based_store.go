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

package ou

import (
	"errors"

	"github.com/asgardeo/thunder/internal/system/file_based_runtime/entity"
	"github.com/asgardeo/thunder/internal/system/immutableresource"
)

var (
	// ErrOUNotFound is the error returned when an organization unit is not found in the file-based store
	ErrOUNotFound = errors.New("organization unit not found")
)

type ouFileBasedStore struct {
	*immutableresource.GenericFileBasedStore
}

// Create implements immutableresource.Storer interface for resource loader
func (f *ouFileBasedStore) Create(id string, data interface{}) error {
	ou := data.(*OrganizationUnit)
	return f.CreateOrganizationUnit(*ou)
}

// CreateOrganizationUnit implements organizationUnitStoreInterface.
func (f *ouFileBasedStore) CreateOrganizationUnit(ou OrganizationUnit) error {
	return f.GenericFileBasedStore.Create(ou.ID, &ou)
}

// DeleteOrganizationUnit implements organizationUnitStoreInterface.
func (f *ouFileBasedStore) DeleteOrganizationUnit(id string) error {
	return errors.New("DeleteOrganizationUnit is not supported in file-based store")
}

// GetOrganizationUnit implements organizationUnitStoreInterface.
func (f *ouFileBasedStore) GetOrganizationUnit(id string) (OrganizationUnit, error) {
	data, err := f.GenericFileBasedStore.Get(id)
	if err != nil {
		return OrganizationUnit{}, ErrOUNotFound
	}
	ou, ok := data.(*OrganizationUnit)
	if !ok {
		immutableresource.LogTypeAssertionError("organization unit", id)
		return OrganizationUnit{}, errors.New("organization unit data corrupted")
	}
	return *ou, nil
}

// GetOrganizationUnitByPath implements organizationUnitStoreInterface.
func (f *ouFileBasedStore) GetOrganizationUnitByPath(handles []string) (OrganizationUnit, error) {
	// For file-based store, we need to traverse the path to find the OU
	// This implementation assumes a hierarchical structure
	list, err := f.GenericFileBasedStore.List()
	if err != nil {
		return OrganizationUnit{}, err
	}

	// Build a path from handles
	// Start from root (parent = nil) and traverse
	var currentParent *string = nil
	var currentOU *OrganizationUnit = nil

	for i, handle := range handles {
		found := false
		for _, item := range list {
			ou, ok := item.Data.(*OrganizationUnit)
			if !ok {
				continue
			}

			// Check if this OU matches current handle and parent
			if ou.Handle == handle {
				if (currentParent == nil && ou.Parent == nil) ||
					(currentParent != nil && ou.Parent != nil && *currentParent == *ou.Parent) {
					found = true
					currentOU = ou
					currentParent = &ou.ID
					break
				}
			}
		}

		if !found {
			return OrganizationUnit{}, ErrOUNotFound
		}

		// If we've reached the last handle, return this OU
		if i == len(handles)-1 {
			return *currentOU, nil
		}
	}

	return OrganizationUnit{}, ErrOUNotFound
}

// GetOrganizationUnitList implements organizationUnitStoreInterface.
func (f *ouFileBasedStore) GetOrganizationUnitList(limit, offset int) ([]OrganizationUnitBasic, error) {
	list, err := f.GenericFileBasedStore.List()
	if err != nil {
		return nil, err
	}

	var ouList []OrganizationUnitBasic
	for _, item := range list {
		if ou, ok := item.Data.(*OrganizationUnit); ok {
			// Only include root OUs (parent = nil)
			if ou.Parent == nil {
				basicOU := OrganizationUnitBasic{
					ID:          ou.ID,
					Handle:      ou.Handle,
					Name:        ou.Name,
					Description: ou.Description,
				}
				ouList = append(ouList, basicOU)
			}
		}
	}

	// Apply pagination
	start := offset
	end := offset + limit
	if start >= len(ouList) {
		return []OrganizationUnitBasic{}, nil
	}
	if end > len(ouList) {
		end = len(ouList)
	}

	return ouList[start:end], nil
}

// GetOrganizationUnitListCount implements organizationUnitStoreInterface.
func (f *ouFileBasedStore) GetOrganizationUnitListCount() (int, error) {
	list, err := f.GenericFileBasedStore.List()
	if err != nil {
		return 0, err
	}

	count := 0
	for _, item := range list {
		if ou, ok := item.Data.(*OrganizationUnit); ok {
			// Only count root OUs (parent = nil)
			if ou.Parent == nil {
				count++
			}
		}
	}

	return count, nil
}

// IsOrganizationUnitExists implements organizationUnitStoreInterface.
func (f *ouFileBasedStore) IsOrganizationUnitExists(id string) (bool, error) {
	_, err := f.GetOrganizationUnit(id)
	if err != nil {
		if errors.Is(err, ErrOUNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// CheckOrganizationUnitNameConflict implements organizationUnitStoreInterface.
func (f *ouFileBasedStore) CheckOrganizationUnitNameConflict(name string, parent *string) (bool, error) {
	list, err := f.GenericFileBasedStore.List()
	if err != nil {
		return false, err
	}

	for _, item := range list {
		if ou, ok := item.Data.(*OrganizationUnit); ok {
			if ou.Name == name {
				// Check if parent matches
				if (parent == nil && ou.Parent == nil) ||
					(parent != nil && ou.Parent != nil && *parent == *ou.Parent) {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

// CheckOrganizationUnitHandleConflict implements organizationUnitStoreInterface.
func (f *ouFileBasedStore) CheckOrganizationUnitHandleConflict(handle string, parent *string) (bool, error) {
	list, err := f.GenericFileBasedStore.List()
	if err != nil {
		return false, err
	}

	for _, item := range list {
		if ou, ok := item.Data.(*OrganizationUnit); ok {
			if ou.Handle == handle {
				// Check if parent matches
				if (parent == nil && ou.Parent == nil) ||
					(parent != nil && ou.Parent != nil && *parent == *ou.Parent) {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

// UpdateOrganizationUnit implements organizationUnitStoreInterface.
func (f *ouFileBasedStore) UpdateOrganizationUnit(ou OrganizationUnit) error {
	return errors.New("UpdateOrganizationUnit is not supported in file-based store")
}

// CheckOrganizationUnitHasChildResources implements organizationUnitStoreInterface.
func (f *ouFileBasedStore) CheckOrganizationUnitHasChildResources(id string) (bool, error) {
	list, err := f.GenericFileBasedStore.List()
	if err != nil {
		return false, err
	}

	for _, item := range list {
		if ou, ok := item.Data.(*OrganizationUnit); ok {
			if ou.Parent != nil && *ou.Parent == id {
				return true, nil
			}
		}
	}

	return false, nil
}

// GetOrganizationUnitChildrenCount implements organizationUnitStoreInterface.
func (f *ouFileBasedStore) GetOrganizationUnitChildrenCount(id string) (int, error) {
	list, err := f.GenericFileBasedStore.List()
	if err != nil {
		return 0, err
	}

	count := 0
	for _, item := range list {
		if ou, ok := item.Data.(*OrganizationUnit); ok {
			if ou.Parent != nil && *ou.Parent == id {
				count++
			}
		}
	}

	return count, nil
}

// GetOrganizationUnitChildrenList implements organizationUnitStoreInterface.
func (f *ouFileBasedStore) GetOrganizationUnitChildrenList(id string, limit, offset int) ([]OrganizationUnitBasic, error) {
	list, err := f.GenericFileBasedStore.List()
	if err != nil {
		return nil, err
	}

	var children []OrganizationUnitBasic
	for _, item := range list {
		if ou, ok := item.Data.(*OrganizationUnit); ok {
			if ou.Parent != nil && *ou.Parent == id {
				basicOU := OrganizationUnitBasic{
					ID:          ou.ID,
					Handle:      ou.Handle,
					Name:        ou.Name,
					Description: ou.Description,
				}
				children = append(children, basicOU)
			}
		}
	}

	// Apply pagination
	start := offset
	end := offset + limit
	if start >= len(children) {
		return []OrganizationUnitBasic{}, nil
	}
	if end > len(children) {
		end = len(children)
	}

	return children[start:end], nil
}

// GetOrganizationUnitUsersCount implements organizationUnitStoreInterface.
func (f *ouFileBasedStore) GetOrganizationUnitUsersCount(id string) (int, error) {
	// File-based store doesn't support user associations
	return 0, nil
}

// GetOrganizationUnitUsersList implements organizationUnitStoreInterface.
func (f *ouFileBasedStore) GetOrganizationUnitUsersList(id string, limit, offset int) ([]User, error) {
	// File-based store doesn't support user associations
	return []User{}, nil
}

// GetOrganizationUnitGroupsCount implements organizationUnitStoreInterface.
func (f *ouFileBasedStore) GetOrganizationUnitGroupsCount(id string) (int, error) {
	// File-based store doesn't support group associations
	return 0, nil
}

// GetOrganizationUnitGroupsList implements organizationUnitStoreInterface.
func (f *ouFileBasedStore) GetOrganizationUnitGroupsList(id string, limit, offset int) ([]Group, error) {
	// File-based store doesn't support group associations
	return []Group{}, nil
}

// newOUFileBasedStore creates a new instance of a file-based store.
func newOUFileBasedStore() organizationUnitStoreInterface {
	return &ouFileBasedStore{
		GenericFileBasedStore: immutableresource.NewGenericFileBasedStore(entity.KeyTypeOU),
	}
}
