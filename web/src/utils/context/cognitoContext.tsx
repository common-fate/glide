import { Auth } from "@aws-amplify/auth";
import { CognitoHostedUIIdentityProvider } from "@aws-amplify/auth/lib-esm/types";
import { Amplify, Hub, HubCallback, ICredentials } from "@aws-amplify/core";
import { Center } from "@chakra-ui/react";
import React, { useEffect, useState } from "react";
import { MakeGenerics, useNavigate, useSearch } from "react-location";
import CFSpinner from "../../components/CFSpinner";
import awsExports from "../aws-exports";
import { AccessRuleStatus } from "../backend-client/types";
import { setAPIURL } from "../custom-instance";
import { createCtx } from "./createCtx";

export interface CognitoContextProps {
  cognitoAuthenticatedUserEmail: string | undefined;
  initiateAuth: () => Promise<ICredentials>;
  initiateSignOut: () => Promise<any>;
}

const [useCognito, CognitoContextProvider] = createCtx<CognitoContextProps>();

interface Props {
  children: React.ReactNode;
}

type MyLocationGenerics = MakeGenerics<{
  Search: {
    state?: string;
  };
}>;

const CognitoProvider: React.FC<Props> = ({ children }) => {
  const [amplifyInitialising, setAmplifyInitializing] = useState(true);
  const [loadingCurrentUser, setLoadingCurrentUser] = useState(true);
  const [loggingOut, setLoggingOut] = useState(false);
  const [cognitoAuthenticatedUserEmail, setCognitoAuthenticatedUserEmail] =
    useState<string>();
  const loading = amplifyInitialising || loadingCurrentUser;
  const navigate = useNavigate();
  const search = useSearch<MyLocationGenerics>();
  // this can be improved in future with a more graceful error page if the AWS config doesn't load.
  // The following effect will run on first load of the app, in production, this will fetch a config file from the server to hydrate the amplify configuration
  // in local dev, this is imported from a local file
  useEffect(() => {
    if (window.location.hostname === "localhost") {
      console.debug({ localExports: awsExports });
      Amplify.configure(awsExports);
      const apiURL = (awsExports as any).API.endpoints[0]?.endpoint;
      if (apiURL == null) {
        console.error("could not load API URL");
      } else {
        setAPIURL(apiURL);
      }
      setAmplifyInitializing(false);
    } else {
      console.debug("using fetch to get aws-exports.json");
      const awsConfigRequestHeaders = new Headers();
      awsConfigRequestHeaders.append("pragma", "no-cache");
      awsConfigRequestHeaders.append("cache-control", "no-cache");
      fetch("/aws-exports.json", {
        headers: awsConfigRequestHeaders,
        method: "GET",
      })
        .then((r) =>
          r.json().then((j) => {
            Amplify.configure(j);
            const apiURL = j.API.endpoints[0]?.endpoint;
            if (apiURL == null) {
              console.error("could not load API URL");
            } else {
              setAPIURL(apiURL);
            }
            setAmplifyInitializing(false);
          })
        )
        .catch((e) => console.error(e));
    }
  }, []);

  const tryGetAuthenticatedUser = () => {
    setLoadingCurrentUser(true);
    Auth.currentAuthenticatedUser()
      .then((data) => {
        console.debug("got current authenticated user", data);
        setCognitoAuthenticatedUserEmail(data.username);
        setLoadingCurrentUser(false);
      })
      .catch(() => {
        console.log("couldn't find current authenticated user");

        setCognitoAuthenticatedUserEmail(undefined);
        setLoadingCurrentUser(false);
      });
  };

  const amplifyListener: HubCallback = async ({ payload: { event, data } }) => {
    console.log("aws-amplify Hub recieved event", { event, data });
    switch (event) {
      case "oAuthSignOut":
        setCognitoAuthenticatedUserEmail(undefined);
        break;
      case "signOut":
        setCognitoAuthenticatedUserEmail(undefined);
        break;
      case "customOAuthState":
        navigate({ to: data });
        break;
      default:
        console.log("getting user in listener", { data });
        tryGetAuthenticatedUser();
    }
  };

  useEffect(() => {
    if (!amplifyInitialising) {
      tryGetAuthenticatedUser();
      console.debug("starting hub listener");
      Hub.listen("auth", amplifyListener);
      return () => Hub.remove("auth", amplifyListener);
    }
  }, [amplifyInitialising]);

  useEffect(() => {
    if (loggingOut && search.state === "loggedOut") {
      setLoggingOut(false);
    }
  }, [search, loggingOut]);

  // spinner when amplify is initialising or when the current user is being fetched and the user is undefined
  if (loading && cognitoAuthenticatedUserEmail === undefined) {
    return (
      <Center h="100vh">
        <CFSpinner />
      </Center>
    );
  }

  // force the ts type for cognitoAuthenticatedUserEmail to be a string in the context return by expricitly checking it
  if (!loading && cognitoAuthenticatedUserEmail === undefined) {
    if (!loggingOut) {
      initiateAuth().catch((e) => console.error(e));
    }
    return (
      <Center h="100vh">
        <CFSpinner />
      </Center>
    );
  }

  function initiateAuth() {
    // prevent unexpected redirects back to login page when redirecting to signin page after logout
    const customState =
      location.pathname === "/logout"
        ? undefined
        : location.pathname + location.search;
    return Auth.federatedSignIn({
      customState: customState,
      provider: CognitoHostedUIIdentityProvider.Cognito,
    });
  }

  function initiateSignOut() {
    setLoggingOut(true);
    return Auth.signOut();
  }

  return (
    <CognitoContextProvider
      value={{
        cognitoAuthenticatedUserEmail,
        initiateAuth,
        initiateSignOut,
      }}
    >
      {children}
    </CognitoContextProvider>
  );
};

export { useCognito, CognitoProvider };
