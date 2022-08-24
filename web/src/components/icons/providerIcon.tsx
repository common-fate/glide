import { IconProps } from "@chakra-ui/react";
import React from "react";
import { Provider } from "../../utils/backend-client/types";
import {
  AWSIcon,
  GrantedKeysIcon,
  OktaIcon,
  AzureIcon,
  KubernetesIcon,
  SnowflakeIcon,
  OnePasswordIcon,
  GitHubIcon,
  DjangoIcon,
  GRPCIcon,
  PythonIcon,
  EKSIcon,
} from "./Icons";

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
    case "azure-ad":
      return <AzureIcon {...rest} />;
    case "kubernetes":
      return <KubernetesIcon {...rest} />;
    case "snowflake":
      return <SnowflakeIcon {...rest} />;
    case "1password":
      return <OnePasswordIcon {...rest} />;
    case "github":
      return <GitHubIcon {...rest} />;
    case "django":
      return <DjangoIcon {...rest} />;
    case "grpc":
      return <GRPCIcon {...rest} />;
    case "python":
      return <PythonIcon {...rest} />;
    case "aws-eks-roles-sso":
      return <EKSIcon {...rest} />;
    default:
      return <GrantedKeysIcon {...rest} />;
  }
};
