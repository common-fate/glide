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
  Code,
  Container,
  Divider,
  Flex,
  IconButton,
  SkeletonCircle,
  SkeletonText,
  Spinner,
  Stack,
  Text,
} from "@chakra-ui/react";
import { formatDistance } from "date-fns";
import { Helmet } from "react-helmet";
import { Link, MakeGenerics, useMatch } from "react-location";
import { ProviderIcon } from "../../components/icons/providerIcon";
import { UserLayout } from "../../components/Layout";
import {
  useUserGetRequest,
  useUserListRequestAccessGroupGrants,
  useUserListRequestAccessGroups,
} from "../../utils/backend-client/default/default";
import { AccessGroup } from "../../utils/backend-client/types";
import { durationString } from "../../utils/durationString";

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
        Prefer: "code=200, example=example_1",
      },
    },
  });

  const groups = useUserListRequestAccessGroups(requestId, {
    swr: { refreshInterval: 10000 },
    request: {
      baseURL: "http://127.0.0.1:3100",
      headers: {
        Prefer: "code=200, example=ex_1",
      },
    },
  });

  const grants = useUserListRequestAccessGroupGrants("TEST", {
    swr: { refreshInterval: 10000 },
    request: {
      baseURL: "http://127.0.0.1:3101",
      headers: {
        Prefer: "code=200, example=ex_1",
      },
    },
  });

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

        <Container
          maxW={{
            md: "container.lg",
          }}
        >
          <Stack spacing={4} mt={8}>
            <Flex px={2}>
              {request.data ? (
                <Flex>
                  <Avatar
                    size="sm"
                    src={request.data.user.picture}
                    name={
                      request.data.user.firstName +
                      " " +
                      request.data.user.lastName
                    }
                    mr={2}
                  />
                  <Box>
                    <Flex>
                      <Text textStyle="Body/Small">
                        {request.data.user.firstName +
                          " " +
                          request.data.user.lastName}
                      </Text>
                      <Text textStyle="Body/Small" ml={1} color="neutrals.500">
                        requested at&nbsp;
                        {request.data.createdAt}
                      </Text>
                    </Flex>
                    <Text textStyle="Body/Small" color="neutrals.500">
                      {request.data.user.email}
                    </Text>
                  </Box>
                </Flex>
              ) : (
                <Flex>
                  <SkeletonCircle size="24px" mr={4} />
                  <Box>
                    <Flex>
                      <SkeletonText noOfLines={1} h="12px" w="12ch" mr="4px" />
                      <SkeletonText noOfLines={1} h="12px" w="12ch" />
                    </Flex>
                    <SkeletonText noOfLines={1} h="12px" w="12ch" />
                  </Box>
                </Flex>
              )}
            </Flex>
            <Divider borderColor="neutrals.300" w="100%" />

            <Stack spacing={4}>
              {groups.data ? (
                groups.data.groups.map((group) => (
                  <AccessGroupItem key={group.id} group={group} />
                ))
              ) : (
                <Flex>
                  <SkeletonCircle size="24px" mr={4} />
                  <Box>
                    <Flex>
                      <SkeletonText noOfLines={1} h="12px" w="12ch" mr="4px" />
                      <SkeletonText noOfLines={1} h="12px" w="12ch" />
                    </Flex>
                    <SkeletonText noOfLines={1} h="12px" w="12ch" />
                  </Box>
                </Flex>
              )}
            </Stack>
          </Stack>

          <Code whiteSpace="pre-wrap" mt={32}>
            {JSON.stringify({ request, groups, grants }, null, 2)}
          </Code>
        </Container>
      </UserLayout>
    </div>
  );
};

type NewType = {
  group: AccessGroup;
};

type AccessGroupProps = NewType;

export const AccessGroupItem = ({ group }: AccessGroupProps) => {
  const grants = useUserListRequestAccessGroupGrants(group.id, {
    swr: { refreshInterval: 10000 },
    request: {
      baseURL: "http://127.0.0.1:3101",
      headers: {
        Prefer: "code=200, example=ex_1",
      },
    },
  });

  return (
    <Box bg="neutrals.100" borderColor="neutrals.300" rounded="lg">
      <Accordion allowToggle>
        <AccordionItem border="none">
          <AccordionButton
            p={2}
            bg="neutrals.100"
            roundedTop="md"
            borderBottomColor="neutrals.300"
            borderWidth="1px"
            sx={{
              "&[aria-expanded='false']": {
                roundedBottom: "md",
              },
            }}
          >
            <AccordionIcon boxSize="6" mr={2} />
            <Box
              as="span"
              flex="1"
              textAlign="left"
              sx={{
                p: { lineHeight: "120%", textStyle: "Body/Extra Small" },
              }}
            >
              <Text color="neutrals.700">Review Required</Text>
              <Text color="neutrals.500">
                Duration&nbsp;
                {durationString(group.time.maxDurationSeconds)}
              </Text>
            </Box>
            <ButtonGroup variant="brandSecondary" spacing={2}>
              <Button size="sm">Approve</Button>
              <Button size="sm">Reject</Button>
            </ButtonGroup>
          </AccordionButton>

          <AccordionPanel
            borderColor="neutrals.300"
            roundedBottom="md"
            borderWidth="1px"
            bg="white"
            p={0}
          >
            <Stack spacing={2} p={2}>
              {grants.data ? (
                grants.data.grants.map((grant) => (
                  <Flex
                    w="100%"
                    borderColor="neutrals.300"
                    rounded="md"
                    borderWidth="1px"
                    bg="white"
                    p={2}
                    pos="relative"
                  >
                    <ProviderIcon shortType="aws" mr={2} />
                    <Text textStyle="Body/Small" color="neutrals.500">
                      {grants.data.grants.length} grants
                    </Text>
                    <Code fontFamily="Roboto Mono"></Code>
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
                  </Flex>
                ))
              ) : (
                <Text>loading</Text>
              )}
            </Stack>
          </AccordionPanel>
        </AccordionItem>
      </Accordion>
    </Box>
  );
};

export default Home;
