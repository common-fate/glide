import { ChakraProvider } from "@chakra-ui/react";
import React from "react";
import ErrorBoundary from "../utils/errorBoundary";
import { CognitoProvider } from "../utils/context/cognitoContext";
import { UserProvider } from "../utils/context/userContext";
import { theme } from "../utils/theme";

export default function App({ children }: { children: React.ReactNode }) {
  return (
    <ChakraProvider theme={theme}>
      <ErrorBoundary>
        <CognitoProvider>
          <UserProvider>{children}</UserProvider>
        </CognitoProvider>
      </ErrorBoundary>
    </ChakraProvider>
  );
}
