#!/usr/bin/env bash
go run $(pwd)/vendor/k8s.io/code-generator/cmd/deepcopy-gen/main.go --input-dirs ./pkg/apis/component/v1alpha1/ -O zz_generated.deepcopy --bounding-dirs github.com/snowdrop/component-api/pkg/apis "component:v1alpha1"