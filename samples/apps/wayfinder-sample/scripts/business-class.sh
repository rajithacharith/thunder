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
# Demo script: lock or unlock Business class flight availability.
#
# Usage:
#   ./business-class.sh lock    — make all Business class flights unavailable (forces CIBA approval)
#   ./business-class.sh unlock  — make all Business class flights available again
#
# To trigger the upgrade scheduler after approving a CIBA notification, run trigger-upgrade.sh separately.

BACKEND_URL="${WAYFINDER_BACKEND_URL:-http://localhost:8787}"

ACTION="${1:-}"

case "$ACTION" in
  lock)
    echo "Locking Business class flights..."
    curl -fsS -X POST "${BACKEND_URL}/api/demo/lock-business-class" \
      -H "Content-Type: application/json" | jq .
    ;;
  unlock)
    echo "Unlocking Business class flights..."
    curl -fsS -X POST "${BACKEND_URL}/api/demo/unlock-business-class" \
      -H "Content-Type: application/json" | jq .
    ;;
  *)
    echo "Usage: $0 {lock|unlock}"
    exit 1
    ;;
esac
