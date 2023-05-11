import { ArrowBackIcon } from "@chakra-ui/icons";
import {
  Box,
  Center,
  Container,
  IconButton,
  Skeleton,
  Text,
  VStack,
} from "@chakra-ui/react";
import { Helmet } from "react-helmet";
import { Link, useMatch } from "react-location";
import UpdateAccessRuleForm from "../../../components/forms/access-rule/UpdateForm";
import { AdminLayout } from "../../../components/Layout";
import { useAdminGetAccessRule } from "../../../utils/backend-client/admin/admin";

const Index = () => {
  const {
    params: { id: ruleId },
  } = useMatch();
  // const ruleId = typeof query?.id == "string" ? query.id : "";
  const { data, isValidating, error } = useAdminGetAccessRule(ruleId);
  return (
    <>
      <AdminLayout>
        <Helmet>
          <title>Access Rule</title>
        </Helmet>
        <Center borderBottom="1px solid" borderColor="neutrals.200" h="80px">
          <IconButton
            as={Link}
            to={"/admin/access-rules"}
            aria-label="Go back"
            pos="absolute"
            left={4}
            icon={<ArrowBackIcon />}
            rounded="full"
            variant="ghost"
          />

          <Text as="h4" textStyle="Heading/H4">
            Edit Access Rule
          </Text>
        </Center>
        <Box mb={24}>
          {!data && (
            <Container pt={12} maxW="container.md">
              <VStack w="100%" spacing={6}>
                <Skeleton h="150px" w="100%" rounded="md" />
                <Skeleton h="150px" w="100%" rounded="md" />
                <Skeleton h="150px" w="100%" rounded="md" />
                <Skeleton h="150px" w="100%" rounded="md" />
                <Skeleton h="150px" w="100%" rounded="md" />
              </VStack>
            </Container>
          )}

          {data && <UpdateAccessRuleForm data={data} />}
        </Box>
      </AdminLayout>
    </>
  );
};

export default Index;
