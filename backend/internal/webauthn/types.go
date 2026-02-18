package webauthn

import (
	"time"

	"github.com/asgardeo/thunder/internal/webauthn/protocol"
)

// User interface as required by the library.
type User interface {
	WebAuthnID() []byte
	WebAuthnName() string
	WebAuthnDisplayName() string
	WebAuthnCredentials() []Credential
}

// Credential represents a stored credential.
type Credential struct {
	ID              []byte
	PublicKey       []byte
	AttestationType string
	Transport       []string
	Authenticator   Authenticator
}

// Authenticator represents authenticator data.
type Authenticator struct {
	AAGUID       []byte
	SignCount    uint32
	CloneWarning bool
	Attachment   protocol.AuthenticatorAttachment
}

// SessionData stores session state between steps.
type SessionData struct {
	Challenge            string
	UserID               []byte
	AllowedCredentialIDs [][]byte
	UserVerification     protocol.UserVerificationRequirement
	Extensions           protocol.AuthenticationExtensions

	// Additional fields used by store
	RelyingPartyID string
	Expires        time.Time
	CredParams     []protocol.CredentialParameter
	Mediation      protocol.CredentialMediationRequirement
}

// RegistrationOption configures registration options.
type RegistrationOption func(*protocol.CredentialCreation)

// LoginOption configures login options.
type LoginOption func(*protocol.CredentialAssertion)
