import { ArrowBackIcon } from "@chakra-ui/icons";
import {
  Avatar,
  Badge,
  Center,
  Container,
  Flex,
  IconButton,
  SkeletonText,
  Spacer,
  Text,
  Tooltip,
  VStack,
  Wrap,
  WrapItem,
} from "@chakra-ui/react";
import { Link, useMatch } from "react-location";

import { AdminLayout } from "../../../components/Layout";
import { useGetGroup } from "../../../utils/backend-client/admin/admin";
import { useGetUser } from "../../../utils/backend-client/end-user/end-user";
const GroupDisplay: React.FC<{ groupId: string }> = ({ groupId }) => {
  const { data } = useGetGroup(encodeURIComponent(groupId));
  return (
    <Tooltip label={data?.description}>
      <Flex
        cursor="help"
        textStyle={"Body/Small"}
        rounded="full"
        bg="neutrals.300"
        py={1}
        px={4}
      >
        {data?.name}
      </Flex>
    </Tooltip>
  );
};
const Index = () => {
  const {
    params: { id: userId },
  } = useMatch();
  const { data: user, isValidating, error } = useGetUser(userId);

  const Content = () => {
    if (user?.id === undefined) {
      return (
        <>
          <VStack>
            <Text>Name</Text>
            <SkeletonText noOfLines={1} />
            <Text>Email</Text>
            <SkeletonText noOfLines={1} />
            <Text>Groups</Text>
            <SkeletonText noOfLines={3} />
          </VStack>
        </>
      );
    }
    return (
      <>
        <VStack align={"left"} spacing={1}>
          <Text textStyle="Body/Medium">Name</Text>
          <Text textStyle="Body/Small">{`${user.firstName} ${user.lastName}`}</Text>
          <Text textStyle="Body/Medium">Email</Text>
          <Text textStyle="Body/Small">{user.email}</Text>
          <Text textStyle="Body/Medium">Groups</Text>
          <Wrap>
            {user.groups.map((g) => {
              return (
                <WrapItem key={g}>
                  <GroupDisplay groupId={g} />
                </WrapItem>
              );
            })}
          </Wrap>
        </VStack>
        <Spacer />
        <Avatar
          src={user.picture}
          name={
            user.firstName ? `${user.firstName} ${user.lastName}` : user.email
          }
          boxSize="200px"
        />
      </>
    );
  };
  return (
    <AdminLayout>
      <Center borderBottom="1px solid" borderColor="neutrals.200" h="80px">
        <IconButton
          as={Link}
          aria-label="Go back"
          pos="absolute"
          left={4}
          icon={<ArrowBackIcon />}
          rounded="full"
          variant="ghost"
          to={"/admin/users"}
        />

        <Text as="h4" textStyle="Heading/H4">
          User Details
        </Text>
      </Center>
      {/* Main content */}
      <Container maxW="container.xl" py={16}>
        <Center>
          <Flex
            direction={["column", "row"]}
            rounded="md"
            bg="neutrals.100"
            w={{ base: "100%", md: "500px", lg: "716px" }}
            p={8}
          >
            <Content />
          </Flex>
        </Center>
      </Container>
    </AdminLayout>
  );
};

export default Index;
