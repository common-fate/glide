import {
  GrantedKeysIcon,
  AzureIcon,
  OktaIcon,
  AWSIcon,
} from "../components/icons/Icons";
import {
  CognitoLogo,
  GoogleLogo,
  OneLoginLogo,
} from "../components/icons/Logos";

type IdpLogoProps = {
  idpType: string;
  size: number;
};

export const GetIDPLogo = (Props: IdpLogoProps) => {
  switch (Props.idpType) {
    case "internal":
      return <GrantedKeysIcon boxSize={Props.size} />;
    case "cognito":
      return <CognitoLogo boxSize={Props.size} />;
    case "azure":
      return <AzureIcon boxSize={Props.size} />;
    case "okta":
      return <OktaIcon boxSize={Props.size} />;
    case "aws-sso":
      return <AWSIcon boxSize={Props.size} />;
    case "google":
      return <GoogleLogo boxSize={Props.size} />;
    case "one-login":
      return <OneLoginLogo boxSize={Props.size} />;

    default:
      break;
  }
};
