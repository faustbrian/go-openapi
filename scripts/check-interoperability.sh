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
go mod edit -require "github.com/faustbrian/go-json-schema/v2@$(awk '$1 == "github.com/faustbrian/go-json-schema/v2" { print $2 }' "$root/go.mod")" \
    -require "$module@$version"
# The v2 producer's public JSON Schema dependency legitimately advances its
# supplier graph. Retain every peer pin; update only already-selected shared
# dependencies to the versions in the producer's resolved graph.
(cd "$root" && go list -m -f '{{if not .Main}}{{.Path}} {{.Version}}{{end}}' all) >producer-modules.txt
while read -r dependency dependency_version; do
    [ -n "$dependency" ] || continue
    if awk -v dependency="$dependency" '$1 == dependency { found=1 } END { exit !found }' go.mod; then
        go mod edit -require "$dependency@$dependency_version"
    fi
done <producer-modules.txt
# Require-block grouping is formatting: edit may insert the new direct module
# into an existing indirect block, and tidy moves it without changing its pin.
awk 'NF && $1 != "github.com/faustbrian/go-openapi/v2" && $1 != ")" && !($1 == "require" && $2 == "(") { sub(/^[ \t]+/, ""); print }' \
    go.mod | LC_ALL=C sort >expected-go.mod
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
awk 'NF && $1 != "github.com/faustbrian/go-openapi/v2" && $1 != ")" && !($1 == "require" && $2 == "(") { sub(/^[ \t]+/, ""); print }' \
    go.mod | LC_ALL=C sort >resolved-go.mod
diff -u expected-go.mod resolved-go.mod
# Public supplier sums are authenticated by tidy/verify above. Preserve exact
# historical sums for the remaining peer graph rather than rewriting evidence.
awk 'NR == FNR { producer[$1]=1; next } $1 !~ /^github\.com\/faustbrian\/go-/ && !producer[$1] { print }' \
    producer-modules.txt "$root/interoperability/go.sum" >"$temporary/tracked-external.sum"
awk 'NR == FNR { producer[$1]=1; next } $1 !~ /^github\.com\/faustbrian\/go-/ && !producer[$1] { print }' \
    producer-modules.txt "$temporary/go.sum" >"$temporary/resolved-external.sum"
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
