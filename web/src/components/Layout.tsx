import { Box, Spacer, Text } from "@chakra-ui/react";
import React from "react";
import { useAdminGetDeploymentVersion } from "../utils/backend-client/admin/admin";
import { AdminNavbar } from "./nav/AdminNavbar";
import { Navbar } from "./nav/Navbar";

export const AdminLayout: React.FC<{ children?: React.ReactNode }> = ({
  children,
}) => {
  const { data } = useAdminGetDeploymentVersion();

  return (
    <>
      <Box as="main" h="100%" minH="100vh">
        <AdminNavbar />
        {children}
      </Box>
      <Box as="footer" marginTop={"-36px"} px={5}>
        <Text textStyle={"Body/ExtraSmall"}>
          {data?.version !== undefined && "Version: " + data.version}
        </Text>
      </Box>
    </>
  );
};

export const UserLayout: React.FC<{ children?: React.ReactNode }> = ({
  children,
}) => {
  return (
    <main>
      <Navbar />
      {children}
    </main>
  );
};
