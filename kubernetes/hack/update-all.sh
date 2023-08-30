#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

bash "./hack/update-codegen.sh"
bash "./hack/update-crdgen.sh"
bash "./hack/update-webhook.sh"
