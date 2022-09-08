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
  Textarea,
  useToast,
} from "@chakra-ui/react";
import axios from "axios";
import { useEffect } from "react";
import { FormProvider, useForm } from "react-hook-form";
import { postApiV1AdminGroups } from "../../utils/backend-client/default/default";
import { CreateGroupRequestBody } from "../../utils/backend-client/types";
type Props = Omit<ModalProps, "children">;

const CreateGroupModal = (props: Props) => {
  const methods = useForm<CreateGroupRequestBody>({});
  const toast = useToast();
  useEffect(() => {
    if (!props.isOpen) {
      methods.reset();
    }
  }, [props.isOpen]);

  const onSubmit = async (data: CreateGroupRequestBody) => {
    console.log({ data });
    try {
      await postApiV1AdminGroups(data);
      toast({
        title: "Group Created",
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
        title: "Error Creating Group",
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
            <Text textStyle="Heading/H3">Create Group</Text>
          </ModalHeader>

          <ModalBody>
            <Stack
              spacing="5"
              //   divider={<StackDivider />}
            >
              <FormControl id="name">
                <Stack>
                  <FormLabel
                    textStyle="Body/Medium"
                    fontWeight="normal"
                    mb={-1}
                  >
                    Name
                  </FormLabel>
                  <Input
                    variant="outline"
                    bg="white"
                    maxW={{ md: "3xl" }}
                    placeholder="Developers"
                    {...methods.register("name", {
                      required: true,
                      minLength: 1,
                    })}
                  />
                </Stack>
              </FormControl>
              <FormControl id="description">
                <Stack>
                  <FormLabel
                    textStyle="Body/Medium"
                    fontWeight="normal"
                    mb={-1}
                  >
                    Description
                  </FormLabel>
                  <Textarea
                    variant="outline"
                    bg="white"
                    maxW={{ md: "3xl" }}
                    placeholder="Developers group"
                    {...methods.register("description")}
                  />
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
              Create group
            </Button>
          </ModalFooter>
        </ModalContent>
      </FormProvider>
    </Modal>
  );
};

export default CreateGroupModal;
