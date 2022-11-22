import { LockIcon } from "@chakra-ui/icons";
import {
  Box,
  ColorModeContext,
  Divider,
  Drawer,
  DrawerBody,
  DrawerCloseButton,
  DrawerContent,
  DrawerOverlay,
  DrawerProps,
  Stack,
} from "@chakra-ui/react";
import * as React from "react";
import { useCognito } from "../../utils/context/cognitoContext";
import { useUser } from "../../utils/context/userContext";
import { useInnerHeight } from "../../utils/hooks/useInnerHeight";
import Counter from "../Counter";
import { DoorIcon } from "../icons/Icons";
import { CommonFateAdminLogo, CommonFateLogo } from "../icons/Logos";
import { UserAvatarDetails } from "../UserAvatar";
import { StyledNextButton } from "./AdminNavbar";
import { StyledButton } from "./Navbar";

type DrawerNavProps = Omit<DrawerProps, "children"> & { isAdmin?: boolean };

export const DrawerNav = ({ isAdmin, ...props }: DrawerNavProps) => {
  const user = useUser();
  const auth = useCognito();
  const innerHeight = useInnerHeight();

  type NavItem = {
    label: string;
    href: string;
    count?: number;
  };

  const navItems: NavItem[] = [
    {
      label: "Access Requests",
      href: "/requests",
    },
    {
      label: "Reviews",
      href: "/reviews?status=pending",
      //   count: reviews?.requests?.length ?? 0,
    },
  ];

  const adminNavItems: NavItem[] = [
    {
      label: "Access Rules",
      href: "/admin/access-rules",
    },
    {
      label: "Requests",
      href: "/admin/requests",
    },
    {
      label: "Users",
      href: "/admin/users",
    },
    {
      label: "Groups",
      href: "/admin/groups",
    },
    {
      label: "Providers",
      href: "/admin/providers",
    },
  ];

  return (
    <ColorModeContext.Provider
      value={{
        colorMode: isAdmin ? "dark" : "light",
        // noop
        toggleColorMode: () => {
          undefined;
        },
        setColorMode: () => {
          undefined;
        },
      }}
    >
      <Drawer
        isOpen={props.isOpen}
        placement="left"
        onClose={props.onClose}
        isFullHeight={true}
        // size="xs"
        // finalFocusRef={btnRef}
      >
        <DrawerOverlay />
        <DrawerCloseButton
          pos="absolute"
          zIndex={9999}
          bg="neutrals.200"
          _hover={{ bg: "neutrals.300" }}
          rounded="full"
        />
        <DrawerContent h={innerHeight}>
          <DrawerBody mt={10} h="100%">
            {isAdmin ? (
              <CommonFateAdminLogo h="32px" w="auto" />
            ) : (
              <CommonFateLogo h="32px" w="auto" />
            )}
            <Stack spacing={1} mt={8}>
              {(isAdmin ? adminNavItems : navItems).map((el) => (
                <StyledNextButton
                  key={el.href}
                  w="100%"
                  justifyContent="flex-start"
                  variant="ghost"
                  href={el.href}
                  pr={10}
                >
                  {el.label}
                  {el.count && (
                    <Counter pos="absolute" right={2} count={el.count} />
                  )}
                </StyledNextButton>
              ))}
            </Stack>

            <Box pos="absolute" bottom={8} w="272px">
              <Divider />
              <Stack spacing={3} my={4}>
                <StyledButton
                  justifyContent="flex-start"
                  leftIcon={
                    <DoorIcon color={isAdmin ? "gray.400" : "gray.700"} />
                  }
                  onClick={auth.initiateSignOut}
                  opacity=".8"
                >
                  Log out
                </StyledButton>
                <StyledNextButton
                  variant="ghost"
                  justifyContent="flex-start"
                  leftIcon={
                    <LockIcon color={isAdmin ? "gray.400" : "gray.700"} />
                  }
                  href={isAdmin ? "/" : "/admin/access-rules"}
                >
                  {isAdmin ? "Switch To User" : "Switch To Admin"}
                </StyledNextButton>
              </Stack>

              <Box mt={2} pl={2}>
                <UserAvatarDetails size="sm" user={user.user?.id} />
              </Box>
            </Box>
          </DrawerBody>
        </DrawerContent>
      </Drawer>
    </ColorModeContext.Provider>
  );
};
