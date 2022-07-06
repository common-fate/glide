import {
  Avatar,
  Box,
  Button,
  ButtonGroup,
  Flex,
  FormControl,
  FormLabel,
  Input,
  Link,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  ModalProps,
  Tab,
  TabList,
  TabPanel,
  TabPanels,
  Tabs,
  Text,
} from "@chakra-ui/react";
import React from "react";
import { Group, User } from "../../utils/backend-client/types";
import { userName } from "../../utils/userName";

type Props = { group?: Group; members: User[] } & Omit<ModalProps, "children">;

const GroupModal = ({ group, members, ...props }: Props) => {
  return (
    <Modal {...props}>
      <ModalOverlay />
      <ModalContent>
        <ModalCloseButton />
        <ModalHeader mt={10}>Group Settings</ModalHeader>
        <ModalBody>
          <Tabs>
            <Tabs variant="brand">
              <TabList>
                <Tab>General</Tab>
                <Tab>Members</Tab>
              </TabList>

              <TabPanels minH="30vh" px={0}>
                <TabPanel px={2} mt={5}>
                  {/* ok */}
                  <FormControl>
                    <FormLabel>Group Name</FormLabel>
                    <Input value={group?.name} />
                  </FormControl>
                </TabPanel>
                <TabPanel px={0} maxH="60vh" overflow="scroll">
                  {members.map((user) => (
                    <Flex
                      py={2}
                      px={4}
                      key={user.id}
                      rounded="md"
                      _hover={{
                        "bg": "neutrals.100",
                        ".hide_show_btn": {
                          display: "block",
                        },
                      }}
                    >
                      <Avatar
                        size="sm"
                        name={userName(user)}
                        src={user.picture}
                        mr={2}
                      />
                      <Flex>
                        <Box>
                          <Text color="neutrals.900">{userName(user)}</Text>
                          <Text color="neutrals.500">{user.email}</Text>
                        </Box>
                      </Flex>
                      <Flex ml="auto" pos="relative">
                        <Link
                          as="button"
                          pos="absolute"
                          right={4}
                          top={2}
                          display="none"
                          className="hide_show_btn"
                          variant="link"
                          textDecor="none"
                          _hover={{ textDecor: "underline" }}
                        >
                          Remove
                        </Link>
                      </Flex>
                    </Flex>
                  ))}
                </TabPanel>
              </TabPanels>
            </Tabs>
          </Tabs>
        </ModalBody>
        <ModalFooter minH={12}>
          <ButtonGroup rounded="full" spacing={2} ml="auto">
            <Button variant="outline" rounded="full" onClick={props.onClose}>
              Cancel
            </Button>
          </ButtonGroup>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
};

export default GroupModal;
