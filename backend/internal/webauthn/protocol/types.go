package protocol

import (
	"fmt"
)

// Enums
type CredentialType string
type ConveyancePreference string
type UserVerificationRequirement string
type AuthenticatorAttachment string
type ResidentKeyRequirement string
type CredentialMediationRequirement string

// CredentialCreation represents PublicKeyCredentialCreationOptions.
type CredentialCreation struct {
	Response CreationResponse `json:"response"`
}

type CreationResponse struct {
	Challenge              []byte                   `json:"challenge"`
	RelyingParty           RelyingPartyEntity       `json:"rp"`
	User                   UserEntity               `json:"user"`
	Parameters             []CredentialParameter    `json:"pubKeyCredParams"`
	AuthenticatorSelection AuthenticatorSelection   `json:"authenticatorSelection,omitempty"`
	Timeout                int                      `json:"timeout,omitempty"`
	CredentialExcludeList  []CredentialDescriptor   `json:"excludeCredentials,omitempty"`
	Extensions             AuthenticationExtensions `json:"extensions,omitempty"`
	Attestation            ConveyancePreference     `json:"attestation,omitempty"`
}

// CredentialAssertion represents PublicKeyCredentialRequestOptions.
type CredentialAssertion struct {
	Response AssertionResponse `json:"response"`
}

type AssertionResponse struct {
	Challenge          []byte                      `json:"challenge"`
	Timeout            int                         `json:"timeout,omitempty"`
	RelyingPartyID     string                      `json:"rpId,omitempty"`
	AllowedCredentials []CredentialDescriptor      `json:"allowCredentials,omitempty"`
	UserVerification   UserVerificationRequirement `json:"userVerification,omitempty"`
	Extensions         AuthenticationExtensions    `json:"extensions,omitempty"`
}

// ParsedCredentialCreationData represents the parsed registration response.
type ParsedCredentialCreationData struct {
	ID       string
	RawID    []byte
	Type     CredentialType
	Response ParsedCreationResponse
	Raw      CredentialCreationResponse
}

type ParsedCreationResponse struct {
	ClientDataJSON    []byte
	AttestationObject []byte
	CollectedClientData CollectedClientData
	AttestationObjectParsed AttestationObject
}

type AttestationObject struct {
	AuthData     []byte                 `cbor:"authData"`
	Format       string                 `cbor:"fmt"`
	AttStatement map[string]interface{} `cbor:"attStmt"`
}

// CredentialCreationResponse is the struct for unmarshalling the JSON response.
type CredentialCreationResponse struct {
	ID       string                           `json:"id"`
	RawID    string                           `json:"rawId"`
	Type     CredentialType                   `json:"type"`
	Response AuthenticatorAttestationResponse `json:"response"`
}

type AuthenticatorAttestationResponse struct {
	ClientDataJSON    string `json:"clientDataJSON"`    // Base64 encoded in JSON
	AttestationObject string `json:"attestationObject"` // Base64 encoded in JSON
}

// ParsedCredentialAssertionData represents the parsed authentication response.
type ParsedCredentialAssertionData struct {
	ParsedPublicKeyCredential
	Response ParsedAssertionResponse
	Raw      CredentialAssertionResponse
}

type ParsedPublicKeyCredential struct {
	RawID            []byte
	ParsedCredential
}

type ParsedCredential struct {
	ID   string
	Type CredentialType
}

type ParsedAssertionResponse struct {
	CollectedClientData CollectedClientData
	AuthenticatorData   AuthenticatorData
	Signature           []byte
	UserHandle          []byte
}

type CredentialAssertionResponse struct {
	PublicKeyCredential PublicKeyCredential            `json:"publicKeyCredential"`
	AssertionResponse   AuthenticatorAssertionResponse `json:"assertionResponse"`
}

type PublicKeyCredential struct {
	Credential Credential `json:"credential"`
	RawID      []byte     `json:"rawId"`
}

type Credential struct {
	ID   string         `json:"id"`
	Type CredentialType `json:"type"`
}

type AuthenticatorAssertionResponse struct {
	AuthenticatorResponse AuthenticatorResponse `json:"authenticatorResponse"`
	AuthenticatorData     []byte                `json:"authenticatorData"`
	Signature             []byte                `json:"signature"`
	UserHandle            []byte                `json:"userHandle"`
}

type AuthenticatorResponse struct {
	ClientDataJSON []byte `json:"clientDataJSON"`
}

// Common types

type RelyingPartyEntity struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type UserEntity struct {
	ID          []byte `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

type CredentialParameter struct {
	Type      CredentialType `json:"type"`
	Algorithm int            `json:"alg"`
}

type AuthenticatorSelection struct {
	AuthenticatorAttachment AuthenticatorAttachment `json:"authenticatorAttachment,omitempty"`
	RequireResidentKey      *bool                   `json:"requireResidentKey,omitempty"`
	ResidentKey             ResidentKeyRequirement  `json:"residentKey,omitempty"`
	UserVerification        UserVerificationRequirement `json:"userVerification,omitempty"`
}

type CredentialDescriptor struct {
	Type       CredentialType `json:"type"`
	ID         []byte         `json:"id"`
	Transports []string       `json:"transports,omitempty"`
}

type AuthenticationExtensions map[string]interface{}

type CollectedClientData struct {
	Type        string `json:"type"`
	Challenge   string `json:"challenge"`
	Origin      string `json:"origin"`
	CrossOrigin bool   `json:"crossOrigin,omitempty"`
}

type AuthenticatorData struct {
	RPIDHash []byte
	Flags    AuthenticatorFlags
	Counter  uint32
	// AttestedCredentialData is optional and variable length
}

func (a *AuthenticatorData) Unmarshal(data []byte) error {
	if len(data) < 37 {
		return fmt.Errorf("authenticator data too short")
	}
	a.RPIDHash = data[:32]
	a.Flags = AuthenticatorFlags(data[32])
	a.Counter = uint32(data[33])<<24 | uint32(data[34])<<16 | uint32(data[35])<<8 | uint32(data[36])
	return nil
}

type AuthenticatorFlags byte

const (
	FlagUserPresent    AuthenticatorFlags = 0x01
	FlagUserVerified   AuthenticatorFlags = 0x04
	FlagAttestedData   AuthenticatorFlags = 0x40
	FlagExtensionData  AuthenticatorFlags = 0x80
)

// Constants for algorithms (COSE)
const (
	ES256 = -7
	RS256 = -257
)
