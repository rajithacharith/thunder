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
# Resets the Wayfinder SQLite database by deleting it and re-running seed.js.
# Run from anywhere — script resolves paths relative to itself.

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$(dirname "$SCRIPT_DIR")/backend"
DB_PATH="${SQLITE_DB_PATH:-$BACKEND_DIR/wayfinder.sqlite}"

for f in "$DB_PATH" "${DB_PATH}-shm" "${DB_PATH}-wal"; do
  if [ -f "$f" ]; then
    rm "$f"
    echo "Deleted $f"
  fi
done

echo "Re-seeding database..."
node "$BACKEND_DIR/scripts/seed.js"
echo "Database reset complete."
