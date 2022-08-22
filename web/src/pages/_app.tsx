import { ChakraProvider } from "@chakra-ui/react";
import React from "react";
import { CognitoProvider } from "../utils/context/cognitoContext";
import { UserProvider } from "../utils/context/userContext";
import { theme } from "../utils/theme";

export default function App({ children }: { children: React.ReactNode }) {
  return (
    <ChakraProvider theme={theme}>
      <CognitoProvider>
        <UserProvider>{children}</UserProvider>
      </CognitoProvider>
    </ChakraProvider>
  );
}
