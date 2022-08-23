import { Container } from "@chakra-ui/react";
import { AdminLayout } from "../../../components/Layout";
import { AdminRequestsTable } from "../../../components/tables/AdminRequestsTable";

const Requests = () => {
  return (
    <AdminLayout>
      <Container maxW="1200px" pb={5}>
        <AdminRequestsTable />
      </Container>
    </AdminLayout>
  );
};

export default Requests;
