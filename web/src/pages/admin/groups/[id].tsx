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
import axios from "axios";
import { groupCollapsed } from "console";
import { data } from "msw/lib/types/context";
import { type } from "os";
import { useEffect, useState } from "react";
import { FormProvider, useForm } from "react-hook-form";
import { Link, useMatch } from "react-location";
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

const Index = () => {
  const methods = useForm<CreateGroupRequestBody>({});

  const {
    params: { id: groupId },
  } = useMatch();
  const { data: group, mutate } = useGetGroup(groupId);

  const [isEditable, setIsEditable] = useState(false);

  const handleSubmit = async (data: CreateGroupRequestBody) => {
    await createGroup(data)
      .then(() => {
        toast({
          title: "Updated Groups",
          status: "success",
          variant: "subtle",
          duration: 2200,
          isClosable: true,
        });
      })
      .catch(() => {
        toast({
          title: "Error Updating Group",
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

    return (
      <VStack w="100%">
        <Flex w="100%">
          <FormProvider {...methods}>
            <VStack
              align={"left"}
              spacing={1}
              flex={1}
              mr={4}
              onSubmit={methods.handleSubmit(handleSubmit)}
            >
              <InputGroup size="lg" bg="white" maxW="700px">
                <VStack align={"left"} w="100%">
                  <VStack align={"left"}>
                    <Text textStyle="Body/Medium">Name</Text>
                    <Input
                      textStyle="Body/Medium"
                      value={group.name}
                      readOnly={!isEditable}
                    />
                  </VStack>

                  <VStack align={"left"}>
                    <Text textStyle="Body/Medium">Description</Text>
                    <Input
                      w="100%"
                      textStyle="Body/Medium"
                      value={group.description}
                      readOnly={!isEditable}
                    />
                  </VStack>
                </VStack>
              </InputGroup>

              <Members
                group={group}
                onSubmit={(u) => mutate(u)}
                isEditing={isEditable}
              />
            </VStack>
            {group.source == "INTERNAL" && (
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
          </FormProvider>
        </Flex>
        {isEditable && <Button>Save</Button>}
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
  onSubmit?: (u: Group) => void;
  isEditing: boolean;
}

type GroupForm = {
  members: string[];
};

const Members: React.FC<MemberProps> = ({ group, isEditing }) => {
  const methods = useForm<GroupForm>({});

  useEffect(() => {
    methods.reset({ members: group.members });
  }, []);

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
};
function toast(arg0: {
  title: string;
  status: string;
  variant: string;
  duration: number;
  isClosable: boolean;
}) {
  throw new Error("Function not implemented.");
}
