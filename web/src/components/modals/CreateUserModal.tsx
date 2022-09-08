import {
  Button,
  FormControl,
  FormLabel,
  Input,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  ModalProps,
  Stack,
  Switch,
  Text,
  useToast,
} from "@chakra-ui/react";
import axios from "axios";
import { useEffect } from "react";
import { FormProvider, useForm } from "react-hook-form";
import {
  postApiV1AdminUsers,
  useGetApiV1AdminIdentity,
} from "../../utils/backend-client/default/default";
import { CreateUserRequestBody } from "../../utils/backend-client/types";
type Props = Omit<ModalProps, "children">;

const CreateUserModal = (props: Props) => {
  const methods = useForm<CreateUserRequestBody>({});
  const toast = useToast();

  useEffect(() => {
    if (!props.isOpen) {
      methods.reset();
    }
  }, [props.isOpen]);

  const onSubmit = async (data: CreateUserRequestBody) => {
    console.log({ data });
    try {
      await postApiV1AdminUsers(data);
      toast({
        title: "Created User",
        status: "success",
        variant: "subtle",
        duration: 2200,
        isClosable: true,
      });
      props.onClose();
    } catch (err) {
      let description: string | undefined;
      if (axios.isAxiosError(err)) {
        // @ts-ignore
        description = err?.response?.data.error;
      }

      toast({
        title: "Error Creating User",
        description,
        status: "error",
        variant: "subtle",
        duration: 2200,
        isClosable: true,
      });
    }
  };

  return (
    <Modal {...props}>
      <ModalOverlay />
      <FormProvider {...methods}>
        <ModalContent as="form" onSubmit={methods.handleSubmit(onSubmit)}>
          <ModalCloseButton />
          <ModalHeader mt={10}>
            <Text textStyle="Heading/H3">Create User</Text>
          </ModalHeader>

          <ModalBody>
            <Stack
              spacing="5"
              //   divider={<StackDivider />}
            >
              <FormControl id="firstName">
                <Stack>
                  <FormLabel
                    textStyle="Body/Medium"
                    fontWeight="normal"
                    mb={-1}
                  >
                    First Name
                  </FormLabel>
                  <Input
                    variant="outline"
                    bg="white"
                    maxW={{ md: "3xl" }}
                    placeholder="Alice"
                    {...methods.register("firstName", {
                      required: true,
                      minLength: 1,
                    })}
                  />
                </Stack>
              </FormControl>
              <FormControl id="lastName">
                <Stack>
                  <FormLabel
                    textStyle="Body/Medium"
                    fontWeight="normal"
                    mb={-1}
                  >
                    Last Name
                  </FormLabel>
                  <Input
                    variant="outline"
                    bg="white"
                    maxW={{ md: "3xl" }}
                    placeholder="Alison"
                    {...methods.register("lastName", {
                      required: true,
                      minLength: 1,
                    })}
                  />
                </Stack>
              </FormControl>
              <FormControl id="email">
                <Stack>
                  <FormLabel
                    textStyle="Body/Medium"
                    fontWeight="normal"
                    mb={-1}
                  >
                    Email
                  </FormLabel>
                  <Input
                    variant="outline"
                    bg="white"
                    maxW={{ md: "3xl" }}
                    placeholder="Email"
                    {...methods.register("email", {
                      required: true,
                      minLength: 1,
                    })}
                  />
                </Stack>
              </FormControl>
              <FormControl id="isAdmin">
                <Stack>
                  <FormLabel
                    textStyle="Body/Medium"
                    fontWeight="normal"
                    mb={-1}
                  >
                    Admin
                  </FormLabel>
                  <Switch {...methods.register("isAdmin")} />
                </Stack>
              </FormControl>
            </Stack>
          </ModalBody>
          <ModalFooter minH={12}>
            <Button
              mr={3}
              isLoading={methods.formState.isSubmitting}
              type="submit"
            >
              Create user
            </Button>
          </ModalFooter>
        </ModalContent>
      </FormProvider>
    </Modal>
  );
};

export default CreateUserModal;
