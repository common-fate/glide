import { Container, Input } from "@chakra-ui/react";
import { useState } from "react";
import { Helmet } from "react-helmet";
import { AdminLayout } from "../../../components/Layout";
import { UsersTable } from "../../../components/tables/UsersTable";
import { useAdminListUsers } from "../../../utils/backend-client/admin/admin";

const Index = () => {
  // text state
  const [text, setText] = useState("");

  const userLookup = useAdminListUsers({
    query: text,
  });

  return (
    <AdminLayout>
      <Helmet>
        <title>Users</title>
      </Helmet>
      <Input
        type="text"
        value={text}
        onChange={(e) => setText(e.target.value)}
      />
      {JSON.stringify(userLookup.data)}
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
