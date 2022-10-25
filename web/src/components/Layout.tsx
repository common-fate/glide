import { Text } from "@chakra-ui/react";
import React from "react";
import { useAdminGetDeploymentVersion } from "../utils/backend-client/admin/admin";
import { AdminNavbar } from "./nav/AdminNavbar";
import { Navbar } from "./nav/Navbar";

export const AdminLayout: React.FC<{ children?: React.ReactNode }> = ({
  children,
}) => {
  const { data } = useAdminGetDeploymentVersion();

  return (
    <main>
      <AdminNavbar />
      {children}
      <Text position="fixed" bottom={4} left={4} textStyle={"Body/ExtraSmall"}>
        {data?.version !== undefined && "Version: " + data.version}
      </Text>
    </main>
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
