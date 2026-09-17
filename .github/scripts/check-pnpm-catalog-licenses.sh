#!/usr/bin/env bash
# GitHub's Dependency Graph reads package.json's literal "catalog:" specifier
# for pnpm workspace catalog dependencies instead of resolving it through
# pnpm-workspace.yaml, so dependency-review-action never sees a version change
# for catalog-managed packages and can't check their licenses. This script
# diffs pnpm-lock.yaml's actually resolved versions between base and head
# directly, and checks the npm registry's license field for anything new or
# changed. Vulnerabilities are unaffected: the separate `pnpm audit` step
# reads the lockfile directly and already covers the full current dependency
# set regardless of this gap.
set -euo pipefail

BASE_SHA="$1"
HEAD_SHA="$2"

if [ -z "${APPROVED_LICENSES:-}" ]; then
  echo "::error::APPROVED_LICENSES environment variable is not set"
  exit 1
fi

APPROVED_JSON=$(printf '%s' "$APPROVED_LICENSES" | tr ',' '\n' | jq -Rn '[inputs | ltrimstr(" ") | rtrimstr(" ")] | map(select(length > 0))')

extract_versions() {
  git show "$1:pnpm-lock.yaml" | yq -o=json '.' - | jq -r '
    .importers // {} | to_entries[] |
    (.value.dependencies // {}, .value.devDependencies // {}, .value.optionalDependencies // {}) |
    to_entries[] |
    select(.value.version != null) |
    "\(.key)@\(.value.version | sub("\\(.*$"; ""))"
  ' | sort -u
}

extract_versions "$BASE_SHA" > /tmp/pnpm-license-check-base.txt
extract_versions "$HEAD_SHA" > /tmp/pnpm-license-check-head.txt

CHANGED=$(comm -13 /tmp/pnpm-license-check-base.txt /tmp/pnpm-license-check-head.txt || true)

if [ -z "$CHANGED" ]; then
  echo "No pnpm dependency version changes detected between base and head."
  exit 0
fi

# SPDX-lite check: an "A OR B" expression passes if any side is approved; an
# "A AND B" expression (dual licensing where both apply) needs every side approved.
is_approved() {
  jq -en --arg license "$1" --argjson approved "$APPROVED_JSON" '
    ($license | gsub("^\\(|\\)$"; "")) as $expr |
    if ($expr | contains(" AND ")) then
      ($expr | split(" AND ") | map(. as $p | ($approved | index($p) != null)) | all)
    else
      ($expr | split(" OR ") | map(. as $p | ($approved | index($p) != null)) | any)
    end
  ' > /dev/null
}

FAILED=0
while IFS= read -r dep; do
  [ -z "$dep" ] && continue
  name="${dep%@*}"
  version="${dep##*@}"
  encoded_name=$(jq -rn --arg n "$name" '$n|@uri')

  registry_response=$(curl -s --retry 2 --retry-connrefused -o /tmp/pnpm-license-check-pkg.json -w '%{http_code}' "https://registry.npmjs.org/${encoded_name}" || echo "000")
  if [ "$registry_response" != "200" ]; then
    echo "::warning::Could not fetch registry metadata for $dep (HTTP $registry_response) - likely a private/workspace-only package, skipping"
    continue
  fi

  license=$(jq -r --arg v "$version" '.versions[$v].license // .versions[$v].licenses[0].type // "UNKNOWN"' /tmp/pnpm-license-check-pkg.json)

  if [ "$license" = "UNKNOWN" ] || [ "$license" = "null" ]; then
    echo "::warning::No license metadata found for $dep, skipping"
    continue
  fi

  if is_approved "$license"; then
    echo "✅ $dep - $license (approved)"
  else
    echo "::error::$dep uses license '$license', which is not on the approved list"
    FAILED=1
  fi
done <<< "$CHANGED"

exit "$FAILED"
