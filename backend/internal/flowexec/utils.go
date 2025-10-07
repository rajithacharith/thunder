package flowexec

import (
	"fmt"

	"github.com/asgardeo/thunder/internal/executor/attributecollect"
	"github.com/asgardeo/thunder/internal/executor/authassert"
	"github.com/asgardeo/thunder/internal/executor/basicauth"
	"github.com/asgardeo/thunder/internal/executor/githubauth"
	"github.com/asgardeo/thunder/internal/executor/googleauth"
	"github.com/asgardeo/thunder/internal/executor/provision"
	"github.com/asgardeo/thunder/internal/executor/smsauth"
	"github.com/asgardeo/thunder/internal/flow"
	"github.com/asgardeo/thunder/internal/idp"
	"github.com/asgardeo/thunder/internal/system/cmodels"
	"github.com/asgardeo/thunder/internal/system/utils"
)

func GetExecutorByName(execConfig *flow.ExecutorConfig) (flow.ExecutorInterface, error) {
	if execConfig == nil {
		return nil, fmt.Errorf("executor configuration cannot be nil")
	}
	if execConfig.Name == "" {
		return nil, fmt.Errorf("executor name cannot be empty")
	}

	var executor flow.ExecutorInterface
	switch execConfig.Name {
	case "BasicAuthExecutor":
		executor = basicauth.NewBasicAuthExecutor("local", "Local", execConfig.Properties)
	case "SMSOTPAuthExecutor":
		if len(execConfig.Properties) == 0 {
			return nil, fmt.Errorf("properties for SMSOTPAuthExecutor cannot be empty")
		}
		senderID, exists := execConfig.Properties["senderId"]
		if !exists || senderID == "" {
			return nil, fmt.Errorf("senderId property is required for SMSOTPAuthExecutor")
		}
		executor = smsauth.NewSMSOTPAuthExecutor("local", "Local", execConfig.Properties)
	case "GithubOAuthExecutor":
		idp, err := getIDP(execConfig.IdpName)
		if err != nil {
			return nil, fmt.Errorf("error while getting IDP for GithubOAuthExecutor: %w", err)
		}

		clientID, clientSecret, redirectURI, scopes, additionalParams, err := getIDPConfigs(
			idp.Properties, execConfig)
		if err != nil {
			return nil, err
		}

		executor = githubauth.NewGithubOAuthExecutor(idp.ID, idp.Name, execConfig.Properties,
			clientID, clientSecret, redirectURI, scopes, additionalParams)
	case "GoogleOIDCAuthExecutor":
		idp, err := getIDP(execConfig.IdpName)
		if err != nil {
			return nil, fmt.Errorf("error while getting IDP for GoogleOIDCAuthExecutor: %w", err)
		}

		clientID, clientSecret, redirectURI, scopes, additionalParams, err := getIDPConfigs(
			idp.Properties, execConfig)
		if err != nil {
			return nil, err
		}

		executor = googleauth.NewGoogleOIDCAuthExecutor(idp.ID, idp.Name, execConfig.Properties,
			clientID, clientSecret, redirectURI, scopes, additionalParams)
	case "AttributeCollector":
		executor = attributecollect.NewAttributeCollector("attribute-collector", "AttributeCollector",
			execConfig.Properties)
	case "ProvisioningExecutor":
		executor = provision.NewProvisioningExecutor("provisioning-executor", "ProvisioningExecutor",
			execConfig.Properties)
	case "AuthAssertExecutor":
		executor = authassert.NewAuthAssertExecutor("auth-assert-executor", "AuthAssertExecutor",
			execConfig.Properties)
	default:
		return nil, fmt.Errorf("executor with name %s not found", execConfig.Name)
	}

	if executor == nil {
		return nil, fmt.Errorf("executor with name %s could not be created", execConfig.Name)
	}
	return executor, nil
}

// getIDP retrieves the IDP by its name. Returns an error if the IDP does not exist or if the name is empty.
func getIDP(idpName string) (*idp.IDPDTO, error) {
	if idpName == "" {
		return nil, fmt.Errorf("IDP name cannot be empty")
	}

	idpSvc := idp.NewIDPService()
	identityProvider, svcErr := idpSvc.GetIdentityProviderByName(idpName)
	if svcErr != nil {
		if svcErr.Code == idp.ErrorIDPNotFound.Code {
			return nil, fmt.Errorf("IDP with name %s does not exist", idpName)
		}
		return nil, fmt.Errorf("error while getting IDP with the name %s: code: %s, error: %s",
			idpName, svcErr.Code, svcErr.ErrorDescription)
	}
	if identityProvider == nil {
		return nil, fmt.Errorf("IDP with name %s does not exist", idpName)
	}

	return identityProvider, nil
}

// getIDPConfigs retrieves the IDP configurations for a given executor configuration.
func getIDPConfigs(idpProperties []cmodels.Property, execConfig *flow.ExecutorConfig) (string,
	string, string, []string, map[string]string, error) {
	if len(idpProperties) == 0 {
		return "", "", "", nil, nil, fmt.Errorf("IDP properties not found for executor with IDP name %s",
			execConfig.IdpName)
	}
	var clientID, clientSecret, redirectURI, scopesStr string
	additionalParams := map[string]string{}
	for _, prop := range idpProperties {
		value, err := prop.GetValue()
		if err != nil {
			return "", "", "", nil, nil, err
		}

		switch prop.GetName() {
		case "client_id":
			clientID = value
		case "client_secret":
			clientSecret = value
		case "redirect_uri":
			redirectURI = value
		case "scopes":
			scopesStr = value
		default:
			additionalParams[prop.GetName()] = value
		}
	}
	if clientID == "" || clientSecret == "" || redirectURI == "" || scopesStr == "" {
		return "", "", "", nil, nil, fmt.Errorf("missing required properties for executor with IDP name %s",
			execConfig.IdpName)
	}
	scopes := utils.ParseStringArray(scopesStr, ",")

	return clientID, clientSecret, redirectURI, scopes, additionalParams, nil
}
