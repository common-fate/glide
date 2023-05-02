import { ArrowBackIcon, CheckCircleIcon, SettingsIcon } from "@chakra-ui/icons";
import {
  Box,
  Button,
  ButtonProps,
  Center,
  CenterProps,
  chakra,
  Container,
  Divider,
  Flex,
  HStack,
  Input,
  Stack,
  TabPanel,
  TabPanels,
  Tabs,
  Text,
  Textarea,
  Tooltip,
  useBoolean,
  useEventListener,
} from "@chakra-ui/react";
import { useEffect, useMemo, useState } from "react";
import { useNavigate } from "react-location";
import { FixedSizeList as List, ListChildComponentProps } from "react-window";
import Counter from "../components/Counter";
// @ts-ignore
import commandScore from "command-score";
import { ProviderIcon, ShortTypes } from "../components/icons/providerIcon";
import { UserLayout } from "../components/Layout";
import {
  userListEntitlementTargets,
  userPostRequests,
  userRequestPreflight,
  useUserListEntitlements,
} from "../utils/backend-client/default/default";
import {
  Preflight,
  Target,
  TargetField,
  UserListEntitlementTargetsParams,
} from "../utils/backend-client/types";
import debounce from "lodash.debounce";
import { Command as CommandNew } from "../utils/cmdk";
const StyledList = chakra(CommandNew.List);
// CONSTANTS
const ACTION_KEY_DEFAULT = ["Ctrl", "Control"];
const ACTION_KEY_APPLE = ["⌘", "Command"];
const TARGET_HEIGHT = 100;
const TARGETS = 5;
// https://erikmartinjordan.com/navigator-platform-deprecated-alternative
const isMac = () =>
  /(Mac|iPhone|iPod|iPad)/i.test(
    // @ts-ignore
    navigator?.userAgentData?.platform || navigator?.platform || "unknown"
  );

const Search = () => {
  // DATA FETCHING

  const entitlements = useUserListEntitlements({
    swr: { refreshInterval: 10000 },
  });

  // HOOKS
  const navigate = useNavigate();

  // STATE
  const [targetKeyMap, setTargetKeyMap] = useState<{
    [key: string]: Target;
  }>({});
  const [inputValue, setInputValue] = useState<string>("");

  const [checked, setChecked] = useState<string[]>([]);
  const [actionKey, setActionKey] = useState<string[]>(ACTION_KEY_APPLE);
  const [tabIndex, setTabIndex] = useState(0);
  const [preflightRes, setPreflightRes] = useState<Preflight>();
  const [accessReason, setAccessReason] = useState<string>("");
  const [allTargets, setAllTargets] = useState<Target[]>([]);
  const [nextToken, setNextToken] = useState<string | undefined>("initial");
  const [submitLoading, submitLoadingToggle] = useBoolean();
  const [showOnlyChecked, showOnlyCheckedToggle] = useBoolean();
  // EFFECTS

  useEffect(() => {
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

  useEffect(() => {
    const fetchData = async () => {
      const params: UserListEntitlementTargetsParams = {
        nextToken: nextToken === "initial" ? undefined : nextToken,
      };
      const result = await userListEntitlementTargets(params);
      setTargetKeyMap((tkm) => {
        result.targets.forEach((t) => {
          // the command palette library casts the value to lowercase, so we need to do the same here
          tkm[t.id.toLowerCase()] = t;
        });
        return tkm;
      });
      setNextToken(result.next);
    };
    if (nextToken !== undefined) {
      fetchData();
    } else {
      setAllTargets(Object.values(targetKeyMap));
    }
  }, [nextToken]);

  // this doesn't update if you deselect a target when in the "selected" filter.
  // this can be considered a desired effect because if you accidentally deselect something you can easily select it again.
  // if you shift back to the "all" view then it will reset this
  const filteredItems = useMemo(() => {
    if (allTargets.length === 0 || !allTargets) return [];
    if (showOnlyChecked)
      return allTargets.filter((t) => checked.includes(t.id.toLowerCase()));
    if (inputValue === "") return allTargets;
    return allTargets.filter((target) => {
      const key = target.id.toLowerCase() + targetFieldsToString(target.fields);
      return commandScore(key, inputValue.toLowerCase()) > 0;
    });
  }, [inputValue, allTargets, showOnlyChecked]);

  const targetComponent: React.FC<ListChildComponentProps> = ({
    index,
    style,
  }) => {
    const target = filteredItems[index];

    if (!target) return <></>;

    return (
      <Flex
        h={TARGET_HEIGHT}
        style={style}
        alignContent="flex-start"
        p={2}
        rounded="md"
        _selected={{
          bg: "neutrals.100",
        }}
        key={target.id}
        // this value is used by the command palette
        value={target.id}
        as={CommandNew.Item}
      >
        <Tooltip
          key={target.id}
          label={`${target.kind.publisher}/${target.kind.name}/${target.kind.kind}`}
          placement="right"
        >
          <Flex p={6} position="relative">
            <CheckCircleIcon
              visibility={
                checked.includes(target.id.toLowerCase()) ? "visible" : "hidden"
              }
              position="absolute"
              top={0}
              left={0}
              boxSize={"12px"}
              color={"brandBlue.300"}
            />
            <ProviderIcon
              boxSize={"24px"}
              shortType={target.kind.icon as ShortTypes}
            />
          </Flex>
        </Tooltip>

        <HStack>
          {target.fields.map((field, i) => (
            <>
              <Divider
                orientation="vertical"
                borderColor={"black"}
                h="80%"
                hidden={i === 0}
              />

              <Tooltip
                key={field.id}
                label={
                  <Stack>
                    <Text color="white" textStyle={"Body/Small"}>
                      {field.fieldTitle}
                    </Text>
                    <Text color="white" textStyle={"Body/Small"}>
                      {field.fieldDescription}
                    </Text>
                    <Text color="white" textStyle={"Body/Small"}>
                      {field.valueLabel}
                    </Text>
                    <Text color="white" textStyle={"Body/Small"}>
                      {field.value}
                    </Text>
                    <Text color="white" textStyle={"Body/Small"}>
                      {field.valueDescription}
                    </Text>
                  </Stack>
                }
                placement="top"
              >
                <Stack>
                  <Text textStyle={"Body/SmallBold"} noOfLines={1}>
                    {field.fieldTitle}
                  </Text>
                  <Text textStyle={"Body/Small"} noOfLines={1}>
                    {field.valueLabel}
                  </Text>
                </Stack>
              </Tooltip>
            </>
          ))}
        </HStack>
      </Flex>
    );
  };

  // HANDLERS
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

    userRequestPreflight({
      targets: entitlementTargets,
    })
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
  };

  const handleRequest = () => {
    preflightRes &&
      // test
      userPostRequests({
        preflightId: preflightRes?.id,
        reason: accessReason,
        groupOptions: preflightRes.accessGroups.map((g) => {
          return {
            id: g.id,
            timing: {
              durationSeconds: g.timeConstraints.maxDurationSeconds,
            },
          };
        }),
      })
        .then((res) => {
          console.log(res);
          submitLoadingToggle.on();
          // redirect to request...
          navigate({ to: `/requests/${res.id}` });
          // clear state
          setChecked([]);
          handleInputChange("");
        })
        .catch((err) => {
          console.log(err);
        });
  };

  // Debounce the setInputValue function with a 500ms delay
  const handleInputChange = debounce((newValue) => {
    setInputValue(newValue);
  }, 500);

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
                  shouldFilter={false}
                  label="Global Command Menu"
                  checked={checked}
                  setChecked={setChecked}
                >
                  <Input
                    size="lg"
                    type="text"
                    placeholder="What do you want to access?"
                    value={inputValue}
                    onValueChange={handleInputChange}
                    autoFocus={true}
                    as={CommandNew.Input}
                  />
                  <HStack mt={2} overflowX="auto">
                    <FilterBlock
                      label="All Resources"
                      total={allTargets.length}
                      onClick={() => {
                        handleInputChange("");
                        showOnlyCheckedToggle.off();
                      }}
                    />
                    <FilterBlock
                      label="Selected"
                      selected={checked.length}
                      onClick={() => {
                        handleInputChange("");
                        showOnlyCheckedToggle.on();
                        document.getElementById(":rd:")?.focus();
                      }}
                    />
                    {entitlements.data?.entitlements.map((kind) => {
                      const key = (
                        kind.publisher +
                        "#" +
                        kind.name +
                        "#" +
                        kind.kind +
                        "#"
                      ).toLowerCase();
                      return (
                        <FilterBlock
                          label={kind.kind}
                          icon={kind.icon as ShortTypes}
                          onClick={() => {
                            handleInputChange(key);
                            showOnlyCheckedToggle.off();
                            // then set the focus back to the input
                            // so that the user can continue typing
                            document.getElementById(":rd:")?.focus();
                          }}
                          selected={
                            checked.filter((c) => c.startsWith(key)).length
                          }
                        />
                      );
                    })}
                  </HStack>
                  <StyledList
                    mt={2}
                    border="1px solid"
                    rounded="md"
                    borderColor="neutrals.300"
                    p={1}
                    pt={2}
                  >
                    <List
                      style={{}}
                      height={TARGETS * TARGET_HEIGHT}
                      itemCount={filteredItems.length}
                      itemSize={TARGET_HEIGHT}
                      width="100%"
                    >
                      {targetComponent}
                    </List>
                  </StyledList>
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
                              {/* <FieldsCodeBlock fields={target.fields} /> */}
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
interface FilterBlockProps extends CenterProps {
  icon?: ShortTypes;
  total?: number;
  selected?: number;
  label: string;
}
const FilterBlock: React.FC<FilterBlockProps> = ({
  label,
  total,
  selected,
  icon,
  ...rest
}) => {
  return (
    <Center
      rounded="md"
      h="84px"
      borderColor="neutrals.300"
      bg="white"
      borderWidth="1px"
      px={2}
      flexDirection="column"
      as={"button"}
      {...rest}
    >
      {icon !== undefined ? (
        <ProviderIcon shortType={icon} />
      ) : (
        <Box boxSize="22px" />
      )}
      <Text textStyle="Body/Small" noOfLines={1} textAlign="center">
        {label}
      </Text>
      {total === undefined ? (
        selected === undefined ? (
          <Box boxSize="22px" />
        ) : (
          <Text
            textStyle="Body/Small"
            noOfLines={1}
            textAlign="center"
            color="neutrals.500"
          >
            {`${selected} selected`}
          </Text>
        )
      ) : (
        <Text
          textStyle="Body/Small"
          noOfLines={1}
          textAlign="center"
          color="neutrals.500"
        >
          {`${total} total`}
        </Text>
      )}
    </Center>
  );
};
function targetFieldsToString(targetFields: TargetField[]): string {
  const strings = targetFields.map((targetField) => {
    // Concatenate the values with a separator
    const values = [targetField.valueLabel, targetField.valueDescription].join(
      "; "
    ); // Use semicolon and space as a separator
    return values;
  });
  return strings.join("; "); // Use newline character as a separator
}
