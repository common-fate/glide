import { ArrowBackIcon, EditIcon } from "@chakra-ui/icons";
import {
  Accordion,
  AccordionButton,
  AccordionIcon,
  AccordionItem,
  AccordionPanel,
  Avatar,
  Box,
  Button,
  ButtonGroup,
  Center,
  Code,
  Container,
  chakra,
  Divider,
  Flex,
  IconButton,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalHeader,
  ModalOverlay,
  SkeletonCircle,
  Skeleton,
  SkeletonText,
  Spinner,
  Stack,
  Text,
  useDisclosure,
  Grid,
  GridItem,
  Popover,
  PopoverArrow,
  PopoverBody,
  PopoverCloseButton,
  PopoverContent,
  PopoverHeader,
  PopoverTrigger,
  Portal,
  FormLabel,
  useBoolean,
  Select,
} from "@chakra-ui/react";
import validateFormData from "@rjsf/core/lib/validate";
import { formatDistance, intervalToDuration } from "date-fns";
import { useState } from "react";
import { Helmet } from "react-helmet";
import { Link, MakeGenerics, useMatch } from "react-location";
import { AuditLog } from "../../components/AuditLog";
import {
  DurationInput,
  Days,
  Hours,
  Minutes,
  Weeks,
} from "../../components/DurationInput";
import FieldsCodeBlock from "../../components/FieldsCodeBlock";
import { EditPencilIcon } from "../../components/icons/Icons";
import { ProviderIcon, ShortTypes } from "../../components/icons/providerIcon";
import { UserLayout } from "../../components/Layout";
import { StatusCell } from "../../components/StatusCell";
import { useUserGetRequest } from "../../utils/backend-client/default/default";
import {
  PreflightAccessGroup,
  RequestAccessGroup,
  RequestAccessGroupTarget,
  Target,
} from "../../utils/backend-client/types";
import {
  durationString,
  durationStringHoursMinutes,
} from "../../utils/durationString";

type MyLocationGenerics = MakeGenerics<{
  Search: {
    action?: "approve" | "close";
  };
}>;

const Home = () => {
  const {
    params: { id: requestId },
  } = useMatch();

  const request = useUserGetRequest(requestId, {
    swr: { refreshInterval: 10000 },
    request: {
      baseURL: "http://127.0.0.1:3100",
      headers: {
        Prefer: "code=200, example=ex_1",
      },
    },
  });

  // request.data?.accessGroups[0].targets[0]

  // const groups = useUserListRequestAccessGroups(requestId, {
  //   swr: { refreshInterval: 10000 },
  //   request: {
  //     baseURL: "http://127.0.0.1:3100",
  //     headers: {
  //       Prefer: "code=200, example=ex_1",
  //     },
  //   },
  // });

  // const grants = useUserListRequestAccessGroupGrants("TEST", {
  //   swr: { refreshInterval: 10000 },
  //   request: {
  //     baseURL: "http://127.0.0.1:3100",
  //     headers: {
  //       Prefer: "code=200, example=ex_1",
  //     },
  //   },
  // });

  // const search = useSearch<MyLocationGenerics>();
  // const { action } = search;

  // const [cachedReq, setCachedReq] = useState(data);
  // useEffect(() => {
  //   if (data !== undefined) setCachedReq(data);
  //   return () => {
  //     setCachedReq(undefined);
  //   };
  // }, [data]);

  // const user = useUser();

  return (
    <div>
      <UserLayout>
        <Helmet>
          <title>Access Request</title>
        </Helmet>
        {/* The header bar */}
        <Center borderBottom="1px solid" borderColor="neutrals.200" h="80px">
          <IconButton
            as={Link}
            aria-label="Go back"
            pos="absolute"
            left={4}
            icon={<ArrowBackIcon />}
            rounded="full"
            variant="ghost"
            // to={data?.canReview ? "/reviews?status=pending" : "/requests"}
          />

          <Text as="h4" textStyle="Heading/H4">
            Request details
          </Text>
        </Center>

        {/* Main content */}
        <Container
          maxW={{
            md: "container.lg",
          }}
        >
          <Grid mt={8} gridTemplateColumns={{ sm: "1fr 240px" }} gap="4">
            <GridItem>
              <>
                <Stack spacing={4}>
                  <Flex px={2}>
                    {request.data ? (
                      <Flex>
                        <Avatar
                          size="sm"
                          src={request.data.requestedBy.picture}
                          name={
                            request.data.requestedBy.firstName +
                            " " +
                            request.data.requestedBy.lastName
                          }
                          mr={2}
                        />
                        <Box>
                          <Flex>
                            <Text textStyle="Body/Small">
                              {request.data.requestedBy.firstName +
                                " " +
                                request.data.requestedBy.lastName}
                            </Text>
                            <Text
                              textStyle="Body/Small"
                              ml={1}
                              color="neutrals.500"
                            >
                              requested at&nbsp;
                              {request.data.requestedAt}
                            </Text>
                          </Flex>
                          <Text textStyle="Body/Small" color="neutrals.500">
                            {request.data.requestedBy.email}
                          </Text>
                        </Box>
                      </Flex>
                    ) : (
                      <Flex h="42px">
                        <SkeletonCircle size="24px" mr={4} />
                        <Box>
                          <Flex>
                            <SkeletonText
                              noOfLines={1}
                              h="12px"
                              w="12ch"
                              mr="4px"
                            />
                            <SkeletonText noOfLines={1} h="12px" w="12ch" />
                          </Flex>
                          <SkeletonText noOfLines={1} h="12px" w="12ch" />
                        </Box>
                      </Flex>
                    )}
                  </Flex>
                  <Divider borderColor="neutrals.300" w="100%" />

                  <Stack spacing={4} w="100%">
                    {request.data
                      ? request.data.accessGroups.map((group) => (
                          <AccessGroupItem key={group.id} group={group} />
                        ))
                      : [
                          <Skeleton key={1} rounded="md" h="282px" w="100%" />,
                          <Skeleton key={2} rounded="md" h="82px" w="100%" />,
                          <Skeleton key={3} rounded="md" h="82px" w="100%" />,
                        ]}
                  </Stack>
                </Stack>

                {/* <Code
                  maxW="60ch"
                  textOverflow="clip"
                  whiteSpace="pre-wrap"
                  mt={32}
                >
                  {JSON.stringify({ request }, null, 2)}
                </Code> */}
              </>
            </GridItem>
            <GridItem>
              <AuditLog />
            </GridItem>
          </Grid>
        </Container>
      </UserLayout>
    </div>
  );
};

type AccessGroupProps = {
  group: RequestAccessGroup;
};

export const HeaderStatusCell = ({ group }: AccessGroupProps) => {
  if (group.status === "PENDING_APPROVAL") {
    return (
      <Box
        as="span"
        flex="1"
        textAlign="left"
        sx={{
          p: { lineHeight: "120%", textStyle: "Body/Extra Small" },
        }}
      >
        <Text color="neutrals.700">Review Required</Text>
        {/* <AvatarGroup size="sm" max={2} ml={-2}>
          {group.reviewers.map((reviewer) => (
            <Avatar
              key={reviewer.id}
              name={reviewer.firstName + " " + reviewer.lastName}
              src={reviewer.picture}
            />
          ))}
        </AvatarGroup> */}
      </Box>
    );
  }

  if (group.status === "APPROVED") {
    return (
      <Flex flex="1">
        <StatusCell
          success="ACTIVE"
          value={group.status}
          replaceValue={
            "Active for the next " +
            durationStringHoursMinutes(
              intervalToDuration({
                start: new Date(),
                end: new Date(),
                // end: new Date(Date.parse(group.grant.end)),
              })
            )
          }
        />
      </Flex>
    );
  }

  return null;
};

export const AccessGroupItem = ({ group }: AccessGroupProps) => {
  // const grants = useUserListRequestAccessGroupGrants(group.id, {
  //   swr: { refreshInterval: 10000 },
  //   request: {
  //     baseURL: "http://127.0.0.1:3100",
  //     headers: {
  //       Prefer: "code=200, example=ex_1",
  //     },
  //   },
  // });

  const [selectedGrant, setSelectedGrant] =
    useState<RequestAccessGroupTarget>();
  const grantModalState = useDisclosure();

  const handleGrantClick = (grant: RequestAccessGroupTarget) => {
    setSelectedGrant(grant);
    grantModalState.onOpen();
  };
  const handleClose = () => {
    setSelectedGrant(undefined);
    grantModalState.onClose();
  };

  const isReviewer = true;

  return (
    <Box bg="neutrals.100" borderColor="neutrals.300" rounded="lg">
      <Accordion
        key={group.id}
        allowToggle
        // we may want to play with how default index works
        defaultIndex={[0]}
      >
        <AccordionItem border="none">
          <AccordionButton
            p={2}
            bg="neutrals.100"
            roundedTop="md"
            borderColor="neutrals.300"
            borderWidth="1px"
            sx={{
              "&[aria-expanded='false']": {
                roundedBottom: "md",
              },
            }}
          >
            <AccordionIcon boxSize="6" mr={2} />
            <HeaderStatusCell group={group} />
            <ApproveRejectDuration group={group} />
          </AccordionButton>

          <AccordionPanel
            borderColor="neutrals.300"
            borderTop="none"
            roundedBottom="md"
            borderWidth="1px"
            bg="white"
            p={0}
          >
            <Stack spacing={2} p={2}>
              {group.targets.map((target) => (
                <Flex
                  w="100%"
                  borderColor="neutrals.300"
                  rounded="md"
                  borderWidth="1px"
                  bg="white"
                  p={2}
                  pos="relative"
                >
                  <ProviderIcon boxSize="24px" shortType="aws-sso" mr={2} />
                  <FieldsCodeBlock fields={target.fields} />
                  {false && (
                    <Spinner
                      thickness="2px"
                      speed="0.65s"
                      emptyColor="neutrals.300"
                      color="neutrals.800"
                      size="sm"
                      top={4}
                      right={4}
                      pos="absolute"
                    />
                  )}
                  <Button
                    variant="brandSecondary"
                    size="xs"
                    top={2}
                    right={2}
                    pos="absolute"
                    onClick={() => handleGrantClick(target)}
                  >
                    View
                  </Button>
                </Flex>
              ))}
            </Stack>
          </AccordionPanel>
        </AccordionItem>
      </Accordion>
      <Modal isOpen={grantModalState.isOpen} onClose={handleClose}>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader></ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            <Box>
              <ProviderIcon
                shortType={selectedGrant?.targetGroupIcon as ShortTypes}
              />
              <FieldsCodeBlock fields={selectedGrant?.fields || []} />
            </Box>
            <Text textStyle="Body/Small">Access Instructions</Text>

            <Code bg="white" whiteSpace="pre-wrap">
              {JSON.stringify(selectedGrant, null, 2)}
            </Code>
            <Text></Text>
          </ModalBody>
        </ModalContent>
      </Modal>
    </Box>
  );
};

// @TODO: sort out state for props.........
type ApproveRejectDurationProps = {
  group: RequestAccessGroup;
};

export const ApproveRejectDuration = ({
  group,
}: ApproveRejectDurationProps) => {
  const isReviewer = true;

  const [state, setState] = useState<"default" | "max" | "custom">("default");

  const handleClickMax = () => {
    setState("max");
    setDurationSeconds(group.accessRule.timeConstraints.maxDurationSeconds);
    // TODO: also set the duration to max
  };

  const handleClickEdit = () => {
    setState("custom");
  };

  const handleCancel = () => {
    setState("default");
    // todo: ensure new duration is set
  };

  console.log({ group });

  // durationSeconds state
  const [durationSeconds, setDurationSeconds] = useState<number>(
    group.time.durationSeconds
  );

  const [isEditing, setIsEditing] = useBoolean();

  return (
    <Flex
      alignSelf="baseline"
      flexDir="row"
      alignItems="center"
      onClick={(e) => {
        e.stopPropagation();
      }}
    >
      <Flex h="32px" alignItems="baseline" flexDir="column" mr={4}>
        <Text textStyle="Body/ExtraSmall" color="neutrals.800">
          {isEditing
            ? "Custom Duration"
            : durationSeconds
            ? durationString(durationSeconds)
            : "No Duration Set"}
        </Text>
        <Popover
          placement="bottom-start"
          isOpen={isEditing}
          onOpen={setIsEditing.on}
          onClose={setIsEditing.off}
        >
          <PopoverTrigger>
            <Button
              pt="4px"
              size="sm"
              textStyle="Body/ExtraSmall"
              fontSize="12px"
              lineHeight="8px"
              color="neutrals.500"
              variant="link"
            >
              Edit Duration
            </Button>
          </PopoverTrigger>
          <Portal>
            <PopoverContent
              minW="256px"
              w="min-content"
              borderColor="neutrals.300"
            >
              <PopoverHeader fontWeight="normal" borderColor="neutrals.300">
                Edit Duration
              </PopoverHeader>
              <PopoverArrow
                sx={{
                  "--popper-arrow-shadow-color": "#E5E5E5",
                }}
              />
              <PopoverCloseButton />
              <PopoverBody py={4}>
                {/* {durationString(durationSeconds)} */}
                {state != "custom" ? (
                  <Flex
                    display="flex"
                    flexDir="row"
                    alignItems="center"
                    sx={{
                      button: {
                        w: "50%",
                        rounded: "md",
                        borderColor: "neutrals.300",
                        color: "neutrals.800",
                        p: 2,
                        _active: {
                          borderColor: "brandBlue.100",
                          color: "brandBlue.300",
                          bg: "white",
                        },
                      },
                    }}
                  >
                    <Flex w="100%">
                      <Button
                        variant="brandSecondary"
                        onClick={(e) => {}}
                        flexDir="column"
                        fontSize="12px"
                        lineHeight="12px"
                        mr={2}
                        isActive={state === "max"}
                        onClick={handleClickMax}
                      >
                        <chakra.span
                          display="block"
                          w="100%"
                          letterSpacing="1.1px"
                        >
                          MAX
                        </chakra.span>
                        {durationString(
                          group.accessRule.timeConstraints.maxDurationSeconds
                        )}
                      </Button>
                      <Button
                        variant="brandSecondary"
                        onClick={handleClickEdit}
                        isActive={state === "custom"}
                        leftIcon={
                          <EditPencilIcon marginTop="-2px" marginRight="-8px" />
                        }
                      />
                    </Flex>
                  </Flex>
                ) : (
                  <Box>
                    {/* <Text textStyle="Body/Small">Session Duration</Text> */}
                    <Box mt={1}>
                      <DurationInput
                        // {...rest}
                        onChange={setDurationSeconds}
                        value={durationSeconds}
                        hideUnusedElements={true}
                        max={
                          group.accessRule.timeConstraints.maxDurationSeconds
                        }
                        min={60}
                        defaultValue={group.overrideTiming.durationSeconds}
                      >
                        <Weeks />
                        <Days />
                        <Hours />
                        <Minutes />
                        <Button
                          variant="brandSecondary"
                          onClick={(e) => {}}
                          flexDir="column"
                          fontSize="12px"
                          lineHeight="12px"
                          mr={2}
                          isActive={state === "max"}
                          onClick={handleClickMax}
                          sx={{
                            w: "50%",
                            rounded: "md",
                            borderColor: "neutrals.300",
                            color: "neutrals.800",
                            p: 2,
                            _active: {
                              borderColor: "brandBlue.100",
                              color: "brandBlue.300",
                              bg: "white",
                            },
                          }}
                        >
                          <chakra.span
                            display="block"
                            w="100%"
                            letterSpacing="1.1px"
                          >
                            MAX
                          </chakra.span>
                          {durationString(
                            group.accessRule.timeConstraints.maxDurationSeconds
                          )}
                        </Button>
                      </DurationInput>
                    </Box>
                  </Box>
                )}
                <Select
                  mt={8}
                  size="xs"
                  variant="brandSecondary"
                  onChange={(e) => setState(e.target.value)}
                >
                  {["default", "max", "custom"].map((option) => (
                    <option value={option}>{option}</option>
                  ))}
                </Select>
              </PopoverBody>
            </PopoverContent>
          </Portal>
        </Popover>
        {/* {durationString(durationSeconds)} */}
      </Flex>
      {isReviewer && (
        <ButtonGroup ml="auto" variant="brandSecondary" spacing={2}>
          <Button
            size="sm"
            onClick={() => {
              console.log("approve");
              // @TODO: add in admin approval API methods
            }}
          >
            Approve
          </Button>
          <Button
            size="sm"
            onClick={() => {
              console.log("reject");
              // @TODO: add in admin approval API methods
            }}
          >
            Reject
          </Button>
        </ButtonGroup>
      )}
    </Flex>
  );
};

export default Home;
