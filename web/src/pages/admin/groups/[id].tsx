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
import { type } from "os";
import { useEffect } from "react";
import { FormProvider, useForm } from "react-hook-form";
import { Link, useMatch } from "react-location";
import {
  GroupSelect,
  UserSelect,
} from "../../../components/forms/access-rule/components/Select";

import { AdminLayout } from "../../../components/Layout";
import {
  updateUser,
  useGetGroup,
} from "../../../utils/backend-client/admin/admin";

import { useGetUser } from "../../../utils/backend-client/end-user/end-user";
import { Group, User } from "../../../utils/backend-client/types";

const UserDisplay: React.FC<{ userId: string }> = ({ userId }) => {
  const { data } = useGetUser(encodeURIComponent(userId));
  return (
    <Flex
      cursor="help"
      textStyle={"Body/Small"}
      rounded="full"
      bg="neutrals.300"
      py={1}
      px={4}
    >
      {data?.email}
    </Flex>
  );
};

const Index = () => {
  const {
    params: { id: groupId },
  } = useMatch();
  const { data: group, isValidating, error, mutate } = useGetGroup(groupId);

  const Content = () => {
    if (group?.id === undefined) {
      return (
        <>
          <VStack>
            <Text>Name</Text>
            <SkeletonText noOfLines={1} />
            <Text>Description</Text>
            <SkeletonText noOfLines={1} />
            <Text>Members</Text>
            <SkeletonText noOfLines={3} />
          </VStack>
        </>
      );
    }

    return (
      <>
        <VStack align={"left"} spacing={1} flex={1} mr={4}>
          <Text textStyle="Body/Medium">Name</Text>
          <Text textStyle="Body/Small">{group?.name}</Text>
          <Text textStyle="Body/Medium">Description</Text>
          <Text textStyle="Body/Small">{group?.description}</Text>
          <Members group={group} onSubmit={(u) => mutate(u)} />
        </VStack>
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
          Group Details
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

interface MemberProps {
  group: Group;
  onSubmit?: (u: Group) => void;
}

type GroupForm = {
  members: string[];
};

const Members: React.FC<MemberProps> = ({ group, onSubmit }) => {
  const methods = useForm<GroupForm>({});
  // const toast = useToast();
  const { onOpen, onClose, isOpen } = useDisclosure();
  useEffect(() => {
    if (!isOpen) {
      methods.reset({ members: group.members });
    }
  }, [isOpen]);

  // const handleSubmit = async (data: UpdateGroupBody) => {
  //   // try {
  //   //   const u = await updateUser(user.id, data);
  //   //   toast({
  //   //     title: "Updated Groups",
  //   //     status: "success",
  //   //     variant: "subtle",
  //   //     duration: 2200,
  //   //     isClosable: true,
  //   //   });
  //   //   onSubmit?.(u);
  //   //   onClose();
  //   // } catch (err) {
  //   //   let description: string | undefined;
  //   //   if (axios.isAxiosError(err)) {
  //   //     // @ts-ignore
  //   //     description = err?.response?.data.error;
  //   //   }
  //   //   toast({
  //   //     title: "Error Updating Groups",
  //   //     description,
  //   //     status: "error",
  //   //     variant: "subtle",
  //   //     duration: 2200,
  //   //     isClosable: true,
  //   //   });
  //   // }
  // };

  if (isOpen) {
    return (
      <FormProvider {...methods}>
        <VStack
          as="form"
          // onSubmit={methods.handleSubmit(handleSubmit)}
          align={"left"}
          spacing={1}
        >
          <FormControl id="members">
            <FormLabel>
              <HStack>
                <Text textStyle="Body/Medium">Members</Text>
                <IconButton
                  isLoading={methods.formState.isSubmitting}
                  size="sm"
                  variant="ghost"
                  icon={<CheckIcon />}
                  aria-label={"save members"}
                  type="submit"
                />
                <IconButton
                  isDisabled={methods.formState.isSubmitting}
                  size="sm"
                  variant="ghost"
                  icon={<CloseIcon />}
                  aria-label={"cancel edit members"}
                  onClick={onClose}
                />
              </HStack>
            </FormLabel>
            <Flex flex={1}>
              <UserSelect
                fieldName="members"
                isDisabled={methods.formState.isSubmitting}
              />
            </Flex>
          </FormControl>
        </VStack>
      </FormProvider>
    );
  }
  return (
    <VStack align={"left"} spacing={1}>
      <HStack>
        <Text textStyle="Body/Medium">Members</Text>
        <IconButton
          size="sm"
          variant="ghost"
          icon={<EditIcon />}
          aria-label={"edit members"}
          onClick={onOpen}
        />
      </HStack>
      <Wrap>
        {group.members.map((u) => {
          return (
            <WrapItem key={u}>
              <UserDisplay userId={u} />
            </WrapItem>
          );
        })}
      </Wrap>
    </VStack>
  );
};
