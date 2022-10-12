import {
  Flex,
  Stack,
  Heading,
  Button,
  Text,
  Link as ChakraLink,
} from "@chakra-ui/react";
import { Link, useNavigate } from "react-location";

import React, { Component, ReactNode } from "react";
import { ApprovalsLogo } from "../components/icons/Logos";

interface Props {
  error: Error;
}

export const UnhandledError = (props: Props) => {
  const navigate = useNavigate();

  return (
    <>
      <Flex
        height="100vh"
        padding="0"
        alignItems="center"
        justifyContent="center"
      >
        <Stack textAlign="center" w="70%">
          <Heading pb="10px">An unexpected error occurred.</Heading>
          {props.error && <Text>{props.error.message}</Text>}

          <Button
            top="40px"
            alignSelf="center"
            size="md"
            w="40%"
            onClick={() => {
              navigate({ to: "./admin" });
            }}
          >
            Back to home page
          </Button>
          <Flex alignItems="center" justifyContent="center">
            <ChakraLink
              as={Link}
              to={".."}
              transition="all .2s ease"
              rounded="sm"
              pt="75px"
            >
              <ApprovalsLogo h="42px" w="auto" />
            </ChakraLink>
          </Flex>
        </Stack>
      </Flex>
    </>
  );
};

export default UnhandledError;
