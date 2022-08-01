import { Auth } from "@aws-amplify/auth";
import { Amplify, Hub, HubCallback, ICredentials } from "@aws-amplify/core";
import { Center } from "@chakra-ui/layout";
import React, { useEffect, useState } from "react";
import NoUser from "../../pages/noUserPage";
import CFSpinner from "../../pages/CFSpinner";
import awsExports from "../aws-exports";
import { getMe } from "../backend-client/end-user/end-user";
import { User } from "../backend-client/types";
import { setAPIURL } from "../custom-instance";
import { createCtx } from "./createCtx";

export interface UserContextProps {
  user?: User;
  initiateAuth: () => Promise<ICredentials>;
  initiateSignOut: () => Promise<any>;
  isAdmin: boolean | undefined;
}

const [useUser, UserContextProvider] = createCtx<UserContextProps>();

interface Props {
  children: React.ReactNode;
}

const UserProvider: React.FC<Props> = ({ children }) => {
  const [user, setUser] = useState<User>();
  const [amplifyUser, setamplifyUser] = useState<any>();
  const [loading, setLoading] = useState<boolean>(true);
  const [amplifyLoggedIn, setamplifyLoggedIn] = useState<boolean>(false);

  const [initialized, setInitialized] = useState(false);
  const [isAdmin, setIsAdmin] = useState<boolean>();

  async function getUser() {
    if (!initialized) {
      console.debug("userContext: not initialized, skipping getUser");
      return;
    }

    // await Auth.currentSession();

    const me = await getMe();
    console.log(me);
    if (me != null) {
      setUser(me.user);
      setIsAdmin(me.isAdmin);
    }
    console.debug({ msg: "getMe response", me });
  }

  const amplifyListener: HubCallback = async ({ payload: { event, data } }) => {
    console.debug("aws-amplify Hub recieved event", { event, data });
    switch (event) {
      case "signIn":
        setamplifyLoggedIn(true);
        setamplifyUser(data);
      case "cognitoHostedUI":
      case "signOut":
        setUser(undefined);
        //setLoading(false);
        //setamplifyLoggedIn(false);
        break;
      case "signIn_failure":

      case "cognitoHostedUI_failure":
        break;
    }
  };

  useEffect(() => {
    Hub.listen("auth", amplifyListener);
    return () => Hub.remove("auth", amplifyListener);
  }, []);

  useEffect(() => {
    if (initialized) {
      getUser()
        .then(() => {
          if (user !== undefined) {
            setLoading(false);
          }
        })
        .catch(() => {
          setLoading(false);
        });
    }
  }, [initialized]);

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
      setInitialized(true);
    } else {
      console.debug("using fetch to get aws-exports.json");
      const awsConfigRequestHeaders = new Headers();
      awsConfigRequestHeaders.append("pragma", "no-cache");
      awsConfigRequestHeaders.append("cache-control", "no-cache");
      fetch("/aws-exports.json", {
        headers: awsConfigRequestHeaders,
        method: "GET",
      }).then((r) =>
        r.json().then((j) => {
          Amplify.configure(j);
          const apiURL = j.API.endpoints[0]?.endpoint;
          if (apiURL == null) {
            console.error("could not load API URL");
          } else {
            setAPIURL(apiURL);
          }
          setInitialized(true);
        })
      );
    }
  }, []);

  if (loading && user === undefined) {
    return (
      <Center h="100vh">
        <CFSpinner />
      </Center>
    );
  }
  console.log(amplifyLoggedIn, user, loading);
  if (amplifyLoggedIn && user === undefined && !loading) {
    return (
      <Center h="100vh">
        <NoUser
          userEmail={amplifyUser.username}
          initiateSignOut={initiateSignOut}
        />
      </Center>
    );
  }

  if (user === undefined && !loading) {
    initiateAuth();
  }

  if (window.location.pathname.startsWith("/admin") && !isAdmin) {
    return <>Sorry, you don&apos;t have access</>;
  }

  return (
    <UserContextProvider
      value={{
        user,
        initiateAuth,
        initiateSignOut,
        isAdmin,
      }}
    >
      {children}
    </UserContextProvider>
  );
};

function initiateAuth() {
  return Auth.federatedSignIn();
}

function initiateSignOut() {
  return Auth.signOut();
}

export { useUser, UserProvider };
