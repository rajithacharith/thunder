/*
 * Copyright (c) 2026, WSO2 LLC. (https://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package consent

// ConsentPromptData holds the structured data needed to render a consent prompt for the user.
// It contains all purposes whose consent is required, with their elements grouped under each purpose.
type ConsentPromptData struct {
	// Purposes is the list of consent purposes that require user consent, along with their elements
	Purposes []ConsentPurposePrompt `json:"purposes"`
	// SessionToken is the signed JWT token that encapsulates the consent session data
	SessionToken string `json:"session_token,omitempty"`
}

// ConsentPurposePrompt holds a single consent purpose's elements that need user consent.
type ConsentPurposePrompt struct {
	// PurposeName is the name of the consent purpose (e.g. "app:my_app:attrs")
	PurposeName string `json:"purpose_name"`
	// PurposeID is the unique identifier of the consent purpose
	PurposeID string `json:"purpose_id"`
	// Description is a human-readable description of the consent purpose
	Description string `json:"description,omitempty"`
	// Essential is the list of mandatory attribute names that require user consent
	Essential []string `json:"essential"`
	// Optional is the list of optional attribute names for which the user can choose
	Optional []string `json:"optional"`
}

// ConsentDecisions holds the user's consent decisions.
type ConsentDecisions struct {
	// Purposes contains the per-purpose element approval decisions
	Purposes []PurposeDecision `json:"purposes"`
}

// PurposeDecision holds the consent decisions for a single purpose
type PurposeDecision struct {
	// PurposeName is the name of the consent purpose
	PurposeName string `json:"purpose_name"`
	// Approved indicates whether the user approved this purpose
	Approved bool `json:"approved"`
	// Elements contains the per-element approval decisions
	Elements []ElementDecision `json:"elements"`
}

// ElementDecision holds the approval decision for a single consent element
type ElementDecision struct {
	// Name is the name of the consent element
	Name string `json:"name"`
	// Approved indicates whether the user approved sharing this element
	Approved bool `json:"approved"`
}

// consentSessionData holds the consent session state that is signed into a JWT token.
// It captures the purposes and their elements from the resolve step so that the record step
// can verify that the user's decisions match exactly what was prompted.
type consentSessionData struct {
	// Purposes holds the per-purpose element information from the resolve step
	Purposes []consentSessionPurpose `json:"purposes"`
}

// consentSessionPurpose represents a single purpose's elements within the consent session.
type consentSessionPurpose struct {
	// PurposeName is the unique name of the consent purpose
	PurposeName string `json:"purpose_name"`
	// Essential holds the names of mandatory elements for this purpose
	Essential []string `json:"essential"`
	// Optional holds the names of optional elements for this purpose
	Optional []string `json:"optional"`
}
