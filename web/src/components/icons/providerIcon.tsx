import { IconProps } from "@chakra-ui/react";
import React from "react";
import { Provider } from "../../utils/backend-client/types";
import {
  AWSIcon,
  GrantedKeysIcon,
  OktaIcon,
  AzureIcon,
  EKSIcon,
} from "./Icons";

export const ProviderIcon = ({
  provider,
  ...rest
}: {
  provider: string | undefined;
} & IconProps): React.ReactElement => {
  if (provider === undefined) {
    // @ts-ignore
    return null;
  }
  switch (provider) {
    case "aws-sso":
      return <AWSIcon {...rest} />;
    case "okta":
      return <OktaIcon {...rest} />;
    case "azure-ad":
      return <AzureIcon {...rest} />;
    case "aws-eks-roles-sso":
      return <EKSIcon {...rest} />;
    default:
      return <GrantedKeysIcon {...rest} />;
  }
};
