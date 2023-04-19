import {
  Box,
  Center,
  chakra,
  Container,
  Flex,
  Input,
  Stack,
  useDisclosure,
  useEventListener,
} from "@chakra-ui/react";
import { Command } from "cmdk";
import React from "react";
import { UserLayout } from "../components/Layout";

const search = () => {
  const modal = useDisclosure();

  // https://erikmartinjordan.com/navigator-platform-deprecated-alternative
  const isMac = () =>
    /(Mac|iPhone|iPod|iPad)/i.test(
      // @ts-ignore
      navigator?.userAgentData?.platform || navigator?.platform || "unknown"
    );

  const ACTION_KEY_DEFAULT = ["Ctrl", "Control"];
  const ACTION_KEY_APPLE = ["⌘", "Command"];
  const [actionKey, setActionKey] = React.useState<string[]>(ACTION_KEY_APPLE);

  React.useEffect(() => {
    if (typeof navigator === "undefined") return;
    if (!isMac()) {
      setActionKey(ACTION_KEY_DEFAULT);
    }
  }, []);

  useEventListener("keydown", (event) => {
    const hotkey = isMac() ? "metaKey" : "ctrlKey";
    if (event?.key?.toLowerCase() === "k" && event[hotkey]) {
      event.preventDefault();
      modal.isOpen ? modal.onClose() : modal.onOpen();
    }
  });

  const ProviderObjExample = {
    aws: [
      // five example aws accounts differnet values (account number, account name, account alias)
      {
        accountNumber: "0123456789012",
        accountName: "Cloud Watch Logs",
        accountAlias: "cloudwatchlogs",
      },
      {
        // new one that isn't cloud watch
        accountNumber: "0123456789012",
      },
    ],
    okta: [{}],
  };

  return (
    <UserLayout>
      <Container mt={24}>
        <Stack spacing={4}>
          {/* <Input
            size="lg"
            type="text"
            placeholder="What do you want to access?"
          /> */}
          <Command
            // open={modal.isOpen}
            // onOpenChange={modal.onToggle}
            label="Global Command Menu"
          >
            <Input
              size="lg"
              type="text"
              placeholder="What do you want to access?"
              as={Command.Input}
            />
            {/* <ChakraInput /> */}
            <Stack as={Command.List} spacing={4}>
              <Center
                as={Command.Empty}
                minH="200px"
                border="1px solid"
                rounded="md"
                borderColor="neutrals.300"
              >
                No results found.
              </Center>

              <Command.Group heading="Letters">
                <Command.Item>
                  a<Box>z</Box>
                </Command.Item>
                <Command.Item>Cloud Watch Logs (0123456789012)</Command.Item>
                <Command.Separator />
                <Command.Item>Cloud Watch Logs (0123456789012)</Command.Item>
                <Command.Item>Okta (0123456789012)</Command.Item>
                <Command.Item>AWS (0123456789012)</Command.Item>
              </Command.Group>

              <Command.Item>Apple</Command.Item>
            </Stack>
          </Command>
          <Flex>{JSON.stringify(modal)}</Flex>
        </Stack>
      </Container>
    </UserLayout>
  );
};

const ChakraInput = chakra(Command.Input);
const ChakraList = chakra(Command.List);
const ChakraEmpty = chakra(Command.Empty);
const ChakraGroup = chakra(Command.Group);
const ChakraItem = chakra(Command.Item);
const ChakraSeparator = chakra(Command.Separator);

export default search;
