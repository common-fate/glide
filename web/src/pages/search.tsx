import { ArrowBackIcon, CheckCircleIcon } from "@chakra-ui/icons";
import {
  Box,
  Button,
  Center,
  Code,
  Container,
  Flex,
  Input,
  Stack,
  Text,
  TabPanel,
  TabPanels,
  Tabs,
  useBoolean,
  useDisclosure,
  useEventListener,
} from "@chakra-ui/react";
import React from "react";
import { Link, useNavigate } from "react-location";
import Counter from "../components/Counter";
import { ProviderIcon, ShortTypes } from "../components/icons/providerIcon";
import { UserLayout } from "../components/Layout";
import {
  userRequestPreflight,
  useUserListEntitlementTargets,
} from "../utils/backend-client/default/default";
import { Preflight, Target } from "../utils/backend-client/types";
import { Command as CommandNew } from "../utils/cmdk";
import { HeaderStatusCell } from "./request/[id]";

const Search = () => {
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

  // Watch keys for cmd Enter submit
  useEventListener("keydown", (event) => {
    const hotkey = isMac() ? "metaKey" : "ctrlKey";

    if (event?.key?.toLowerCase() === "enter" && event[hotkey]) {
      event.preventDefault();
      checked.length > 0 && handleSubmit();
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

  const navigate = useNavigate();

  const [submitLoading, submitLoadingToggle] = useBoolean();

  // tabIndex
  const [tabIndex, setTabIndex] = React.useState(0);

  // preflightRes
  const [preflightRes, setPreflightRes] = React.useState<Preflight>();

  const handleSubmit = () => {
    submitLoadingToggle.on();

    let entitlementTargerts = [];

    // map over keys and get the ids of the targets
    checked.forEach((key) => {
      const target = targetKeyMap[key];
      target && entitlementTargerts.push(target?.id);
    });

    userRequestPreflight(
      {
        targets: entitlementTargerts,
      },
      {
        baseURL: "http://127.0.0.1:3100",
        headers: {
          Prefer: "code=200, example=ex_1",
        },
      }
    )
      .then((res) => {
        setPreflightRes(res);
        setTabIndex(1);
      })
      .catch((err) => {
        setTabIndex(1);
        console.log(err);
      })
      .finally(() => {
        submitLoadingToggle.off();
      });

    // navigate({ to: "/search2" });
  };

  return (
    <UserLayout>
      <Container mt={24}>
        {/* set index */}
        <Tabs index={tabIndex}>
          <TabPanels>
            <TabPanel>
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
                    {[
                      "aws",
                      "cloudwatch",
                      "okta",
                      "azure",
                      "github",
                      "gcp",
                    ].map((key) => {
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
                    })}
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
                                  // onClick={() => {
                                  //   console.log("click");
                                  //   if (checked.includes(key)) {
                                  //     setChecked(
                                  //       checked.filter(
                                  //         (checkedKey) => checkedKey !== key
                                  //       )
                                  //     );
                                  //   } else {
                                  //     setChecked([...checked, key]);
                                  //   }
                                  // }}
                                />
                                <Box>
                                  <Box>{target.fields[0].value}</Box>
                                  {/* then map over the proceeding fields */}
                                  <Box color="neutrals.500" minH="1em">
                                    {target.fields
                                      .map(
                                        (field, index) => index && field.value
                                      )
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
                    Next ({actionKey[0]}+Enter)
                  </Button>
                </Flex>
              </Box>
            </TabPanel>
            {/* PAGE 2 */}
            <TabPanel>
              {/* main */}
              <Text textStyle="Body/Medium">Access</Text>
              <Stack spacing={1} w="100%">
                {preflightRes?.accessGroups.map((group) => {
                  return (
                    <Box
                      p={2}
                      w="100%"
                      borderColor="neutrals.300"
                      borderWidth="1px"
                      rounded="md"
                    >
                      <HeaderStatusCell group={group} />
                      <Stack spacing={2}>
                        {group.targets.map((target) => {
                          return (
                            <Box
                              p={2}
                              borderColor="neutrals.300"
                              borderWidth="1px"
                              rounded="md"
                            >
                              <ProviderIcon
                                shortType={
                                  target.targetGroupFrom.name as ShortTypes
                                }
                                mr={2}
                              />
                              <Code bg="white">
                                {target.fields.map((f) => f.value).join("\n")}
                              </Code>
                              {/* {JSON.stringify(target)} */}
                            </Box>
                          );
                        })}
                      </Stack>
                    </Box>
                  );
                })}
              </Stack>

              {/* buttons */}
              <Flex w="100%" mt={4}>
                <Button
                  // ml="auto"
                  // disabled={checked.length == 0}
                  // onClick={handleSubmit}
                  variant="brandSecondary"
                  leftIcon={<ArrowBackIcon />}
                  // to="/search"
                  // as={Link}
                  onClick={() => setTabIndex(0)}
                >
                  Go back
                </Button>
                <Button
                  ml="auto"
                  // disabled={checked.length == 0}
                  // onClick={handleSubmit}
                  // isLoading={submitLoading}
                  loadingText="Processing request..."
                >
                  Next (⌘+Enter)
                </Button>
              </Flex>
            </TabPanel>
          </TabPanels>
        </Tabs>

        {/* <Code bg="gray.50" whiteSpace="pre-wrap">
          {JSON.stringify({ preflightRes, targets }, null, 2)}
        </Code> */}
      </Container>
    </UserLayout>
  );
};

export default Search;
