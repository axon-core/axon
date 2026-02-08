#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

make image
make push
go install github.com/axon-core/axon/cmd/axon

axon install
kubectl rollout restart deployment/axon-controller-manager -n axon-system
kubectl rollout status deployment/axon-controller-manager -n axon-system
