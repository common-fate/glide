import {
  ArrowBackIcon,
  CheckIcon,
  CloseIcon,
  EditIcon,
} from "@chakra-ui/icons";
import {
  Avatar,
  Center,
  Container,
  Flex,
  FormControl,
  FormLabel,
  HStack,
  IconButton,
  SkeletonText,
  Spacer,
  Text,
  Tooltip,
  useDisclosure,
  useToast,
  VStack,
  Wrap,
  WrapItem,
} from "@chakra-ui/react";
import axios from "axios";
import { useEffect } from "react";
import { FormProvider, useForm } from "react-hook-form";
import { Link, useMatch } from "react-location";
import { GroupSelect } from "../../../components/forms/access-rule/components/Select";

import { AdminLayout } from "../../../components/Layout";
import { useGetGroup } from "../../../utils/backend-client/admin/admin";
import { postApiV1AdminUsersUserId } from "../../../utils/backend-client/default/default";
import { useGetUser } from "../../../utils/backend-client/end-user/end-user";
import {
  PostApiV1AdminUsersUserIdBody,
  User,
} from "../../../utils/backend-client/types";
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
          <Groups user={user} />
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

interface GroupsProps {
  user: User;
}
const Groups: React.FC<GroupsProps> = ({ user }) => {
  const methods = useForm<PostApiV1AdminUsersUserIdBody>({});
  const toast = useToast();
  const { onOpen, onClose, isOpen } = useDisclosure();
  useEffect(() => {
    if (!isOpen) {
      methods.reset({
        groups: user.groups,
      });
    }
  }, [isOpen]);

  const onSubmit = async (data: PostApiV1AdminUsersUserIdBody) => {
    console.log({ data });
    try {
      await postApiV1AdminUsersUserId(user.id, data);
      toast({
        title: "Updated Groups",
        status: "success",
        variant: "subtle",
        duration: 2200,
        isClosable: true,
      });
      onClose();
    } catch (err) {
      let description: string | undefined;
      if (axios.isAxiosError(err)) {
        // @ts-ignore
        description = err?.response?.data.error;
      }

      toast({
        title: "Error Updating Groups",
        description,
        status: "error",
        variant: "subtle",
        duration: 2200,
        isClosable: true,
      });
    }
  };
  // make it an editable with a form component when editing
  //
  //
  if (isOpen) {
    return (
      <FormProvider {...methods}>
        <VStack
          as="form"
          onSubmit={methods.handleSubmit(onSubmit)}
          align={"left"}
          spacing={1}
        >
          <FormControl id="groups">
            <FormLabel>
              <HStack>
                <Text textStyle="Body/Medium">Groups</Text>
                <IconButton
                  size="sm"
                  variant="ghost"
                  icon={<CheckIcon />}
                  aria-label={"save groups"}
                  type="submit"
                />
                <IconButton
                  size="sm"
                  variant="ghost"
                  icon={<CloseIcon />}
                  aria-label={"cancel edit groups"}
                  onClick={onClose}
                />
              </HStack>
            </FormLabel>
            <GroupSelect fieldName="groups" />
          </FormControl>
        </VStack>
      </FormProvider>
    );
  }
  return (
    <VStack align={"left"} spacing={1}>
      <HStack>
        <Text textStyle="Body/Medium">Groups</Text>
        <IconButton
          size="sm"
          variant="ghost"
          icon={<EditIcon />}
          aria-label={"edit groups"}
          onClick={onOpen}
        />
      </HStack>
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
  );
};
