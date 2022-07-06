import { ChevronDownIcon, HamburgerIcon, SearchIcon } from "@chakra-ui/icons";
import {
  Avatar,
  Box,
  Button,
  ButtonGroup,
  ButtonProps,
  ColorModeContext,
  Container,
  Divider,
  Flex,
  HStack,
  IconButton,
  Input,
  InputGroup,
  InputLeftElement,
  Link as ChakraLink,
  Menu,
  MenuButton,
  MenuItem,
  MenuList,
  useBreakpointValue,
  useColorModeValue,
  useDisclosure,
} from "@chakra-ui/react";
import * as React from "react";
import { Link } from "react-location";
import { useUser } from "../../utils/context/userContext";
import { DoorIcon } from "../icons/Icons";
import { ApprovalsLogoAdmin } from "../icons/Logos";
import { DrawerNav } from "./DrawerNav";

export const AdminNavbar: React.FC<{}> = () => {
  const isDesktop = useBreakpointValue({ base: false, lg: true }, "800px");

  // Keep track of whether app has mounted (hydration is done)
  const [hasMounted, setHasMounted] = React.useState(false);
  React.useEffect(() => {
    setHasMounted(true);
  }, []);

  const auth = useUser();

  const { isOpen, onOpen, onClose } = useDisclosure();

  return (
    <Box as="section">
      <ColorModeContext.Provider
        value={{
          colorMode: "dark",
          // noop
          toggleColorMode: () => {},
          setColorMode: () => {},
        }}
      >
        {/* <DarkMode /> */}
        <Box
          as="nav"
          bg="gray.800"
          color="white"
          boxShadow={useColorModeValue("sm", "sm-dark")}
        >
          <Container maxW="container.xl">
            <Flex justify="space-between" py={{ base: "3", lg: "4" }}>
              <HStack spacing="4">
                <ChakraLink
                  as={Link}
                  to="/admin/access-rules"
                  transition="all .2s ease"
                  rounded="sm"
                >
                  <ApprovalsLogoAdmin h="32px" w="auto" />
                </ChakraLink>
                {isDesktop && (
                  <ButtonGroup variant="ghost" spacing="1">
                    <StyledNextButton href="/admin/access-rules">
                      Access Rules
                    </StyledNextButton>
                    <StyledNextButton href="/admin/requests">
                      Requests
                    </StyledNextButton>
                    <StyledNextButton href="/admin/users">
                      Users
                    </StyledNextButton>
                    <StyledNextButton href="/admin/groups">
                      Groups
                    </StyledNextButton>
                    <StyledNextButton href="/admin/settings">
                      Settings
                    </StyledNextButton>
                    {/* <StyledNextButton href="/admin/auditlog">
                      Audit Log
                    </StyledNextButton> */}
                  </ButtonGroup>
                )}
              </HStack>
              {isDesktop ? (
                <HStack spacing="4">
                  <Button
                    as={Link}
                    to={"/requests"}
                    size="md"
                    variant="link"
                    px={3}
                    aria-label="Admin"
                  >
                    User
                  </Button>
                  <Menu>
                    <MenuButton
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
                        name={auth?.user?.email}
                      />
                    </MenuButton>
                    <MenuList _dark={{ borderColor: "gray.500" }}>
                      <MenuItem
                        icon={<DoorIcon color={"gray.400"} />}
                        onClick={auth.initiateSignOut}
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
          {isDesktop && false && (
            <>
              <Divider />
              <Container maxW="container.xl" py={3}>
                <Flex justify="space-between">
                  <ButtonGroup variant="ghost" spacing="1">
                    <Button>Overview</Button>
                    <StyledNextButton href="/admin/access-rules/create">
                      Create
                    </StyledNextButton>
                    <Button>Key Metrics</Button>
                    <Button>Risks</Button>
                    <Button>Alerts</Button>
                  </ButtonGroup>

                  <InputGroup maxW="xs">
                    <InputLeftElement pointerEvents="none">
                      <SearchIcon color={"gray.200"} boxSize="5" />
                    </InputLeftElement>
                    <Input placeholder="Search" />
                  </InputGroup>
                </Flex>
              </Container>
            </>
          )}
        </Box>
      </ColorModeContext.Provider>
      <DrawerNav isAdmin={true} isOpen={isOpen} onClose={onClose} />
    </Box>
  );
};

export const StyledNextButton = (props: any) => (
  <NextButton
    opacity={0.8}
    _activeLink={{ fontWeight: "bold", opacity: 1 }}
    {...props}
  />
);

export const NextButton: React.FC<
  {
    href: string;
  } & ButtonProps
> = ({ href, ...buttonProps }) => {
  return (
    <Button
      as={Link}
      to={href}
      preload={5000}
      getActiveProps={() => ({ style: { fontWeight: "bold", opacity: 1 } })}
      {...buttonProps}
    />
  );
};
