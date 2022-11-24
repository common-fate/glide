import { CloseIcon, SmallAddIcon } from "@chakra-ui/icons";
import {
  Button,
  Center,
  CircularProgress,
  Code,
  Container,
  Flex,
  HStack,
  IconButton,
  LinkBox,
  LinkOverlay,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  Stack,
  Text,
  useDisclosure,
} from "@chakra-ui/react";
import { useMemo, useState } from "react";
import { Helmet } from "react-helmet";
import { Link } from "react-location";
import { Column } from "react-table";
import { CFCode } from "../../../components/CodeInstruction";
import { AdminLayout } from "../../../components/Layout";
import { TableRenderer } from "../../../components/tables/TableRenderer";
import {
  useListProviders,
  useListProvidersetups,
} from "../../../utils/backend-client/admin/admin";
import { deleteProvidersetup } from "../../../utils/backend-client/default/default";
import { Provider, ProviderSetup } from "../../../utils/backend-client/types";

const AdminProvidersTable = () => {
  const { data } = useListProviders();

  const cols: Column<Provider>[] = useMemo(
    () => [
      {
        accessor: "id",
        Header: "ID",
      },
      {
        accessor: "type",
        Header: "Type",
      },
    ],
    []
  );

  return TableRenderer<Provider>({
    columns: cols,
    data: data,
    emptyText: "No providers have been set up yet.",
    linkTo: false,
  });
};

const Providers = () => {
  const { data } = useListProvidersetups();

  const setups = data?.providerSetups ?? [];

  return (
    <AdminLayout>
      <Helmet>
        <title>Providers</title>
      </Helmet>
      <Container
        my={12}
        // This prevents unbounded widths for small screen widths
        minW={{ base: "100%", xl: "container.xl" }}
        overflowX="auto"
      >
        {setups.length > 0 && (
          <Stack p={1}>
            {setups.map((s) => (
              <ProviderSetupBanner setup={s} key={s.id} />
            ))}
          </Stack>
        )}
        <Button
          my={5}
          size="sm"
          variant="ghost"
          leftIcon={<SmallAddIcon />}
          as={Link}
          to="/admin/providers/setup"
          id="new-provider-button"
        >
          New Access Provider
        </Button>
        <AdminProvidersTable />
        <HStack mt={2} spacing={1} w="100%" justify={"center"}>
          <Text textStyle={"Body/ExtraSmall"}>
            View the full configuration of each access provider in your{" "}
          </Text>
          <Code fontSize={"12px"}>granted-deployment.yml</Code>
          <Text textStyle={"Body/ExtraSmall"}>file.</Text>
        </HStack>
      </Container>
    </AdminLayout>
  );
};

interface ProviderSetupBannerProps {
  setup: ProviderSetup;
}

const ProviderSetupBanner: React.FC<ProviderSetupBannerProps> = ({ setup }) => {
  const stepsOverview = setup.steps ?? [];
  const { data, mutate } = useListProvidersetups();
  const { onOpen, isOpen, onClose } = useDisclosure();
  const [loading, setLoading] = useState(false);

  const handleCancelSetup = async () => {
    setLoading(true);
    await deleteProvidersetup(setup.id);
    const oldSetups = data?.providerSetups ?? [];
    void mutate({
      providerSetups: [...oldSetups.filter((s) => s.id !== setup.id)],
    });
    setLoading(false);
    onClose();
  };

  const completedSteps = stepsOverview.filter((s) => s.complete).length;

  const completedPercentage =
    stepsOverview.length ?? 0 > 0
      ? (completedSteps / stepsOverview.length) * 100
      : 0;

  return (
    <LinkBox
      as={Flex}
      position="relative"
      justify="space-between"
      bg="neutrals.100"
      rounded="md"
      p={8}
      flexDirection={{ base: "column", md: "row" }}
    >
      <LinkOverlay as={Link} to={"/admin/providers/setup/" + setup.id}>
        <Stack>
          <Text textStyle={"Body/Medium"}>Continue setting up {setup.id}</Text>
          <Text>
            {setup.type}@{setup.version}
          </Text>
        </Stack>
      </LinkOverlay>
      <HStack spacing={3}>
        <Text>
          {completedSteps} of {setup.steps.length} steps complete
        </Text>
        <CircularProgress value={completedPercentage} color="#449157" />
      </HStack>
      <IconButton
        position="absolute"
        top={1}
        right={1}
        size="xs"
        variant={"unstyled"}
        onClick={(e) => {
          e.stopPropagation();
          onOpen();
        }}
        icon={<CloseIcon />}
        aria-label="Cancel setup"
      />
      <Modal isOpen={isOpen} onClose={onClose}>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>Cancel setting up {setup.id}</ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            Are you sure you want to stop setting up this provider? You'll lose
            any configuration values that we've stored.
          </ModalBody>

          <ModalFooter>
            <Button
              variant={"solid"}
              colorScheme="red"
              rounded="full"
              mr={3}
              onClick={handleCancelSetup}
              isLoading={loading}
            >
              Stop setup
            </Button>
            <Button
              variant={"brandSecondary"}
              onClick={onClose}
              isDisabled={loading}
            >
              Go back
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </LinkBox>
  );
};

export default Providers;
