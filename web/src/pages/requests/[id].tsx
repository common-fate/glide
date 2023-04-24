import { ArrowBackIcon } from "@chakra-ui/icons";
import { Center, Container, IconButton, Stack, Text } from "@chakra-ui/react";
import { MakeGenerics, useMatch, useSearch, Link } from "react-location";

import { UserLayout } from "../../components/Layout";

import { useUser } from "../../utils/context/userContext";

import { Helmet } from "react-helmet";
import { useEffect, useMemo, useState } from "react";
import { useUserGetRequest } from "../../utils/backend-client/default/default";

type MyLocationGenerics = MakeGenerics<{
  Search: {
    action?: "approve" | "close";
  };
}>;

const Home = () => {
  const {
    params: { id: requestId },
  } = useMatch();
  const { data, mutate, isValidating } = useUserGetRequest(requestId, {
    swr: { refreshInterval: 10000 },
  });
  const search = useSearch<MyLocationGenerics>();
  const { action } = search;

  const [cachedReq, setCachedReq] = useState(data);
  useEffect(() => {
    if (data !== undefined) setCachedReq(data);
    return () => {
      setCachedReq(undefined);
    };
  }, [data]);

  const user = useUser();

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
        {/* <Container maxW="container.xl" py={16}>
          <Stack spacing={12} direction={{ base: "column", md: "row" }}>
            {data?.canReview && data.status == "PENDING" ? (
              <RequestDisplay request={data} isValidating={isValidating}>
                <RequestDetails>
                  <RequestTime canReview={cachedReq?.canReview} />
                  <RequestRequestor />
                </RequestDetails>
                <RequestReview
                  onSubmitReview={mutate}
                  focus={action}
                  canReview={!!data.canReview}
                />
              </RequestDisplay>
            ) : (
              <RequestDisplay request={data} isValidating={isValidating}>
                <RequestDetails>
                  <RequestTime />

                  {user.user?.id === data?.requestor && (
                    <RequestAccessInstructions />
                  )}
                  {user.user?.id === data?.requestor && <RequestAccessToken />}
                  <RequestCancelButton />
                  <RequestRevoke onSubmitRevoke={mutate} />
                </RequestDetails>
              </RequestDisplay>
            )}
            <AuditLog request={data} />
          </Stack>
        </Container> */}
      </UserLayout>
    </div>
  );
};

export default Home;
