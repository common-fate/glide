import {
  Box,
  Button,
  Center,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  ModalProps,
  Skeleton,
  Stack,
  Text,
} from "@chakra-ui/react";
import { Link } from "react-location";
import { useListUserAccessRules } from "../../utils/backend-client/end-user/end-user";
import { ProviderIcon } from "../icons/providerIcon";

type Props = Omit<ModalProps, "children">;

const AcessRulesMobileModal = (props: Props) => {
  const { data: rules, isValidating } = useListUserAccessRules();

  return (
    <Modal size="full" scrollBehavior="inside" {...props}>
      <ModalOverlay />
      <ModalContent top="20vh">
        <ModalCloseButton />
        <ModalHeader mt={10}>All Access</ModalHeader>
        <ModalBody>
          <Stack gap={6}>
            {rules ? (
              rules.accessRules.length > 0 ? (
                rules.accessRules.map((r) => (
                  <Box
                    key={r.id}
                    className="group"
                    textAlign="center"
                    bg="neutrals.100"
                    p={6}
                    as="a"
                    h="172px"
                    w="100%"
                    rounded="md"
                  >
                    <ProviderIcon
                      provider={r.target.provider}
                      mb={3}
                      h="8"
                      w="8"
                    />

                    <Text textStyle="Body/SmallBold" color="neutrals.700">
                      {r.name}
                    </Text>
                    <Link to={"/access/request/" + r.id} key={r.id}>
                      <Button mt={4} variant="brandSecondary" size="sm">
                        Request
                      </Button>
                    </Link>
                  </Box>
                ))
              ) : (
                <Center
                  bg="neutrals.100"
                  p={6}
                  as="a"
                  h="193px"
                  w="100%"
                  rounded="md"
                  flexDir="column"
                  textAlign="center"
                >
                  <Text textStyle="Heading/H3" color="neutrals.500">
                    No Access
                  </Text>
                  <Text textStyle="Body/Medium" color="neutrals.400" mt={2}>
                    You don’t have access to anything yet. Ask your Granted
                    administrator to finish setting up Granted.
                  </Text>
                </Center>
              )
            ) : (
              // Otherwise loading state
              [1, 2, 3, 4].map((i) => (
                <Skeleton
                  key={i}
                  p={6}
                  h="172px"
                  w={{ base: "100%", md: "232px" }}
                  rounded="sm"
                />
              ))
            )}
          </Stack>
        </ModalBody>
        <ModalFooter minH={12}></ModalFooter>
      </ModalContent>
    </Modal>
  );
};

export default AcessRulesMobileModal;
