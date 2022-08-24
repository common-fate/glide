export interface RegisteredProvider {
  type: string;
  name: string;
}

export const registeredProviders: RegisteredProvider[] = [
  {
    type: "aws-sso",
    name: "AWS SSO",
  },
  {
    type: "okta",
    name: "Okta Groups",
  },
  {
    type: "azure-ad",
    name: "Azure AD Groups",
  },
  {
    type: "aws-eks-roles-sso",
    name: "EKS (with AWS SSO)",
  },
];
