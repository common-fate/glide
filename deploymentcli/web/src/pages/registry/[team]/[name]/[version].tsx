import { Button, Container, Heading, Text } from "@chakra-ui/react";
import { useState } from "react";
import { Helmet } from "react-helmet";
import { useMatch } from "react-location";
import { UserLayout } from "../../../../components/Layout";
import { postDeployment } from "../../../../utils/local-client/orval";
import { useGetProvider } from "../../../../utils/registry-client/orval";
import { adminCreateProviderv2 } from "../../../..//utils/common-fate-client/default/default";

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
    postDeployment({ name, team, version })
      .then(({ stackId }) => {
        adminCreateProviderv2({
          team,
          name,
          version,
          stackId,
          alias: "",
        })
          .then(() => {
            setLoading(false);
            // navigate to the provider page
          })
          .catch((e) => {
            console.log(e);
          });
      })
      .finally(() => {
        setLoading(false);
      });

    // now call CF to create the provider
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
        <Button onClick={handleClick}>Create/Deploy Provider</Button>
        <Text>Schema</Text>
        <Text whiteSpace={"pre-wrap"} as={"pre"}>
          {JSON.stringify(provider.data?.schema, undefined, 2)}
        </Text>
      </Container>
    </UserLayout>
  );
};
export default RegistryProvider;
