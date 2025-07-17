#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

PROJECT_MODULE="github.com/clbiggs/cyberark-to-k8s"

HACK_DIR="$(dirname ${BASH_SOURCE[0]})"
REPO_DIR="${HACK_DIR}/.."
CODEGEN_PKG="/tmp/go/src/k8s.io/code-generator"

BOILER="${HACK_DIR}/boilerplate.go.txt"

source "${CODEGEN_PKG}/kube_codegen.sh"

kube::codegen::gen_helpers \
  --boilerplate "${BOILER}" \
  "${REPO_DIR}"

if [[ -n "${API_KNOWN_VIOLATIONS_DIR:-}" ]]; then
    report_filename="${API_KNOWN_VIOLATIONS_DIR}/codegen_violation_exceptions.list"
    if [[ "${UPDATE_API_KNOWN_VIOLATIONS:-}" == "true" ]]; then
        update_report="--update-report"
    fi
fi

kube::codegen::gen_client \
  --with-watch \
  --boilerplate "${BOILER}" \
  --output-dir "${REPO_DIR}/pkg/k8s/client" \
  --output-pkg "${PROJECT_MODULE}/pkg/k8s/client" \
  "${REPO_DIR}/pkg/k8s/apis"

