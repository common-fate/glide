import {
  Avatar,
  Box,
  Container,
  Flex,
  Heading,
  Skeleton,
} from "@chakra-ui/react";
import { useMatch } from "react-location";
import { CFCard } from "../../../components/CFCard";
import { AdminLayout } from "../../../components/Layout";
import { useGetUser } from "../../../utils/backend-client/end-user/end-user";
import { userName } from "../../../utils/userName";

const Index = () => {
  const {
    params: { id: userId },
  } = useMatch();
  const { data: user, isValidating, error } = useGetUser(userId);

  return (
    <AdminLayout>
      <Box bg="neutrals.100" minH="90vh">
        <Container pt={12}>
          {user && (
            <CFCard w="500px">
              <Flex>
                <Avatar
                  mr={4}
                  size="xl"
                  name={userName(user)}
                  src={user.picture}
                />
                <Box>
                  <Heading>{userName(user)}</Heading>
                  <Flex mt={1}>
                    Email:
                    <a href={"mailto:" + user.email}>&nbsp;{user.email}</a>
                  </Flex>
                </Box>
              </Flex>
            </CFCard>
          )}
          {!user && !isValidating && error && <>User not found</>}
          {isValidating && <Skeleton h="200px" w="500px" rounded="md" />}
        </Container>
      </Box>
    </AdminLayout>
  );
};

export default Index;
