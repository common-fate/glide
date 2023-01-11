import { ArrowBackIcon } from "@chakra-ui/icons";
import {
  Badge,
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
import { ProviderIcon } from "../../../../components/icons/providerIcon";
import { AdminLayout } from "../../../../components/Layout";
import { adminCreateProvidersetup } from "../../../../utils/backend-client/admin/admin";
import { registeredProviders } from "../../../../utils/providerRegistry";

const Page = () => {
  const navigate = useNavigate();
  // used to show a spinner when the provider is being initialised.
  const [providerLoading, setProviderLoading] = useState<string>();

  const createProvider = async (providerType: string) => {
    setProviderLoading(providerType);
    const res = await adminCreateProvidersetup({
      providerType,
    });
    navigate({ to: `/admin/providers/setup/${res.id}` });
  };

  return (
    <AdminLayout>
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
              <ProviderIcon type={provider.type} mb={3} h="8" w="8" />

              <Text textStyle="Body/SmallBold" color="neutrals.700">
                {provider.name}
              </Text>
              {provider?.alpha && (
                <Badge
                  variant="outline"
                  position="absolute"
                  top={4}
                  right={4}
                  colorScheme="gray"
                >
                  ALPHA
                </Badge>
              )}
            </Box>
          ))}
        </SimpleGrid>
      </Container>
    </AdminLayout>
  );
};

export default Page;
