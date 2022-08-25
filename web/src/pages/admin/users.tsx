import { Container, Stack } from "@chakra-ui/react";
import { AdminLayout } from "../../components/Layout";
import { UsersTable } from "../../components/tables/UsersTable";

const Index = () => {
  return (
    <AdminLayout>
      <Container maxW="1200px" pb={5}>
        <Stack padding={{ base: 2, md: "50px" }}>
          <UsersTable />
        </Stack>
      </Container>
    </AdminLayout>
  );
};

export default Index;
