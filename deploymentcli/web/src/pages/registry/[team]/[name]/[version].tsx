import { Box, Button, Container, Flex, Stack, Text } from "@chakra-ui/react";
import React, { useState } from "react";
import { Helmet } from "react-helmet";
import { useMatch, useNavigate } from "react-location";
import { ProviderIcon } from "../../../../components/icons/providerIcon";
import { UserLayout } from "../../../../components/Layout";
import { createProvider } from "../../../../utils/local-client/deploymentcli/deploymentcli";
import { useGetProvider } from "../../../../utils/registry-client/orval";

const RegistryProvider = () => {
  const {
    params: { team, name, version },
  } = useMatch();

  const { data: provider } = useGetProvider(team, name, version);

  // deployment loading state
  const [loading, setLoading] = useState(false);

  const navigate = useNavigate();

  const handleClick = async () => {
    // call deployment CLI, get the stack id...
    // stackId = await deployCLI.create()
    setLoading(true);
    createProvider({ name, team, version, alias: "" })
      .then((prov) => {
        // redirect to the /setup page using react-location
        navigate({ to: `/providers/${prov.stackId}/setup` });
      })
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
        minW={{ base: "100%", md: "container.sm" }}
        overflowX="auto"
      >
        <Stack
          className="group"
          textAlign="center"
          bg="neutrals.100"
          p={6}
          rounded="md"
          position="relative"
          _disabled={{
            opacity: "0.5",
          }}
          spacing={4}
          flexDir="column"
          justifyItems="flex-start"
        >
          <Flex flexDir="row" alignItems="center">
            <ProviderIcon type={name} mr={3} h="8" w="8" />
            <Text textStyle="Body/SmallBold" color="neutrals.700">
              {`${team}/${name}@${version}`}
            </Text>
          </Flex>

          <Button
            isLoading={loading}
            loadingText="Deploying..."
            onClick={handleClick}
            w="min-content"
          >
            Deploy Provider
          </Button>
        </Stack>

        <Box mt={12}>
          <Text textStyle="Heading/H1">Docs</Text>
        </Box>
      </Container>
    </UserLayout>
  );
};
export default RegistryProvider;
