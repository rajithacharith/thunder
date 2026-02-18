// Package webauthn implements the WebAuthn API.
package webauthn

// Config is the configuration for the WebAuthn service.
type Config struct {
	RPDisplayName string
	RPID          string
	RPOrigins     []string
}
