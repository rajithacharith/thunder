#!/bin/bash
# Copyright (c) 2026, WSO2 LLC. (https://www.wso2.com).
#
# WSO2 LLC. licenses this file to you under the Apache License,
# Version 2.0 (the "License"); you may not use this file except
# in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.

set -euo pipefail
# Demo script: makes all Business class flights available and triggers upgrade processing.
# By default, Business class flights require CIBA async approval (available = 0).

BACKEND_URL="${WAYFINDER_BACKEND_URL:-http://localhost:8787}"
AGENT_URL="${WAYFINDER_AGENT_URL:-http://localhost:8790}"

echo "Unlocking Business class flights..."
curl -fsS -X POST "${BACKEND_URL}/api/demo/unlock-business-class" \
  -H "Content-Type: application/json" | jq .

curl -fsS -X POST "${AGENT_URL}/api/demo/process-upgrades" \
  -H "Content-Type: application/json" | jq .
