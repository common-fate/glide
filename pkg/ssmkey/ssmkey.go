package ssmkey

import "path"

type SSMKeyOpts struct {
	HandlerID    string
	Key          string
	Publisher    string
	ProviderName string
}

// this will create a unique identifier for AWS System Manager Parameter Store
// for configuration field "api_url" this will result: 'publisher/provider-name/version/configuration/api_url'
func SSMKey(opts SSMKeyOpts) string {
	return "/" + path.Join("common-fate", "provider", opts.Publisher, opts.ProviderName, opts.HandlerID, opts.Key)
}
