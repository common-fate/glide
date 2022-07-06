import { Center } from "@chakra-ui/react";
import React from "react";
import { useUser } from "../utils/context/userContext";

type Props = {};

const Logout = () => {
  const auth = useUser();
  auth.initiateSignOut();

  return <Center h="80vh">Logging you out 👋</Center>;
};

export default Logout;
