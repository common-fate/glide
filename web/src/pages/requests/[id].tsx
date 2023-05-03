import { ArrowBackIcon } from "@chakra-ui/icons";
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
  chakra,
  Code,
  Container,
  Divider,
  Flex,
  Grid,
  GridItem,
  IconButton,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalHeader,
  ModalOverlay,
  Popover,
  PopoverArrow,
  PopoverBody,
  PopoverCloseButton,
  PopoverContent,
  PopoverHeader,
  PopoverTrigger,
  Portal,
  Skeleton,
  SkeletonCircle,
  SkeletonText,
  Spinner,
  Stack,
  Text,
  useBoolean,
  useDisclosure,
  useToast,
} from "@chakra-ui/react";
import { intervalToDuration, format } from "date-fns";
import { useState } from "react";
import { Helmet } from "react-helmet";
import { Link, MakeGenerics, useMatch } from "react-location";
import { AuditLog } from "../../components/AuditLog";
import {
  Days,
  DurationInput,
  Hours,
  Minutes,
  Weeks,
} from "../../components/DurationInput";

import { ProviderIcon, ShortTypes } from "../../components/icons/providerIcon";
import { UserLayout } from "../../components/Layout";
import { GrantStatusCell, StatusCell } from "../../components/StatusCell";
import {
  getGroupTargetStatus,
  useGetGroupTargetStatus,
  useUserGetRequest,
  useUserListRequests,
} from "../../utils/backend-client/default/default";
import {
  userCancelRequest,
  userReviewRequest,
  userRevokeRequest,
} from "../../utils/backend-client/end-user/end-user";
import {
  RequestAccessGroup,
  RequestAccessGroupTarget,
  RequestStatus,
} from "../../utils/backend-client/types";
import {
  durationString,
  durationStringHoursMinutes,
  getEndTimeWithDuration,
} from "../../utils/durationString";
import { request } from "http";
import FieldsCodeBlock from "../../components/FieldsCodeBlock";
import { TargetDetail } from "../../components/Target";

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
  });
  const toast = useToast();

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
                  <Flex direction="row">
                    {request.data ? (
                      <Flex
                        direction="row"
                        justifyContent="space-between"
                        alignItems="flex-start"
                      >
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
                          <Stack>
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
                                {format(
                                  new Date(
                                    Date.parse(request.data.requestedAt)
                                  ),
                                  "p dd/MM/yy"
                                )}
                              </Text>
                            </Flex>

                            <Text textStyle="Body/Small" color="neutrals.500">
                              {request.data.requestedBy.email}
                            </Text>
                          </Stack>
                        </Flex>
                        <Flex alignSelf="flex-start" alignContent="flex-end">
                          {request.data.status == "ACTIVE" && (
                            <ButtonGroup variant="brandSecondary">
                              <Button
                                size="sm"
                                onClick={() => {
                                  console.log("revoke");
                                  //@ts-ignore
                                  userRevokeRequest(request.data?.id)
                                    .then((e) => {
                                      toast({
                                        title: "Revoke Initiated",
                                        status: "success",
                                        variant: "subtle",
                                        duration: 2200,
                                        isClosable: true,
                                      });
                                    })
                                    .catch((e) => {
                                      toast({
                                        title: "Error Revoking",
                                        status: "error",
                                        variant: "subtle",
                                        duration: 2200,
                                        isClosable: true,
                                      });
                                    });
                                }}
                              >
                                Revoke
                              </Button>
                            </ButtonGroup>
                          )}
                          {request.data.status == "PENDING" && (
                            <ButtonGroup variant="brandSecondary">
                              <Button
                                size="sm"
                                onClick={() => {
                                  console.log("cancel");
                                  //@ts-ignore
                                  userCancelRequest(request.data?.id)
                                    .then((e) => {
                                      console.log(e);
                                    })
                                    .catch((e) => {
                                      console.log(e);
                                    });
                                }}
                              >
                                Cancel
                              </Button>
                            </ButtonGroup>
                          )}
                        </Flex>
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
    if (group.requestStatus === "CANCELLED") {
      return (
        <Box
          as="span"
          flex="1"
          textAlign="left"
          sx={{
            p: { lineHeight: "120%", textStyle: "Body/Extra Small" },
          }}
        >
          <Text color="neutrals.700">Cancelled</Text>
        </Box>
      );
    }

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
      </Box>
    );
  }

  if (group.status === "APPROVED") {
    switch (group.requestStatus) {
      case "ACTIVE":
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
                    end: getEndTimeWithDuration(
                      group.requestedTiming.startTime
                        ? group.requestedTiming.startTime
                        : "",
                      group.requestedTiming.durationSeconds
                    ),
                  })
                )
              }
            />
          </Flex>
        );
      case "PENDING":
        return (
          <Box
            as="span"
            flex="1"
            textAlign="left"
            sx={{
              p: { lineHeight: "120%", textStyle: "Body/Extra Small" },
            }}
          >
            <Text color="neutrals.700">Pending</Text>
          </Box>
        );
      case "REVOKING":
        return (
          <Box
            as="span"
            flex="1"
            textAlign="left"
            sx={{
              p: { lineHeight: "120%", textStyle: "Body/Extra Small" },
            }}
          >
            <Text color="neutrals.700">Revoking</Text>
          </Box>
        );
      case "REVOKED":
        return (
          <Box
            as="span"
            flex="1"
            textAlign="left"
            sx={{
              p: { lineHeight: "120%", textStyle: "Body/Extra Small" },
            }}
          >
            <Text color="neutrals.700">Revoked</Text>
          </Box>
        );
      case "COMPLETE":
        return (
          <Box
            as="span"
            flex="1"
            textAlign="left"
            sx={{
              p: { lineHeight: "120%", textStyle: "Body/Extra Small" },
            }}
          >
            <Text color="neutrals.700">Complete</Text>
          </Box>
        );
      default:
        break;
    }
  }

  return null;
};

export const AccessGroupItem = ({ group }: AccessGroupProps) => {
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
                  <TargetDetail
                    target={{
                      fields: target.fields,
                      id: target.id,
                      kind: target.targetKind,
                    }}
                  />

                  <Stack>
                    <Flex
                      onClick={(e) => {
                        e.stopPropagation();
                      }}
                    >
                      <GrantStatusCell
                        requestId={group.requestId}
                        groupId={group.id}
                        targetId={target.id}
                      />
                    </Flex>
                    <Button
                      variant="brandSecondary"
                      size="xs"
                      onClick={() => handleGrantClick(target)}
                    >
                      View
                    </Button>
                  </Stack>
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
                shortType={selectedGrant?.targetKind.icon as ShortTypes}
              />
              {/* <FieldsCodeBlock fields={selectedGrant?.fields || []} /> */}
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
  const isReviewer = false;

  const handleClickMax = () => {
    setDurationSeconds(group.accessRule.timeConstraints.maxDurationSeconds);
  };

  // console.log({ group });

  // durationSeconds state
  const [durationSeconds, setDurationSeconds] = useState<number>(
    group.requestedTiming.durationSeconds
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
                <Box>
                  <Box mt={1}>
                    <DurationInput
                      // {...rest}
                      onChange={setDurationSeconds}
                      value={durationSeconds}
                      hideUnusedElements={true}
                      max={group.accessRule.timeConstraints.maxDurationSeconds}
                      min={60}
                      defaultValue={group.overrideTiming?.durationSeconds}
                    >
                      <Weeks />
                      <Days />
                      <Hours />
                      <Minutes />
                      <Button
                        variant="brandSecondary"
                        flexDir="column"
                        fontSize="12px"
                        lineHeight="12px"
                        mr={2}
                        isActive={
                          durationSeconds ==
                          group.accessRule.timeConstraints.maxDurationSeconds
                        }
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
                {/*    <Select
                  mt={8}
                  size="xs"
                  variant="brandSecondary"
                  onChange={(e) => setState(e.target.value)}
                >
                  {["default", "max", "custom"].map((option) => (
                    <option value={option}>{option}</option>
                  ))}
                </Select> */}
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
              userReviewRequest(group.requestId, group.id, {
                decision: "APPROVED",
              });
            }}
          >
            Approve
          </Button>
          <Button
            size="sm"
            onClick={() => {
              console.log("reject");
              // @TODO: add in admin approval API methods
              userReviewRequest(group.requestId, group.id, {
                decision: "DECLINED",
              });
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
