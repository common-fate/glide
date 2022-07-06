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
} from "@chakra-ui/react";
import { useState } from "react";

type Props = Omit<ModalProps, "children">;

const CreateUserModal = (props: Props) => {
  const [name, setName] = useState<string>("");
  const [email, setEmail] = useState<string>("");
  const [isAdmin, setIsAdmin] = useState(false);
  const [loading, setLoading] = useState(false);

  const createUserHandler = () => {
    setLoading(true);
  };

  return (
    <Modal {...props}>
      <ModalOverlay />
      <ModalContent>
        <ModalCloseButton />
        <ModalHeader mt={10}>Create User</ModalHeader>
        <ModalBody>
          <Stack
            spacing="5"
            //   divider={<StackDivider />}
          >
            <FormControl id="name">
              <Stack>
                <FormLabel mb={-1}>Name</FormLabel>
                <Input
                  variant="outline"
                  bg="white"
                  maxW={{ md: "3xl" }}
                  placeholder="name"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                />
              </Stack>
            </FormControl>
            <FormControl id="email">
              <Stack>
                <FormLabel mb={-1}>Email</FormLabel>
                <Input
                  variant="outline"
                  bg="white"
                  maxW={{ md: "3xl" }}
                  placeholder="Email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                />
              </Stack>
            </FormControl>
            <FormControl id="email">
              <Stack>
                <FormLabel mb={-1}>Is Admin?</FormLabel>
                <Switch
                  isChecked={isAdmin}
                  onChange={(e) => setIsAdmin(e.target.checked)}
                />
              </Stack>
            </FormControl>
          </Stack>
        </ModalBody>
        <ModalFooter minH={12}>
          <Button mr={3} isLoading={loading} onClick={createUserHandler}>
            Create user
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
};

export default CreateUserModal;
