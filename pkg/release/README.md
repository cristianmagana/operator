# Release Package

> ## About

This package will get the latest EKS optimized release based on the `EKS_VERSION` and if using arm64 instance you will explicitly need the `NODE_ARCHITECTURE` set to `arm64` or the default of amd64 will be pulled. 

> ### Run tests

1. go test -cover
2. go test -coverprofile=coverage.out
3. go tool cover -html="coverage.out"
