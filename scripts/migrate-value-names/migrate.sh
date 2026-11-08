#!/bin/sh
# Copyright 2026 The ThunderID Authors
# SPDX-License-Identifier: Apache-2.0

# Runs the migration of captured values to ID-based names on each Data Plane given: checks the naming rules,
# reads what the Data Plane holds, shows what it would copy, and copies when asked. It does this for each
# deployment ID asked for, or for every deployment the Data Plane holds values for.
#
# It only needs the Data Plane's databases, so it can run from a host that cannot reach the Control Plane.
#
# Plain POSIX sh: it runs under dash, ash (Alpine, busybox), bash and zsh, on Linux and macOS. It needs only
# psql 10 or later and the usual tools (sed, tr, awk, wc, mktemp).
#
# Without --apply it only looks: nothing is written.

set -eu

HERE=$(cd "$(dirname "$0")" && pwd)

usage() {
    cat <<'EOF'
Usage: migrate.sh (--deployment-id ID ... | --all-deployments) --gateway-dsn DSN [--entity-dsn DSN] ... [options]

  --gateway-dsn DSN     a Data Plane's CONFIGURATION database; repeat for each Data Plane
  --entity-dsn DSN      that Data Plane's ENTITY database. Give it right after its --gateway-dsn. Without it,
                        applications, users and agents of that Data Plane are not migrated.
  --deployment-id ID    a DEPLOYMENT_ID to migrate; repeat for several       (env: DEPLOYMENT_IDS, one
                                                                              per line)
  --all-deployments     migrate every deployment ID the Data Planes hold values for
  --apply               copy the values. Without it nothing is written.
  --yes                 with --apply, do not ask before copying
  -h, --help            show this help

One of --deployment-id and --all-deployments is required, so a run never widens its scope by accident.

A DSN is anything psql accepts, such as "host=db dbname=thunder user=ops" or a postgresql:// URL. Put
passwords in PGPASSWORD or ~/.pgpass so they do not appear in the process list.

Example, two Data Planes, the first with its entity database:
  migrate.sh --all-deployments \
    --gateway-dsn "host=dp1 dbname=configdb user=ops" --entity-dsn "host=dp1 dbname=entitydb user=ops" \
    --gateway-dsn "host=dp2 dbname=configdb user=ops"
EOF
}

die() { echo "error: $*" >&2; exit 1; }

# A newline and a tab, kept intact by the substitution. Lists below hold one value per line.
NL=$(printf '\nx')
NL=${NL%x}
TAB=$(printf '\t')

# count_lines TEXT prints how many non-empty lines TEXT has.
count_lines() { printf '%s\n' "$1" | awk 'NF { n++ } END { print n + 0 }'; }

DEPLOYMENT_IDS=${DEPLOYMENT_IDS:-}
GATEWAYS=""   # one line per Data Plane: its configuration DSN, a tab, and its entity DSN or a dash
ALL_DEPLOYMENTS=false
APPLY=false
ASSUME_YES=false

while [ $# -gt 0 ]; do
    case "$1" in
        --deployment-id)
            [ $# -ge 2 ] || die "--deployment-id needs a value"
            DEPLOYMENT_IDS="${DEPLOYMENT_IDS:+$DEPLOYMENT_IDS$NL}$2"; shift 2 ;;
        --all-deployments) ALL_DEPLOYMENTS=true; shift ;;
        --gateway-dsn)
            [ $# -ge 2 ] || die "--gateway-dsn needs a value"
            GATEWAYS="${GATEWAYS:+$GATEWAYS$NL}$2${TAB}-"; shift 2 ;;
        --entity-dsn)
            [ $# -ge 2 ] || die "--entity-dsn needs a value"
            [ -n "$GATEWAYS" ] || die "--entity-dsn must come right after the --gateway-dsn it belongs to"
            last=$(printf '%s\n' "$GATEWAYS" | tail -n 1)
            case "$last" in *"${TAB}-") ;; *) die "that --gateway-dsn already has an --entity-dsn" ;; esac
            rest=$(printf '%s\n' "$GATEWAYS" | sed '$d')
            GATEWAYS="${rest:+$rest$NL}${last%"${TAB}-"}${TAB}$2"; shift 2 ;;
        --apply) APPLY=true; shift ;;
        --yes) ASSUME_YES=true; shift ;;
        -h|--help) usage; exit 0 ;;
        *) usage >&2; die "unknown option: $1" ;;
    esac
done

if [ "$ALL_DEPLOYMENTS" = true ] && [ "$(count_lines "$DEPLOYMENT_IDS")" -gt 0 ]; then
    die "use either --all-deployments or --deployment-id, not both"
fi
if [ "$ALL_DEPLOYMENTS" != true ] && [ "$(count_lines "$DEPLOYMENT_IDS")" -eq 0 ]; then
    die "give --deployment-id ID (repeatable) or --all-deployments"
fi
[ "$(count_lines "$GATEWAYS")" -gt 0 ] || die "at least one --gateway-dsn is required"

command -v psql >/dev/null 2>&1 || die "psql is not installed"
PSQL_MAJOR=$(psql --version | sed 's/[^0-9]*\([0-9][0-9]*\).*/\1/')
[ "$PSQL_MAJOR" -ge 10 ] || die "psql 10 or later is required, found $PSQL_MAJOR"

# The scripts read and write their files in the current directory, so each run works in a private one.
WORK=$(mktemp -d)
cleanup() { rm -rf "$WORK"; }
trap cleanup EXIT
trap 'exit 130' INT TERM

# run_psql DIR ARGS... runs psql inside DIR.
run_psql() {
    _dir=$1; shift
    (cd "$_dir" && psql -X -q -v ON_ERROR_STOP=1 "$@")
}

# gateway_part N FIELD prints the configuration (1) or entity (2) DSN of the Nth Data Plane.
gateway_part() {
    printf '%s\n' "$GATEWAYS" | sed -n "${1}p" | awk -F"$TAB" -v f="$2" '{ print $f }'
}

FIRST_CONFIG=$(gateway_part 1 1)

echo "== Checking the naming rules"
mkdir -p "$WORK/check"
CHECK_OUTPUT=$(run_psql "$WORK/check" -tA "$FIRST_CONFIG" -f "$HERE/check_names.sql")
MISMATCHES=$(printf '%s\n' "$CHECK_OUTPUT" | grep -v -e '^Checked ' -e '^$' || true)
if [ -n "$MISMATCHES" ]; then
    echo "$MISMATCHES" >&2
    die "the naming rules in names.sql differ from the server's. Do not migrate with them."
fi
echo "   The rules match the server."

if [ "$ALL_DEPLOYMENTS" = true ]; then
    FOUND=""
    while IFS= read -r line; do
        [ -n "$line" ] || continue
        config=${line%%"$TAB"*}
        ids=$(run_psql "$WORK/check" -tA "$config" -c \
            'SELECT DEPLOYMENT_ID FROM "SECRET" UNION SELECT DEPLOYMENT_ID FROM "VARIABLE" ORDER BY 1')
        FOUND="${FOUND:+$FOUND$NL}$ids"
    done <<EOF
$GATEWAYS
EOF
    DEPLOYMENT_IDS=$(printf '%s\n' "$FOUND" | awk 'NF && !seen[$0]++')
    [ "$(count_lines "$DEPLOYMENT_IDS")" -gt 0 ] || die "the Data Planes hold no variables or secrets"
    echo "== Deployments found: $(count_lines "$DEPLOYMENT_IDS")"
fi

index=0
PLANNED=""   # one line per run that would copy something: its directory
while IFS= read -r dep; do
    [ -n "$dep" ] || continue
    gateway=0
    while IFS= read -r line; do
        [ -n "$line" ] || continue
        gateway=$((gateway + 1))
        index=$((index + 1))
        config=${line%%"$TAB"*}
        entity=${line#*"$TAB"}
        dir="$WORK/r$index"
        mkdir -p "$dir"
        printf '%s' "$dep" >"$dir/deployment-id"
        printf '%s' "$config" >"$dir/config-dsn"

        echo
        echo "== Deployment '$dep', Data Plane $gateway: reading what it holds"
        run_psql "$dir" "$config" -v "dep=$dep" -f "$HERE/1a_resources_config.sql"
        : >"$dir/resources-entity.csv"
        if [ "$entity" != "-" ]; then
            run_psql "$dir" "$entity" -v "dep=$dep" -f "$HERE/1b_resources_entity.sql"
        else
            echo "   No --entity-dsn for this Data Plane: its applications, users and agents are not migrated."
        fi
        cat "$dir/resources-config.csv" "$dir/resources-entity.csv" >"$dir/resources.csv"

        echo "-- Deployment '$dep', Data Plane $gateway: what would be copied"
        look=$(run_psql "$dir" "$config" -v "dep=$dep" -v apply=false -f "$HERE/2_apply.sql")
        printf '%s\n' "$look"
        would=$(printf '%s\n' "$look" | sed -n 's/^Would copy \([0-9][0-9]*\) values.*/\1/p')
        if [ "${would:-0}" -gt 0 ]; then
            PLANNED="${PLANNED:+$PLANNED$NL}$dir"
        fi
    done <<EOF
$GATEWAYS
EOF
done <<EOF
$DEPLOYMENT_IDS
EOF

if [ "$(count_lines "$PLANNED")" -eq 0 ]; then
    echo
    echo "Nothing to migrate."
    exit 0
fi

if [ "$APPLY" != true ]; then
    echo
    echo "Nothing was written. Run again with --apply to copy."
    exit 0
fi

if [ "$ASSUME_YES" != true ]; then
    [ -t 0 ] || die "--apply asks before copying and there is no terminal. Add --yes to go ahead."
    printf 'Copy these values in %s place(s)? [y/N] ' "$(count_lines "$PLANNED")"
    read -r answer
    case "$answer" in
        y|Y) ;;
        *) echo "Cancelled. Nothing was written."; exit 0 ;;
    esac
fi

while IFS= read -r dir; do
    [ -n "$dir" ] || continue
    dep=$(cat "$dir/deployment-id")
    config=$(cat "$dir/config-dsn")
    echo "== Deployment '$dep': copying"
    # Captured rather than piped: a plain sh pipeline reports only its last command's status.
    result=$(run_psql "$dir" "$config" -v "dep=$dep" -v apply=true -f "$HERE/2_apply.sql")
    printf '%s\n' "$result" | tail -n 3
done <<EOF
$PLANNED
EOF

cat <<'EOF'

Done. Next:
  1. Capture a new configuration version on the Control Plane.
  2. Apply it to each Data Plane and confirm nothing is reported as missing.
  3. Set the values again for the resources listed as ambiguous.
  4. The old rows are still there and harmless. To remove them later, see 3_cleanup.sql: it is a manual step
     that this script never runs.
EOF
