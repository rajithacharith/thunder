package protocol

const (
	// Authenticator Attachment
	Platform      AuthenticatorAttachment = "platform"
	CrossPlatform AuthenticatorAttachment = "cross-platform"

	// User Verification Requirement
	VerificationRequired    UserVerificationRequirement = "required"
	VerificationPreferred   UserVerificationRequirement = "preferred"
	VerificationDiscouraged UserVerificationRequirement = "discouraged"

	// Resident Key Requirement
	ResidentKeyRequired    ResidentKeyRequirement = "required"
	ResidentKeyPreferred   ResidentKeyRequirement = "preferred"
	ResidentKeyDiscouraged ResidentKeyRequirement = "discouraged"

	// Aliases to match go-webauthn library naming
	ResidentKeyRequirementRequired = ResidentKeyRequired

	// Conveyance Preference
	PreferNoAttestation ConveyancePreference = "none"
	PreferIndirect      ConveyancePreference = "indirect"
	PreferDirect        ConveyancePreference = "direct"
	PreferEnterprise    ConveyancePreference = "enterprise"

	// Credential Type
	PublicKeyCredentialType CredentialType = "public-key"
)
