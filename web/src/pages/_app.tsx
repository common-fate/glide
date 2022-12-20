import {
  ChakraProvider,
  useDisclosure,
  useEventListener,
} from "@chakra-ui/react";
import React from "react";
import CommandPalette from "../components/CommandPalette";
import CommandPalette2 from "../components/CommandPalette2";
import { CognitoProvider } from "../utils/context/cognitoContext";
import { UserProvider } from "../utils/context/userContext";
import ErrorBoundary from "../utils/errorBoundary";
import { theme } from "../utils/theme";

export default function App({ children }: { children: React.ReactNode }) {
  const modal = useDisclosure();

  useEventListener("keydown", (event) => {
    const isMac = /(Mac|iPhone|iPod|iPad)/i.test(navigator?.platform);
    const hotkey = isMac ? "metaKey" : "ctrlKey";
    if (event?.key?.toLowerCase() === "k" && event[hotkey]) {
      event.preventDefault();
      modal.isOpen ? modal.onClose() : modal.onOpen();
    }
  });

  return (
    <ChakraProvider theme={theme}>
      <ErrorBoundary>
        <CognitoProvider>
          <UserProvider>{children}</UserProvider>
          <CommandPalette isOpen={modal.isOpen} onClose={modal.onClose} />
          {/* <CommandPalette2 isOpen={modal.isOpen} onClose={modal.onClose} /> */}
        </CognitoProvider>
      </ErrorBoundary>
    </ChakraProvider>
  );
}
