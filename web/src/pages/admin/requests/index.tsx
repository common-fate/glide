import { Container } from "@chakra-ui/react";
import { Helmet } from "react-helmet";
import { AdminLayout } from "../../../components/Layout";
import { AdminRequestsTable } from "../../../components/tables/AdminRequestsTable";

const Requests = () => {
  return (
    <AdminLayout>
      <Helmet>
        <title>Requests</title>
      </Helmet>
      <Container
        my={12}
        // This prevents unbounded widths for small screen widths
        minW={{ base: "100%", xl: "container.xl" }}
        overflowX="auto"
      >
        <AdminRequestsTable />
      </Container>
    </AdminLayout>
  );
};

export default Requests;
