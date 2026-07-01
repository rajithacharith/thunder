-- ----------------------------------------------------------------------------
-- Copyright (c) 2026, WSO2 LLC. (https://www.wso2.com).
--
-- WSO2 LLC. licenses this file to you under the Apache License,
-- Version 2.0 (the "License"); you may not use this file except
-- in compliance with the License. You may obtain a copy of the License at
--
-- http://www.apache.org/licenses/LICENSE-2.0
--
-- Unless required by applicable law or agreed to in writing,
-- software distributed under the License is distributed on an
-- "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
-- KIND, either express or implied. See the License for the
-- specific language governing permissions and limitations
-- under the License.
-- ----------------------------------------------------------------------------

-- Table to store revoked token JTIs (single-token revocation deny list).
-- Part of the database.operation classification: authoritative authorization
-- enforcement state that must survive a runtime database flush.
CREATE TABLE "REVOKED_TOKEN" (
    DEPLOYMENT_ID VARCHAR(255) NOT NULL,
    ID VARCHAR(36) NOT NULL PRIMARY KEY,
    JTI VARCHAR(255) NOT NULL,
    REVOCATION_REASON VARCHAR(30) NOT NULL CHECK (REVOCATION_REASON IN ('explicit', 'refresh_rotation')),
    REVOKED_AT TIMESTAMP NOT NULL,
    EXPIRY_TIME TIMESTAMP NOT NULL
);

-- Unique index backs the hot deny-list lookup by (deployment, jti) and enforces idempotent revocation writes.
CREATE UNIQUE INDEX idx_revoked_token_jti_deployment ON "REVOKED_TOKEN" (DEPLOYMENT_ID, JTI);

-- Index for expiry time on REVOKED_TOKEN (supports cleanup and expiry checks).
CREATE INDEX idx_revoked_token_expiry_time ON "REVOKED_TOKEN" (EXPIRY_TIME);
