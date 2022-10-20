/**
 * The IdentityProviderRegistry matches the registry keys from the identitysync package in go
 * These keys will be set in the deployment config and should be used as options for idp type
 */
export const IdentityProviderRegistry = {
  Cognito: "cognito",
  Okta: "okta",
  AzureAD: "azure",
  Google: "google",
  AWSSSO: "aws-sso",
} as const;

export type IdentityProviderTypes = typeof IdentityProviderRegistry[keyof typeof IdentityProviderRegistry];
