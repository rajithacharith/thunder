package protocol

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// ParseCredentialCreationResponseBytes parses the JSON response from the client.
func ParseCredentialCreationResponseBytes(data []byte) (*ParsedCredentialCreationData, error) {
	// 1. Unmarshal into map or struct to get base64 encoded strings
	var rawResponse struct {
		ID       string         `json:"id"`
		RawID    string         `json:"rawId"`
		Type     CredentialType `json:"type"`
		Response struct {
			ClientDataJSON    string `json:"clientDataJSON"`
			AttestationObject string `json:"attestationObject"`
		} `json:"response"`
	}

	if err := json.Unmarshal(data, &rawResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	// 2. Decode base64 strings
	// WebAuthn spec says RawID is usually base64url encoded in JSON.
	rawIDBytes, err := decodeBase64(rawResponse.RawID)
	if err != nil {
		// Just in case it's not base64, use the ID string as bytes if rawId decoding fails?
		// No, standard is base64url.
		// However, webauthn_service.go passes `credentialID` (string) as `rawId`.
		// If credentialID is already base64, then it works.
		// If it's hex or something else, we might have issues.
		// But let's assume standard base64url.
		return nil, fmt.Errorf("failed to decode rawId: %w", err)
	}

	clientDataBytes, err := decodeBase64(rawResponse.Response.ClientDataJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to decode clientDataJSON: %w", err)
	}

	attestationObjectBytes, err := decodeBase64(rawResponse.Response.AttestationObject)
	if err != nil {
		return nil, fmt.Errorf("failed to decode attestationObject: %w", err)
	}

	// 3. Parse ClientDataJSON
	var clientData CollectedClientData
	if err := json.Unmarshal(clientDataBytes, &clientData); err != nil {
		return nil, fmt.Errorf("failed to parse clientDataJSON: %w", err)
	}

	// 4. Parse AttestationObject (CBOR)
	var attestationObject AttestationObject
	if err := UnmarshalAttestationObject(attestationObjectBytes, &attestationObject); err != nil {
		return nil, fmt.Errorf("failed to parse attestationObject: %w", err)
	}

	// 5. Return ParsedCredentialCreationData
	return &ParsedCredentialCreationData{
		ID:    rawResponse.ID,
		RawID: rawIDBytes,
		Type:  rawResponse.Type,
		Response: ParsedCreationResponse{
			ClientDataJSON:          clientDataBytes,
			AttestationObject:       attestationObjectBytes,
			CollectedClientData:     clientData,
			AttestationObjectParsed: attestationObject,
		},
	}, nil
}

func decodeBase64(s string) ([]byte, error) {
	// Try RawURLEncoding first
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err == nil {
		return b, nil
	}
	// Try StdEncoding
	b, err = base64.StdEncoding.DecodeString(s)
	if err == nil {
		return b, nil
	}
	// Try URLEncoding (with padding)
	return base64.URLEncoding.DecodeString(s)
}
