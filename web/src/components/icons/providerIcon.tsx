import { IconProps } from "@chakra-ui/react";
import React from "react";
import {
  AWSIcon,
  AzureIcon,
  EKSIcon,
  GrantedKeysIcon,
  OktaIcon,
} from "./Icons";

interface Props extends IconProps {
  /**
   * The short type of the provider,
   * e.g. "aws-sso".
   * @deprecated use `type` instead (which uses the namespaced `commonfate/aws-sso` type).
   */
  shortType?: string;

  /**
   * The type of the provider, including the namespace, e.g. `commonfate/aws-sso`.
   */
  type?: string;
}

export const ProviderIcon: React.FC<Props> = ({
  shortType,
  type,
  ...rest
}): React.ReactElement => {
  if (shortType === undefined && type === undefined) {
    // @ts-ignore
    return null;
  }
  switch (type) {
    case "commonfate/aws-sso":
      return <AWSIcon {...rest} />;
    case "commonfate/okta":
      return <OktaIcon {...rest} />;
    case "commonfate/azure-ad":
      return <AzureIcon {...rest} />;
    case "commonfate/aws-eks-roles-sso":
      return <EKSIcon {...rest} />;
  }

  switch (shortType) {
    case "aws-sso":
      return <AWSIcon {...rest} />;
    case "okta":
      return <OktaIcon {...rest} />;
    case "azure-ad":
      return <AzureIcon {...rest} />;
    case "aws-eks-roles-sso":
      return <EKSIcon {...rest} />;
  }
  return <GrantedKeysIcon {...rest} />;
};
