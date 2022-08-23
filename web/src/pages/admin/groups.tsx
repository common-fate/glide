import { Container, Stack } from "@chakra-ui/react";
import { AdminLayout } from "../../components/Layout";
import { GroupsTable } from "../../components/tables/GroupsTable";

const Index = () => {
  return (
    <AdminLayout>
      <Container maxW="1200px" pb={5}>
        <Stack padding={{ base: 2, md: "50px" }}>
          {/* <CFSearchBar
            placeholderMessage={"Search for a user"}
            setSearchVal={(val) => setInput(val)}
          /> */}
          <GroupsTable />
        </Stack>
      </Container>
    </AdminLayout>
  );
};

export default Index;
