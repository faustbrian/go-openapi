#!/bin/sh
set -eu

root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
mode=${OPENAPI_INTEROPERABILITY_MODE:-public}
case "$mode" in
    candidate|public) ;;
    *) echo 'interoperability mode must be candidate or public' >&2; exit 1 ;;
esac
temporary=$(mktemp -d "${GOTMPDIR:-${TMPDIR:-/tmp}}/go-openapi-interoperability.XXXXXX")
cleanup() {
    chmod -R u+w "$temporary" 2>/dev/null || true
    [ ! -e "$temporary" ] || find "$temporary" -depth -delete
}
trap cleanup EXIT HUP INT TERM

cp "$root/interoperability/go.mod" "$temporary/go.mod"
awk '$1 !~ /^github\.com\/faustbrian\/go-/ { print }' \
    "$root/interoperability/go.sum" >"$temporary/go.sum"
cp "$root/interoperability/runner.go" "$temporary/runner.go"
cd "$temporary"
export GOWORK=off
module=github.com/faustbrian/go-openapi/v2
if [ "$mode" = candidate ]; then
    go mod edit -replace "$module=$root"
else
    if [ -n "${OPENAPI_INTEROPERABILITY_VERSION:-}" ]; then
        go mod edit -require "$module=$OPENAPI_INTEROPERABILITY_VERSION"
    fi
    if awk '$1 == "github.com/faustbrian/go-openapi/v2" && $2 == "v2.0.0-00010101000000-000000000000" { found=1 } END { exit !found }' go.mod; then
        echo 'public interoperability requires a published v2 version; the tracked version is a candidate placeholder' >&2
        exit 1
    fi
    export GOPROXY=https://proxy.golang.org GOSUMDB=sum.golang.org GOPRIVATE= GONOSUMDB= GONOPROXY=
fi
go mod tidy
go mod verify >/dev/null

if [ "$mode" = candidate ]; then
    go mod edit -dropreplace "$module"
    diff -u "$root/interoperability/go.mod" "$temporary/go.mod"
fi
awk '$1 !~ /^github\.com\/faustbrian\/go-/ { print }' \
    "$root/interoperability/go.sum" >"$temporary/tracked-external.sum"
awk '$1 !~ /^github\.com\/faustbrian\/go-/ { print }' \
    "$temporary/go.sum" >"$temporary/resolved-external.sum"
diff -u "$temporary/tracked-external.sum" "$temporary/resolved-external.sum"
if [ "$mode" = candidate ]; then
    go mod edit -replace "$module=$root"
fi

report="$temporary/report.tsv"
OPENAPI_INTEROPERABILITY_ROOT="$root" go run -mod=readonly . \
    "$root"/interoperability/fixtures/* \
    "$root"/specification/independent/swagger-petstore/openapi.yaml \
    "$root"/specification/independent/github-rest-api/api.github.com.2022-11-28.json \
    >"$report"

if [ "${INTEROP_UPDATE:-false}" = true ]; then
    cp "$report" "$root/interoperability/expected.tsv"
    exit 0
fi

diff -u "$root/interoperability/expected.tsv" "$report"
