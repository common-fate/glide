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
  useAdminListProviders,
  useAdminListProvidersetups,
  adminDeleteProvidersetup,
  useAdminListProvidersV2,
} from "../../../utils/backend-client/admin/admin";
import {
  Provider,
  ProviderSetup,
  ProviderV2,
} from "../../../utils/backend-client/types";

const AdminProvidersTable = () => {
  const { data } = useAdminListProvidersV2();

  const cols: Column<ProviderV2>[] = useMemo(
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

  return TableRenderer<ProviderV2>({
    columns: cols,
    data: data,
    emptyText: "No providers have been set up yet.",
    linkTo: false,
  });
};

const Providers = () => {
  const { data } = useAdminListProvidersetups();

  const setups = data?.providerSetups ?? [];

  return (
    <AdminLayout>
      <Helmet>
        <title>Providers V2</title>
      </Helmet>
      <Container
        my={12}
        // This prevents unbounded widths for small screen widths
        minW={{ base: "100%", xl: "container.xl" }}
        overflowX="auto"
      >
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
          <Code fontSize={"12px"}>deployment.yml</Code>
          <Text textStyle={"Body/ExtraSmall"}>file.</Text>
        </HStack>
      </Container>
    </AdminLayout>
  );
};

export default Providers;
