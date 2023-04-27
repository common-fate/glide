import { ArrowBackIcon, CheckCircleIcon, SettingsIcon } from "@chakra-ui/icons";
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
  chakra,
  Textarea,
  Spinner,
  Divider,
  Icon,
  Tooltip,
  HStack,
} from "@chakra-ui/react";
import React from "react";
import { Link, useNavigate } from "react-location";
import Counter from "../components/Counter";
import FieldsCodeBlock from "../components/FieldsCodeBlock";
import { ProviderIcon, ShortTypes } from "../components/icons/providerIcon";
import { UserLayout } from "../components/Layout";
import {
  userPostRequests,
  userRequestPreflight,
  useUserListEntitlements,
  useUserListEntitlementTargets,
} from "../utils/backend-client/default/default";
import { Preflight, Target } from "../utils/backend-client/types";
import { Command as CommandNew } from "../utils/cmdk";
import { HeaderStatusCell } from "./request/[id]";
function TargetComponent({ target }: { target: Target }) {
  return (
    <Flex>
      <ProviderIcon shortType={target.kind.icon as ShortTypes} />
      <HStack>
        {target.fields.map((field, i) => (
          <Box borderLeftWidth={i > 0 ? 2 : undefined} paddingRight={10}>
            <Tooltip
              key={field.id}
              label={
                <>
                  <Box fontWeight="bold">{field.fieldTitle}</Box>
                  {field.fieldDescription && (
                    <Box>{field.fieldDescription}</Box>
                  )}
                  {field.valueDescription && (
                    <Box mt={2}>{field.valueDescription}</Box>
                  )}
                </>
              }
              placement="top"
            >
              <Box display="inline-block" verticalAlign="top">
                <Box>{field.valueLabel}</Box>
                <Box>{field.value}</Box>
              </Box>
            </Tooltip>
          </Box>
        ))}
      </HStack>
    </Flex>
  );
}
const Search = () => {
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

  const targets = useUserListEntitlementTargets(undefined, {
    swr: { refreshInterval: 10000 },
    // request: {
    //   baseURL: "http://127.0.0.1:3100",
    //   headers: {
    //     Prefer: "code=200, example=example_targets",
    //   },
    // },
  });
  const entitlements = useUserListEntitlements({
    swr: { refreshInterval: 10000 },
    // request: {
    //   baseURL: "http://127.0.0.1:3100",
    //   headers: {
    //     Prefer: "code=200, example=example_targets",
    //   },
    // },
  });
  const [targetKeyMap, setTargetKeyMap] = React.useState<{
    [key: string]: Target;
  }>({});

  const [inputValue, setInputValue] = React.useState<string>("");

  const [checked, setChecked] = React.useState<string[]>([]);

  React.useEffect(() => {
    if (!targets.data) return;
    const map: { [key: string]: Target } = {};
    targets.data.targets.forEach((target) => {
      map[target.id] = target;
    });
    setTargetKeyMap(map);
  }, [targets.data]);

  const targetsLoading =
    targets.isValidating && Object.entries(targetKeyMap).length == 0;

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
    if (tabIndex == 0) handlePreflight();
    if (tabIndex == 1) handleRequest();
  };

  const handlePreflight = () => {
    submitLoadingToggle.on();

    const entitlementTargets: string[] = [];

    // map over keys and get the ids of the targets
    checked.forEach((key) => {
      const target = targetKeyMap[key];
      target && entitlementTargets.push(target?.id);
    });

    userRequestPreflight(
      {
        targets: entitlementTargets,
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

  const handleRequest = () => {
    preflightRes &&
      // test
      userPostRequests(
        {
          preflightId: preflightRes?.id,
          reason: accessReason,
          groupOptions: [
            // @TODO: add in group options, after dynamic UI is supported (next scope of work)
            //   {
            //     id: preflightRes?.id,
            //     timing: { startTime: new Date().toISOString() },
            //   },
          ],
        },
        {
          baseURL: "http://127.0.0.1:3100",
          headers: {
            Prefer: "code=200, example=ex_1",
          },
        }
      )
        .then((res) => {
          console.log(res);
          submitLoadingToggle.on();
          // redirect to request...
          navigate({ to: `/request/${"REPLACEME"}` });
          // clear state
          setChecked([]);
          setInputValue("");
        })
        .catch((err) => {
          console.log(err);
        });
  };

  const [accessReason, setAccessReason] = React.useState<string>("");

  return (
    <UserLayout>
      <Container
        mt={24}
        maxW={{
          md: tabIndex == 0 ? "container.lg" : "container.sm",
        }}
      >
        {/* set index */}
        <Tabs index={tabIndex}>
          <TabPanels>
            <TabPanel>
              <Box minH="200px">
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
                    autoFocus={true}
                    as={CommandNew.Input}
                  />
                  <Flex mt={2} direction="row" overflowX="scroll">
                    <Center
                      // boxSize="90px"
                      rounded="md"
                      // w="90px !important"
                      borderColor="neutrals.300"
                      bg="white"
                      borderWidth="1px"
                      flexDir="column"
                      textStyle="Body/Small"
                      onClick={() => {
                        setInputValue("");
                      }}
                      px={2}
                      mr={2}
                      as="button"
                    >
                      <Counter size="md" count={checked.length} />
                      <Text
                        textStyle="Body/Small"
                        noOfLines={1}
                        textOverflow="clip"
                        w="90px"
                        textAlign="center"
                      >
                        All resources
                      </Text>
                      <Flex color="neutrals.500">
                        {targets.data?.targets.length}&nbsp;total
                      </Flex>
                    </Center>
                    {entitlements.data?.entitlements.map((kind) => {
                      const key =
                        kind.publisher +
                        "#" +
                        kind.name +
                        "#" +
                        kind.kind +
                        "#";
                      return (
                        <Center
                          boxSize="90px"
                          rounded="md"
                          borderColor="neutrals.300"
                          bg="white"
                          borderWidth="1px"
                          textStyle="Body/Small"
                          flexDir="column"
                          onClick={() => {
                            setInputValue(key);
                            // then set the focus back to the input
                            // so that the user can continue typing
                            document.getElementById(":rd:")?.focus();
                          }}
                          px={8}
                          mr={2}
                          as="button"
                        >
                          <ProviderIcon shortType={kind.icon as ShortTypes} />
                          {kind.kind}
                        </Center>
                      );
                    })}
                  </Flex>
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
                        {targetsLoading ? (
                          <Center as={CommandNew.Loading} minH="200px">
                            <Spinner />
                          </Center>
                        ) : (
                          targets.data &&
                          targets.data.targets.map((target) => {
                            return (
                              <Flex
                                alignContent="flex-start"
                                p={2}
                                rounded="md"
                                _selected={{
                                  "bg": "neutrals.100",
                                  "#description": {
                                    display: "block",
                                  },
                                }}
                                _checked={{
                                  "#checked": {
                                    display: "block",
                                  },
                                }}
                                pos="relative"
                                key={target.id}
                                // this value is used by the command palette
                                value={target.id}
                                as={CommandNew.Item}
                              >
                                <TargetComponent target={target} />
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
                                {/* @TODO: review me as a part of CF-1028 */}
                                {/* <Box
                                  rounded="md"
                                  w="24ch"
                                  zIndex={9999}
                                  pos="absolute"
                                  top={4}
                                  right={4}
                                  id="description"
                                  display="none"
                                  textStyle="Body/ExtraSmall"
                                  p={1}
                                >
                                  Admin access to {target.fields[0].value}{" "}
                                  account
                                </Box> */}
                                <Box
                                  rounded="md"
                                  w="24ch"
                                  bg="white"
                                  border="1px solid"
                                  borderColor="neutrals.300"
                                  zIndex={9999}
                                  pos="absolute"
                                  bottom={-4}
                                  right={0}
                                  id="description"
                                  display="none"
                                  textStyle="Body/ExtraSmall"
                                  p={1}
                                >
                                  Admin access to {target.fields[0].value}{" "}
                                  account
                                </Box>
                              </Flex>
                            );
                          })
                        )}
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

              <Flex mt={12}>
                <Button
                  leftIcon={<SettingsIcon />}
                  size="xs"
                  ml="auto"
                  variant="brandSecondary"
                >
                  Settings
                </Button>
              </Flex>
            </TabPanel>
            {/* PAGE 2 */}
            <TabPanel>
              {/* main */}
              <Text textStyle="Body/Medium">Access</Text>
              <Stack spacing={2} w="100%">
                {preflightRes?.accessGroups.map((group) => {
                  return (
                    <Box
                      p={2}
                      w="100%"
                      borderColor="neutrals.300"
                      borderWidth="1px"
                      rounded="md"
                    >
                      {/* <HeaderStatusCell group={group} /> */}
                      <Stack spacing={2}>
                        {group.targets.map((target) => {
                          return (
                            <Flex
                              p={2}
                              borderColor="neutrals.300"
                              borderWidth="1px"
                              rounded="md"
                              flexDir="row"
                            >
                              <ProviderIcon
                                shortType={target.kind.icon as ShortTypes}
                                mr={2}
                              />
                              <FieldsCodeBlock fields={target.fields} />
                            </Flex>
                          );
                        })}
                      </Stack>
                    </Box>
                  );
                })}
              </Stack>
              <Box mt={4}>
                <Text textStyle="Body/Medium">Why do you need access?</Text>
                <Textarea
                  value={accessReason}
                  onChange={(e) => setAccessReason(e.target.value)}
                  placeholder="Deploying initial Terraform infrastructure..."
                />
              </Box>

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
                  onClick={handleSubmit}
                  isLoading={submitLoading}
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
