import {
  Box,
  Center,
  chakra,
  Code,
  Container,
  Flex,
  Input,
  Stack,
  useDisclosure,
  useEventListener,
} from "@chakra-ui/react";
import { Command } from "cmdk";
import React from "react";
import { ProviderIcon, ShortTypes } from "../components/icons/providerIcon";
import { UserLayout } from "../components/Layout";
import {
  useUserListEntitlements,
  useUserListEntitlementTargets,
} from "../utils/backend-client/default/default";

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

  // const entitlements = useUserListEntitlements({
  //   swr: { refreshInterval: 10000 },
  //   request: {
  //     baseURL: "http://127.0.0.1:3100",
  //     headers: {
  //       Prefer: "code=200, example=ex_1",
  //     },
  //   },
  // });

  //   I need a query param.......
  //   const resources = useUserListRequestAccessGroupGrants(group.id, {
  //     swr: { refreshInterval: 10000 },
  //     request: {
  //       baseURL: "http://127.0.0.1:3100",
  //       headers: {
  //         Prefer: "code=200, example=ex_1",
  //       },
  //     },
  //   });

  const targets = useUserListEntitlementTargets(
    {},
    {
      swr: { refreshInterval: 10000 },
      request: {
        baseURL: "http://127.0.0.1:3100",
        headers: {
          Prefer: "code=200, example=example_targets",
        },
      },
    }
  );

  // @TODO:
  // Actually use the fixture data, maybe write it with actual values.
  // Add in page responses etc.

  return (
    <UserLayout>
      <Container mt={24}>
        <Box spacing={4} minH="200px">
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
            <Command.List>
              <Stack
                // as={Command.List}
                mt={2}
                spacing={4}
                border="1px solid"
                rounded="md"
                borderColor="neutrals.300"
                p={1}
                pt={2}
              >
                <Center as={Command.Empty} minH="200px">
                  No results found.
                </Center>

                <Command.Group heading="Permissions">
                  <Command.Separator />
                  {targets.data &&
                    targets.data.targets.map((target) => {
                      return (
                        <Flex
                          alignContent="flex-start"
                          p={2}
                          rounded="md"
                          sx={{
                            " &[data-selected='true']": {
                              bg: "neutrals.100",
                            },
                          }}
                          pos="relative"
                          as={Command.Item}
                          value={target.fields
                            .map((field) => field.value)
                            .join(", ")}
                        >
                          <ProviderIcon
                            mr={2}
                            shortType={
                              target.targetGroupFrom.name as ShortTypes
                            }
                          />
                          <Box>
                            <Box>{target.fields[0].value}</Box>
                            {/* then map over the proceeding fields */}
                            <Box color="neutrals.500" minH="1em">
                              {target.fields
                                .map((field, index) => index && field.value)
                                .filter((field) => field)
                                .join(", ")}
                            </Box>
                          </Box>
                        </Flex>
                      );
                    })}
                </Command.Group>
              </Stack>
            </Command.List>
          </Command>
        </Box>

        <Flex my={12}>{JSON.stringify(modal)}</Flex>

        <Code bg="gray.50" whiteSpace="pre-wrap">
          {JSON.stringify({ targets }, null, 2)}
        </Code>
      </Container>
    </UserLayout>
  );
};

// const ChakraInput = chakra(Command.Input);
// const ChakraList = chakra(Command.List);
// const ChakraEmpty = chakra(Command.Empty);
// const ChakraGroup = chakra(Command.Group);
// const ChakraItem = chakra(Command.Item);
// const ChakraSeparator = chakra(Command.Separator);

export default search;
