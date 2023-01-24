import { Box } from "@chakra-ui/react";
import React from "react";
import { Navbar } from "./Navbar";

export const UserLayout: React.FC<{ children?: React.ReactNode }> = ({
  children,
}) => {
  return (
    <Box as="main" h="100%" minH="100vh">
      <Navbar />
      {children}
    </Box>
  );
};
