package webauthn

import "github.com/asgardeo/thunder/internal/webauthn/protocol"

// WithAuthenticatorSelection sets the authenticator selection criteria.
func WithAuthenticatorSelection(selection protocol.AuthenticatorSelection) RegistrationOption {
	return func(opts *protocol.CredentialCreation) {
		opts.Response.AuthenticatorSelection = selection
	}
}

// WithConveyancePreference sets the attestation conveyance preference.
func WithConveyancePreference(preference protocol.ConveyancePreference) RegistrationOption {
	return func(opts *protocol.CredentialCreation) {
		opts.Response.Attestation = preference
	}
}

// WithUserVerification sets the user verification requirement for login.
func WithUserVerification(verification protocol.UserVerificationRequirement) LoginOption {
	return func(opts *protocol.CredentialAssertion) {
		opts.Response.UserVerification = verification
	}
}
