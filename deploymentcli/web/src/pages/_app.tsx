import { ChakraProvider } from "@chakra-ui/react";
import React from "react";
import ErrorBoundary from "../utils/errorBoundary";
import { theme } from "../utils/theme";

export default function App({ children }: { children: React.ReactNode }) {
  return (
    <ChakraProvider theme={theme}>
      <ErrorBoundary>{children}</ErrorBoundary>
    </ChakraProvider>
  );
}
