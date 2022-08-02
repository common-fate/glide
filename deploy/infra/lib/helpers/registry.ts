
/**
 * The IdentityProviderRegistry matches the registry keys from the identitysync package in go
 * These keys will be set in the deployment config and should be used as options for idp type
 */
export const IdentityProviderRegistry = {
    CognitoV1Key : "commonfate/identity/cognito@v1",
    OktaV1Key    : "commonfate/identity/okta@v1",
    AzureADV1Key : "commonfate/identity/azure-ad@v1",
    GoogleV1Key  : "commonfate/identity/google@v1"
} as const 

export type IdentityProviderRegistryValues = typeof IdentityProviderRegistry[keyof typeof IdentityProviderRegistry]