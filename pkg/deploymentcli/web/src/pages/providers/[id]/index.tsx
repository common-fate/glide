import { Button, Container, Heading, Text } from "@chakra-ui/react";
import { Helmet } from "react-helmet";
import { useMatch } from "react-location";
import { UserLayout } from "../../../components/Layout";
import { useAdminGetProvidersetupv2 } from "../../../utils/common-fate-client/admin/admin";
import { Link } from "react-location";

const Provider = () => {
  const {
    params: { id },
  } = useMatch();

  const provider = useAdminGetProvidersetupv2(id);

  return (
    <UserLayout>
      <Helmet>
        <title>Provider</title>
      </Helmet>
      <Container
        my={12}
        // This prevents unbounded widths for small screen widths
        minW={{ base: "100%", lg: "container.lg" }}
        overflowX="auto"
      >
        <Link>
          <Button as={Link} to={`/providers/${id}/setup`}>
            Configure
          </Button>
        </Link>
        <Heading>{provider.data?.name}</Heading>
        <Heading>{provider.data?.team}</Heading>
        <Heading>{provider.data?.status}</Heading>
        <Heading>{provider.data?.version}</Heading>
      </Container>
    </UserLayout>
  );
};
export default Provider;
