package protocol

import (
	"fmt"
)

// CredentialType represents the type of credential.
type CredentialType string

// ConveyancePreference represents the preference for attestation conveyance.
type ConveyancePreference string

// UserVerificationRequirement represents the requirement for user verification.
type UserVerificationRequirement string

// AuthenticatorAttachment represents the attachment mechanism of the authenticator.
type AuthenticatorAttachment string

// ResidentKeyRequirement represents the requirement for a resident key.
type ResidentKeyRequirement string

// CredentialMediationRequirement represents the requirement for credential mediation.
type CredentialMediationRequirement string

// CredentialCreation represents PublicKeyCredentialCreationOptions.
type CredentialCreation struct {
	Response CreationResponse `json:"response"`
}

// CreationResponse represents the response parameters for credential creation.
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

// AssertionResponse represents the response parameters for credential assertion.
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

// ParsedCreationResponse represents the parsed components of a creation response.
type ParsedCreationResponse struct {
	ClientDataJSON          []byte
	AttestationObject       []byte
	CollectedClientData     CollectedClientData
	AttestationObjectParsed AttestationObject
}

// AttestationObject represents the decoded attestation object.
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

// AuthenticatorAttestationResponse represents the raw response from the authenticator during registration.
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

// ParsedPublicKeyCredential represents a parsed public key credential.
type ParsedPublicKeyCredential struct {
	RawID []byte
	ParsedCredential
}

// ParsedCredential represents a basic parsed credential.
type ParsedCredential struct {
	ID   string
	Type CredentialType
}

// ParsedAssertionResponse represents the parsed components of an assertion response.
type ParsedAssertionResponse struct {
	CollectedClientData CollectedClientData
	AuthenticatorData   AuthenticatorData
	Signature           []byte
	UserHandle          []byte
}

// CredentialAssertionResponse represents the raw JSON response for authentication.
type CredentialAssertionResponse struct {
	PublicKeyCredential PublicKeyCredential            `json:"publicKeyCredential"`
	AssertionResponse   AuthenticatorAssertionResponse `json:"assertionResponse"`
}

// PublicKeyCredential represents a public key credential.
type PublicKeyCredential struct {
	Credential Credential `json:"credential"`
	RawID      []byte     `json:"rawId"`
}

// Credential represents a credential with an ID and type.
type Credential struct {
	ID   string         `json:"id"`
	Type CredentialType `json:"type"`
}

// AuthenticatorAssertionResponse represents the authenticator's response to an assertion request.
type AuthenticatorAssertionResponse struct {
	AuthenticatorResponse AuthenticatorResponse `json:"authenticatorResponse"`
	AuthenticatorData     []byte                `json:"authenticatorData"`
	Signature             []byte                `json:"signature"`
	UserHandle            []byte                `json:"userHandle"`
}

// AuthenticatorResponse represents the authenticator response data.
type AuthenticatorResponse struct {
	ClientDataJSON []byte `json:"clientDataJSON"`
}

// Common types

// RelyingPartyEntity represents the relying party.
type RelyingPartyEntity struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// UserEntity represents the user entity.
type UserEntity struct {
	ID          []byte `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

// CredentialParameter represents a credential parameter (type and algorithm).
type CredentialParameter struct {
	Type      CredentialType `json:"type"`
	Algorithm int            `json:"alg"`
}

// AuthenticatorSelection represents options for authenticator selection.
type AuthenticatorSelection struct {
	AuthenticatorAttachment AuthenticatorAttachment `json:"authenticatorAttachment,omitempty"`
	RequireResidentKey      *bool                   `json:"requireResidentKey,omitempty"`
	ResidentKey             ResidentKeyRequirement  `json:"residentKey,omitempty"`
	UserVerification        UserVerificationRequirement `json:"userVerification,omitempty"`
}

// CredentialDescriptor represents a descriptor for a credential.
type CredentialDescriptor struct {
	Type       CredentialType `json:"type"`
	ID         []byte         `json:"id"`
	Transports []string       `json:"transports,omitempty"`
}

// AuthenticationExtensions represents authentication extensions.
type AuthenticationExtensions map[string]interface{}

// CollectedClientData represents the client data collected during the ceremony.
type CollectedClientData struct {
	Type        string `json:"type"`
	Challenge   string `json:"challenge"`
	Origin      string `json:"origin"`
	CrossOrigin bool   `json:"crossOrigin,omitempty"`
}

// AuthenticatorData represents the authenticator data structure.
type AuthenticatorData struct {
	RPIDHash []byte
	Flags    AuthenticatorFlags
	Counter  uint32
	// AttestedCredentialData is optional and variable length
}

// Unmarshal decodes bytes into AuthenticatorData.
func (a *AuthenticatorData) Unmarshal(data []byte) error {
	if len(data) < 37 {
		return fmt.Errorf("authenticator data too short")
	}
	a.RPIDHash = data[:32]
	a.Flags = AuthenticatorFlags(data[32])
	a.Counter = uint32(data[33])<<24 | uint32(data[34])<<16 | uint32(data[35])<<8 | uint32(data[36])
	return nil
}

// AuthenticatorFlags represents the flags in authenticator data.
type AuthenticatorFlags byte

// Authenticator Flags
const (
	FlagUserPresent   AuthenticatorFlags = 0x01
	FlagUserVerified  AuthenticatorFlags = 0x04
	FlagAttestedData  AuthenticatorFlags = 0x40
	FlagExtensionData AuthenticatorFlags = 0x80
)

// Constants for algorithms (COSE)
const (
	ES256 = -7
	RS256 = -257
)
