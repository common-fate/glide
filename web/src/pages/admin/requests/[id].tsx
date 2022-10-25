import { ArrowBackIcon } from "@chakra-ui/icons";
import { Center, Container, IconButton, Stack, Text } from "@chakra-ui/react";
import { Helmet } from "react-helmet";
import { Link, useMatch } from "react-location";
import { AuditLog } from "../../../components/AuditLog";
import { AdminLayout } from "../../../components/Layout";
import {
  RequestDetails,
  RequestDisplay,
  RequestOverridableTime,
  RequestRequestor,
  RequestReview,
  RequestRevoke,
  RequestTime,
} from "../../../components/Request";
import { useAdminGetRequest } from "../../../utils/backend-client/end-user/end-user";

const Home = () => {
  const {
    params: { id: requestId },
  } = useMatch();

  const { data, mutate } = useAdminGetRequest(requestId, {
    swr: { refreshInterval: 10000 },
  });
  return (
    <div>
      <Helmet>
        <title>{requestId}</title>
      </Helmet>
      <AdminLayout>
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
            to="/admin/requests"
          />

          <Text as="h4" textStyle="Heading/H4">
            Request details
          </Text>
        </Center>
        {/* Main content */}
        <Container maxW="container.xl" py={16}>
          <Stack spacing={12} direction={{ base: "column", md: "row" }}>
            <RequestDisplay request={data}>
              <RequestDetails>
                {data?.canReview ? <RequestOverridableTime /> : <RequestTime />}
                <RequestRequestor />
              </RequestDetails>
              <RequestReview
                onSubmitReview={mutate}
                canReview={!!data?.canReview}
              />
              <RequestRevoke onSubmitRevoke={mutate} />
            </RequestDisplay>
            <AuditLog request={data} />
          </Stack>
        </Container>
      </AdminLayout>
    </div>
  );
};

export default Home;
