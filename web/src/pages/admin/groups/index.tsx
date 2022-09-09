import { Container, Stack } from "@chakra-ui/react";
import { AdminLayout } from "../../../components/Layout";
import { GroupsTable } from "../../../components/tables/GroupsTable";

const Index = () => {
  return (
    <AdminLayout>
      <Container
        my={12}
        // This prevents unbounded widths for small screen widths
        minW={{ base: "100%", xl: "container.xl" }}
        overflowX="auto"
      >
        <GroupsTable />
      </Container>
    </AdminLayout>
  );
};

export default Index;
