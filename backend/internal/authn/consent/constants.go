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

const (
	// consentSessionTokenAudience is the JWT audience for consent session tokens
	consentSessionTokenAudience = "consent-svc"

	// consentSessionTokenValidityPeriod is the validity period of consent session tokens in seconds
	consentSessionTokenValidityPeriod = int64(3600)

	// consentSessionClaimKey is the JWT claim key for consent session data
	consentSessionClaimKey = "consent_session"
)
