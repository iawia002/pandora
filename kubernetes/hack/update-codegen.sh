#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

GO111MODULE=on go install k8s.io/code-generator/cmd/deepcopy-gen
GO111MODULE=on go install k8s.io/code-generator/cmd/register-gen
GO111MODULE=on go install k8s.io/code-generator/cmd/client-gen
GO111MODULE=on go install k8s.io/code-generator/cmd/lister-gen
GO111MODULE=on go install k8s.io/code-generator/cmd/informer-gen
GO111MODULE=on go install k8s.io/kube-openapi/cmd/openapi-gen

# All API group names in the pkg/apis directory that need code generation
PKGS=(foo/v1alpha1 foo/v1alpha2)

CLIENT_PATH=github.com/iawia002/pandora/kubernetes
CLIENT_APIS=${CLIENT_PATH}/apis

FULL_PKGS=()
for path in "${PKGS[@]}"; do
  FULL_PKGS+=("$CLIENT_APIS/$path")
done

echo "Generating for API group:" "${PKGS[@]}"

echo "Generating with deepcopy-gen"
deepcopy-gen \
  --go-header-file=hack/boilerplate/boilerplate.go.txt \
  "${FULL_PKGS[@]}"

echo "Generating with register-gen"
register-gen \
  --go-header-file=hack/boilerplate/boilerplate.go.txt \
  "${FULL_PKGS[@]}"

echo "Generating with client-gen"
client-gen \
  --go-header-file=hack/boilerplate/boilerplate.go.txt \
  --clientset-name="clientset" \
  --input-base="" \
  --output-pkg="${CLIENT_PATH}/generated" \
  --output-dir=generated \
  $(printf -- " --input %s" "${FULL_PKGS[@]}")

echo "Generating with lister-gen"
lister-gen \
  --go-header-file=hack/boilerplate/boilerplate.go.txt \
  --output-pkg="${CLIENT_PATH}/generated/listers" \
  --output-dir=generated/listers \
  "${FULL_PKGS[@]}"

echo "Generating with informer-gen"
informer-gen \
  --go-header-file=hack/boilerplate/boilerplate.go.txt \
  --single-directory \
  --versioned-clientset-package="${CLIENT_PATH}/generated/clientset" \
  --listers-package="${CLIENT_PATH}/generated/listers" \
  --output-pkg="${CLIENT_PATH}/generated/informers" \
  --output-dir=generated/informers \
  "${FULL_PKGS[@]}"

echo "Generating with openapi-gen"
openapi-gen \
  --go-header-file=hack/boilerplate/boilerplate.go.txt \
  --output-pkg="${CLIENT_PATH}/generated/openapi" \
  --output-dir=generated/openapi \
  "${FULL_PKGS[@]}" \
  "k8s.io/api/core/v1" \
  "k8s.io/apimachinery/pkg/api/resource" \
  "k8s.io/apimachinery/pkg/apis/meta/v1" \
  "k8s.io/apimachinery/pkg/runtime" \
  "k8s.io/apimachinery/pkg/version" \
  "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1" \
  "k8s.io/api/admissionregistration/v1" \
  "k8s.io/api/networking/v1"
