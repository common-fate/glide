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
      return <GrantedKeysIcon h={Props.size} w="auto" />;
    case "cognito":
      return <CognitoLogo h={Props.size} w="auto" />;
    case "azure":
      return <AzureIcon h={Props.size} w="auto" />;
    case "okta":
      return <OktaIcon h={Props.size} w="auto" />;
    case "aws-sso":
      return <AWSIcon h={Props.size} w="auto" />;
    case "google":
      return <GoogleLogo h={Props.size} w="auto" />;
    case "one-login":
      //wide rectangular logos require being halfed to fit the page better
      return <OneLoginLogo h={Props.size / 2} w="auto" />;

    default:
      break;
  }
};
