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

package design

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/asgardeo/thunder/tests/integration/testutils"
	"github.com/stretchr/testify/suite"
)

const (
	usagesBasePath = "/design/usages"
)

// UsagesAPITestSuite tests the GET /design/usages endpoint end-to-end.
type UsagesAPITestSuite struct {
	suite.Suite
	client *http.Client

	// shared resources created in SetupSuite and cleaned up in TearDownSuite
	themeID  string
	layoutID string
	flowID   string
	ouID     string

	// applications that reference the shared resources
	appWithThemeID  string
	appWithLayoutID string
	appWithFlowID   string
	appWithAllID    string
	appWithNoneID   string
}

func TestUsagesAPITestSuite(t *testing.T) {
	suite.Run(t, new(UsagesAPITestSuite))
}

// -------------------------------------------------------------------------
// Suite lifecycle
// -------------------------------------------------------------------------

func (suite *UsagesAPITestSuite) SetupSuite() {
	suite.client = testutils.GetHTTPClient()

	var err error

	// Create an OU to own the applications
	suite.ouID, err = testutils.CreateOrganizationUnit(testutils.OrganizationUnit{
		Handle:      "test-usages-ou",
		Name:        "Test OU for Usages",
		Description: "OU created for design usages integration tests",
	})
	suite.Require().NoError(err, "SetupSuite: failed to create OU")

	// Create a theme
	suite.themeID, err = suite.createTheme()
	suite.Require().NoError(err, "SetupSuite: failed to create theme")

	// Create a layout
	suite.layoutID, err = suite.createLayout()
	suite.Require().NoError(err, "SetupSuite: failed to create layout")

	// Retrieve the default auth flow
	suite.flowID, err = testutils.GetFlowIDByHandle("default-basic-flow", "AUTHENTICATION")
	suite.Require().NoError(err, "SetupSuite: failed to get default auth flow ID")

	// Create test applications
	suite.appWithThemeID, err = suite.createApplication(applicationPayload{
		OUID:    suite.ouID,
		Name:    "Usages Test App - Theme Only",
		ThemeID: suite.themeID,
	})
	suite.Require().NoError(err, "SetupSuite: failed to create app with theme")

	suite.appWithLayoutID, err = suite.createApplication(applicationPayload{
		OUID:     suite.ouID,
		Name:     "Usages Test App - Layout Only",
		LayoutID: suite.layoutID,
	})
	suite.Require().NoError(err, "SetupSuite: failed to create app with layout")

	suite.appWithFlowID, err = suite.createApplication(applicationPayload{
		OUID:       suite.ouID,
		Name:       "Usages Test App - Flow Only",
		AuthFlowID: suite.flowID,
	})
	suite.Require().NoError(err, "SetupSuite: failed to create app with flow")

	suite.appWithAllID, err = suite.createApplication(applicationPayload{
		OUID:       suite.ouID,
		Name:       "Usages Test App - All Resources",
		ThemeID:    suite.themeID,
		LayoutID:   suite.layoutID,
		AuthFlowID: suite.flowID,
	})
	suite.Require().NoError(err, "SetupSuite: failed to create app with all resources")

	suite.appWithNoneID, err = suite.createApplication(applicationPayload{
		OUID: suite.ouID,
		Name: "Usages Test App - No Design Resources",
	})
	suite.Require().NoError(err, "SetupSuite: failed to create app with no design resources")
}

func (suite *UsagesAPITestSuite) TearDownSuite() {
	for _, id := range []string{
		suite.appWithThemeID,
		suite.appWithLayoutID,
		suite.appWithFlowID,
		suite.appWithAllID,
		suite.appWithNoneID,
	} {
		if id != "" {
			_ = suite.deleteApplication(id)
		}
	}

	if suite.themeID != "" {
		_ = suite.deleteTheme(suite.themeID)
	}
	if suite.layoutID != "" {
		_ = suite.deleteLayout(suite.layoutID)
	}
	if suite.ouID != "" {
		_ = testutils.DeleteOrganizationUnit(suite.ouID)
	}
}

// -------------------------------------------------------------------------
// Validation error tests (no server-side resources required)
// -------------------------------------------------------------------------

// TestUsages_MissingType verifies that omitting the ?type= parameter returns 400 DSU-1001.
func (suite *UsagesAPITestSuite) TestUsages_MissingType() {
	url := fmt.Sprintf("%s%s?id=00000000-0000-0000-0000-000000000000", testServerURL, usagesBasePath)
	statusCode, errResp := suite.doRawGet(url)
	suite.Equal(http.StatusBadRequest, statusCode)
	suite.Require().NotNil(errResp)
	suite.Equal("DSU-1001", errResp.Code)
}

// TestUsages_MissingID verifies that omitting the ?id= parameter returns 400 DSU-1002.
func (suite *UsagesAPITestSuite) TestUsages_MissingID() {
	url := fmt.Sprintf("%s%s?type=THEME", testServerURL, usagesBasePath)
	statusCode, errResp := suite.doRawGet(url)
	suite.Equal(http.StatusBadRequest, statusCode)
	suite.Require().NotNil(errResp)
	suite.Equal("DSU-1002", errResp.Code)
}

// TestUsages_UnsupportedType verifies that an unrecognised type returns 400 DSU-1003.
func (suite *UsagesAPITestSuite) TestUsages_UnsupportedType() {
	url := fmt.Sprintf("%s%s?type=UNKNOWN&id=00000000-0000-0000-0000-000000000000", testServerURL, usagesBasePath)
	statusCode, errResp := suite.doRawGet(url)
	suite.Equal(http.StatusBadRequest, statusCode)
	suite.Require().NotNil(errResp)
	suite.Equal("DSU-1003", errResp.Code)
}

// TestUsages_TypeCaseInsensitive verifies that lowercase type names are accepted
// (the handler normalises with strings.ToUpper, so "theme" == "THEME").
func (suite *UsagesAPITestSuite) TestUsages_TypeCaseInsensitive() {
	url := fmt.Sprintf("%s%s?type=theme&id=00000000-0000-0000-0000-000000000001", testServerURL, usagesBasePath)
	// A non-existent theme ID → 404, not 400
	statusCode, _ := suite.doRawGet(url)
	suite.Equal(http.StatusNotFound, statusCode)
}

// -------------------------------------------------------------------------
// 404 not-found tests
// -------------------------------------------------------------------------

// TestUsages_ThemeNotFound verifies that a non-existent theme ID returns 404 DSU-1004.
func (suite *UsagesAPITestSuite) TestUsages_ThemeNotFound() {
	url := fmt.Sprintf("%s%s?type=THEME&id=00000000-0000-0000-0000-000000000000", testServerURL, usagesBasePath)
	statusCode, errResp := suite.doRawGet(url)
	suite.Equal(http.StatusNotFound, statusCode)
	suite.Require().NotNil(errResp)
	suite.Equal("DSU-1004", errResp.Code)
}

// TestUsages_LayoutNotFound verifies that a non-existent layout ID returns 404 DSU-1004.
func (suite *UsagesAPITestSuite) TestUsages_LayoutNotFound() {
	url := fmt.Sprintf("%s%s?type=LAYOUT&id=00000000-0000-0000-0000-000000000000", testServerURL, usagesBasePath)
	statusCode, errResp := suite.doRawGet(url)
	suite.Equal(http.StatusNotFound, statusCode)
	suite.Require().NotNil(errResp)
	suite.Equal("DSU-1004", errResp.Code)
}

// TestUsages_FlowNotFound verifies that a non-existent flow ID returns 404 DSU-1004.
func (suite *UsagesAPITestSuite) TestUsages_FlowNotFound() {
	url := fmt.Sprintf("%s%s?type=FLOW&id=00000000-0000-0000-0000-000000000000", testServerURL, usagesBasePath)
	statusCode, errResp := suite.doRawGet(url)
	suite.Equal(http.StatusNotFound, statusCode)
	suite.Require().NotNil(errResp)
	suite.Equal("DSU-1004", errResp.Code)
}

// -------------------------------------------------------------------------
// Theme usages tests
// -------------------------------------------------------------------------

// TestUsages_Theme_ReturnsAppsUsingTheme verifies that querying by theme ID returns the
// applications that were assigned that theme, and excludes apps that were not.
func (suite *UsagesAPITestSuite) TestUsages_Theme_ReturnsAppsUsingTheme() {
	resp, err := suite.getUsages("THEME", suite.themeID)
	suite.Require().NoError(err)

	suite.GreaterOrEqual(resp.TotalResults, 2, "expected at least 2 apps using the theme")
	suite.Equal(resp.TotalResults, resp.Count)
	suite.Len(resp.Applications, resp.Count)

	ids := appIDs(resp.Applications)
	suite.Contains(ids, suite.appWithThemeID, "app with theme only should be included")
	suite.Contains(ids, suite.appWithAllID, "app with all resources should be included")
	suite.NotContains(ids, suite.appWithLayoutID, "app with layout only should be excluded")
	suite.NotContains(ids, suite.appWithFlowID, "app with flow only should be excluded")
	suite.NotContains(ids, suite.appWithNoneID, "app with no design resources should be excluded")
}

// TestUsages_Theme_ApplicationRefShape verifies that each ApplicationRef has the required fields.
func (suite *UsagesAPITestSuite) TestUsages_Theme_ApplicationRefShape() {
	resp, err := suite.getUsages("THEME", suite.themeID)
	suite.Require().NoError(err)
	suite.Require().NotEmpty(resp.Applications)

	for _, ref := range resp.Applications {
		suite.NotEmpty(ref.ID, "application ref must have an ID")
		suite.NotEmpty(ref.Name, "application ref must have a name")
		// clientId may be absent for apps without an OAuth profile; no assertion needed
	}
}

// TestUsages_Theme_NoAppsUsing verifies that a freshly-created theme with no apps returns
// an empty list with 0 totalResults.
func (suite *UsagesAPITestSuite) TestUsages_Theme_NoAppsUsing() {
	emptyThemeID, err := suite.createTheme()
	suite.Require().NoError(err)
	defer func() { _ = suite.deleteTheme(emptyThemeID) }()

	resp, err := suite.getUsages("THEME", emptyThemeID)
	suite.Require().NoError(err)
	suite.Equal(0, resp.TotalResults)
	suite.Equal(0, resp.Count)
	suite.Empty(resp.Applications)
}

// -------------------------------------------------------------------------
// Layout usages tests
// -------------------------------------------------------------------------

// TestUsages_Layout_ReturnsAppsUsingLayout verifies that querying by layout ID returns
// the applications assigned to that layout.
func (suite *UsagesAPITestSuite) TestUsages_Layout_ReturnsAppsUsingLayout() {
	resp, err := suite.getUsages("LAYOUT", suite.layoutID)
	suite.Require().NoError(err)

	suite.GreaterOrEqual(resp.TotalResults, 2, "expected at least 2 apps using the layout")
	ids := appIDs(resp.Applications)
	suite.Contains(ids, suite.appWithLayoutID)
	suite.Contains(ids, suite.appWithAllID)
	suite.NotContains(ids, suite.appWithThemeID)
	suite.NotContains(ids, suite.appWithNoneID)
}

// TestUsages_Layout_NoAppsUsing verifies that a freshly-created layout with no apps returns
// an empty list.
func (suite *UsagesAPITestSuite) TestUsages_Layout_NoAppsUsing() {
	emptyLayoutID, err := suite.createLayout()
	suite.Require().NoError(err)
	defer func() { _ = suite.deleteLayout(emptyLayoutID) }()

	resp, err := suite.getUsages("LAYOUT", emptyLayoutID)
	suite.Require().NoError(err)
	suite.Equal(0, resp.TotalResults)
	suite.Empty(resp.Applications)
}

// -------------------------------------------------------------------------
// Flow usages tests
// -------------------------------------------------------------------------

// TestUsages_Flow_ReturnsAppsUsingFlow verifies that querying by flow ID returns the
// applications that reference that flow as their auth flow.
func (suite *UsagesAPITestSuite) TestUsages_Flow_ReturnsAppsUsingFlow() {
	resp, err := suite.getUsages("FLOW", suite.flowID)
	suite.Require().NoError(err)

	// The default-basic-flow may be used by other pre-existing apps too, so just assert
	// that our explicitly configured apps appear in the list.
	ids := appIDs(resp.Applications)
	suite.Contains(ids, suite.appWithFlowID)
	suite.Contains(ids, suite.appWithAllID)
	suite.NotContains(ids, suite.appWithNoneID)
}

// -------------------------------------------------------------------------
// Mutation: assign / un-assign theme and verify usages update
// -------------------------------------------------------------------------

// TestUsages_Theme_AfterUpdate verifies that updating an application to remove its theme
// removes it from the usages response, and that re-assigning it brings it back.
func (suite *UsagesAPITestSuite) TestUsages_Theme_AfterUpdate() {
	// Create an isolated app that starts with the shared theme
	appID, err := suite.createApplication(applicationPayload{
		OUID:    suite.ouID,
		Name:    "Usages Mutation Test App",
		ThemeID: suite.themeID,
	})
	suite.Require().NoError(err)
	defer func() { _ = suite.deleteApplication(appID) }()

	// Confirm it is listed
	resp, err := suite.getUsages("THEME", suite.themeID)
	suite.Require().NoError(err)
	suite.Contains(appIDs(resp.Applications), appID, "app should be listed before theme removal")

	// Remove the theme from the application
	err = suite.updateApplicationTheme(appID, "")
	suite.Require().NoError(err)

	// Confirm it is no longer listed
	resp, err = suite.getUsages("THEME", suite.themeID)
	suite.Require().NoError(err)
	suite.NotContains(appIDs(resp.Applications), appID, "app should not be listed after theme removal")

	// Re-assign the theme
	err = suite.updateApplicationTheme(appID, suite.themeID)
	suite.Require().NoError(err)

	// Confirm it is listed again
	resp, err = suite.getUsages("THEME", suite.themeID)
	suite.Require().NoError(err)
	suite.Contains(appIDs(resp.Applications), appID, "app should be listed again after theme re-assignment")
}

// -------------------------------------------------------------------------
// OPTIONS pre-flight
// -------------------------------------------------------------------------

// TestUsages_OPTIONS verifies that the OPTIONS pre-flight request returns 204.
func (suite *UsagesAPITestSuite) TestUsages_OPTIONS() {
	url := fmt.Sprintf("%s%s", testServerURL, usagesBasePath)
	req, err := http.NewRequest("OPTIONS", url, nil)
	suite.Require().NoError(err)

	resp, err := suite.client.Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusNoContent, resp.StatusCode)
}

// -------------------------------------------------------------------------
// Helpers
// -------------------------------------------------------------------------

// getUsages calls GET /design/usages?type=<t>&id=<id> and decodes the success body.
func (suite *UsagesAPITestSuite) getUsages(resourceType, id string) (*DesignUsagesResponse, error) {
	url := fmt.Sprintf("%s%s?type=%s&id=%s", testServerURL, usagesBasePath, resourceType, id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := suite.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.Unmarshal(bodyBytes, &errResp); err == nil {
			return nil, fmt.Errorf("expected status 200, got %d. Code: %s, Message: %s",
				resp.StatusCode, errResp.Code, errResp.Message.DefaultValue)
		}
		return nil, fmt.Errorf("expected status 200, got %d. Response: %s", resp.StatusCode, string(bodyBytes))
	}

	var usagesResp DesignUsagesResponse
	if err := json.Unmarshal(bodyBytes, &usagesResp); err != nil {
		return nil, fmt.Errorf("failed to parse response body: %w. Response: %s", err, string(bodyBytes))
	}
	return &usagesResp, nil
}

// doRawGet sends a GET request and returns the status code along with a decoded error response
// (nil if the body is not an error response).
func (suite *UsagesAPITestSuite) doRawGet(url string) (int, *ErrorResponse) {
	req, err := http.NewRequest("GET", url, nil)
	suite.Require().NoError(err)

	resp, err := suite.client.Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	suite.Require().NoError(err)

	var errResp ErrorResponse
	if err := json.Unmarshal(bodyBytes, &errResp); err == nil && errResp.Code != "" {
		return resp.StatusCode, &errResp
	}
	return resp.StatusCode, nil
}

// applicationPayload holds the fields needed to create a minimal test application.
type applicationPayload struct {
	OUID       string `json:"ouId,omitempty"`
	Name       string `json:"name"`
	ThemeID    string `json:"themeId,omitempty"`
	LayoutID   string `json:"layoutId,omitempty"`
	AuthFlowID string `json:"authFlowId,omitempty"`
}

// createApplication creates a minimal application and returns its ID.
func (suite *UsagesAPITestSuite) createApplication(app applicationPayload) (string, error) {
	payload, err := json.Marshal(app)
	if err != nil {
		return "", fmt.Errorf("failed to marshal application: %w", err)
	}

	req, err := http.NewRequest("POST", testServerURL+"/applications", bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("expected status 201, got %d. Response: %s", resp.StatusCode, string(bodyBytes))
	}

	var created map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &created); err != nil {
		return "", fmt.Errorf("failed to parse response body: %w", err)
	}

	id, ok := created["id"].(string)
	if !ok || id == "" {
		return "", fmt.Errorf("response does not contain id")
	}
	return id, nil
}

// deleteApplication deletes an application by ID.
func (suite *UsagesAPITestSuite) deleteApplication(appID string) error {
	req, err := http.NewRequest("DELETE", testServerURL+"/applications/"+appID, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete request: %w", err)
	}

	resp, err := suite.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send delete request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotFound {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("expected status 204 or 404, got %d. Response: %s", resp.StatusCode, string(bodyBytes))
	}
	return nil
}

// updateApplicationTheme patches an application's themeId field. Passing an empty themeID
// clears the field.
func (suite *UsagesAPITestSuite) updateApplicationTheme(appID, themeID string) error {
	// Fetch current application state first
	req, err := http.NewRequest("GET", testServerURL+"/applications/"+appID, nil)
	if err != nil {
		return fmt.Errorf("failed to create GET request: %w", err)
	}

	resp, err := suite.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send GET request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read GET response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GET application returned %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var appMap map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &appMap); err != nil {
		return fmt.Errorf("failed to decode application: %w", err)
	}

	// Update the themeId field (or remove it)
	if themeID == "" {
		delete(appMap, "themeId")
	} else {
		appMap["themeId"] = themeID
	}

	// Remove read-only / response-only fields that the server would reject
	delete(appMap, "clientSecret")

	payload, err := json.Marshal(appMap)
	if err != nil {
		return fmt.Errorf("failed to marshal updated application: %w", err)
	}

	updateReq, err := http.NewRequest("PUT", testServerURL+"/applications/"+appID, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("failed to create PUT request: %w", err)
	}
	updateReq.Header.Set("Content-Type", "application/json")

	updateResp, err := suite.client.Do(updateReq)
	if err != nil {
		return fmt.Errorf("failed to send PUT request: %w", err)
	}
	defer updateResp.Body.Close()

	if updateResp.StatusCode != http.StatusOK {
		updateBody, _ := io.ReadAll(updateResp.Body)
		return fmt.Errorf("PUT application returned %d: %s", updateResp.StatusCode, string(updateBody))
	}
	return nil
}

// createTheme creates a minimal theme and returns its ID.
func (suite *UsagesAPITestSuite) createTheme() (string, error) {
	payload, err := json.Marshal(map[string]interface{}{
		"handle":      fmt.Sprintf("usages-test-theme-%d", time.Now().UnixNano()),
		"displayName": "Usages Test Theme",
		"theme":       json.RawMessage(`{"colorSchemes":{"light":{},"dark":{}}}`),
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal theme: %w", err)
	}

	req, err := http.NewRequest("POST", testServerURL+"/design/themes", bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("expected 201, got %d. Response: %s", resp.StatusCode, string(bodyBytes))
	}

	var created map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &created); err != nil {
		return "", fmt.Errorf("failed to decode theme response: %w", err)
	}

	id, ok := created["id"].(string)
	if !ok || id == "" {
		return "", fmt.Errorf("theme response missing id")
	}
	return id, nil
}

// deleteTheme deletes a theme by ID.
func (suite *UsagesAPITestSuite) deleteTheme(id string) error {
	req, err := http.NewRequest("DELETE", testServerURL+"/design/themes/"+id, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete request: %w", err)
	}

	resp, err := suite.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send delete request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotFound {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("expected 204 or 404, got %d. Response: %s", resp.StatusCode, string(bodyBytes))
	}
	return nil
}

// createLayout creates a minimal layout and returns its ID.
func (suite *UsagesAPITestSuite) createLayout() (string, error) {
	payload, err := json.Marshal(map[string]interface{}{
		"handle":      fmt.Sprintf("usages-test-layout-%d", time.Now().UnixNano()),
		"displayName": "Usages Test Layout",
		"layout":      json.RawMessage(`{"components":[]}`),
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal layout: %w", err)
	}

	req, err := http.NewRequest("POST", testServerURL+"/design/layouts", bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("expected 201, got %d. Response: %s", resp.StatusCode, string(bodyBytes))
	}

	var created map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &created); err != nil {
		return "", fmt.Errorf("failed to decode layout response: %w", err)
	}

	id, ok := created["id"].(string)
	if !ok || id == "" {
		return "", fmt.Errorf("layout response missing id")
	}
	return id, nil
}

// deleteLayout deletes a layout by ID.
func (suite *UsagesAPITestSuite) deleteLayout(id string) error {
	req, err := http.NewRequest("DELETE", testServerURL+"/design/layouts/"+id, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete request: %w", err)
	}

	resp, err := suite.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send delete request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotFound {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("expected 204 or 404, got %d. Response: %s", resp.StatusCode, string(bodyBytes))
	}
	return nil
}

// appIDs extracts just the ID strings from a slice of ApplicationRef.
func appIDs(refs []ApplicationRef) []string {
	ids := make([]string, len(refs))
	for i, r := range refs {
		ids[i] = r.ID
	}
	return ids
}
