#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

GO111MODULE=on go install sigs.k8s.io/controller-tools/cmd/controller-gen

echo "Generating CRDs"
controller-gen crd paths=./apis/... output:dir=./crds
