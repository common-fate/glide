import {
  ArrowBackIcon,
  CheckIcon,
  CloseIcon,
  EditIcon,
} from "@chakra-ui/icons";
import {
  Avatar,
  Button,
  Center,
  Container,
  Flex,
  FormControl,
  FormLabel,
  HStack,
  IconButton,
  Input,
  InputGroup,
  InputRightElement,
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

import { useEffect, useState } from "react";
import { FormProvider, useForm, UseFormReturn } from "react-hook-form";
import { Link, useMatch } from "react-location";
import { useGetUser } from "../../../utils/backend-client/end-user/end-user";
import { UserSelect } from "../../../components/forms/access-rule/components/Select";

import { AdminLayout } from "../../../components/Layout";
import {
  createGroup,
  useGetGroup,
} from "../../../utils/backend-client/admin/admin";

import {
  CreateGroupRequestBody,
  Group,
} from "../../../utils/backend-client/types";
import {
  GrantedKeysIcon,
  AzureIcon,
  OktaIcon,
  AWSIcon,
} from "../../../components/icons/Icons";
import { CognitoLogo, GoogleLogo } from "../../../components/icons/Logos";
import { GetIDPLogo } from "../../../utils/idp-logo";

const Index = () => {
  const methods = useForm<CreateGroupRequestBody>({});
  const [loading, setLoading] = useState(false);

  const {
    params: { id: groupId },
  } = useMatch();
  const { data: group } = useGetGroup(groupId);
  const toast = useToast();

  useEffect(() => {
    if (group) {
      const formValues: CreateGroupRequestBody = {
        id: group.id,
        name: group?.name ? group.name : "",
        description: group?.description,
        members: group?.members ? group.members : [],
      };
      methods.reset(formValues);
    }
  }, [group]);

  const [isEditable, setIsEditable] = useState(false);

  const handleSubmit = async (data: CreateGroupRequestBody) => {
    setLoading(true);

    await createGroup(data)
      .then(() => {
        toast({
          title: "Updated Group",
          status: "success",
          variant: "subtle",
          duration: 2200,
          isClosable: true,
        });
        setIsEditable(false);
        setLoading(false);
      })
      .catch(() => {
        setLoading(false);

        toast({
          title: "Error updating group",
          status: "error",
          variant: "subtle",
          duration: 2200,
          isClosable: true,
        });
      });
  };

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

    if (!isEditable) {
      return (
        <>
          <VStack align={"left"} spacing={1} flex={1} mr={4}>
            <Text textStyle="Body/Medium">Name</Text>
            <Text textStyle="Body/Small">{group.name}</Text>
            <Text textStyle="Body/Medium">Description</Text>
            <Text textStyle="Body/Small">{group.description}</Text>
            <Members group={group} isEditing={isEditable} methods={methods} />
          </VStack>
          {group.source == "internal" && (
            <Button
              variant="brandSecondary"
              size="sm"
              onClick={() => {
                setIsEditable(true);
              }}
            >
              Edit
            </Button>
          )}

          {GetIDPLogo({ idpType: group.source, size: 200 })}
        </>
      );
    }

    return (
      <VStack w="100%">
        <Flex w="100%">
          <FormProvider {...methods}>
            <VStack
              align={"left"}
              spacing={5}
              flex={1}
              as="form"
              onSubmit={methods.handleSubmit(handleSubmit)}
            >
              <VStack spacing={5} flex={1} align={"left"} w="100%">
                <FormControl>
                  <VStack align={"left"}>
                    <FormLabel display="inline">
                      <Text textStyle="Body/Medium">Name</Text>
                    </FormLabel>
                    <Input
                      textStyle="Body/Medium"
                      readOnly={!isEditable}
                      {...methods.register("name", {
                        required: true,
                        minLength: 1,
                      })}
                    />
                  </VStack>
                </FormControl>
                <FormControl>
                  <VStack align={"left"}>
                    <FormLabel display="inline">
                      <Text textStyle="Body/Medium">Description</Text>
                    </FormLabel>
                    <Input
                      w="100%"
                      textStyle="Body/Medium"
                      readOnly={!isEditable}
                      {...methods.register("description", {
                        required: true,
                        minLength: 1,
                      })}
                    />
                  </VStack>
                </FormControl>
              </VStack>

              <Members group={group} isEditing={isEditable} methods={methods} />

              {isEditable && (
                <HStack>
                  <Button w="20%" mr={3} type="submit" isLoading={loading}>
                    Save
                  </Button>
                  <Button
                    variant="brandSecondary"
                    w="20%"
                    mr={3}
                    onClick={() => {
                      setIsEditable(false);
                      setLoading(false);
                    }}
                    isLoading={loading}
                  >
                    Cancel
                  </Button>
                </HStack>
              )}
            </VStack>
          </FormProvider>
        </Flex>
      </VStack>
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
          to={"/admin/groups"}
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

  isEditing: boolean;
  methods: UseFormReturn<CreateGroupRequestBody>;
}

const Members: React.FC<MemberProps> = ({ isEditing, methods, group }) => {
  if (isEditing) {
    return (
      <FormProvider {...methods}>
        <VStack as="form" align={"left"} spacing={1}>
          <FormControl id="members">
            <FormLabel>
              <HStack>
                <Text textStyle="Body/Medium">Members</Text>
              </HStack>
            </FormLabel>
            <Flex flex={1}>
              <UserSelect
                fieldName="members"
                isDisabled={methods.formState.isSubmitting || !isEditing}
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
      </HStack>
      <Wrap>
        {group.members.map((g) => {
          return (
            <WrapItem key={g}>
              <UserDisplay userId={g} />
            </WrapItem>
          );
        })}
      </Wrap>
    </VStack>
  );
};

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
