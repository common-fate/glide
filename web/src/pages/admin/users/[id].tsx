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
  Text,
  Tooltip,
  useDisclosure,
  useToast,
  VStack,
  Wrap,
  WrapItem,
} from "@chakra-ui/react";
import axios from "axios";
import { useEffect, useState } from "react";
import { FormProvider, useForm } from "react-hook-form";
import { Link, useMatch } from "react-location";
import { GroupSelect } from "../../../components/forms/access-rule/components/Select";

import { AdminLayout } from "../../../components/Layout";
import {
  getGroup,
  updateUser,
  useIdentityConfiguration,
} from "../../../utils/backend-client/admin/admin";

import { useGetUser } from "../../../utils/backend-client/end-user/end-user";
import {
  Group,
  ListGroupsSource,
  UpdateUserBody,
  User,
} from "../../../utils/backend-client/types";
import { GetIDPName } from "../../../utils/idp-logo";

const GroupDisplay: React.FC<{ group: Group }> = ({ group }) => {
  return (
    <Tooltip label={group.description}>
      <Flex
        textStyle={"Body/Small"}
        rounded="full"
        bg="neutrals.300"
        py={1}
        data-testid={group.name}
        px={4}
      >
        {group.name}
      </Flex>
    </Tooltip>
  );
};
const Index = () => {
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

const Content: React.FC = () => {
  const {
    params: { id: userId },
  } = useMatch();
  const { data: user, mutate } = useGetUser(userId);
  const [userGroups, setUserGroups] = useState<Group[]>();
  const toast = useToast();
  useEffect(() => {
    if (user) {
      const groups = Promise.all(
        user.groups.map((g) => getGroup(encodeURIComponent(g)))
      );
      groups
        .then((g) => {
          setUserGroups(g);
        })
        .catch((err) => {
          let description: string | undefined;
          if (axios.isAxiosError(err)) {
            // @ts-ignore
            description = err?.response?.data.error;
          }
          toast({
            title: "Failed to load users groups",
            description,
            status: "error",
            variant: "subtle",
            duration: 2200,
            isClosable: true,
          });
        });
    }
  }, [user]);
  if (user?.id === undefined || userGroups === undefined) {
    return (
      <>
        <VStack>
          <Text>Name</Text>
          <SkeletonText noOfLines={1} />
          <Text>Email</Text>
          <SkeletonText noOfLines={1} />
          <ExternalGroupsLabel />
          <SkeletonText noOfLines={3} />
          <InternalGroupsLabel />
          <SkeletonText noOfLines={3} />
        </VStack>
      </>
    );
  }

  return (
    <>
      <VStack align={"left"} spacing={1} flex={1} mr={4}>
        <Text textStyle="Body/Medium">Name</Text>
        <Text textStyle="Body/Small">{`${user.firstName} ${user.lastName}`}</Text>
        <Text textStyle="Body/Medium">Email</Text>
        <Text textStyle="Body/Small">{user.email}</Text>
        <ExternalGroups userGroups={userGroups} />
        <InternalGroups
          user={user}
          onSubmit={(u) => mutate(u)}
          userGroups={userGroups}
        />
      </VStack>

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
interface InternalGroupsProps {
  userGroups: Group[];
  user: User;
  onSubmit?: (u: User) => void;
}
const InternalGroups: React.FC<InternalGroupsProps> = ({
  user,
  onSubmit,
  userGroups,
}) => {
  const methods = useForm<UpdateUserBody>({});
  const toast = useToast();
  const { onOpen, onClose, isOpen } = useDisclosure();

  useEffect(() => {
    if (isOpen) {
      methods.reset({
        groups: userGroups
          .filter((g) => g.source === "internal")
          .map((g) => g.id),
      });
    }
  }, [isOpen]);

  const handleSubmit = async (data: UpdateUserBody) => {
    try {
      const u = await updateUser(user.id, data);
      toast({
        title: "Updated Groups",
        status: "success",
        variant: "subtle",
        duration: 2200,
        isClosable: true,
      });
      onSubmit?.(u);
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

  if (isOpen) {
    return (
      <FormProvider {...methods}>
        <VStack
          as="form"
          onSubmit={methods.handleSubmit(handleSubmit)}
          align={"left"}
          spacing={1}
        >
          <FormControl id="groups">
            <FormLabel>
              <HStack>
                <InternalGroupsLabel />
                <IconButton
                  isLoading={methods.formState.isSubmitting}
                  size="sm"
                  variant="ghost"
                  icon={<CheckIcon />}
                  data-testid="save-group-submit"
                  aria-label={"save groups"}
                  type="submit"
                />
                <IconButton
                  isDisabled={methods.formState.isSubmitting}
                  size="sm"
                  variant="ghost"
                  icon={<CloseIcon />}
                  aria-label={"cancel edit groups"}
                  onClick={onClose}
                />
              </HStack>
            </FormLabel>
            <Flex flex={1}>
              <GroupSelect
                fieldName="groups"
                isDisabled={methods.formState.isSubmitting}
                source={ListGroupsSource.INTERNAL}
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
        <InternalGroupsLabel />
        <IconButton
          size="sm"
          variant="ghost"
          icon={<EditIcon />}
          data-testid="edit-groups-icon"
          aria-label={"edit groups"}
          onClick={onOpen}
        />
      </HStack>
      <Wrap>
        {userGroups
          .filter((g) => g.source === "internal")
          .map((g) => {
            return (
              <WrapItem key={g.id}>
                <GroupDisplay group={g} />
              </WrapItem>
            );
          })}
      </Wrap>
    </VStack>
  );
};

interface ExternalGroupsProps {
  userGroups: Group[];
}
const ExternalGroups: React.FC<ExternalGroupsProps> = ({ userGroups }) => {
  return (
    <VStack align={"left"} spacing={1}>
      <ExternalGroupsLabel />
      <Wrap>
        {userGroups
          .filter((g) => g.source !== "internal")
          .map((g) => {
            return (
              <WrapItem key={g.id}>
                <GroupDisplay group={g} />
              </WrapItem>
            );
          })}
      </Wrap>
    </VStack>
  );
};

const InternalGroupsLabel = () => {
  return (
    <Tooltip label="Internal groups are managed by Granted Approvals, use them when you need more granular access control than you have defined by groups in your external identity provider.">
      <Text textStyle="Body/Medium">Groups</Text>
    </Tooltip>
  );
};
const ExternalGroupsLabel = () => {
  const { data } = useIdentityConfiguration();
  return (
    <Tooltip label="External groups are managed by your identity provider. You can use your identity providers management console to update group memberships. These groups are synced automatically every 5 minutes.">
      <Text textStyle="Body/Medium" textTransform={"capitalize"}>
        {GetIDPName(data?.identityProvider || "")} Groups
      </Text>
    </Tooltip>
  );
};
