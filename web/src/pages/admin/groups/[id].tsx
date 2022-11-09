import { ArrowBackIcon } from "@chakra-ui/icons";
import {
  Button,
  Center,
  Container,
  Flex,
  FormControl,
  FormErrorMessage,
  FormLabel,
  HStack,
  IconButton,
  Input,
  SkeletonText,
  Text,
  useToast,
  VStack,
  Wrap,
  WrapItem,
} from "@chakra-ui/react";

import { useEffect, useState } from "react";
import { FormProvider, useForm } from "react-hook-form";
import { Link, useMatch } from "react-location";
import { UserSelect } from "../../../components/forms/access-rule/components/Select";
import { useGetUser } from "../../../utils/backend-client/end-user/end-user";

import { AdminLayout } from "../../../components/Layout";

import {
  adminUpdateGroup,
  useGetGroup,
} from "../../../utils/backend-client/admin/admin";
import { CreateGroupRequestBody } from "../../../utils/backend-client/types";
import { GetIDPLogo } from "../../../utils/idp-logo";

const Index = () => {
  const methods = useForm<CreateGroupRequestBody>({});
  const [loading, setLoading] = useState(false);

  const {
    params: { id: groupId },
  } = useMatch();
  const { data: group, mutate } = useGetGroup(groupId);
  const toast = useToast();
  const [isEditable, setIsEditable] = useState(false);

  useEffect(() => {
    if (group) {
      const formValues: CreateGroupRequestBody = {
        name: group?.name ? group.name : "",
        description: group?.description,
        members: group?.members ? group.members : [],
      };
      methods.reset(formValues);
    }
  }, [group]);

  const handleSubmit = (data: CreateGroupRequestBody) => {
    setLoading(true);
    mutate(adminUpdateGroup(groupId, data))
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
        <HStack align={"flex-start"} w="100%">
          {GetIDPLogo({ idpType: group.source, size: 150 })}
          <VStack align={"left"} spacing={1} flex={1} mr={4}>
            <Text textStyle="Body/Medium">Name</Text>
            <Text textStyle="Body/Small">{group.name}</Text>
            <Text textStyle="Body/Medium">Description</Text>
            <Text textStyle="Body/Small">{group.description}</Text>
            <Text textStyle="Body/Medium">Members</Text>
            <Wrap>
              {group.members.length === 0 ? (
                <WrapItem>
                  <Text textStyle="Body/Small">No members</Text>
                </WrapItem>
              ) : (
                group.members.map((g) => {
                  return (
                    <WrapItem key={g}>
                      <UserDisplay userId={g} />
                    </WrapItem>
                  );
                })
              )}
            </Wrap>
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
        </HStack>
      );
    }

    return (
      <VStack
        spacing={6}
        align={"left"}
        w="100%"
        as="form"
        onSubmit={methods.handleSubmit(handleSubmit)}
      >
        <VStack>
          <FormProvider {...methods}>
            <FormControl isInvalid={!!methods.formState.errors.name}>
              <FormLabel fontWeight={"normal"}>Name</FormLabel>
              <Input
                background={"white"}
                {...methods.register("name", {
                  required: true,
                  minLength: 1,
                })}
                onBlur={() => {
                  void methods.trigger("name");
                }}
              />
              <FormErrorMessage>Name is required</FormErrorMessage>
            </FormControl>
            <FormControl isInvalid={!!methods.formState.errors.description}>
              <FormLabel fontWeight={"normal"}>Description</FormLabel>
              <Input
                background={"white"}
                {...methods.register("description", {
                  required: "Description is required",
                  minLength: 1,
                })}
                onBlur={() => {
                  void methods.trigger("description");
                }}
              />
              <FormErrorMessage>Description is required</FormErrorMessage>
            </FormControl>

            <FormControl id="members">
              <FormLabel>
                <HStack>
                  <Text textStyle="Body/Medium">Members</Text>
                </HStack>
              </FormLabel>
              <UserSelect
                fieldName="members"
                isDisabled={methods.formState.isSubmitting}
              />
            </FormControl>
          </FormProvider>
        </VStack>
        <HStack justify={"right"}>
          <Button type="submit" isLoading={loading}>
            Save
          </Button>
          <Button
            variant="brandSecondary"
            onClick={() => {
              setIsEditable(false);
              setLoading(false);
            }}
            isDisabled={loading}
          >
            Cancel
          </Button>
        </HStack>
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
