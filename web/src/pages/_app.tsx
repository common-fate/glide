import { ChakraProvider } from "@chakra-ui/react";
import React from "react";
import { UserProvider } from "../utils/context/userContext";
import { theme } from "../utils/theme";

export default function App({ children }: { children: React.ReactNode }) {
  return (
    <ChakraProvider theme={theme}>
      <UserProvider>{children}</UserProvider>
    </ChakraProvider>
  );
}
