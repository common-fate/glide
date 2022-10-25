import { ArrowBackIcon } from "@chakra-ui/icons";
import {
  Box,
  Center,
  IconButton,
  Text,
  Container,
  Skeleton,
  useToast,
  VStack,
} from "@chakra-ui/react";
import { Link } from "react-location";
import UpdateAccessRuleForm from "../../../components/forms/access-rule/UpdateForm";
import { AdminLayout } from "../../../components/Layout";
import { useMatch } from "react-location";
import { useAdminGetAccessRule } from "../../../utils/backend-client/admin/admin";
import { Helmet } from "react-helmet";

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
          <title>{ruleId}</title>
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
            {data?.status === "ACTIVE" ? "Edit" : "View"} Access Rule
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

          {!data?.isCurrent ||
            (data.status === "ARCHIVED" && (
              <UpdateAccessRuleForm data={data} readOnly={true} />
            ))}
          {data?.isCurrent && data?.status === "ACTIVE" && (
            <UpdateAccessRuleForm data={data} />
          )}
        </Box>
      </AdminLayout>
    </>
  );
};

export default Index;
