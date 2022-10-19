import { ArrowBackIcon } from "@chakra-ui/icons";
import { Center, Container, IconButton, Stack, Text } from "@chakra-ui/react";
import { MakeGenerics, useMatch, useSearch, Link } from "react-location";
import { AuditLog } from "../../components/AuditLog";
import { UserLayout } from "../../components/Layout";
import {
  RequestAccessInstructions,
  RequestAccessToken,
  RequestCancelButton,
  RequestDetails,
  RequestDisplay,
  RequestOverridableTime,
  RequestRequestor,
  RequestReview,
  RequestRevoke,
  RequestTime,
} from "../../components/Request";
import { useUser } from "../../utils/context/userContext";

import { useUserGetRequest } from "../../utils/backend-client/end-user/end-user";

type MyLocationGenerics = MakeGenerics<{
  Search: {
    action?: "approve" | "close";
  };
}>;

const Home = () => {
  const {
    params: { id: requestId },
  } = useMatch();
  const { data, mutate } = useUserGetRequest(requestId, {
    swr: { refreshInterval: 10000 },
  });
  const search = useSearch<MyLocationGenerics>();
  const { action } = search;
  const Content = () => {
    if (data?.canReview && data.status == "PENDING") {
      return (
        <RequestDisplay request={data}>
          <RequestDetails>
            {data?.canReview ? <RequestOverridableTime /> : <RequestTime />}
            <RequestRequestor />
          </RequestDetails>
          <RequestReview
            onSubmitReview={mutate}
            focus={action}
            canReview={!!data.canReview}
          />
        </RequestDisplay>
      );
    }

    const user = useUser();
    return (
      <RequestDisplay request={data}>
        <RequestDetails>
          <RequestTime />

          {user.user?.id === data?.requestor && <RequestAccessInstructions />}
          {user.user?.id === data?.requestor && (
            <RequestAccessToken reqId={data ? data.id : ""} />
          )}
          <RequestCancelButton />
          <RequestRevoke onSubmitRevoke={mutate} />
        </RequestDetails>
      </RequestDisplay>
    );
  };
  return (
    <div>
      <UserLayout>
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
            to={data?.canReview ? "/reviews?status=pending" : "/requests"}
          />

          <Text as="h4" textStyle="Heading/H4">
            Request details
          </Text>
        </Center>
        {/* Main content */}
        <Container maxW="container.xl" py={16}>
          <Stack spacing={12} direction={{ base: "column", md: "row" }}>
            <Content />
            <AuditLog request={data} />
          </Stack>
        </Container>
      </UserLayout>
    </div>
  );
};

export default Home;
