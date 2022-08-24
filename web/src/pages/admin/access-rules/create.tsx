import { Box, Flex, Text } from "@chakra-ui/react";
import CreateAccessRuleForm from "../../../components/forms/access-rule/CreateForm";

import { AdminLayout } from "../../../components/Layout";
const Index = () => {
  return (
    <AdminLayout>
      <Box minH="90vh" pb={12}>
        <Flex
          h={20}
          borderBottomWidth="1px"
          position="relative"
          justify="center"
          align="center"
        >
          <Text textStyle={"Heading/H4"}>New Access Rule</Text>
        </Flex>
        <CreateAccessRuleForm />
      </Box>
    </AdminLayout>
  );
};

export default Index;
