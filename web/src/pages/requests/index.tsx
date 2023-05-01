import {
  Box,
  Button,
  Center,
  CenterProps,
  Container,
  Flex,
  Grid,
  Link as ChakraLink,
  LinkBox,
  LinkBoxProps,
  LinkOverlay,
  Skeleton,
  SkeletonCircle,
  SkeletonText,
  Spinner,
  Stack,
  Tab,
  TabList,
  TabPanel,
  TabPanels,
  Tabs,
  Text,
  useDisclosure,
  VStack,
} from "@chakra-ui/react";
import React from "react";
import { Helmet } from "react-helmet";
import { Link, MakeGenerics, useNavigate, useSearch } from "react-location";
import { ProviderIcon, ShortTypes } from "../../components/icons/providerIcon";
import { UserLayout } from "../../components/Layout";
import { useUserListRequests } from "../../utils/backend-client/default/default";
import {} from "../../utils/backend-client/end-user/end-user";
import { AccessRule, Request } from "../../utils/backend-client/types";
import { useUser } from "../../utils/context/userContext";
import { renderTiming } from "../../utils/renderTiming";
import { useInfiniteScrollApi } from "../../utils/useInfiniteScrollApi";

type MyLocationGenerics = MakeGenerics<{
  Search: {
    filter?: "upcoming" | "past";
  };
}>;

const Home = () => {
  const search = useSearch<MyLocationGenerics>();
  const navigate = useNavigate<MyLocationGenerics>();

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

  const { isOpen, onClose, onToggle } = useDisclosure();

  return (
    <>
      <UserLayout>
        <Helmet>
          <title>Common Fate</title>
        </Helmet>
        <Box overflow="auto">
          <Container maxW="container.xl" pt={{ base: 12, lg: 32 }}>
            <Stack
              direction={["column", "column", "column", "row", "row"]}
              justifyContent="center"
              spacing={12}
            >
              <VStack spacing={8}>
                {/* <Favorites /> */}
                <Flex flexDirection="column" w="100%">
                  <Flex>
                    <Text
                      as="h3"
                      textStyle="Heading/H3"
                      mt="6px" // this minor adjustment aligns heading with Tabbed content on XL screen widths
                    >
                      New Request
                    </Text>
                    <Button
                      display={{ base: "flex", lg: "none" }}
                      variant="brandSecondary"
                      size="sm"
                      ml="auto"
                      onClick={onToggle}
                    >
                      View All
                    </Button>
                  </Flex>
                  <Rules />
                </Flex>
              </VStack>

              <Tabs
                variant="brand"
                w="100%"
                index={search.filter === "past" ? 1 : 0}
                onChange={(i: any) => {
                  const tab = i === 1 ? "past" : "upcoming";
                  navigate({ search: (old) => ({ ...old, filter: tab }) });
                }}
              >
                <TabList>
                  <Tab>Upcoming</Tab>
                  <Tab>Past</Tab>
                </TabList>
                <TabPanels>
                  <TabPanel overflowY="auto">
                    <Stack spacing={5} maxH="80vh">
                      {reqsUpcoming?.requests?.map((request, i) => (
                        <UserAccessCard
                          type="upcoming"
                          key={request.id + i}
                          req={request}
                          index={i}
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
                  <TabPanel overflowY="auto">
                    <Stack
                      spacing={5}
                      maxH="80vh"
                      // ref={pastRef}
                      // onScroll={() => handleScroll("past")}
                    >
                      {reqsPast?.requests.map((request, i) => (
                        <UserAccessCard
                          index={i}
                          type="past"
                          key={request.id}
                          req={request}
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
            </Stack>
          </Container>
        </Box>
        {/* <AcessRulesMobileModal isOpen={isOpen} onClose={onClose} /> */}
      </UserLayout>
    </>
  );
};

const Rules = () => {
  const rules: { accessRules: Array<AccessRule> } = {
    accessRules: [],
  };
  const user = useUser();

  // loading/standard state needs to be rendered in a Grid container
  if ((rules && rules.accessRules.length > 0) || typeof rules == "undefined") {
    return (
      <Grid
        mt={4}
        templateColumns={{
          base: "repeat(20, 1fr)",
          lg: "repeat(1, 1fr)",
          xl: "repeat(2, 1fr)",
        }}
        templateRows={{ base: "repeat(1, 1fr)", xl: "unset" }}
        minW={{ base: "unset", xl: "488px" }}
        gap={6}
        overflowX={["scroll", "auto"]}
      >
        {
          // Loading state
          typeof rules === "undefined"
            ? [1, 2, 3, 4].map((i) => (
                <Skeleton key={i} p={6} h="172px" w="232px" rounded="sm" />
              ))
            : rules.accessRules.map((r, i) => (
                <Link
                  style={{ display: "flex" }}
                  to={"/access/requests/" + r.id}
                  key={r.id}
                >
                  <Box
                    className="group"
                    textAlign="center"
                    bg="neutrals.100"
                    p={6}
                    h="172px"
                    w="232px"
                    rounded="md"
                    data-testid={"r_" + i}
                    pos="relative"
                    overflow="hidden"
                  >
                    {/* <ProviderIcon
                      shortType={r.target.provider.id as ShortTypes}
                      mb={3}
                      h="8"
                      w="8"
                    /> */}

                    <Text
                      textStyle="Body/SmallBold"
                      color="neutrals.700"
                      noOfLines={3}
                    >
                      {r.name}
                    </Text>

                    <Button
                      mt={4}
                      variant="brandSecondary"
                      size="sm"
                      opacity={0}
                      pos="absolute"
                      bottom={6}
                      sx={{
                        // This media query ensure always visible for touch screens
                        "@media (hover: none)": {
                          opacity: 1,
                        },
                      }}
                      transition="all .2s ease-in-out"
                      transform="translate(-50%, 8px)"
                      _groupHover={{
                        bg: "white",
                        opacity: 1,
                        transform: "translateY(-50%, 0px)",
                      }}
                      left="50%"
                    >
                      Request
                    </Button>
                  </Box>
                </Link>
              ))
        }
      </Grid>
    );
  }
  // empty state
  if (rules?.accessRules.length === 0) {
    return (
      <Center
        bg="neutrals.100"
        p={6}
        as="a"
        h="193px"
        w={{ base: "100%", md: "488px" }}
        rounded="md"
        flexDir="column"
        textAlign="center"
        mt={4}
      >
        <Text textStyle="Heading/H3" color="neutrals.500">
          No Access
        </Text>
        <Text textStyle="Body/Medium" color="neutrals.400" mt={2}>
          You don’t have access to anything yet.{" "}
          {user?.isAdmin ? (
            <ChakraLink
              as={Link}
              to="/admin/access-rules/create"
              textDecor="none"
              _hover={{ textDecor: "underline" }}
            >
              Click here to create a new access rule.
            </ChakraLink>
          ) : (
            "Ask your Common Fate administrator to finish setting up Common Fate."
          )}
        </Text>
      </Center>
    );
  }
  // should never be reached; but needed for type safety
  return null;
};

export default Home;

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

const UserAccessCard: React.FC<
  {
    req: Request;
    type: "upcoming" | "past";
    index: number;
  } & LinkBoxProps
> = ({ req, type, index, ...rest }) => {
  //@ts-ignore
  const rule: AccessRule = {};

  return (
    <LinkBox {...rest}>
      <Link to={"/requests/" + req.id}>
        <LinkOverlay>
          <Flex
            rounded="md"
            bg="neutrals.100"
            flexDir="column"
            key={req.id}
            pos="relative"
            data-testid={"req_" + req.id}
          >
            <Stack flexDir="column" p={8} pos="relative" spacing={2}>
              {/* <RequestStatusDisplay request={req} /> */}

              <Flex justify="space-between">
                <Box>
                  {rule ? (
                    <Flex align="center" mr="auto">
                      {/* <ProviderIcon
                        shortType={rule?.target.provider.id as ShortTypes}
                        h={10}
                        w={10}
                      /> */}
                      <Text
                        ml={2}
                        textStyle="Body/LargeBold"
                        color="neutrals.700"
                      >
                        {rule?.name}
                      </Text>
                    </Flex>
                  ) : (
                    <Flex align="center" h="40px">
                      <SkeletonCircle h={10} w={10} mr={2} />
                      <SkeletonText noOfLines={1} width="6ch" />
                    </Flex>
                  )}
                </Box>
              </Flex>
            </Stack>
          </Flex>
        </LinkOverlay>
      </Link>
    </LinkBox>
  );
};

// const Favorites: React.FC = () => {
//   const { data: favorites } = useUserListFavorites();

//   if (favorites?.favorites.length === 0) {
//     return null;
//   }

//   return (
//     <Box w="100%">
//       <Flex>
//         <Text
//           as="h3"
//           textStyle="Heading/H3"
//           mt="6px" // this minor adjustment aligns heading with Tabbed content on XL screen widths
//         >
//           Favorites
//         </Text>
//       </Flex>
//       <Grid
//         mt={4}
//         templateColumns={{
//           base: "repeat(20, 1fr)",
//           lg: "repeat(1, 1fr)",
//           xl: "repeat(2, 1fr)",
//         }}
//         templateRows={{ base: "repeat(1, 1fr)", xl: "unset" }}
//         minW={{ base: "unset", xl: "488px" }}
//         gap={6}
//       >
//         {favorites
//           ? favorites.favorites.map((r, i) => (
//               <Link
//                 style={{ display: "flex" }}
//                 to={"/access/requests/" + r.ruleId + "?favorite=" + r.id}
//                 key={r.id}
//               >
//                 <Box
//                   className="group"
//                   textAlign="center"
//                   bg="neutrals.100"
//                   p={6}
//                   h="172px"
//                   w="232px"
//                   rounded="md"
//                   data-testid={`fav-request-item-${r.name}`}
//                 >
//                   <Text textStyle="Body/SmallBold" color="neutrals.700">
//                     {r.name}
//                   </Text>

//                   <Button
//                     mt={4}
//                     variant="brandSecondary"
//                     size="sm"
//                     opacity={0}
//                     sx={{
//                       // This media query ensure always visible for touch screens
//                       "@media (hover: none)": {
//                         opacity: 1,
//                       },
//                     }}
//                     transition="all .2s ease-in-out"
//                     transform="translateY(8px)"
//                     _groupHover={{
//                       bg: "white",
//                       opacity: 1,
//                       transform: "translateY(0px)",
//                     }}
//                   >
//                     Request
//                   </Button>
//                 </Box>
//               </Link>
//             ))
//           : // Otherwise loading state
//             [1, 2].map((i) => (
//               <Skeleton key={i} p={6} h="172px" w="232px" rounded="sm" />
//             ))}
//       </Grid>
//     </Box>
//   );
// };
