import { Button, Container, Heading, Text } from "@chakra-ui/react";
import { useState } from "react";
import { Helmet } from "react-helmet";
import { useMatch } from "react-location";
import { UserLayout } from "../../../../components/Layout";
import { createProvider } from "../../../../utils/local-client/deploymentcli/deploymentcli";
import { useGetProvider } from "../../../../utils/registry-client/orval";

const RegistryProvider = () => {
  const {
    params: { team, name, version },
  } = useMatch();

  const provider = useGetProvider(team, name, version);

  // deployment loading state
  const [loading, setLoading] = useState(false);

  const handleClick = async () => {
    // call deployment CLI, get the stack id...
    // stackId = await deployCLI.create()
    setLoading(true);
    createProvider({ name, team, version, alias: "" })
      .then(() => {})
      .finally(() => {
        setLoading(false);
      });
  };

  return (
    <UserLayout>
      <Helmet>
        <title>Registry</title>
      </Helmet>
      <Container
        my={12}
        // This prevents unbounded widths for small screen widths
        minW={{ base: "100%", lg: "container.lg" }}
        overflowX="auto"
      >
        <Heading>
          {team}/{name}/{version}
        </Heading>
        <Button isLoading={loading} onClick={handleClick}>
          Create/Deploy Provider
        </Button>
        <Text>Schema</Text>
        <Text whiteSpace={"pre-wrap"} as={"pre"}>
          {JSON.stringify(provider.data?.schema, undefined, 2)}
        </Text>
      </Container>
    </UserLayout>
  );
};
export default RegistryProvider;
