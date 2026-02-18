package webauthn

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"

	"github.com/asgardeo/thunder/internal/webauthn/protocol"
)

// BeginRegistration initiates the registration ceremony.
func (w *WebAuthn) BeginRegistration(user User, opts ...RegistrationOption) (*protocol.CredentialCreation, *SessionData, error) {
	// Generate challenge
	challenge := make([]byte, 32)
	if _, err := rand.Read(challenge); err != nil {
		return nil, nil, fmt.Errorf("failed to generate challenge: %w", err)
	}

	creation := &protocol.CredentialCreation{
		Response: protocol.CreationResponse{
			RelyingParty: protocol.RelyingPartyEntity{
				ID:   w.Config.RPID,
				Name: w.Config.RPDisplayName,
			},
			User: protocol.UserEntity{
				ID:          user.WebAuthnID(),
				Name:        user.WebAuthnName(),
				DisplayName: user.WebAuthnDisplayName(),
			},
			Challenge: challenge,
			Parameters: []protocol.CredentialParameter{
				{Type: protocol.PublicKeyCredentialType, Algorithm: protocol.ES256},
				{Type: protocol.PublicKeyCredentialType, Algorithm: protocol.RS256},
			},
			Timeout:     60000, // 60 seconds
			Attestation: protocol.PreferNoAttestation,
		},
	}

	// Apply options
	for _, opt := range opts {
		opt(creation)
	}

	// Create session data
	session := &SessionData{
		Challenge:        base64.RawURLEncoding.EncodeToString(challenge),
		UserID:           user.WebAuthnID(),
		UserVerification: protocol.VerificationPreferred,
	}

	return creation, session, nil
}

// CreateCredential verifies the registration response and creates a credential.
func (w *WebAuthn) CreateCredential(user User, session SessionData, response *protocol.ParsedCredentialCreationData) (*Credential, error) {
	// 1. Verify Challenge
	// The response challenge is base64url encoded in ClientDataJSON.
	// ParsedCredentialCreationData has ParsedCreationResponse with CollectedClientData.
	// CollectedClientData has Challenge as string (base64url encoded).

	if response.Response.CollectedClientData.Challenge != session.Challenge {
		return nil, fmt.Errorf("challenge mismatch")
	}

	// 2. Verify Origin
	originFound := false
	for _, origin := range w.Config.RPOrigins {
		if response.Response.CollectedClientData.Origin == origin {
			originFound = true
			break
		}
	}
	if !originFound {
		return nil, fmt.Errorf("origin mismatch")
	}

	// 3. Verify RPID Hash
	authData := response.Response.AttestationObjectParsed.AuthData
	if len(authData) < 37 {
		return nil, fmt.Errorf("authenticator data too short")
	}

	rpIDHash := sha256.Sum256([]byte(w.Config.RPID))
	if !bytes.Equal(rpIDHash[:], authData[:32]) {
		return nil, fmt.Errorf("RPID hash mismatch")
	}

	// 4. Parse Attestation Object to get Public Key and Credential ID
	// attestationObject is parsed into AttestationObject struct in protocol package.
	// It has AuthData []byte.
	// We need to parse AuthData to get AttestedCredentialData.

	flags := protocol.AuthenticatorFlags(authData[32])
	if flags&protocol.FlagAttestedData == 0 {
		return nil, fmt.Errorf("attested credential data flag not set")
	}

	// Attested Credential Data starts at index 37
	// AAGUID (16) + CredentialIDLen (2) + CredentialID (len) + Public Key (variable)
	if len(authData) < 37+16+2 {
		return nil, fmt.Errorf("authenticator data too short for attested credential data")
	}

	offset := 37
	aaguid := authData[offset : offset+16]
	offset += 16

	credIDLen := binary.BigEndian.Uint16(authData[offset : offset+2])
	offset += 2

	if len(authData) < offset+int(credIDLen) {
		return nil, fmt.Errorf("authenticator data too short for credential ID")
	}

	credentialID := authData[offset : offset+int(credIDLen)]
	offset += int(credIDLen)

	// Parse the public key (COSE Key) to determine its length
	// We need the raw bytes of the public key for verification
	_, bytesConsumed, err := protocol.UnmarshalNext(authData[offset:])
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	publicKey := authData[offset : offset+bytesConsumed]

	signCount := binary.BigEndian.Uint32(authData[33:37])

	return &Credential{
		ID:              credentialID,
		PublicKey:       publicKey,
		AttestationType: "none", // Simplified
		Authenticator: Authenticator{
			AAGUID:       aaguid,
			SignCount:    signCount,
			CloneWarning: false,
		},
	}, nil
}
