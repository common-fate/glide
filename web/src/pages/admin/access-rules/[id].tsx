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
import EditAccessRuleForm from "../../../components/forms/access-rule/AccessRuleForm";
import { AdminLayout } from "../../../components/Layout";
import { useMatch } from "react-location";
import { useAdminGetAccessRule } from "../../../utils/backend-client/admin/admin";
type Props = {};

const Index = (props: Props) => {
  const {
    params: { id: ruleId },
  } = useMatch();
  // const ruleId = typeof query?.id == "string" ? query.id : "";
  const { data, isValidating, error } = useAdminGetAccessRule(ruleId);
  return (
    <AdminLayout>
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
            <EditAccessRuleForm data={data} readOnly={true} />
          ))}
        {data?.isCurrent && data?.status === "ACTIVE" && (
          <EditAccessRuleForm data={data} />
        )}
      </Box>
    </AdminLayout>
  );
};

export default Index;
