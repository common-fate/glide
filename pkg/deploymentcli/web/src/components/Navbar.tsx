import Auth from "@aws-amplify/auth";
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
  LightMode,
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
import { Link, useNavigate } from "react-location";
import { DoorIcon } from "./icons/Icons";
import { CommonFateAdminLogo, CommonFateLogo } from "./icons/Logos";

export const Navbar: React.FC = () => {
  const isDesktop = useBreakpointValue({ base: false, sm: true }, "800px");

  // const { isOpen, onOpen, onClose } = useDisclosure();
  // const modal = useDisclosure();

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
            <HStack spacing="4" pos="relative" w="100%" position="relative">
              <ChakraLink
                as={Link}
                to={"/"}
                transition="all .2s ease"
                rounded="sm"
              >
                <CommonFateLogo h="32px" w="auto" />
              </ChakraLink>
              {/* <Button
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
              </Button> */}
              {isDesktop && (
                <ButtonGroup
                  variant="ghost"
                  spacing="0"
                  mb={"-32px !important;"}
                  // ml="auto !important;"
                  // mr="auto !important;"
                  pos="absolute"
                  left="50%"
                  transform="translateX(-50%)"
                  bottom="15px"
                >
                  {/* I've hardcoded widths here to prevent the bold/unbold text from 
                  altering the containing divs width. Reduces *jittering* */}
                  <TabsStyledButton
                    href="/"
                    // w={showReqCount ? "190px" : "142px"}
                    // pr={showReqCount ? 10 : undefined}
                  >
                    Providers
                  </TabsStyledButton>
                  <TabsStyledButton
                    href="/registry"
                    w="125px"
                    // pr={showReqCount ? 10 : undefined}
                  >
                    Registry
                  </TabsStyledButton>
                </ButtonGroup>
              )}
            </HStack>
          </Flex>
        </Container>
      </Box>
      {/* <DrawerNav isOpen={isOpen} onClose={onClose} /> */}
      {/* <CommandPalette isOpen={modal.isOpen} onClose={modal.onClose} /> */}
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
