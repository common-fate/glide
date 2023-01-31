import { Button, Container, Heading, Text } from "@chakra-ui/react";
import { useState } from "react";
import { Helmet } from "react-helmet";
import { useMatch } from "react-location";
import { UserLayout } from "../../components/Layout";
import {
  deleteProvider,
  useGetProvider,
} from "../../utils/local-client/deploymentcli/deploymentcli";

const Provider = () => {
  const {
    params: { id },
  } = useMatch();

  const provider = useGetProvider(id);

  const [loading, setLoading] = useState(false);

  const handleDelete = () => {
    setLoading(true);
    deleteProvider(id).finally(() => {
      setLoading(false);
    });
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
        <Button onClick={handleDelete} isLoading={loading}>
          Delete
        </Button>
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
