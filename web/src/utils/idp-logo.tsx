import { IconProps } from "@chakra-ui/icons";
import { AzureIcon, OktaIcon, AWSIcon } from "../components/icons/Icons";
import {
  CognitoLogo,
  CommonFateIcon,
  CommonFateLogo,
  GoogleLogo,
  OneLoginLogo,
} from "../components/icons/Logos";

type IdpLogoProps = {
  idpType: string;
} & IconProps;

export const IDPLogo: React.FC<IdpLogoProps> = ({ idpType, ...rest }) => {
  switch (idpType) {
    case "internal":
      return <CommonFateIcon {...rest} />;
    case "cognito":
      return <CognitoLogo {...rest} />;
    case "azure":
      return <AzureIcon {...rest} />;
    case "okta":
      return <OktaIcon {...rest} />;
    case "aws-sso":
      return <AWSIcon {...rest} />;
    case "google":
      return <GoogleLogo {...rest} />;
    case "one-login":
      return <OneLoginLogo {...rest} />;
    default:
      return null;
  }
};

export const GetIDPName = (idpType: string): string => {
  switch (idpType) {
    case "internal":
      return "Internal";
    case "cognito":
      return "Cognito";
    case "azure":
      return "Azure AD";
    case "okta":
      return "Okta";
    case "aws-sso":
      return "AWS SSO";
    case "google":
      return "Google Workspace";
    case "one-login":
      return "One Login";
    default:
      return idpType;
  }
};
