package commonfate

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=.api-codegen.yaml openapi.yml

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --old-config-style --package client --o pkg/client/api.gen.go --generate "types,client,spec,skip-prune" --include-tags=ProviderSetup --import-mapping=./accesshandler/openapi.yml:github.com/common-fate/common-fate/accesshandler/pkg/types openapi.yml

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=pkg/deploymentcli/.api-codegen.yaml ./pkg/deploymentcli/openapi.yml
