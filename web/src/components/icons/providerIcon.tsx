import { IconProps } from "@chakra-ui/react";
import React from "react";
import { CommonFateIcon } from "./Logos";
import {
  AWSCloudwatch,
  AWSIcon,
  AzureIcon,
  ECSIcon,
  EKSIcon,
  FlaskIcon,
  GCPIcon,
  GoogleIcon,
  OktaIcon,
  SnowflakeIcon,
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

  /**
   * temporary hack for PDK provider icons not showing
   */
  id?: string;
}

export const ProviderIcon: React.FC<Props> = ({
  shortType,
  type,
  id,
  ...rest
}): React.ReactElement => {
  if (shortType === undefined && type === undefined) {
    // @ts-ignore
    return null;
  }

  switch (id) {
    case "azure":
      return <AzureIcon {...rest} />;
    case "gcp":
      return <GCPIcon {...rest} />;
    case "snowflake":
      return <SnowflakeIcon {...rest} />;
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
    case "commonfate/ecs-exec-sso":
      return <ECSIcon {...rest} />;
    case "commonfate/flask":
      return <FlaskIcon {...rest} />;
  }

  switch (shortType) {
    case "aws-sso" || "aws":
      return <AWSIcon {...rest} />;
    case "okta":
      return <OktaIcon {...rest} />;
    case "azure-ad" || "azure":
      return <AzureIcon {...rest} />;
    case "aws-eks-roles-sso":
      return <EKSIcon {...rest} />;
    case "ecs-exec-sso":
      return <ECSIcon {...rest} />;
    case "flask":
      return <FlaskIcon {...rest} />;
    case "aws-cloudwatch":
      return <AWSCloudwatch {...rest} />;
    case "google":
      return <GoogleIcon {...rest} />;
  }

  return <CommonFateIcon {...rest} />;
};
