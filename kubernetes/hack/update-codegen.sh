#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

GO111MODULE=on go install k8s.io/code-generator/cmd/deepcopy-gen
GO111MODULE=on go install k8s.io/code-generator/cmd/register-gen
GO111MODULE=on go install k8s.io/code-generator/cmd/conversion-gen
GO111MODULE=on go install k8s.io/code-generator/cmd/client-gen
GO111MODULE=on go install k8s.io/code-generator/cmd/lister-gen
GO111MODULE=on go install k8s.io/code-generator/cmd/informer-gen
GO111MODULE=on go install k8s.io/code-generator/cmd/openapi-gen

# All API group names in the pkg/apis directory that need code generation
PKGS=(foo/v1alpha1)

CLIENT_PATH=github.com/iawia002/pandora/kubernetes
CLIENT_APIS=${CLIENT_PATH}/apis

ALL_PKGS=""
for path in "${PKGS[@]}"
do
  ALL_PKGS=$ALL_PKGS" $CLIENT_APIS/$path"
done

function join {
  local IFS="$1"
  shift
  result="$*"
}

join "," $ALL_PKGS
FULL_PKGS=$result

echo "Generating for API group:" "${PKGS[@]}"

echo "Generating with deepcopy-gen"
deepcopy-gen \
  --go-header-file=hack/boilerplate/boilerplate.go.txt \
  --input-dirs="${FULL_PKGS}" \
  --output-file-base=zz_generated.deepcopy

echo "Generating with register-gen"
register-gen \
  --go-header-file=hack/boilerplate/boilerplate.go.txt \
  --input-dirs="${FULL_PKGS}" \
  --output-file-base=zz_generated.register

echo "Generating with conversion-gen"
conversion-gen \
  --go-header-file=hack/boilerplate/boilerplate.go.txt \
  --input-dirs="${FULL_PKGS}" \
  --output-file-base=zz_generated.conversion

echo "Generating with client-gen"
client-gen \
  --go-header-file=hack/boilerplate/boilerplate.go.txt \
  --clientset-name="clientset" \
  --input-base="" \
  --input="${FULL_PKGS}" \
  --output-package="${CLIENT_PATH}/generated"

echo "Generating with lister-gen"
lister-gen \
  --go-header-file=hack/boilerplate/boilerplate.go.txt \
  --input-dirs="${FULL_PKGS}" \
  --output-package="${CLIENT_PATH}/generated/listers"

echo "Generating with informer-gen"
informer-gen \
  --go-header-file=hack/boilerplate/boilerplate.go.txt \
  --single-directory \
  --input-dirs="${FULL_PKGS}" \
  --versioned-clientset-package="${CLIENT_PATH}/generated/clientset" \
  --listers-package="${CLIENT_PATH}/generated/listers" \
  --output-package="${CLIENT_PATH}/generated/informers"

echo "Generating with openapi-gen"
openapi-gen \
  --go-header-file=hack/boilerplate/boilerplate.go.txt \
  --input-dirs="${FULL_PKGS}" \
  --input-dirs="k8s.io/api/core/v1,k8s.io/apimachinery/pkg/api/resource" \
  --input-dirs="k8s.io/apimachinery/pkg/apis/meta/v1,k8s.io/apimachinery/pkg/runtime,k8s.io/apimachinery/pkg/version" \
  --input-dirs="k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1,k8s.io/api/admissionregistration/v1,k8s.io/api/networking/v1" \
  --output-package="${CLIENT_PATH}/generated/openapi" \
  --output-file-base=zz_generated.openapi
