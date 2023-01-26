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

import {
  Provider,
  useListAllProviders,
} from "../../utils/registry-client/orval";
import { useAdminListProvidersetupsv2 } from "../../utils/common-fate-client/provider-setup/provider-setup";

/** `${provider.team}/${provider.name}` is the format that will be used for detail lookup on /provider/[id] routes */
export const uniqueProviderKey = (provider: Provider) =>
  `${provider.team}/${provider.name}/${provider.version}`;

const Providers = () => {
  const { data: providers } = useAdminListProvidersetupsv2();

  return (
    <UserLayout>
      <Helmet>
        <title>My Providers</title>
      </Helmet>
      <Container
        my={12}
        // This prevents unbounded widths for small screen widths
        minW={{ base: "100%", lg: "container.lg" }}
        overflowX="auto"
      >
        <Heading>My Providers</Heading>
        <SimpleGrid columns={2} spacing={4} p={0} mt={6}>
          {providers?.providerSetups?.map((provider) => {
            return (
              <Box
                key={provider.id}
                as="button"
                className="group"
                textAlign="center"
                bg="neutrals.100"
                p={6}
                rounded="md"
                data-testid={"provider_" + provider.id}
                position="relative"
                _disabled={{
                  opacity: "0.5",
                }}
              >
                <LinkOverlay
                  href={`/providers/${provider.id}`}
                  as={Link}
                  to={`/providers/${provider.id}`}
                >
                  {/* <ProviderIcon type={provider.name} mb={3} h="8" w="8" /> */}

                  <Text textStyle="Body/SmallBold" color="neutrals.700">
                    {provider.id}
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
