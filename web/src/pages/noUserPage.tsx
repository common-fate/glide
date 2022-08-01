import { Flex, Text, Stack, Heading, Button } from "@chakra-ui/react";
import { ICredentials } from "@aws-amplify/core";

import { useNavigate } from "react-location";

interface Props {
  userEmail?: string;
  initiateSignOut: () => Promise<any>;
}

export const NoUser = (props: Props) => {
  const navigate = useNavigate();
  return (
    <Flex
      height="100vh"
      padding="0"
      alignItems="center"
      justifyContent="center"
    >
      <Stack textAlign="center" w="70%">
        <Heading pb="50px">An error occurred signing you in</Heading>
        <Text>
          You&apos;ve successfully logged in, but we couldn&apos;t find a
          matching user account for you in our database. (
          {props.userEmail
            ?.split("_")
            .slice(1, props.userEmail?.split("_").length)
            .join()}
          ){/* Removes prefixed idp provider that amplify adds */}
        </Text>
        <Text>
          This is likely because your user directory settings are misconfigured.
          Check that your configuration variables (including client IDs and
          secrets) in your granted-deployment.yml file are correct.
        </Text>

        <Text>
          If you need help debugging this, contact us at:{" "}
          <a href="mailto:hello@commonfate.io">hello@commonfate.io</a>
        </Text>

        <Button
          onClick={() => {
            console.log("clicked");
            props.initiateSignOut().then(() => {
              window.location.reload();
            });
          }}
          top="40px"
          alignSelf="center"
          size="md"
          w="40%"
        >
          Back to login
        </Button>
      </Stack>
    </Flex>
  );
};

export default NoUser;
