#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

HACK_DIR="$(dirname ${BASH_SOURCE})"
REPO_DIR="${HACK_DIR}/.."

PROJECT_MODULE="github.com/clbiggs/cyberark-to-k8s"
IMAGE_NAME="kubernetes-codegen:latest"

echo "Building codegen Docker image..."
docker build -f "${HACK_DIR}/Dockerfile" \
  -t "${IMAGE_NAME}" \
  "${REPO_DIR}"

CMD="${HACK_DIR}/project-codegen.sh"

echo "Generating kube code..."
echo "$CMD"
docker run --rm \
  -v "$(readlink -e ${REPO_DIR}):/tmp/go/src/${PROJECT_MODULE}" \
  "${IMAGE_NAME}" $CMD

echo ""
echo "Resetting file ownership to '${USER}:${USER}'"
echo "Password may be required"
echo ""

sudo chown ${USER}:${USER} -R ${REPO_DIR}/pkg/k8s
