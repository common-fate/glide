import {
  Box,
  Button,
  Center,
  chakra,
  Code,
  Container,
  Flex,
  Input,
  Stack,
  useBoolean,
  useDisclosure,
  useEventListener,
} from "@chakra-ui/react";
import { Command } from "cmdk";
import { Command as CommandNew } from "../utils/cmdk";
import React from "react";
import { ProviderIcon, ShortTypes } from "../components/icons/providerIcon";
import { UserLayout } from "../components/Layout";
import {
  userPostRequests,
  useUserListEntitlements,
  useUserListEntitlementTargets,
} from "../utils/backend-client/default/default";
import { CheckCircleIcon } from "@chakra-ui/icons";
import { useRouter } from "react-location";
import Counter from "../components/Counter";

const search = () => {
  const modal = useDisclosure();

  // https://erikmartinjordan.com/navigator-platform-deprecated-alternative
  const isMac = () =>
    /(Mac|iPhone|iPod|iPad)/i.test(
      // @ts-ignore
      navigator?.userAgentData?.platform || navigator?.platform || "unknown"
    );

  // const ACTION_KEY_DEFAULT = ["Ctrl", "Control"];
  // const ACTION_KEY_APPLE = ["⌘", "Command"];
  // const [actionKey, setActionKey] = React.useState<string[]>(ACTION_KEY_APPLE);

  // React.useEffect(() => {
  //   if (typeof navigator === "undefined") return;
  //   if (!isMac()) {
  //     setActionKey(ACTION_KEY_DEFAULT);
  //   }
  // }, []);

  // useEventListener("keydown", (event) => {
  //   const hotkey = isMac() ? "metaKey" : "ctrlKey";
  //   if (event?.key?.toLowerCase() === "k" && event[hotkey]) {
  //     event.preventDefault();
  //     modal.isOpen ? modal.onClose() : modal.onOpen();
  //   }
  // });

  // Watch keys for cmd Enter submit
  useEventListener("keydown", (event) => {
    const hotkey = isMac() ? "metaKey" : "ctrlKey";
    if (event?.key?.toLowerCase() === "Enter" && event[hotkey]) {
      event.preventDefault();
      handleSubmit();
    }
  });

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

  const [targetKeyMap, setTargetKeyMap] = React.useState<{
    [key: string]: Target;
  }>({});

  const [inputValue, setInputValue] = React.useState<string>("");

  const [checked, setChecked] = React.useState<string[]>([]);

  React.useEffect(() => {
    if (!targets.data) return;
    const map: { [key: string]: Target } = {};
    targets.data.targets.forEach((target) => {
      const key =
        target.targetGroupFrom.name +
        " " +
        target.fields.map((f) => f.value).join(" ");
      map[key] = target;
    });
    setTargetKeyMap(map);
  }, [targets.data]);

  // @TODO:
  // Actually use the fixture data, maybe write it with actual values.
  // Add in page responses etc.

  const router = useRouter();

  const [submitLoading, submitLoadingToggle] = useBoolean();

  const handleSubmit = () => {
    submitLoadingToggle.on();

    let entitlementTargerts = [];

    // map over keys and get the ids of the targets
    checked.forEach((key) => {
      const target = targetKeyMap[key];
      entitlementTargerts.push(target.id);
    });

    // preflight request
    router.push();
  };

  return (
    <UserLayout>
      <Container mt={24}>
        <Box spacing={4} minH="200px">
          <CommandNew
            // open={modal.isOpen}
            // onOpenChange={modal.onToggle}
            label="Global Command Menu"
            checked={checked}
            setChecked={setChecked}
          >
            <Input
              size="lg"
              type="text"
              placeholder="What do you want to access?"
              value={inputValue}
              onValueChange={setInputValue}
              as={CommandNew.Input}
            />
            <Stack mt={2} spacing={2} direction="row">
              <Center
                boxSize="90px"
                rounded="md"
                borderColor="neutrals.300"
                bg="white"
                borderWidth="1px"
                flexDir="column"
                onClick={() => {
                  setInputValue("");
                }}
              >
                <Counter count={checked.length} />
                All resources
                <Flex>{targets.data?.targets.length}&nbsp;total</Flex>
              </Center>
              {["aws", "cloudwatch", "okta", "azure", "github", "gcp"].map(
                (key) => {
                  return (
                    <Center
                      boxSize="90px"
                      rounded="md"
                      borderColor="neutrals.300"
                      bg="white"
                      borderWidth="1px"
                      flexDir="column"
                      onClick={() => {
                        setInputValue(key);
                      }}
                    >
                      <ProviderIcon shortType={key} />
                      {key}
                    </Center>
                  );
                }
              )}
            </Stack>
            <CommandNew.List>
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
                <Center as={CommandNew.Empty} minH="200px">
                  No results found.
                </Center>

                <CommandNew.Group
                // heading="Permissions"
                >
                  <CommandNew.Separator />
                  {targets.data &&
                    Object.entries(targetKeyMap).map(([key, target]) => {
                      return (
                        <Flex
                          alignContent="flex-start"
                          p={2}
                          rounded="md"
                          _selected={{
                            bg: "neutrals.100",
                          }}
                          _checked={{
                            "#checked": {
                              display: "block",
                            },
                          }}
                          pos="relative"
                          key={key}
                          value={key}
                          as={CommandNew.Item}
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
                          <CheckCircleIcon
                            id="checked"
                            position="absolute"
                            display="none"
                            top={2}
                            right={2}
                            h="12px"
                            w="12px"
                            color={"brandBlue.300"}
                          />
                        </Flex>
                      );
                    })}
                </CommandNew.Group>
              </Stack>
            </CommandNew.List>
          </CommandNew>
          <Flex w="100%" mt={4}>
            <Button
              disabled={checked.length == 0}
              ml="auto"
              onClick={handleSubmit}
              isLoading={submitLoading}
              loadingText="Processing request..."
            >
              Next (⌘+Enter)
            </Button>
          </Flex>
        </Box>

        {/* <Code bg="gray.50" whiteSpace="pre-wrap">
          {JSON.stringify({ targets }, null, 2)}
        </Code> */}
      </Container>
    </UserLayout>
  );
};

export default search;
