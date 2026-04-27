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

package inboundclient

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/asgardeo/thunder/internal/cert"
	inboundmodel "github.com/asgardeo/thunder/internal/inboundclient/model"
	sysconfig "github.com/asgardeo/thunder/internal/system/config"
	"github.com/asgardeo/thunder/internal/system/error/serviceerror"
	"github.com/asgardeo/thunder/internal/system/transaction"
	"github.com/asgardeo/thunder/tests/mocks/certmock"
)

type InboundClientServiceTestSuite struct {
	suite.Suite
}

func TestInboundClientServiceTestSuite(t *testing.T) {
	suite.Run(t, new(InboundClientServiceTestSuite))
}

func (suite *InboundClientServiceTestSuite) SetupTest() {
	sysconfig.ResetThunderRuntime()
	suite.Require().NoError(sysconfig.InitializeThunderRuntime("/tmp/test", &sysconfig.Config{}))
}

func newServiceForTest(store inboundClientStoreInterface) InboundClientServiceInterface {
	return newInboundClientService(store, transaction.NewNoOpTransactioner(), nil, nil, nil, nil, nil, nil, nil)
}

func newServiceWithCert(certService cert.CertificateServiceInterface) *inboundClientService {
	svc := newInboundClientService(
		nil, transaction.NewNoOpTransactioner(), certService, nil, nil, nil, nil, nil, nil,
	)
	return svc.(*inboundClientService)
}

func validInboundClient(id string) inboundmodel.InboundClient {
	return inboundmodel.InboundClient{
		ID:                        id,
		AuthFlowID:                "flow-1",
		RegistrationFlowID:        "reg-1",
		IsRegistrationFlowEnabled: true,
	}
}

func ptrInboundClient() *inboundmodel.InboundClient {
	c := validInboundClient("p1")
	return &c
}

func validOAuthProfileData() *inboundmodel.OAuthProfileData {
	return &inboundmodel.OAuthProfileData{
		RedirectURIs:            []string{"https://app.example.com/cb"},
		GrantTypes:              []string{"authorization_code"},
		ResponseTypes:           []string{"code"},
		TokenEndpointAuthMethod: "client_secret_basic",
	}
}

// ----- Inbound client CRUD -----

func (suite *InboundClientServiceTestSuite) TestCreateInboundClient_RunsValidationBeforePersist() {
	store := newInboundClientStoreInterfaceMock(suite.T())
	store.EXPECT().IsDeclarative(mock.Anything, "p1").Return(false)
	svc := newServiceForTest(store)

	p := validOAuthProfileData()
	p.GrantTypes = []string{"not_a_real_grant"}

	err := svc.CreateInboundClient(context.Background(), ptrInboundClient(), nil, p, false, "")

	assert.ErrorIs(suite.T(), err, ErrOAuthInvalidGrantType)
}

func (suite *InboundClientServiceTestSuite) TestCreateInboundClient_PersistsBoth() {
	store := newInboundClientStoreInterfaceMock(suite.T())
	store.EXPECT().IsDeclarative(mock.Anything, "p1").Return(false)
	store.EXPECT().CreateInboundClient(mock.Anything, mock.Anything).Return(nil)
	store.EXPECT().CreateOAuthProfile(mock.Anything, "p1", mock.Anything).Return(nil)

	svc := newServiceForTest(store)
	err := svc.CreateInboundClient(context.Background(), ptrInboundClient(),
		nil, validOAuthProfileData(), true, "")

	assert.NoError(suite.T(), err)
}

func (suite *InboundClientServiceTestSuite) TestCreateInboundClient_PersistsClientOnlyWhenOAuthNil() {
	store := newInboundClientStoreInterfaceMock(suite.T())
	store.EXPECT().IsDeclarative(mock.Anything, "p1").Return(false)
	store.EXPECT().CreateInboundClient(mock.Anything, mock.Anything).Return(nil)

	svc := newServiceForTest(store)
	err := svc.CreateInboundClient(context.Background(), ptrInboundClient(), nil, nil, false, "")

	assert.NoError(suite.T(), err)
}

func (suite *InboundClientServiceTestSuite) TestCreateInboundClient_RefusesDeclarative() {
	store := newInboundClientStoreInterfaceMock(suite.T())
	store.EXPECT().IsDeclarative(mock.Anything, "p1").Return(true)

	svc := newServiceForTest(store)
	err := svc.CreateInboundClient(context.Background(), ptrInboundClient(), nil, nil, false, "")

	assert.ErrorIs(suite.T(), err, ErrCannotModifyDeclarative)
}

func (suite *InboundClientServiceTestSuite) TestUpdateInboundClient_RefusesDeclarative() {
	store := newInboundClientStoreInterfaceMock(suite.T())
	store.EXPECT().IsDeclarative(mock.Anything, "p1").Return(true)

	svc := newServiceForTest(store)
	err := svc.UpdateInboundClient(context.Background(), ptrInboundClient(), nil, nil, false, "", "")

	assert.ErrorIs(suite.T(), err, ErrCannotModifyDeclarative)
}

func (suite *InboundClientServiceTestSuite) TestDeleteInboundClient_RefusesDeclarative() {
	store := newInboundClientStoreInterfaceMock(suite.T())
	store.EXPECT().IsDeclarative(mock.Anything, "p1").Return(true)

	svc := newServiceForTest(store)
	err := svc.DeleteInboundClient(context.Background(), "p1")

	assert.ErrorIs(suite.T(), err, ErrCannotModifyDeclarative)
}

func (suite *InboundClientServiceTestSuite) TestDelegatesPlainReads() {
	store := newInboundClientStoreInterfaceMock(suite.T())
	store.EXPECT().GetInboundClientList(mock.Anything, mock.Anything).
		Return([]inboundmodel.InboundClient{validInboundClient("p1")}, nil)
	store.EXPECT().IsDeclarative(mock.Anything, "p1").Return(true)

	svc := newServiceForTest(store)
	list, err := svc.GetInboundClientList(context.Background())
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), list, 1)

	assert.True(suite.T(), svc.IsDeclarative(context.Background(), "p1"))
}

func (suite *InboundClientServiceTestSuite) TestDeleteInboundClient() {
	store := newInboundClientStoreInterfaceMock(suite.T())
	store.EXPECT().IsDeclarative(mock.Anything, "p1").Return(false)
	store.EXPECT().DeleteInboundClient(mock.Anything, "p1").Return(nil)

	svc := newServiceForTest(store)
	assert.NoError(suite.T(), svc.DeleteInboundClient(context.Background(), "p1"))
}

func (suite *InboundClientServiceTestSuite) TestStorePropagatesErrors() {
	storeErr := errors.New("db error")
	store := newInboundClientStoreInterfaceMock(suite.T())
	store.EXPECT().IsDeclarative(mock.Anything, "p1").Return(false)
	store.EXPECT().CreateInboundClient(mock.Anything, mock.Anything).Return(storeErr)

	svc := newServiceForTest(store)
	err := svc.CreateInboundClient(context.Background(), ptrInboundClient(), nil, nil, false, "")

	assert.ErrorIs(suite.T(), err, storeErr)
}

// ----- ValidateCertificateInput -----

func (suite *InboundClientServiceTestSuite) TestValidateCertificateInput_Empty() {
	c, err := validateCertificateInput(cert.CertificateReferenceTypeOAuthApp, "ref-1", "", nil)

	suite.Nil(c)
	suite.Nil(err)
}

func (suite *InboundClientServiceTestSuite) TestValidateCertificateInput_JWKS_Success() {
	c, err := validateCertificateInput(cert.CertificateReferenceTypeOAuthApp, "ref-1", "existing",
		&inboundmodel.Certificate{Type: cert.CertificateTypeJWKS, Value: `{"keys":[]}`})

	suite.Nil(err)
	suite.NotNil(c)
	suite.Equal("existing", c.ID)
	suite.Equal(cert.CertificateTypeJWKS, c.Type)
	suite.Equal(cert.CertificateReferenceTypeOAuthApp, c.RefType)
	suite.Equal("ref-1", c.RefID)
}

func (suite *InboundClientServiceTestSuite) TestValidateCertificateInput_JWKS_MissingValue() {
	c, err := validateCertificateInput(cert.CertificateReferenceTypeOAuthApp, "ref-1", "",
		&inboundmodel.Certificate{Type: cert.CertificateTypeJWKS, Value: ""})

	suite.Nil(c)
	suite.ErrorIs(err, ErrCertValueRequired)
}

func (suite *InboundClientServiceTestSuite) TestValidateCertificateInput_JWKSURI_Success() {
	c, err := validateCertificateInput(cert.CertificateReferenceTypeOAuthApp, "ref-1", "",
		&inboundmodel.Certificate{Type: cert.CertificateTypeJWKSURI, Value: "https://example.com/jwks"})

	suite.Nil(err)
	suite.Equal(cert.CertificateTypeJWKSURI, c.Type)
}

func (suite *InboundClientServiceTestSuite) TestValidateCertificateInput_JWKSURI_Invalid() {
	c, err := validateCertificateInput(cert.CertificateReferenceTypeOAuthApp, "ref-1", "",
		&inboundmodel.Certificate{Type: cert.CertificateTypeJWKSURI, Value: "not-a-uri"})

	suite.Nil(c)
	suite.ErrorIs(err, ErrCertInvalidJWKSURI)
}

func (suite *InboundClientServiceTestSuite) TestValidateCertificateInput_InvalidType() {
	c, err := validateCertificateInput(cert.CertificateReferenceTypeOAuthApp, "ref-1", "",
		&inboundmodel.Certificate{Type: "bogus", Value: "x"})

	suite.Nil(c)
	suite.ErrorIs(err, ErrCertInvalidType)
}

// ----- CreateCertificate -----

func (suite *InboundClientServiceTestSuite) TestCreateCertificate_Nil() {
	svc := newServiceWithCert(certmock.NewCertificateServiceInterfaceMock(suite.T()))

	out, vErr, opErr := svc.createCertificate(context.Background(),
		cert.CertificateReferenceTypeOAuthApp, "ref-1", nil)

	suite.Nil(out)
	suite.Nil(vErr)
	suite.Nil(opErr)
}

func (suite *InboundClientServiceTestSuite) TestCreateCertificate_Success() {
	mockCert := certmock.NewCertificateServiceInterfaceMock(suite.T())
	mockCert.EXPECT().CreateCertificate(mock.Anything, mock.Anything).
		Return(&cert.Certificate{}, nil)
	svc := newServiceWithCert(mockCert)

	in := &inboundmodel.Certificate{Type: cert.CertificateTypeJWKS, Value: `{}`}
	out, vErr, opErr := svc.createCertificate(context.Background(),
		cert.CertificateReferenceTypeOAuthApp, "ref-1", in)

	suite.Nil(vErr)
	suite.Nil(opErr)
	suite.Equal(cert.CertificateTypeJWKS, out.Type)
	suite.Equal(`{}`, out.Value)
}

func (suite *InboundClientServiceTestSuite) TestCreateCertificate_InvalidInput() {
	svc := newServiceWithCert(certmock.NewCertificateServiceInterfaceMock(suite.T()))

	in := &inboundmodel.Certificate{Type: cert.CertificateTypeJWKSURI, Value: "not-a-uri"}
	out, vErr, opErr := svc.createCertificate(context.Background(),
		cert.CertificateReferenceTypeOAuthApp, "ref-1", in)

	suite.Nil(out)
	suite.Nil(opErr)
	suite.ErrorIs(vErr, ErrCertInvalidJWKSURI)
}

func (suite *InboundClientServiceTestSuite) TestCreateCertificate_ServiceError() {
	mockCert := certmock.NewCertificateServiceInterfaceMock(suite.T())
	clientErr := &serviceerror.ServiceError{Type: serviceerror.ClientErrorType, Code: "C-1"}
	mockCert.EXPECT().CreateCertificate(mock.Anything, mock.Anything).Return(nil, clientErr)
	svc := newServiceWithCert(mockCert)

	in := &inboundmodel.Certificate{Type: cert.CertificateTypeJWKS, Value: `{}`}
	out, vErr, opErr := svc.createCertificate(context.Background(),
		cert.CertificateReferenceTypeOAuthApp, "ref-1", in)

	suite.Nil(out)
	suite.Nil(vErr)
	suite.Equal(CertOpCreate, opErr.Operation)
	suite.Same(clientErr, opErr.Underlying)
	suite.True(opErr.IsClientError())
}

// ----- GetCertificate -----

func (suite *InboundClientServiceTestSuite) TestGetCertificate_NotFound() {
	mockCert := certmock.NewCertificateServiceInterfaceMock(suite.T())
	mockCert.EXPECT().
		GetCertificateByReference(mock.Anything, cert.CertificateReferenceTypeApplication, "ref-1").
		Return(nil, &cert.ErrorCertificateNotFound)
	svc := newServiceWithCert(mockCert)

	out, err := svc.GetCertificate(context.Background(), cert.CertificateReferenceTypeApplication, "ref-1")

	suite.Nil(out)
	suite.Nil(err)
}

func (suite *InboundClientServiceTestSuite) TestGetCertificate_Success() {
	mockCert := certmock.NewCertificateServiceInterfaceMock(suite.T())
	mockCert.EXPECT().
		GetCertificateByReference(mock.Anything, cert.CertificateReferenceTypeApplication, "ref-1").
		Return(&cert.Certificate{Type: cert.CertificateTypeJWKS, Value: `{}`}, nil)
	svc := newServiceWithCert(mockCert)

	out, err := svc.GetCertificate(context.Background(), cert.CertificateReferenceTypeApplication, "ref-1")

	suite.Nil(err)
	suite.Equal(cert.CertificateTypeJWKS, out.Type)
}

func (suite *InboundClientServiceTestSuite) TestGetCertificate_ServerError() {
	mockCert := certmock.NewCertificateServiceInterfaceMock(suite.T())
	srvErr := &serviceerror.ServiceError{Type: serviceerror.ServerErrorType, Code: "S-1"}
	mockCert.EXPECT().
		GetCertificateByReference(mock.Anything, cert.CertificateReferenceTypeApplication, "ref-1").
		Return(nil, srvErr)
	svc := newServiceWithCert(mockCert)

	out, err := svc.GetCertificate(context.Background(), cert.CertificateReferenceTypeApplication, "ref-1")

	suite.Nil(out)
	suite.Equal(CertOpRetrieve, err.Operation)
	suite.False(err.IsClientError())
}

// ----- DeleteCertificate -----

func (suite *InboundClientServiceTestSuite) TestDeleteCertificate_Success() {
	mockCert := certmock.NewCertificateServiceInterfaceMock(suite.T())
	mockCert.EXPECT().
		DeleteCertificateByReference(mock.Anything, cert.CertificateReferenceTypeApplication, "ref-1").
		Return(nil)
	svc := newServiceWithCert(mockCert)

	err := svc.deleteCertificate(context.Background(), cert.CertificateReferenceTypeApplication, "ref-1")

	suite.Nil(err)
}

func (suite *InboundClientServiceTestSuite) TestDeleteCertificate_Error() {
	mockCert := certmock.NewCertificateServiceInterfaceMock(suite.T())
	clientErr := &serviceerror.ServiceError{Type: serviceerror.ClientErrorType, Code: "D-1"}
	mockCert.EXPECT().
		DeleteCertificateByReference(mock.Anything, cert.CertificateReferenceTypeOAuthApp, "ref-1").
		Return(clientErr)
	svc := newServiceWithCert(mockCert)

	err := svc.deleteCertificate(context.Background(), cert.CertificateReferenceTypeOAuthApp, "ref-1")

	suite.NotNil(err)
	suite.Equal(CertOpDelete, err.Operation)
	suite.Equal(cert.CertificateReferenceTypeOAuthApp, err.RefType)
}

// ----- SyncCertificate -----

func (suite *InboundClientServiceTestSuite) TestSyncCertificate_NoOp_NoExistingNoInput() {
	mockCert := certmock.NewCertificateServiceInterfaceMock(suite.T())
	mockCert.EXPECT().
		GetCertificateByReference(mock.Anything, mock.Anything, mock.Anything).
		Return(nil, &cert.ErrorCertificateNotFound)
	svc := newServiceWithCert(mockCert)

	out, vErr, opErr := svc.syncCertificate(context.Background(),
		cert.CertificateReferenceTypeApplication, "ref-1", nil)

	suite.Nil(out)
	suite.Nil(vErr)
	suite.Nil(opErr)
}

func (suite *InboundClientServiceTestSuite) TestSyncCertificate_CreateWhenAbsent() {
	mockCert := certmock.NewCertificateServiceInterfaceMock(suite.T())
	mockCert.EXPECT().
		GetCertificateByReference(mock.Anything, mock.Anything, mock.Anything).
		Return(nil, &cert.ErrorCertificateNotFound)
	mockCert.EXPECT().CreateCertificate(mock.Anything, mock.Anything).
		Return(&cert.Certificate{}, nil)
	svc := newServiceWithCert(mockCert)

	out, vErr, opErr := svc.syncCertificate(context.Background(),
		cert.CertificateReferenceTypeApplication, "ref-1",
		&inboundmodel.Certificate{Type: cert.CertificateTypeJWKS, Value: `{}`})

	suite.Nil(vErr)
	suite.Nil(opErr)
	suite.NotNil(out)
}

func (suite *InboundClientServiceTestSuite) TestSyncCertificate_UpdateWhenPresent() {
	mockCert := certmock.NewCertificateServiceInterfaceMock(suite.T())
	mockCert.EXPECT().
		GetCertificateByReference(mock.Anything, mock.Anything, mock.Anything).
		Return(&cert.Certificate{ID: "cert-1"}, nil)
	mockCert.EXPECT().UpdateCertificateByID(mock.Anything, "cert-1", mock.Anything).
		Return(&cert.Certificate{}, nil)
	svc := newServiceWithCert(mockCert)

	out, vErr, opErr := svc.syncCertificate(context.Background(),
		cert.CertificateReferenceTypeApplication, "ref-1",
		&inboundmodel.Certificate{Type: cert.CertificateTypeJWKS, Value: `{}`})

	suite.Nil(vErr)
	suite.Nil(opErr)
	suite.NotNil(out)
}

func (suite *InboundClientServiceTestSuite) TestSyncCertificate_DeleteWhenInputEmpty() {
	mockCert := certmock.NewCertificateServiceInterfaceMock(suite.T())
	mockCert.EXPECT().
		GetCertificateByReference(mock.Anything, mock.Anything, mock.Anything).
		Return(&cert.Certificate{ID: "cert-1"}, nil)
	mockCert.EXPECT().
		DeleteCertificateByReference(mock.Anything, mock.Anything, mock.Anything).
		Return(nil)
	svc := newServiceWithCert(mockCert)

	out, vErr, opErr := svc.syncCertificate(context.Background(),
		cert.CertificateReferenceTypeApplication, "ref-1", nil)

	suite.Nil(out)
	suite.Nil(vErr)
	suite.Nil(opErr)
}

func (suite *InboundClientServiceTestSuite) TestSyncCertificate_ValidationError() {
	mockCert := certmock.NewCertificateServiceInterfaceMock(suite.T())
	mockCert.EXPECT().
		GetCertificateByReference(mock.Anything, mock.Anything, mock.Anything).
		Return(nil, &cert.ErrorCertificateNotFound)
	svc := newServiceWithCert(mockCert)

	out, vErr, opErr := svc.syncCertificate(context.Background(),
		cert.CertificateReferenceTypeApplication, "ref-1",
		&inboundmodel.Certificate{Type: "bogus", Value: "x"})

	suite.Nil(out)
	suite.Nil(opErr)
	suite.ErrorIs(vErr, ErrCertInvalidType)
}

func (suite *InboundClientServiceTestSuite) TestGetInboundClientByEntityID_Delegates() {
	store := newInboundClientStoreInterfaceMock(suite.T())
	want := &inboundmodel.InboundClient{ID: "p1"}
	store.EXPECT().GetInboundClientByEntityID(mock.Anything, "p1").Return(want, nil)

	svc := newServiceForTest(store)
	got, err := svc.GetInboundClientByEntityID(context.Background(), "p1")

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "p1", got.ID)
}

func (suite *InboundClientServiceTestSuite) TestGetOAuthProfileByEntityID_Delegates() {
	store := newInboundClientStoreInterfaceMock(suite.T())
	want := &inboundmodel.OAuthProfile{AppID: "p1"}
	store.EXPECT().GetOAuthProfileByEntityID(mock.Anything, "p1").Return(want, nil)

	svc := newServiceForTest(store)
	got, err := svc.GetOAuthProfileByEntityID(context.Background(), "p1")

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "p1", got.AppID)
}

func (suite *InboundClientServiceTestSuite) TestUpdateInboundClient_ValidationFails() {
	store := newInboundClientStoreInterfaceMock(suite.T())
	store.EXPECT().IsDeclarative(mock.Anything, "p1").Return(false)
	svc := newServiceForTest(store)

	p := validOAuthProfileData()
	p.GrantTypes = []string{"not_a_real_grant"}

	err := svc.UpdateInboundClient(context.Background(), ptrInboundClient(), nil, p, false, "", "")
	assert.ErrorIs(suite.T(), err, ErrOAuthInvalidGrantType)
}

func (suite *InboundClientServiceTestSuite) TestUpdateInboundClient_Succeeds() {
	store := newInboundClientStoreInterfaceMock(suite.T())
	store.EXPECT().IsDeclarative(mock.Anything, "p1").Return(false)
	store.EXPECT().UpdateInboundClient(mock.Anything, mock.Anything).Return(nil)
	// syncOAuthProfile path: GetOAuthProfileByEntityID returns not found → CreateOAuthProfile
	store.EXPECT().GetOAuthProfileByEntityID(mock.Anything, "p1").Return(nil, ErrInboundClientNotFound)
	store.EXPECT().CreateOAuthProfile(mock.Anything, "p1", mock.Anything).Return(nil)

	mockCert := certmock.NewCertificateServiceInterfaceMock(suite.T())
	// syncCertificate for app cert (nil input): gets existing (not found), no update needed
	mockCert.EXPECT().GetCertificateByReference(mock.Anything, mock.Anything, mock.Anything).
		Return(nil, &cert.ErrorCertificateNotFound)

	svc := newInboundClientService(store, transaction.NewNoOpTransactioner(), mockCert, nil, nil, nil, nil, nil, nil)
	err := svc.UpdateInboundClient(context.Background(), ptrInboundClient(), nil, validOAuthProfileData(), true, "", "")
	assert.NoError(suite.T(), err)
}

func (suite *InboundClientServiceTestSuite) TestValidate_ValidProfile() {
	store := newInboundClientStoreInterfaceMock(suite.T())
	svc := newServiceForTest(store)

	err := svc.Validate(context.Background(), ptrInboundClient(), validOAuthProfileData(), true)
	assert.NoError(suite.T(), err)
}

func (suite *InboundClientServiceTestSuite) TestValidate_InvalidGrantType() {
	store := newInboundClientStoreInterfaceMock(suite.T())
	svc := newServiceForTest(store)

	p := validOAuthProfileData()
	p.GrantTypes = []string{"bogus_grant"}

	err := svc.Validate(context.Background(), ptrInboundClient(), p, false)
	assert.ErrorIs(suite.T(), err, ErrOAuthInvalidGrantType)
}

func (suite *InboundClientServiceTestSuite) TestValidateRedirectURIs_WildcardInHost_Rejected() {
	p := &inboundmodel.OAuthProfileData{
		RedirectURIs: []string{"https://*.example.com/cb"},
		GrantTypes:   []string{"authorization_code"},
	}
	err := validateRedirectURIs(p)
	assert.ErrorIs(suite.T(), err, ErrOAuthInvalidRedirectURI)
}

func (suite *InboundClientServiceTestSuite) TestValidateRedirectURIs_WildcardInQuery_Rejected() {
	p := &inboundmodel.OAuthProfileData{
		RedirectURIs: []string{"https://app.example.com/cb?foo=*"},
		GrantTypes:   []string{"authorization_code"},
	}
	err := validateRedirectURIs(p)
	assert.ErrorIs(suite.T(), err, ErrOAuthInvalidRedirectURI)
}

func (suite *InboundClientServiceTestSuite) TestValidatePublicClient_PKCENotRequired_Fails() {
	p := &inboundmodel.OAuthProfileData{
		PublicClient:            true,
		PKCERequired:            false,
		TokenEndpointAuthMethod: "none",
	}
	err := validatePublicClient(p)
	assert.ErrorIs(suite.T(), err, ErrOAuthPublicClientMustHavePKCE)
}

func (suite *InboundClientServiceTestSuite) TestValidateTokenEndpointAuthMethod_InvalidMethod() {
	p := &inboundmodel.OAuthProfileData{
		TokenEndpointAuthMethod: "bogus_method",
	}
	err := validateTokenEndpointAuthMethod(p, false)
	assert.ErrorIs(suite.T(), err, ErrOAuthInvalidTokenEndpointAuthMethod)
}
