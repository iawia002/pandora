#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

GO111MODULE=on go install sigs.k8s.io/controller-tools/cmd/controller-gen

echo "Generating Webhook"
controller-gen webhook paths=./apis/... output:dir=./webhook
