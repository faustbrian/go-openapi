#!/usr/bin/env bash
set -euo pipefail

root="$(git rev-parse --show-toplevel)"
task="$(mktemp -d "${TMPDIR:-/tmp}/golib-codeql-build.XXXXXX")"
cleanup() {
    find "${task}" -depth -delete 2>/dev/null || true
}
trap cleanup EXIT HUP INT TERM

root_version="$(jq -er '
    .modules[] | select(.directory == ".") | .version
' "${root}/modules.json")"
self_proxy="${task}/self-proxy"
"${root}/.golib/scripts/build-local-proxy.sh" \
    "${self_proxy}" "v${root_version}" "."
upstream_proxy="${GOPROXY:-$(go env GOPROXY)}"
export GOPROXY="file://${self_proxy},${upstream_proxy}"
current_no_sum_db="$(go env GONOSUMDB)"
export GONOSUMDB="github.com/faustbrian/go-*${current_no_sum_db:+,${current_no_sum_db}}"

clean_flags=""
for flag in ${GOFLAGS:-}; do
    case "${flag}" in
        -mod=*|-modfile=*) ;;
        *) clean_flags="${clean_flags:+${clean_flags} }${flag}" ;;
    esac
done

while IFS= read -r module; do
    [[ -n "${module}" ]] || continue
    module_root="${root}"
    if [[ "${module}" != "." ]]; then
        module_root="${root}/${module}"
    fi
    module_slug="$(printf '%s' "${module}" | tr '/.' '--')"
    module_state="${task}/modfiles/${module_slug}"
    modfile="${module_state}/isolated.mod"
    mkdir -p "${module_state}"
    cp "${module_root}/go.mod" "${modfile}"
    if [[ -f "${module_root}/go.sum" ]]; then
        awk '$1 !~ /^github\.com\/faustbrian\/go-/ { print }' \
            "${module_root}/go.sum" >"${module_state}/isolated.sum"
    else
        : >"${module_state}/isolated.sum"
    fi
    module_flags="${clean_flags:+${clean_flags} }-modfile=${modfile}"
    (cd "${module_root}" && GOWORK=off GOFLAGS="${module_flags}" \
        go mod download all)
    build_flags="${module_flags} -mod=readonly"
    while IFS= read -r package; do
        [[ -n "${package}" ]] || continue
        package_tags="$(
            jq -r \
                --arg module "${module}" \
                --arg package "${package}" '
                    .modules[]
                    | select(.directory == $module)
                    | .packages[]
                    | select(.import_path == $package)
                    | .build_tags[]?
                ' "${root}/modules.json"
        )"
        slug="$(printf '%s' "${package}" | tr '/.' '--')"
        if [[ "${module}" == "benchmarks/platform" && -n "${package_tags}" ]]; then
            package_tags=benchmark_disabled
        fi
        if [[ -z "${package_tags}" ]]; then
            (cd "${module_root}" && GOWORK=off GOFLAGS="${build_flags}" \
                go build -o "${task}/${slug}" "${package}")
            continue
        fi
        variant=0
        while IFS= read -r tag; do
            [[ -n "${tag}" ]] || continue
            (cd "${module_root}" && GOWORK=off GOFLAGS="${build_flags}" go build \
                -tags="${tag}" -o "${task}/${slug}-${variant}" "${package}")
            variant=$((variant + 1))
        done <<<"${package_tags}"
    done < <(
        jq -r --arg module "${module}" '
            .modules[]
            | select(.directory == $module)
            | .packages[]
            | select(.build_required == true)
            | .import_path
        ' "${root}/modules.json"
    )
done < <(jq -r '.modules[].directory' "${root}/modules.json")
