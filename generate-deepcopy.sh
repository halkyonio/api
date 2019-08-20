#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

TMP_DIR=$(mktemp -d)

cleanup() {
  rm -rf "${TMP_DIR}"
}
trap "cleanup" EXIT SIGINT

cleanup

# The following solution for making code generation work with go modules is
# borrowed and modified from https://github.com/heptio/contour/pull/1010.
export GO111MODULE=on
VERSION=$(go list -m all | grep k8s.io/code-generator | rev | cut -d"-" -f1 | cut -d" " -f1 | rev)
echo "Using ${VERSION} of k8s code generator"
git clone https://github.com/kubernetes/code-generator.git "${TMP_DIR}"
(cd "${TMP_DIR}" && git reset --hard "${VERSION}")

# Usage: generate-groups.sh <generators> <output-package> <apis-package> <groups-versions> ...
#
#   <generators>        the generators comma separated to run (deepcopy,defaulter,client,lister,informer) or "all".
#   <output-package>    the output package name (e.g. github.com/example/project/pkg/generated).
#   <apis-package>      the external types dir (e.g. github.com/example/api or github.com/example/project/pkg/apis).
#   <groups-versions>   the groups and their versions in the format "groupA:v1,v2 groupB:v1 groupC:v2", relative
#                       to <api-package>.
#   ...                 arbitrary flags passed to all generator binaries.
#
#
# Examples:
#   generate-groups.sh all             github.com/example/project/pkg/client github.com/example/project/pkg/apis "foo:v1 bar:v1alpha1,v1beta1"
#   generate-groups.sh deepcopy,client github.com/example/project/pkg/client github.com/example/project/pkg/apis "foo:v1 bar:v1alpha1,v1beta1"

for group in capability component link; do
  "${TMP_DIR}"/generate-groups.sh "deepcopy,client" \
    halkyon.io/api/${group} \
    halkyon.io/api \
    "${group}:v1beta1" \
    --go-header-file "${TMP_DIR}"/hack/boilerplate.go.txt
done
