export interface RegisteredProvider {
  type: string;
  /**
   * @deprecated
   * this needs to be removed in favour of a namespaced type
   * (e.g. `commonfate/aws-sso` rather than `aws-sso`)
   */
  shortType: string;
  name: string;
  alpha?: boolean;
}

export const registeredProviders: RegisteredProvider[] = [
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
    alpha: true,
  },
  {
    type: "commonfate/ecs-exec-sso",
    shortType: "ecs-exec-sso",
    name: "ECS Exec (with AWS SSO)",
    alpha: true,
  },
  {
    type: "commonfate/testvault",
    shortType: "testvault",
    name: "TestVault",
  },
  {
    type: "commonfate/shell",
    shortType: "shell",
    name: "Shell",
    alpha: true,
  },
  {
    type: "commonfate/actions",
    shortType: "actions",
    name: "Actions",
    alpha: true,
  },
];

/**
 * If we type registeredProviders with a const assertion i.e. `registeredProviders = [...] as const;`
 * it is possible to strongly type the shortType key-values (could be beneficial)
 */
export type RegisteredShortTypes = typeof registeredProviders[number]["shortType"];

export type RegisteredTypes = typeof registeredProviders[number]["type"];
