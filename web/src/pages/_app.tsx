import {
  ChakraProvider,
  useDisclosure,
  useEventListener,
} from "@chakra-ui/react";
import React from "react";
import { useRouter } from "react-location";
import CommandPalette from "../components/CommandPalette";
import { CognitoProvider } from "../utils/context/cognitoContext";
import { UserProvider } from "../utils/context/userContext";
import ErrorBoundary from "../utils/errorBoundary";
import { theme } from "../utils/theme";

export default function App({ children }: { children: React.ReactNode }) {
  const modal = useDisclosure();

  const router = useRouter();

  const endUser = !router.state.location.href.includes("/admin");

  useEventListener("keydown", (event) => {
    const isMac = /(Mac|iPhone|iPod|iPad)/i.test(navigator?.platform);
    const hotkey = isMac ? "metaKey" : "ctrlKey";
    if (event?.key?.toLowerCase() === "k" && event[hotkey]) {
      event.preventDefault();
      if (endUser) {
        modal.isOpen ? modal.onClose() : modal.onOpen();
      }
    }
  });

  return (
    <ChakraProvider theme={theme}>
      <ErrorBoundary>
        <CognitoProvider>
          <UserProvider>{children}</UserProvider>
          {endUser && (
            <CommandPalette isOpen={modal.isOpen} onClose={modal.onClose} />
          )}
        </CognitoProvider>
      </ErrorBoundary>
    </ChakraProvider>
  );
}
