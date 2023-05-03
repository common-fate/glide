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

export const RecentRequests: React.FC = () => {
  const {
    data: reqsUpcoming,
    isValidating,
    ...upcomingApi
  } = useInfiniteScrollApi<typeof useUserListRequests>({
    swrHook: useUserListRequests,
    hookProps: { filter: "UPCOMING" },
    swrProps: { swr: { refreshInterval: 10000 } },
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
          <SimpleGrid
            p={1}
            rounded="lg"
            bg="neutrals.100"
            columns={2}
            borderWidth={1}
            borderColor="neutrals.300"
          >
            <VStack spacing={2}>
              {request.accessGroups.map((group) => {
                return (
                  <VStack
                    key={group.id}
                    borderWidth={1}
                    borderColor="neutrals.300"
                    rounded="lg"
                    h="80px"
                    w="100%"
                    background="white"
                  >
                    <Text>{group.status}</Text>
                    <Text>{`${group.targets.length} target${
                      group.targets.length > 1 ? "s" : ""
                    }`}</Text>
                  </VStack>
                );
              })}
            </VStack>
            <Flex direction={"column"}>
              <Spacer />
              <Text textAlign={"right"}>{request.purpose.reason}</Text>
            </Flex>
          </SimpleGrid>
        </LinkOverlay>
      </Link>
    </LinkBox>
  );
};
