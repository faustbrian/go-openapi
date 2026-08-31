#!/bin/sh
set -eu

root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
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
go mod edit -replace \
    "github.com/faustbrian/go-openapi=$root"
go mod tidy
go mod verify >/dev/null

go mod edit -dropreplace github.com/faustbrian/go-openapi
diff -u "$root/interoperability/go.mod" "$temporary/go.mod"
awk '$1 !~ /^github\.com\/faustbrian\/go-/ { print }' \
    "$root/interoperability/go.sum" >"$temporary/tracked-external.sum"
awk '$1 !~ /^github\.com\/faustbrian\/go-/ { print }' \
    "$temporary/go.sum" >"$temporary/resolved-external.sum"
diff -u "$temporary/tracked-external.sum" "$temporary/resolved-external.sum"
go mod edit -replace \
    "github.com/faustbrian/go-openapi=$root"

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
