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

package attestation

// appleAppAttestRootPEM is Apple's public "Apple App Attestation Root CA" certificate, published at
// https://www.apple.com/certificateauthority/Apple_App_Attestation_Root_CA.pem. It is a public trust
// anchor, not a secret.
const appleAppAttestRootPEM = `-----BEGIN CERTIFICATE-----
MIICITCCAaegAwIBAgIQC/O+DvHN0uD7jG5yH2IXmDAKBggqhkjOPQQDAzBSMSYw
JAYDVQQDDB1BcHBsZSBBcHAgQXR0ZXN0YXRpb24gUm9vdCBDQTETMBEGA1UECgwK
QXBwbGUgSW5jLjETMBEGA1UECAwKQ2FsaWZvcm5pYTAeFw0yMDAzMTgxODMyNTNa
Fw00NTAzMTUwMDAwMDBaMFIxJjAkBgNVBAMMHUFwcGxlIEFwcCBBdHRlc3RhdGlv
biBSb290IENBMRMwEQYDVQQKDApBcHBsZSBJbmMuMRMwEQYDVQQIDApDYWxpZm9y
bmlhMHYwEAYHKoZIzj0CAQYFK4EEACIDYgAERTHhmLW07ATaFQIEVwTtT4dyctdh
NbJhFs/Ii2FdCgAHGbpphY3+d8qjuDngIN3WVhQUBHAoMeQ/cLiP1sOUtgjqK9au
Yen1mMEvRq9Sk3Jm5X8U62H+xTD3FE9TgS41o0IwQDAPBgNVHRMBAf8EBTADAQH/
MB0GA1UdDgQWBBSskRBTM72+aEH/pwyp5frq5eWKoTAOBgNVHQ8BAf8EBAMCAQYw
CgYIKoZIzj0EAwMDaAAwZQIwQgFGnByvsiVbpTKwSga0kP0e8EeDS4+sQmTvb7vn
53O5+FRXgeLhpJ06ysC5PrOyAjEAp5U4xDgEgllF7En3VcE3iexZZtKeYnpqtijV
oyFraWVIyd/dganmrduC1bmTBGwD
-----END CERTIFICATE-----
`

// appAttestFormat is the only attestation statement format Apple App Attest produces.
const appAttestFormat = "apple-appattest"

// authData field layout, per the WebAuthn authenticator data format Apple App Attest reuses.
const (
	authDataRPIDHashLen  = 32
	authDataFlagsLen     = 1
	authDataSignCountLen = 4
	authDataAAGUIDLen    = 16
	authDataCredIDLenLen = 2
	// authDataMinLen is the smallest valid authData: all fixed fields present with a zero-length
	// credential ID.
	authDataMinLen = authDataRPIDHashLen + authDataFlagsLen + authDataSignCountLen +
		authDataAAGUIDLen + authDataCredIDLenLen
	authDataFlagAttestedCD = 0x40
)

// Apple App Attest AAGUID values identifying the environment a key was attested in. Both are
// accepted; this verifier does not restrict to a single environment.
var (
	aaguidProduction  = [16]byte{'a', 'p', 'p', 'a', 't', 't', 'e', 's', 't', 0, 0, 0, 0, 0, 0, 0}
	aaguidDevelopment = [16]byte{'a', 'p', 'p', 'a', 't', 't', 'e', 's', 't', 'd', 'e', 'v', 'e', 'l', 'o', 'p'}
)
