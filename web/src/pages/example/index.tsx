import { Box, Container, Heading, Select, VStack } from "@chakra-ui/react";
import { Helmet } from "react-helmet";
import { UserLayout } from "../../components/Layout";
import { UserReviewsTable } from "../../components/tables/UserReviewsTable";
import {
  useUserListEntitlementResources,
  useUserListEntitlements,
} from "../../utils/backend-client/default/default";

const Home = () => {
  const { data } = useUserListEntitlements();
  const { data: resources } = useUserListEntitlementResources({
    kind: "",
    name: "",
    publisher: "",
    resourceType: "Account",
    version: "",
  });
  return (
    <div>
      <Helmet>
        <title>Example</title>
      </Helmet>
      <UserLayout>
        <Box overflow="auto">
          <Container minW="864px" maxW="container.xl">
            <VStack>
              <Heading>Select an entitlement</Heading>
              <Select>
                {data?.map((e) => {
                  return (
                    <option>{`${e.Kind.publisher} ${e.Kind.name} ${e.Kind.version} ${e.Kind.kind}`}</option>
                  );
                })}
              </Select>
            </VStack>
            <VStack>
              <Heading>Select an option for Account</Heading>
              <Select>
                {resources?.resources.map((e) => {
                  return <option>{`${e.name} ${e.value} `}</option>;
                })}
              </Select>
            </VStack>
          </Container>
        </Box>
      </UserLayout>
    </div>
  );
};

export default Home;
