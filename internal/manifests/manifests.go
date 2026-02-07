package manifests

import _ "embed"

//go:embed install-crd.yaml
var InstallCRD []byte

//go:embed install.yaml
var InstallController []byte
