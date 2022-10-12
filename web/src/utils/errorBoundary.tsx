import {
  Flex,
  Stack,
  Heading,
  Button,
  Text,
  Link as ChakraLink,
} from "@chakra-ui/react";
import { Link } from "react-location";

import React, { Component, ReactNode } from "react";
import { ApprovalsLogo } from "../components/icons/Logos";

interface Props {
  children?: ReactNode;
}

interface State {
  error: Error | undefined;
}

class ErrorBoundary extends Component<Props, State> {
  public state: State = {
    error: undefined,
  };

  public static getDerivedStateFromError(error: Error): State {
    // Update state so the next render will show the fallback UI.
    return { error: error };
  }

  public render() {
    if (this.state.error != undefined) {
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
              {this.state.error && <Text>{this.state.error.message}</Text>}
              <Link to={"/requests"}>
                <Button top="40px" alignSelf="center" size="md" w="40%">
                  Back to home page
                </Button>
              </Link>
              <Flex alignItems="center" justifyContent="center">
                <ChakraLink
                  as={Link}
                  to={"/requests"}
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
    }

    return this.props.children;
  }
}

export default ErrorBoundary;
