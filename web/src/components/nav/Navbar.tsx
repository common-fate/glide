import {
  ChevronDownIcon,
  ExternalLinkIcon,
  HamburgerIcon,
  SearchIcon,
} from "@chakra-ui/icons";
import {
  Avatar,
  Box,
  Button,
  ButtonGroup,
  ButtonProps,
  Container,
  Flex,
  HStack,
  IconButton,
  Kbd,
  Link as ChakraLink,
  Menu,
  MenuButton,
  MenuItem,
  MenuList,
  useBreakpointValue,
  useColorModeValue,
  useDisclosure,
  useEventListener,
  VStack,
  textDecoration,
} from "@chakra-ui/react";

import * as React from "react";
import { useMemo } from "react";
import { Link, useNavigate } from "react-location";
import { useUserListRequests } from "../../utils/backend-client/end-user/end-user";
import { useCognito } from "../../utils/context/cognitoContext";
import { useUser } from "../../utils/context/userContext";
import CommandPalette from "../CommandPalette";
import Counter from "../Counter";
import { DoorIcon } from "../icons/Icons";
import { CommonFateLogo } from "../icons/Logos";
import { DrawerNav } from "./DrawerNav";

export const Navbar: React.FC = () => {
  const isDesktop = useBreakpointValue({ base: false, lg: true }, "800px");

  const user = useUser();
  const auth = useCognito();
  const { data: requests } = useUserListRequests({ status: "PENDING" });
  const { data: reviews } = useUserListRequests({ reviewer: true });

  const { isOpen, onOpen, onClose } = useDisclosure();
  const modal = useDisclosure();

  // https://erikmartinjordan.com/navigator-platform-deprecated-alternative
  const isMac = () =>
    /(Mac|iPhone|iPod|iPad)/i.test(
      // @ts-ignore
      navigator?.userAgentData?.platform || navigator?.platform || "unknown"
    );

  const ACTION_KEY_DEFAULT = ["Ctrl", "Control"];
  const ACTION_KEY_APPLE = ["⌘", "Command"];
  const [actionKey, setActionKey] = React.useState<string[]>(ACTION_KEY_APPLE);

  React.useEffect(() => {
    if (typeof navigator === "undefined") return;
    if (!isMac()) {
      setActionKey(ACTION_KEY_DEFAULT);
    }
  }, []);

  useEventListener("keydown", (event) => {
    const hotkey = isMac() ? "metaKey" : "ctrlKey";
    if (event?.key?.toLowerCase() === "k" && event[hotkey]) {
      event.preventDefault();
      modal.isOpen ? modal.onClose() : modal.onOpen();
    }
  });

  const showReqCount = useMemo(
    () => requests?.requests && requests?.requests.length > 0,
    [requests]
  );
  const showRewCount = useMemo(
    () =>
      reviews?.requests &&
      reviews.requests.filter((r) => r.status === "PENDING").length > 0,
    [reviews]
  );
  const navigate = useNavigate();

  return (
    <Box as="section">
      <Box
        as="nav"
        bg={"white"}
        color={"gray.900"}
        boxShadow={useColorModeValue("sm", "sm-dark")}
      >
        <Container maxW="container.xl">
          <Flex justify="space-between" py={{ base: "3", lg: "4" }}>
            <HStack spacing="4" pos="relative" w="100%">
              <VStack>
                <ChakraLink
                  href={"https://docs.commonfate.io/common-fate/introduction/"}
                  isExternal
                  transition="all .2s ease"
                  rounded="sm"
                >
                  <CommonFateLogo h="32px" w="auto" />
                </ChakraLink>
                {/* <ChakraLink
                  href={"https://docs.commonfate.io/common-fate/introduction/"}
                  isExternal
                  transition="all .2s ease"
                  rounded="sm"
                >
                  <Button
                    size="xs"
                    variant="link"
                    px={3}
                    aria-label="docs"
                    id="docs-button"
                    color="neutrals.600"
                  >
                    Documentation
                  </Button>
                </ChakraLink> */}
              </VStack>
              <Button
                size="md"
                variant="unstyled"
                bg="neutrals.100"
                rounded="full"
                border="1px solid"
                borderColor="neutrals.200"
                aria-label="Search"
                // w="113px"
                // px={2}
                textAlign="left"
                color="neutrals.600"
                onClick={modal.onOpen}
              >
                <SearchIcon color="neutrals.600" boxSize="15px" ml={3} mr={2} />
                Search
                <Flex ml={2} mr={3} display="inline-flex">
                  <Kbd>{actionKey[0]}</Kbd>
                  <Kbd>k</Kbd>
                </Flex>
              </Button>
              {isDesktop && (
                <ButtonGroup
                  variant="ghost"
                  spacing="0"
                  mb={"-32px !important;"}
                  ml="auto !important;"
                  mr="auto !important;"
                  // pos="absolute"
                  // bottom={-4}
                  // left="50%"
                >
                  {/* I've hardcoded widths here to prevent the bold/unbold text from 
                  altering the containing divs width. Reduces *jittering* */}
                  <TabsStyledButton
                    href="/requests"
                    w={showReqCount ? "190px" : "142px"}
                    pr={showReqCount ? 10 : undefined}
                  >
                    Access Requests
                    {showReqCount && (
                      <Counter
                        pos="absolute"
                        right={2}
                        count={requests?.requests?.length ?? 0}
                      />
                    )}
                  </TabsStyledButton>
                  <TabsStyledButton
                    href="/reviews?status=pending"
                    w="125px"
                    pr={showReqCount ? 10 : undefined}
                  >
                    Reviews
                    {showRewCount && (
                      <Counter
                        pos="absolute"
                        right={2}
                        count={
                          reviews?.requests?.filter(
                            (rw) => rw.status == "PENDING"
                          ).length ?? 0
                        }
                      />
                    )}
                  </TabsStyledButton>
                </ButtonGroup>
              )}
            </HStack>
            {isDesktop ? (
              <HStack spacing="4">
                {user.isAdmin && (
                  <ButtonGroup variant="ghost" spacing="1">
                    <Button
                      as={Link}
                      to={"/admin/access-rules"}
                      size="md"
                      variant="link"
                      px={3}
                      aria-label="Admin"
                      id="admin-button"
                      color="neutrals.600"
                    >
                      Switch To Admin
                    </Button>
                  </ButtonGroup>
                )}

                <Menu>
                  <MenuButton
                    data-testid="logout-icon"
                    as={Button}
                    variant="ghost"
                    rounded="full"
                    rightIcon={
                      <ChevronDownIcon
                        color="gray.500"
                        ml={-1}
                        boxSize="24px"
                      />
                    }
                    py={2}
                    pl={0}
                  >
                    <Avatar
                      variant="withBorder"
                      boxSize="10"
                      name={user.user?.email}
                    />
                  </MenuButton>
                  <MenuList _dark={{ borderColor: "gray.500" }}>
                    <MenuItem
                      _hover={{ cursor: "auto", backgroundColor: "white" }}
                      onClick={(e) => {
                        e.preventDefault();
                      }}
                      borderBottom="2px"
                      borderColor="neutrals.300"
                    >
                      {user.user?.email}
                    </MenuItem>
                    <MenuItem icon={<ExternalLinkIcon boxSize="15px" />}>
                      <ChakraLink
                        isExternal
                        href="https://docs.commonfate.io/common-fate/introduction/"
                        textDecoration="none"
                        _hover={{ textDecoration: "none" }}
                      >
                        Documentation
                      </ChakraLink>
                    </MenuItem>

                    <MenuItem
                      data-testid="logout-button"
                      icon={<DoorIcon color={"gray.700"} />}
                      onClick={async () =>
                        await auth.initiateSignOut().then((e) => console.log(e))
                      }
                    >
                      Sign out
                    </MenuItem>
                  </MenuList>
                </Menu>
              </HStack>
            ) : (
              <IconButton
                variant="ghost"
                icon={<HamburgerIcon fontSize="1.25rem" />}
                aria-label="Open Menu"
                onClick={onOpen}
              />
            )}
          </Flex>
        </Container>
      </Box>
      <DrawerNav isOpen={isOpen} onClose={onClose} />
      <CommandPalette isOpen={modal.isOpen} onClose={modal.onClose} />
    </Box>
  );
};
interface TabsStyledButtonProps extends ButtonProps {
  href: string;
}
export const TabsStyledButton: React.FC<TabsStyledButtonProps> = ({
  href,
  ...rest
}) => {
  const navigate = useNavigate();
  return (
    <Button
      opacity={0.8}
      roundedTop="md"
      onClick={() => {
        navigate({ to: href });
      }}
      isActive={location.pathname === href.split("?")[0]}
      _active={{
        fontWeight: "bold",
        opacity: 1,
        borderColor: "#2E7FFF",
        borderBottomWidth: "2px",
      }}
      sx={{
        rounded: "none",
        // paddingBottom: "10px",
        borderBottom: "2px solid",
        borderColor: "neutrals.300",
        color: "neutrals.700",
        px: 4,
        // hover state
        _hover: {
          borderColor: "neutrals.500",
        },
        // 'Current' state
        _selected: {
          fontWeight: 500,
          borderColor: "#2E7FFF",
        },
        // Disabled state
        _disabled: {
          opacity: 0.3,
        },
      }}
      {...rest}
    />
  );
};

export const StyledButton = (props: ButtonProps) => (
  <Button w="100%" justifyContent="flex-start" variant="ghost" {...props} />
);
