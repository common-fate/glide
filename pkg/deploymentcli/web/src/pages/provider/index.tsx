import { ArrowBackIcon } from "@chakra-ui/icons";
import {
  Box,
  Center,
  Container,
  IconButton,
  SimpleGrid,
  Spinner,
  Text,
} from "@chakra-ui/react";
import { useState } from "react";
import { Helmet } from "react-helmet";
import { Link, useNavigate } from "react-location";
import { ProviderIcon } from "../../components/icons/providerIcon";
import { UserLayout } from "../../components/Layout";
import { createProvidersetup } from "../../utils/backend-client/local/orval";
import {
  Provider,
  useListAllProviders,
} from "../../utils/backend-client/registry/orval";

export const providerKey = (provider: Provider): string => {
  return provider.team + provider.name + provider.version;
};
const Page = () => {
  const navigate = useNavigate();
  const { data } = useListAllProviders();

  // used to show a spinner when the provider is being initialised.
  const [providerLoading, setProviderLoading] = useState<string>();

  // @TODO: hook into this logic, but on the detail/[id] page
  const createProvider = async (provider: Provider) => {
    setProviderLoading(providerKey(provider));
    const res = await createProvidersetup({
      name: provider.name,
      team: provider.team,
      version: provider.version,
    });

    navigate({ to: `/setup/${res.id}` });
  };

  return (
    <UserLayout>
      <Helmet>
        <title>New Access Provider</title>
      </Helmet>
      <Center borderBottom="1px solid" borderColor="neutrals.200" h="80px">
        <IconButton
          as={Link}
          aria-label="Go back"
          pos="absolute"
          left={4}
          icon={<ArrowBackIcon />}
          rounded="full"
          variant="ghost"
          to="/"
        />
        <Text as="h4" textStyle="Heading/H4">
          New Access Provider
        </Text>
      </Center>
      <Container
        my={12}
        // This prevents unbounded widths for small screen widths
        minW={{ base: "100%", xl: "container.xl" }}
        overflowX="auto"
      >
        <SimpleGrid columns={2} spacing={4} p={1}>
          {data?.providers.map((provider) => {
            const key = providerKey(provider);
            return (
              <Box
                key={key}
                as="button"
                className="group"
                textAlign="center"
                bg="neutrals.100"
                p={6}
                rounded="md"
                data-testid={"provider_" + key}
                onClick={() => createProvider(provider)}
                position="relative"
                disabled={providerLoading !== undefined}
                _disabled={{
                  opacity: "0.5",
                }}
              >
                {providerLoading === key && (
                  <Spinner size="xs" position="absolute" right={2} top={2} />
                )}
                <ProviderIcon type="commonfate/aws-sso" mb={3} h="8" w="8" />

                <Text textStyle="Body/SmallBold" color="neutrals.700">
                  {`${provider.team}/${provider.name}@${provider.version}`}
                </Text>
                {/* {provider?.alpha && (
                  <Badge
                    variant="outline"
                    position="absolute"
                    top={4}
                    right={4}
                    colorScheme="gray"
                  >
                    ALPHA
                  </Badge>
                )} */}
              </Box>
            );
          })}
        </SimpleGrid>
      </Container>
    </UserLayout>
  );
};

export default Page;
