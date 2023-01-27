import { ChakraProvider } from "@chakra-ui/react";
import React from "react";
import { ErrorBoundary } from "react-error-boundary";
import { CognitoProvider } from "../utils/context/cognitoContext";
import { UserProvider } from "../utils/context/userContext";
import UnhandledError from "../utils/errPage";
import { theme } from "../utils/theme";

export default function App({ children }: { children: React.ReactNode }) {
  return (
    <ChakraProvider theme={theme}>
      <ErrorBoundary
        fallbackRender={({ error }) => <UnhandledError error={error} />}
      >
        <CognitoProvider>
          <UserProvider>{children}</UserProvider>
        </CognitoProvider>
      </ErrorBoundary>
    </ChakraProvider>
  );
}
