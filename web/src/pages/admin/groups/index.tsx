import { Container, Stack } from "@chakra-ui/react";
import { Helmet } from "react-helmet";
import { ManagementGroupsTable } from "../../../components/tables/ManagementGroupsTable";
import { AdminLayout } from "../../../components/Layout";
import { GroupsTable } from "../../../components/tables/GroupsTable";

const Index = () => {
  return (
    <AdminLayout>
      <Helmet>
        <title>Groups</title>
      </Helmet>
      <Container
        my={12}
        // This prevents unbounded widths for small screen widths
        minW={{ base: "100%", xl: "container.xl" }}
        overflowX="auto"
      >
        Identity Provider Groups
        <GroupsTable />
        Management Groups
        <ManagementGroupsTable />
      </Container>
    </AdminLayout>
  );
};

export default Index;
