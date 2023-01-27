import { Button, Container, Heading, Text } from "@chakra-ui/react";
import { Helmet } from "react-helmet";
import { useMatch } from "react-location";
import { UserLayout } from "../../components/Layout";
import { useAdminGetProvider } from "../../utils/common-fate-client/admin/admin";

const Provider = () => {
  const {
    params: { id },
  } = useMatch();

  const provider = useAdminGetProvider(id);

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
        <Heading>{provider.data?.name}</Heading>
        <Heading>{provider.data?.team}</Heading>
        <Heading>{provider.data?.status}</Heading>
        <Heading>{provider.data?.version}</Heading>
      </Container>
    </UserLayout>
  );
};
export default Provider;
