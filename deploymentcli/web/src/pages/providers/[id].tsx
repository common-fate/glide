import { Button, Container, Heading, Text } from "@chakra-ui/react";
import { useState } from "react";
import { Helmet } from "react-helmet";
import { useMatch } from "react-location";
import { UserLayout } from "../../components/Layout";
import {
  useAdminGetProvider,
  useAdminGetProviderv2,
} from "../../utils/common-fate-client/admin/admin";

const Provider = () => {
  const {
    params: { id },
  } = useMatch();

  const provider = useAdminGetProviderv2(id);

  // handleDelete
  // loading state
  const [loading, setLoading] = useState(false);

  const handleDelete = () => {
    setLoading(true);

    // adminDelete? ⭐

    setLoading(false);
  };

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
        {/* delete button */}
        <Button onClick={handleDelete}>Delete</Button>
        <Heading>{provider.data?.name}</Heading>
        <Heading>{provider.data?.team}</Heading>
        <Heading>{provider.data?.status}</Heading>
        <Heading>{provider.data?.version}</Heading>
        <Heading>{provider.data?.stackId}</Heading>
      </Container>
    </UserLayout>
  );
};
export default Provider;
