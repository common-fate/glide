import { Code, Container, HStack, Text } from "@chakra-ui/react";
import { useMemo } from "react";
import { Helmet } from "react-helmet";
import { Column } from "react-table";
import { AdminLayout } from "../../../components/Layout";
import { TableRenderer } from "../../../components/tables/TableRenderer";
import { useAdminListProviders } from "../../../utils/backend-client/admin/admin";
import { Provider } from "../../../utils/backend-client/types";

const AdminProvidersTable = () => {
  const { data } = useAdminListProviders();

  const cols: Column<Provider>[] = useMemo(
    () => [
      {
        accessor: "name",
        Header: "Name",
      },
      {
        accessor: "version",
        Header: "Version",
      },
      {
        accessor: "schema",
        Header: "Schema",
      },
      {
        accessor: "url",
        Header: "ARN URL",
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
      <Helmet>
        <title>Providers</title>
      </Helmet>
      <Container
        my={12}
        // This prevents unbounded widths for small screen widths
        minW={{ base: "100%", xl: "container.xl" }}
        overflowX="auto"
      >
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
