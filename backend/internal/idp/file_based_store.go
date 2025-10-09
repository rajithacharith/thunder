package idp

import (
	"errors"

	filebackedruntime "github.com/asgardeo/thunder/internal/file_backed_runtime"
	"gopkg.in/yaml.v3"
)

type fileBasedStore struct {
	store map[string]IDPDTO
}

// CreateIdentityProvider implements idpStoreInterface.
func (f *fileBasedStore) CreateIdentityProvider(idp IDPDTO) error {
	panic("unimplemented")
}

// DeleteIdentityProvider implements idpStoreInterface.
func (f *fileBasedStore) DeleteIdentityProvider(idpID string) error {
	panic("unimplemented")
}

// GetIdentityProvider implements idpStoreInterface.
func (f *fileBasedStore) GetIdentityProvider(idpID string) (*IDPDTO, error) {
	panic("unimplemented")
}

// GetIdentityProviderByName implements idpStoreInterface.
func (f *fileBasedStore) GetIdentityProviderByName(idpName string) (*IDPDTO, error) {
	panic("unimplemented")
}

// GetIdentityProviderList implements idpStoreInterface.
func (f *fileBasedStore) GetIdentityProviderList() ([]BasicIDPDTO, error) {
	panic("unimplemented")
}

// UpdateIdentityProvider implements idpStoreInterface.
func (f *fileBasedStore) UpdateIdentityProvider(idp *IDPDTO) error {
	panic("unimplemented")
}

var _ idpStoreInterface = (*fileBasedStore)(nil)

func newFileBasedStore() *fileBasedStore {
	fileConfigs := filebackedruntime.GetConfig().IDPs
	store := make(map[string]IDPDTO)
	for _, fileConfig := range fileConfigs {
		idp, err := convertFileIDPConfigToDTO(fileConfig)
		if err != nil {
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
