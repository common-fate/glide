import {
  Button,
  Container,
  Heading,
  Input,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  Text,
  useDisclosure,
} from "@chakra-ui/react";
import { useState } from "react";
import { Helmet } from "react-helmet";
import { useMatch } from "react-location";
import { UserLayout } from "../../components/Layout";
import {
  deleteProvider,
  updateProvider,
  useGetProvider,
} from "../../utils/local-client/deploymentcli/deploymentcli";

const Provider = () => {
  const {
    params: { id },
  } = useMatch();

  const provider = useGetProvider(id);
  const { isOpen, onOpen, onClose } = useDisclosure();
  const [loading, setLoading] = useState(false);
  const [ver, setVer] = useState<string>("");
  const handleDelete = () => {
    setLoading(true);
    deleteProvider(id).finally(() => {
      setLoading(false);
    });
  };

  return (
    <UserLayout>
      <Helmet>
        <title>Provider</title>
      </Helmet>
      <Container
        my={12}
        // This prevents unbounded widths for small screen widths
        minW={{ base: "100%", lg: "container.lg" }}
        overflowX="auto"
      >
        <Button onClick={handleDelete} isLoading={loading}>
          Delete
        </Button>
        <Button onClick={onOpen} isLoading={loading}>
          Update
        </Button>
        <Modal isOpen={isOpen} onClose={onClose}>
          <ModalOverlay />
          <ModalContent>
            <ModalHeader>Version</ModalHeader>
            <ModalCloseButton />
            <ModalBody>
              <Input onChange={(e) => setVer(e.target.value)} />
            </ModalBody>

            <ModalFooter>
              <Button colorScheme="blue" mr={3} onClick={onClose}>
                Close
              </Button>
              <Button
                variant="ghost"
                onClick={() =>
                  updateProvider(provider.data?.id ?? "", {
                    alias: "example",
                    version: ver,
                  })
                }
              >
                Update
              </Button>
            </ModalFooter>
          </ModalContent>
        </Modal>
        <Heading>{provider.data?.name}</Heading>
        <Heading>{provider.data?.team}</Heading>
        <Heading>{provider.data?.status}</Heading>
        <Heading>{provider.data?.version}</Heading>
        <Heading>{provider.data?.stackId}</Heading>
        <Text
          as={"pre"}
          textStyle="Body/SmallBold"
          color="neutrals.700"
          whiteSpace={"pre-wrap"}
        >
          {JSON.stringify(provider, undefined, 2)}
        </Text>
      </Container>
    </UserLayout>
  );
};
export default Provider;
