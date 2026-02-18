package protocol

const (
	// Platform represents a platform authenticator.
	Platform AuthenticatorAttachment = "platform"
	// CrossPlatform represents a roaming authenticator.
	CrossPlatform AuthenticatorAttachment = "cross-platform"

	// VerificationRequired means user verification is required.
	VerificationRequired UserVerificationRequirement = "required"
	// VerificationPreferred means user verification is preferred.
	VerificationPreferred UserVerificationRequirement = "preferred"
	// VerificationDiscouraged means user verification is discouraged.
	VerificationDiscouraged UserVerificationRequirement = "discouraged"

	// ResidentKeyRequired means a resident key is required.
	ResidentKeyRequired ResidentKeyRequirement = "required"
	// ResidentKeyPreferred means a resident key is preferred.
	ResidentKeyPreferred ResidentKeyRequirement = "preferred"
	// ResidentKeyDiscouraged means a resident key is discouraged.
	ResidentKeyDiscouraged ResidentKeyRequirement = "discouraged"

	// ResidentKeyRequirementRequired is an alias for ResidentKeyRequired.
	ResidentKeyRequirementRequired = ResidentKeyRequired

	// PreferNoAttestation means no attestation is preferred.
	PreferNoAttestation ConveyancePreference = "none"
	// PreferIndirect means indirect attestation is preferred.
	PreferIndirect ConveyancePreference = "indirect"
	// PreferDirect means direct attestation is preferred.
	PreferDirect ConveyancePreference = "direct"
	// PreferEnterprise means enterprise attestation is preferred.
	PreferEnterprise ConveyancePreference = "enterprise"

	// PublicKeyCredentialType represents the public key credential type.
	PublicKeyCredentialType CredentialType = "public-key"
)
