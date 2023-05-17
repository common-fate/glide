import {
  Box,
  Center,
  CenterProps,
  Flex,
  HStack,
  LinkBox,
  LinkBoxProps,
  LinkOverlay,
  SimpleGrid,
  Skeleton,
  chakra,
  SkeletonCircle,
  SkeletonText,
  Spacer,
  Spinner,
  Stack,
  Tab,
  TabList,
  TabPanel,
  TabPanels,
  Tabs,
  Text,
  VStack,
} from "@chakra-ui/react";
import React from "react";
import { Link } from "react-location";
import { useUserListRequests } from "../utils/backend-client/default/default";
import { useInfiniteScrollApi } from "../utils/useInfiniteScrollApi";
import { Request } from "../utils/backend-client/types";

import { StatusCell } from "./StatusCell";
import { ProviderIcon, ShortTypes } from "./icons/providerIcon";

export const RecentRequests: React.FC = () => {
  const {
    data: reqsUpcoming,
    isValidating,
    ...upcomingApi
  } = useInfiniteScrollApi<typeof useUserListRequests>({
    swrHook: useUserListRequests,
    hookProps: { filter: "UPCOMING" },
    swrProps: {
      // @ts-ignore; type discrepancy with latest SWR client
      swr: { refreshInterval: 10000 },
    },
    listObjKey: "requests",
  });

  const { data: reqsPast, ...pastApi } = useInfiniteScrollApi<
    typeof useUserListRequests
  >({
    swrHook: useUserListRequests,
    hookProps: { filter: "PAST" },
    listObjKey: "requests",
  });
  return (
    <Tabs variant="brand" w="100%">
      <TabList>
        <Tab>Upcoming</Tab>
        <Tab>Past</Tab>
      </TabList>
      <TabPanels>
        <TabPanel overflowY="auto" px="0px">
          <Stack spacing={5} maxH="80vh">
            {reqsUpcoming?.requests?.map((request, i) => (
              <RequestCard
                type="upcoming"
                key={"card" + request.id}
                request={request}
              />
            ))}
            {reqsUpcoming === undefined && (
              <>
                <Skeleton h="224px" w="100%" rounded="md" />
                <Skeleton h="224px" w="100%" rounded="md" />
                <Skeleton h="224px" w="100%" rounded="md" />
              </>
            )}
            {!isValidating && reqsUpcoming?.requests.length === 0 && (
              <Center
                bg="neutrals.100"
                p={6}
                as="a"
                h="310px"
                w="100%"
                rounded="md"
                // flexDir="column"
                // textAlign="center"
              >
                <Text textStyle="Heading/H3" color="neutrals.500">
                  No upcoming requests{" "}
                  <Text as="span" opacity={0.5}>
                    ☀️
                  </Text>
                </Text>
              </Center>
            )}
            <LoadMoreButton
              // dont apply ref when validating
              // ref={upcomingRef}
              // isDisabled={!upcomingApi.canNextPage}
              onClick={upcomingApi.incrementPage}
            >
              {isValidating && reqsUpcoming?.requests ? (
                <Spinner />
              ) : upcomingApi.canNextPage ? (
                "Load more"
              ) : (reqsUpcoming?.requests?.length ?? 0) > 4 ? (
                "That's it!"
              ) : (
                ""
              )}
            </LoadMoreButton>
          </Stack>
        </TabPanel>
        <TabPanel overflowY="auto" px="0px">
          <Stack spacing={5} maxH="80vh">
            {reqsPast?.requests.map((request) => (
              <RequestCard
                type="past"
                key={"card" + request.id}
                request={request}
              />
            ))}
            {reqsPast?.requests === undefined && (
              <>
                <Skeleton h="224px" w="100%" rounded="md" />
                <Skeleton h="224px" w="100%" rounded="md" />
                <Skeleton h="224px" w="100%" rounded="md" />
              </>
            )}
            {reqsPast?.requests.length === 0 && (
              <Center
                bg="neutrals.100"
                p={6}
                as="a"
                h="310px"
                w="100%"
                rounded="md"
              >
                <Text textStyle="Heading/H3" color="neutrals.500">
                  No past requests{" "}
                  <Text as="span" opacity={0.5}>
                    ☀️
                  </Text>
                </Text>
              </Center>
            )}
            <LoadMoreButton
              // dont apply ref when validating
              // ref={isValidating ? null : pastRef}
              // disabled={!pastApi.canNextPage}
              onClick={pastApi.incrementPage}
            >
              {pastApi.isValidating && reqsPast?.requests ? (
                <Spinner />
              ) : pastApi.canNextPage ? (
                "Load more"
              ) : (reqsPast?.requests?.length ?? 0) > 4 ? (
                "That's it!"
              ) : (
                ""
              )}
            </LoadMoreButton>
          </Stack>
        </TabPanel>
      </TabPanels>
    </Tabs>
  );
};
const LoadMoreButton = (props: CenterProps) => (
  <Center
    minH={12}
    as="button"
    color="neutrals.500"
    h={10}
    w="100%"
    _hover={{
      _disabled: {
        textDecor: "none",
      },
      textDecor: "underline",
    }}
    {...props}
  />
);

const RequestCard: React.FC<
  {
    request: Request;
    type: "upcoming" | "past";
  } & LinkBoxProps
> = ({ request, type, ...rest }) => {
  return (
    <LinkBox {...rest}>
      <Link to={"/requests/" + request.id}>
        <LinkOverlay>
          <Box
            p={1}
            rounded="lg"
            bg="neutrals.100"
            // columns={2}
            borderWidth={1}
            borderColor="neutrals.300"
          >
            <Flex px={1} align="center">
              <StatusCell
                sx={{
                  span: {
                    textStyle: "Body/Small",
                  },
                }}
                value={request.status}
                replaceValue={request.status.toLowerCase()}
                success={["COMPLETE", "ACTIVE"]}
                warning={["PENDING", "REVOKING"]}
                danger={["REVOKED", "CANCELLED"]}
              />
              <Text textStyle="Body/ExtraSmall">{request.purpose.reason}</Text>
            </Flex>
            <VStack spacing={1}>
              {request.accessGroups.map((group) => {
                return (
                  <Flex
                    key={group.id}
                    borderWidth={1}
                    borderColor="neutrals.300"
                    rounded="lg"
                    // minH="64px"
                    w="100%"
                    background="white"
                    p={3}
                    position="relative"
                    flexDir="column"
                    pos="relative"
                  >
                    {/* ABS. POSITIONED REQ STATUS TOP RIGHT */}
                    <StatusCell
                      sx={{
                        span: {
                          textStyle: "Body/Medium",
                        },
                      }}
                      top={2}
                      right={1}
                      minW="8px"
                      position="absolute"
                      value={group.status}
                      replaceValue=" " // this makes the status cell not render the value
                      success={["APPROVED"]}
                      warning={["PENDING_APPROVAL"]}
                      danger="DECLINED"
                    />
                    {group.targets.slice(0, 6).map((target) => (
                      <Flex align="center">
                        <ProviderIcon
                          h="18px"
                          w="18px"
                          shortType={target.targetKind.name as ShortTypes}
                          mr={2}
                        />
                        {/* <Text textStyle="Body/Medium">{`${
                        group.targets.length
                      } Target${group.targets.length > 1 ? "s" : ""}`}</Text> */}

                        {target.fields.map((field) => {
                          return (
                            <chakra.span
                              textStyle="Body/Small"
                              fontFamily="mono"
                              maxW="24ch"
                              noOfLines={1}
                              textOverflow="ellipsis"
                              wordBreak="break-all"
                            >
                              {field.valueLabel}
                            </chakra.span>
                          );
                        })}
                      </Flex>
                    ))}
                    {group.targets.length > 4 && [
                      <Box
                        position="absolute"
                        bottom={0}
                        left={0}
                        right={0}
                        h="100%"
                        bg="linear-gradient(-15deg, rgba(255,255,255,1) 15%, rgba(255,255,255,0) 100%)"
                        rounded="md"
                      />,
                      <Text
                        textStyle="Body/Small"
                        color="neutrals.500"
                        position="absolute"
                        bottom={2}
                        right={2}
                      >
                        +{group.targets.length - 4} more
                      </Text>,
                    ]}
                  </Flex>
                );
              })}
            </VStack>
          </Box>
        </LinkOverlay>
      </Link>
    </LinkBox>
  );
};
