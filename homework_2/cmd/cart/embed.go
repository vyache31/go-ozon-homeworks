package main

import _ "embed"

//go:embed docs.html
var docsHTML []byte

//go:embed open-api.yaml
var openAPIyaml []byte
