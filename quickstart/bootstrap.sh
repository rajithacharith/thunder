#!/bin/bash
# Copyright 2025-2026 The ThunderID Authors
# SPDX-License-Identifier: Apache-2.0

# Sourced by start-quickstart.sh. Holds what the quickstart pack's bootstrap resources (the
# *.yaml files under quickstart/bootstrap/, layered onto each plane by build.sh's
# package_quickstart) need from the environment, so the orchestration script does not carry it.
#
# Expects AUTH_DP_PORT, DP_PORT and CP_PORT to be set by the caller.

# Management key the managed Data Plane trusts; the same fixed, dev-only value the Control Plane's
# gateway resource (quickstart/bootstrap/cp/00-quickstart-gateway.yaml) registers it with.
QUICKSTART_GATEWAY_KEY="7d9c7f5adadaef3f965cd483ffb0aa4c9993486c3674a1fa1d0a5bf21e32fe46"

# user_password_secret_name prints the name a user's password is held under in a gateway's secret
# store, which the Control Plane's export derives from the username the same way:
# USER_<USERNAME>_PASSWORD, with a word boundary wherever the case changes (ServiceDesk becomes
# SERVICE_DESK) and every character that is not a letter or digit as an underscore. See
# DeriveVariableName in backend/internal/system/varname.
user_password_secret_name() {
    local snake
    snake="$(awk -v s="$1" 'BEGIN {
        n = length(s); out = ""
        for (i = 1; i <= n; i++) {
            c = substr(s, i, 1)
            if (i > 1 && c ~ /[A-Z]/) {
                p = substr(s, i - 1, 1)
                if (p ~ /[a-z]/) out = out "_"
                else if (p != "_" && i < n && substr(s, i + 1, 1) ~ /[a-z]/) out = out "_"
            }
            out = out c
        }
        print toupper(out)
    }')"
    printf '%s' "USER_${snake}_PASSWORD" | sed 's/[^A-Za-z0-9]/_/g; s/__*/_/g; s/^_//; s/_$//'
}

# export_console_urls exports the console redirect URIs the bootstrap resources are templated with:
# CP_CONSOLE_URL for the CP_CONSOLE application on auth-dp, and ADMIN_CONSOLE_URL for the
# ADMIN_CONSOLE application on the Data Plane and the Control Plane. Call it before the first
# plane's setup runs, since auth-dp's bootstrap already needs CP_CONSOLE_URL.
export_console_urls() {
    export CP_CONSOLE_URL="https://localhost:$CP_PORT/console"
    export ADMIN_CONSOLE_URL="https://localhost:$DP_PORT/console"
}

# export_dp_ca_certificate exports DP_CA_CERTIFICATE, the managed Data Plane's certificate, which
# the Control Plane's gateway resource trusts it by. Real newlines are escaped to a literal
# two-character "\n" so the multi-line PEM survives as one line of bootstrap-template text; the
# YAML parser turns them back into real newlines once the file is templated and decoded (a plain
# double-quoted YAML scalar interprets \n as a line break). Needs dp/config/certs/server.cert, so
# call it once the key material has been shared.
export_dp_ca_certificate() {
    local dp_cert
    dp_cert="$(cat dp/config/certs/server.cert)"
    export DP_CA_CERTIFICATE="${dp_cert//$'\n'/\\n}"
}

# seed_dp_values sets the values auto-apply needs on the managed Data Plane before any write has
# captured them.
#
# Auto-apply (and a manual Capture/Apply, for that matter) refuses every sweep with "the gateway
# does not hold every value this version refers to" until each value a version names exists on the
# gateway. Two come from the Control Plane's own bootstrap, which runs while the Data Plane is not
# up to capture them to:
#
#   - the ClientID and redirect URIs of the ADMIN_CONSOLE application
#     (quickstart/bootstrap/dp-cp/04-quickstart-admin-console.yaml). The redirect URIs are a list,
#     so the variable holds a JSON array;
#   - the identifier of the System resource server, which a gateway holds for itself and is
#     the Data Plane's own MCP URL;
#   - the Service Desk Admin's password. A version carries that user to the Data Plane like any
#     other, and the Data Plane fills its password in from this secret.
#
# All are set directly on the gateway's own stores. Needs SERVICE_DESK_USERNAME and
# SERVICE_DESK_PASSWORD to be exported.
seed_dp_values() {
    for _ in 1 2 3 4 5 6 7 8 9 10; do
        curl -sk -o /dev/null "https://localhost:$DP_PORT/health/readiness" && break
        sleep 0.5
    done
    curl -sk -o /dev/null -X PUT "https://localhost:$DP_PORT/variables/APPLICATION_ADMIN_CONSOLE_CLIENT_ID" \
        -H "API-Key: $QUICKSTART_GATEWAY_KEY" \
        -H "Content-Type: application/json" \
        -d '{"value":"ADMIN_CONSOLE"}' || true
    curl -sk -o /dev/null -X PUT "https://localhost:$DP_PORT/variables/APPLICATION_ADMIN_CONSOLE_REDIRECT_URIS" \
        -H "API-Key: $QUICKSTART_GATEWAY_KEY" \
        -H "Content-Type: application/json" \
        -d "{\"value\":\"[\\\"$ADMIN_CONSOLE_URL\\\"]\"}" || true
    curl -sk -o /dev/null -X PUT "https://localhost:$DP_PORT/variables/RESOURCE_SERVER_SYSTEM_IDENTIFIER" \
        -H "API-Key: $QUICKSTART_GATEWAY_KEY" \
        -H "Content-Type: application/json" \
        -d "{\"value\":\"https://localhost:$DP_PORT/mcp\"}" || true
    local password_json
    password_json="$(printf '%s' "$SERVICE_DESK_PASSWORD" | sed 's/\\/\\\\/g; s/"/\\"/g')"
    curl -sk -o /dev/null -X PUT "https://localhost:$DP_PORT/secrets/$(user_password_secret_name "$SERVICE_DESK_USERNAME")" \
        -H "API-Key: $QUICKSTART_GATEWAY_KEY" \
        -H "Content-Type: application/json" \
        -d "{\"value\":\"$password_json\"}" || true
}
