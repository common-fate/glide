import { IconProps } from "@chakra-ui/react";
import React from "react";
import { Provider } from "../../utils/backend-client/types";
import { AWSIcon, GrantedKeysIcon, OktaIcon } from "./Icons";

export const getProviderIcon = (
  provider: Provider | undefined
): React.ReactElement => {
  if (provider === undefined) {
    // @ts-ignore
    return null;
  }
  switch (provider.type) {
    case "aws-sso":
      return <AWSIcon />;
    case "okta":
      return <OktaIcon />;
    default:
      return <GrantedKeysIcon />;
  }
};

export const ProviderIcon = ({
  provider,
  ...rest
}: {
  provider: Provider | undefined;
} & IconProps): React.ReactElement => {
  if (provider === undefined) {
    // @ts-ignore
    return null;
  }
  switch (provider.type) {
    case "aws-sso":
      return <AWSIcon {...rest} />;
    case "okta":
      return <OktaIcon {...rest} />;
    default:
      return <GrantedKeysIcon {...rest} />;
  }
};
