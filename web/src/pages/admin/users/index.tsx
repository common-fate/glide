import { Container } from "@chakra-ui/react";
import { AdminLayout } from "../../../components/Layout";
import { UsersTable } from "../../../components/tables/UsersTable";

const Index = () => {
  return (
    <AdminLayout>
      <Container
        my={12}
        // This prevents unbounded widths for small screen widths
        minW={{ base: "100%", xl: "container.xl" }}
        overflowX="auto"
      >
        <UsersTable />
      </Container>
    </AdminLayout>
  );
};

export default Index;
