import { SmallAddIcon } from "@chakra-ui/icons";
import { Button, Center, Code, Container, Text } from "@chakra-ui/react";
import { useMemo } from "react";
import { Link } from "react-location";
import { Column } from "react-table";
import { AdminLayout } from "../../../components/Layout";
import { TableRenderer } from "../../../components/tables/TableRenderer";
import { useListProviders } from "../../../utils/backend-client/admin/admin";
import { Provider } from "../../../utils/backend-client/types";

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
  return (
    <AdminLayout>
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
        <Center>
          <Text textStyle={"Body/ExtraSmall"}>
            View the full configuration of each access provider in your
            <Code fontSize={"12px"}>granted-deployment.yml</Code> file.
          </Text>
        </Center>
      </Container>
    </AdminLayout>
  );
};

export default Providers;
