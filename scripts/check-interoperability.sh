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
sed 's|"github.com/faustbrian/go-openapi\(["/]\)|"github.com/faustbrian/go-openapi/v2\1|g' \
    "$root/interoperability/runner.go" >"$temporary/runner.go"
cd "$temporary"
export GOWORK=off
module=github.com/faustbrian/go-openapi/v2
if [ "$mode" = public ]; then
    if awk '$1 == "replace" { found=1 } END { exit !found }' go.mod; then
        echo 'public interoperability forbids module replacements' >&2
        exit 1
    fi
    version=${OPENAPI_INTEROPERABILITY_VERSION:-}
    case "$version" in
        v2.*) ;;
        *) echo 'public interoperability requires an explicit published v2 version' >&2; exit 1 ;;
    esac
    if [ "$version" = v2.0.0-00010101000000-000000000000 ]; then
        echo 'public interoperability rejects the unpublished candidate placeholder' >&2
        exit 1
    fi
else
    version=v2.0.0-00010101000000-000000000000
fi
go mod edit -droprequire github.com/faustbrian/go-openapi \
    -require "$module@$version"
awk '$1 != "github.com/faustbrian/go-openapi" { print }' \
    "$root/interoperability/go.mod" >expected-go.mod
if [ "$mode" = candidate ]; then
    go mod edit -replace "$module=$root"
else
    export GOPROXY=https://proxy.golang.org GOSUMDB=sum.golang.org GOPRIVATE= GONOSUMDB= GONOPROXY=
fi
go mod tidy
go mod verify >/dev/null

if [ "$mode" = candidate ]; then
    go mod edit -dropreplace "$module"
fi
awk '$1 != "github.com/faustbrian/go-openapi/v2" { print }' \
    go.mod >resolved-go.mod
diff -u expected-go.mod resolved-go.mod
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
