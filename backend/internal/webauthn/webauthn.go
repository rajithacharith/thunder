package webauthn

// WebAuthn is the main entry point for the WebAuthn service.
type WebAuthn struct {
	Config *Config
}

// New creates a new WebAuthn service.
func New(config *Config) (*WebAuthn, error) {
	// Validate config if necessary
	if config.RPID == "" {
		return nil, &ErrorInvalidConfig{"RPID is required"}
	}
	return &WebAuthn{Config: config}, nil
}

type ErrorInvalidConfig struct {
	Message string
}

func (e *ErrorInvalidConfig) Error() string {
	return e.Message
}
