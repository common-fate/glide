export interface RegisteredProvider {
  type: string;
  /**
   * @deprecated
   * this needs to be removed in favour of a namespaced type
   * (e.g. `commonfate/aws-sso` rather than `aws-sso`)
   */
  shortType: string;
  name: string;
}

export const registeredProviders = [
  {
    type: "commonfate/aws-sso",
    shortType: "aws-sso",
    name: "AWS SSO",
  },
  {
    type: "commonfate/okta",
    shortType: "okta",
    name: "Okta Groups",
  },
  {
    type: "commonfate/azure-ad",
    shortType: "azure-ad",
    name: "Azure AD Groups",
  },
  {
    type: "commonfate/aws-eks-roles-sso",
    shortType: "aws-eks-roles-sso",
    name: "EKS (with AWS SSO)",
  },
  {
    type: "commonfate/flask",
    shortType: "flask",
    name: "Flask (with AWS SSO)",
  },
] as const;

export type RegisteredShortTypes = typeof registeredProviders[number]["shortType"];

export type RegisteredTypes = typeof registeredProviders[number]["type"];
