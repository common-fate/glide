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
import { Link, useNavigate } from "react-location";
import { ProviderIcon } from "../../../../components/icons/providerIcon";
import { AdminLayout } from "../../../../components/Layout";
import { createProvidersetup } from "../../../../utils/backend-client/default/default";
import { registeredProviders } from "../../../../utils/providerRegistry";

const Page = () => {
  const navigate = useNavigate();
  // used to show a spinner when the provider is being initialised.
  const [providerLoading, setProviderLoading] = useState<string>();

  const createProvider = async (providerType: string) => {
    setProviderLoading(providerType);
    const res = await createProvidersetup({
      providerType,
    });
    navigate({ to: `/admin/providers/setup/${res.id}` });
  };

  return (
    <AdminLayout>
      <Center borderBottom="1px solid" borderColor="neutrals.200" h="80px">
        <IconButton
          as={Link}
          aria-label="Go back"
          pos="absolute"
          left={4}
          icon={<ArrowBackIcon />}
          rounded="full"
          variant="ghost"
          to="/admin/providers"
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
          {registeredProviders.map((provider) => (
            <Box
              key={provider.type}
              as="button"
              className="group"
              textAlign="center"
              bg="neutrals.100"
              p={6}
              rounded="md"
              data-testid={"provider_" + provider.type}
              onClick={() => createProvider(provider.type)}
              position="relative"
              disabled={providerLoading !== undefined}
              _disabled={{
                opacity: "0.5",
              }}
            >
              {providerLoading === provider.type && (
                <Spinner size="xs" position="absolute" right={2} top={2} />
              )}
              <ProviderIcon shortType={provider.shortType} mb={3} h="8" w="8" />

              <Text textStyle="Body/SmallBold" color="neutrals.700">
                {provider.name}
              </Text>
            </Box>
          ))}
        </SimpleGrid>
      </Container>
    </AdminLayout>
  );
};

export default Page;
