import {
  Box,
  Container,
  Heading,
  LinkOverlay,
  SimpleGrid,
  Text,
} from "@chakra-ui/react";
import { Helmet } from "react-helmet";
import { Link } from "react-location";
import { UserLayout } from "../../components/Layout";
import { ProviderIcon } from "../../components/icons/providerIcon";
import { useListProviders } from "../../utils/local-client/deploymentcli/deploymentcli";
import { ProviderV2 } from "../../utils/local-client/types/openapi.yml";

/** `${provider.team}/${provider.name}` is the format that will be used for detail lookup on /provider/[id] routes */
export const uniqueProviderKey = (provider: ProviderV2) =>
  `${provider.team}/${provider.name}/${provider.version}`;

const Providers = () => {
  const { data: providers } = useListProviders();

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
        <Heading>Deploy a provider</Heading>
        <SimpleGrid columns={2} spacing={4} p={0} mt={6}>
          {providers &&
            providers.map((provider) => {
              const id = uniqueProviderKey(provider);
              return (
                <Box
                  key={id}
                  as="button"
                  className="group"
                  textAlign="center"
                  bg="neutrals.100"
                  p={6}
                  rounded="md"
                  data-testid={"provider_" + id}
                  position="relative"
                  _disabled={{
                    opacity: "0.5",
                  }}
                >
                  <LinkOverlay
                    href={`/registry/${id}`}
                    as={Link}
                    to={`/registry/${id}`}
                  >
                    <ProviderIcon type={provider.name} mb={3} h="8" w="8" />
                    <Text textStyle="Body/SmallBold" color="neutrals.700">
                      {`${provider.team}/${provider.name}@${provider.version}`}
                    </Text>
                  </LinkOverlay>
                </Box>
              );
            })}
        </SimpleGrid>
      </Container>
    </UserLayout>
  );
};

export default Providers;
