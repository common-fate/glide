import { Container } from "@chakra-ui/react";

import { AdminLayout } from "../../../components/Layout";
import ProviderSetupNotice from "../../../components/ProviderSetupNotice";
import { AccessRuleTable } from "../../../components/tables/AccessRuleTable";

const Index = () => {
  return (
    <AdminLayout>
      <Container
        my={12}
        // This prevents unbounded widths for small screen widths
        minW={{ base: "100%", xl: "container.xl" }}
        overflowX="auto"
      >
        <ProviderSetupNotice />
        <AccessRuleTable />
      </Container>
    </AdminLayout>
  );
};

export default Index;
