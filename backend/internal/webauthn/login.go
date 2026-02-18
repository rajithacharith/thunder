package webauthn

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/asgardeo/thunder/internal/webauthn/protocol"
)

// BeginLogin initiates the authentication ceremony.
func (w *WebAuthn) BeginLogin(user User, opts ...LoginOption) (*protocol.CredentialAssertion, *SessionData, error) {
	// Generate challenge
	challenge := make([]byte, 32)
	if _, err := rand.Read(challenge); err != nil {
		return nil, nil, fmt.Errorf("failed to generate challenge: %w", err)
	}

	assertion := &protocol.CredentialAssertion{
		Response: protocol.AssertionResponse{
			Challenge:      challenge,
			Timeout:        60000, // 60 seconds
			RelyingPartyID: w.Config.RPID,
			AllowedCredentials: []protocol.CredentialDescriptor{},
			UserVerification: protocol.VerificationPreferred,
		},
	}

	// Populate allowed credentials from user
	for _, cred := range user.WebAuthnCredentials() {
		assertion.Response.AllowedCredentials = append(assertion.Response.AllowedCredentials, protocol.CredentialDescriptor{
			Type: protocol.PublicKeyCredentialType,
			ID:   cred.ID,
		})
	}

	// Apply options
	for _, opt := range opts {
		opt(assertion)
	}

	// Create session data
	session := &SessionData{
		Challenge:        base64.RawURLEncoding.EncodeToString(challenge),
		UserID:           user.WebAuthnID(),
		UserVerification: assertion.Response.UserVerification,
	}

	return assertion, session, nil
}

// BeginDiscoverableLogin initiates usernameless authentication ceremony.
func (w *WebAuthn) BeginDiscoverableLogin(opts ...LoginOption) (*protocol.CredentialAssertion, *SessionData, error) {
	// Generate challenge
	challenge := make([]byte, 32)
	if _, err := rand.Read(challenge); err != nil {
		return nil, nil, fmt.Errorf("failed to generate challenge: %w", err)
	}

	assertion := &protocol.CredentialAssertion{
		Response: protocol.AssertionResponse{
			Challenge:      challenge,
			Timeout:        60000,
			RelyingPartyID: w.Config.RPID,
			UserVerification: protocol.VerificationPreferred,
		},
	}

	// Apply options
	for _, opt := range opts {
		opt(assertion)
	}

	// Create session data - no UserID for discoverable flow
	session := &SessionData{
		Challenge:        base64.RawURLEncoding.EncodeToString(challenge),
		UserVerification: assertion.Response.UserVerification,
	}

	return assertion, session, nil
}

// ValidateLogin validates the authentication response.
func (w *WebAuthn) ValidateLogin(user User, session SessionData, response *protocol.ParsedCredentialAssertionData) (*Credential, error) {
	// 1. Verify Challenge
	if response.Response.CollectedClientData.Challenge != session.Challenge {
		return nil, fmt.Errorf("challenge mismatch")
	}

	// 2. Verify Origin
	originValid := false
	for _, origin := range w.Config.RPOrigins {
		if response.Response.CollectedClientData.Origin == origin {
			originValid = true
			break
		}
	}
	if !originValid {
		return nil, fmt.Errorf("origin mismatch")
	}

	// 3. Verify RPID Hash
	authData := response.Raw.AssertionResponse.AuthenticatorData
	if len(authData) < 37 {
		return nil, fmt.Errorf("authenticator data too short")
	}

	rpIDHash := sha256.Sum256([]byte(w.Config.RPID))
	if !bytes.Equal(rpIDHash[:], authData[:32]) {
		return nil, fmt.Errorf("RPID hash mismatch")
	}

	// 4. Verify Flags
	flags := protocol.AuthenticatorFlags(authData[32])
	if flags&protocol.FlagUserPresent == 0 {
		return nil, fmt.Errorf("user present flag not set")
	}

	if session.UserVerification == protocol.VerificationRequired && flags&protocol.FlagUserVerified == 0 {
		return nil, fmt.Errorf("user verification required but not verified")
	}

	// 5. Find Credential
	var credential *Credential
	for _, c := range user.WebAuthnCredentials() {
		if bytes.Equal(c.ID, response.ParsedPublicKeyCredential.RawID) {
			credential = &c
			break
		}
	}
	if credential == nil {
		return nil, fmt.Errorf("credential not found")
	}

	// 6. Verify Signature
	clientDataJSON := response.Raw.AssertionResponse.AuthenticatorResponse.ClientDataJSON
	clientDataHash := sha256.Sum256(clientDataJSON)

	signatureBase := append(authData, clientDataHash[:]...)

	// Parse Public Key
	pubKey, alg, err := protocol.ParseCOSEKey(credential.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	// Verify Signature
	signature := response.Raw.AssertionResponse.Signature

	switch k := pubKey.(type) {
	case *ecdsa.PublicKey:
		// ES256 (-7) uses SHA-256
		if alg == 0 || alg == -7 {
			hash := sha256.Sum256(signatureBase)
			if !ecdsa.VerifyASN1(k, hash[:], signature) {
				return nil, fmt.Errorf("invalid signature")
			}
		} else {
			return nil, fmt.Errorf("unsupported algorithm: %d", alg)
		}

	case *rsa.PublicKey:
		// RS256 (-257) uses SHA-256 and PKCS1v15
		if alg == 0 || alg == -257 {
			hash := sha256.Sum256(signatureBase)
			if err := rsa.VerifyPKCS1v15(k, crypto.SHA256, hash[:], signature); err != nil {
				return nil, fmt.Errorf("invalid signature: %w", err)
			}
		} else {
			return nil, fmt.Errorf("unsupported algorithm: %d", alg)
		}

	default:
		return nil, fmt.Errorf("unsupported public key type")
	}

	// 4. Verify Sign Count
	authDataParsed := response.Response.AuthenticatorData
	if authDataParsed.Counter > 0 && authDataParsed.Counter <= credential.Authenticator.SignCount {
		// Clone warning logic
		credential.Authenticator.CloneWarning = true
		// Fail or warn? Library usually fails.
		return nil, fmt.Errorf("sign count error")
	}
	credential.Authenticator.SignCount = authDataParsed.Counter

	return credential, nil
}

// ValidatePasskeyLogin validates discoverable credential authentication.
func (w *WebAuthn) ValidatePasskeyLogin(
	userHandler func(rawID, userHandle []byte) (User, error),
	session SessionData,
	response *protocol.ParsedCredentialAssertionData,
) (User, *Credential, error) {
	// 1. Verify Challenge
	if response.Response.CollectedClientData.Challenge != session.Challenge {
		return nil, nil, fmt.Errorf("challenge mismatch")
	}

	// 2. Resolve User using userHandler
	user, err := userHandler(response.ParsedPublicKeyCredential.RawID, response.Response.UserHandle)
	if err != nil {
		return nil, nil, fmt.Errorf("user not found: %w", err)
	}

	// 3. Validate Login for the resolved user
	credential, err := w.ValidateLogin(user, session, response)
	if err != nil {
		return nil, nil, err
	}

	return user, credential, nil
}
