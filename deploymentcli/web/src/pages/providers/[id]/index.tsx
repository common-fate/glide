import {
  Box,
  Button,
  Container,
  Flex,
  Heading,
  Input,
  LinkOverlay,
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
import React, { useState } from "react";
import { Helmet } from "react-helmet";
import { Link, useMatch } from "react-location";
import { ProviderIcon } from "../../../components/icons/providerIcon";
import { UserLayout } from "../../../components/Layout";
import {
  deleteProvider,
  updateProvider,
  useGetProvider,
} from "../../../utils/local-client/deploymentcli/deploymentcli";

const Provider = () => {
  const {
    params: { id },
  } = useMatch();

  const { data: provider } = useGetProvider(id);
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
        minW={{ base: "100%", md: "container.md" }}
        overflowX="auto"
      >
        <Box
          key={id}
          as="button"
          className="group"
          textAlign="center"
          bg="neutrals.100"
          p={6}
          rounded="md"
          data-testid={"provider_" + id}
          position="relative"
          _disabled={{
            opacity: "0.5",
          }}
          w="100%"
        >
          <LinkOverlay
            href={`/registry/${id}`}
            as={Link}
            to={`/registry/${id}`}
          >
            <Flex flexDir="row" alignItems="center" my={6}>
              <ProviderIcon type={provider.name} mr={3} h="8" w="8" />
              <Text textStyle="Body/SmallBold" color="neutrals.700">
                {`${provider.team}/${provider.name}@${provider.version}`}
              </Text>
            </Flex>
            <Flex>
              <Button onClick={handleDelete} isLoading={loading}>
                Delete
              </Button>
              <Button onClick={onOpen} isLoading={loading}>
                Update
              </Button>
              <Button as={Link} to={`/providers/${id}/setup`}>
                Setup Configuration
              </Button>
            </Flex>
          </LinkOverlay>
        </Box>
        {/*  */}

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
                  updateProvider(provider?.id ?? "", {
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
        {/* <Heading>{provider.data?.name}</Heading>
        <Heading>{provider.data?.team}</Heading>
        <Heading>{provider.data?.status}</Heading>
        <Heading>{provider.data?.version}</Heading>
        <Heading>{provider.data?.stackId}</Heading> */}
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
