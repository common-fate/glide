import { Center } from "@chakra-ui/layout";
import React, { useEffect, useState } from "react";
import CFSpinner from "../../pages/CFSpinner";
import NoUser from "../../pages/noUserPage";
import { getMe } from "../backend-client/end-user/end-user";
import { User } from "../backend-client/types";
import { createCtx } from "./createCtx";

export interface UserContextProps {
  user?: User;
  isAdmin?: boolean;
}

const [useUser, UserContextProvider] = createCtx<UserContextProps>();

interface Props {
  children: React.ReactNode;
}

const UserProvider: React.FC<Props> = ({ children }) => {
  const [loadingMe, setLoadingMe] = useState<boolean>(true);
  const [user, setUser] = useState<User>();
  const [isAdmin, setIsAdmin] = useState<boolean>();

  useEffect(() => {
    setLoadingMe(true);
    getMe().then((u) => {
      if (u) {
        setUser(u.user);
        setIsAdmin(u.isAdmin);
        setLoadingMe(false);
      } else {
        setUser(undefined);
        setIsAdmin(undefined);
        setLoadingMe(false);
      }
    });
  }, []);

  if (loadingMe && user === undefined) {
    return (
      <Center h="100vh">
        <CFSpinner />
      </Center>
    );
  }

  // if loading has finished, and there is not user, report that something went wrong
  if (!loadingMe && user === undefined) {
    return (
      <Center h="100vh">
        <NoUser />
      </Center>
    );
  }

  if (window.location.pathname.startsWith("/admin") && !isAdmin) {
    return <>Sorry, you don&apos;t have access</>;
  }

  return (
    <UserContextProvider
      value={{
        user,

        isAdmin,
      }}
    >
      {children}
    </UserContextProvider>
  );
};

export { useUser, UserProvider };
