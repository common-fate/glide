import { IconProps } from "@chakra-ui/react";
import React from "react";
import { CommonFateIcon } from "./Logos";
import {
  AWSCloudwatch,
  AWSIcon,
  AzureIcon,
  DataDog,
  ECSIcon,
  EKSIcon,
  FlaskIcon,
  GoogleCloudIcon,
  GoogleIcon,
  JenkinsLogo,
  JiraIcon,
  OktaIcon,
  OnePasswordIcon,
  SnowflakeIcon,
} from "./Icons";
import { GitHubIcon } from "./SocialIcons";

export type ShortTypes =
  | "aws-sso"
  | "aws"
  | "okta"
  | "azure-ad"
  | "azure"
  | "aws-eks-roles-sso"
  | "ecs-exec-sso"
  | "flask"
  | "aws-cloudwatch"
  | "cloudwatch"
  | "google"
  | "gcp"
  | "1pass"
  | "github"
  | "snowflake"
  | "jira"
  | "okta"
  | "jenkins"
  | "datadog";

export const shortTypesArr: ShortTypes[] = [
  "aws-sso",
  "aws",
  "okta",
  "azure-ad",
  "azure",
  "aws-eks-roles-sso",
  "ecs-exec-sso",
  "flask",
  "aws-cloudwatch",
  "cloudwatch",
  "google",
  "gcp",
  "1pass",
  "github",
  "snowflake",
  "jira",
  "okta",
  "jenkins",
  "datadog",
];

// use english title case
export const shortTypeValues: { [key in ShortTypes]: string } = {
  "aws": "AWS",
  "okta": "Okta",
  "azure": "Azure",
  "1pass": "1Password",
  "github": "GitHub",
  "snowflake": "Snowflake",
  "jira": "Jira",
  "jenkins": "Jenkins",
  "datadog": "Datadog",
  "gcp": "Google Cloud",
  "ecs-exec-sso": "ECS",
  "aws-eks-roles-sso": "EKS",
  "aws-cloudwatch": "Cloudwatch",
  "flask": "Flask",
};

interface Props extends IconProps {
  /**
   * The short type of the provider,
   * e.g. "aws-sso".
   */
  shortType?: ShortTypes;

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
    case "commonfate/ecs-exec-sso":
      return <ECSIcon {...rest} />;
    case "commonfate/flask":
      return <FlaskIcon {...rest} />;
  }

  switch (shortType) {
    case "aws-sso":
      return <AWSIcon {...rest} />;
    case "aws":
      return <AWSIcon {...rest} />;
    case "okta":
      return <OktaIcon {...rest} />;
    case "jira":
      return <JiraIcon {...rest} />;
    case "jenkins":
      return <JenkinsLogo {...rest} />;
    case "azure-ad":
      return <AzureIcon {...rest} />;
    case "azure":
      return <AzureIcon {...rest} />;
    case "aws-eks-roles-sso":
      return <EKSIcon {...rest} />;
    case "ecs-exec-sso":
      return <ECSIcon {...rest} />;
    case "flask":
      return <FlaskIcon {...rest} />;
    case "aws-cloudwatch":
      return <AWSCloudwatch {...rest} />;
    case "cloudwatch":
      return <AWSCloudwatch {...rest} />;
    case "google":
      return <GoogleIcon {...rest} />;
    case "1pass":
      return <OnePasswordIcon {...rest} />;
    case "datadog":
      return <DataDog {...rest} />;
    case "github":
      return <GitHubIcon {...rest} />;
    case "snowflake":
      return <SnowflakeIcon {...rest} />;
    case "gcp":
      return <GoogleCloudIcon {...rest} />;
  }

  return <CommonFateIcon {...rest} />;
};
