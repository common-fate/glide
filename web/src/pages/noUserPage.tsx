import { Button, Flex, Heading, Stack, Text } from "@chakra-ui/react";

import { useCognito } from "../utils/context/cognitoContext";

export const NoUser = () => {
  const { cognitoAuthenticatedUserEmail, initiateSignOut } = useCognito();
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
          {cognitoAuthenticatedUserEmail
            ?.split("_")
            .slice(1, cognitoAuthenticatedUserEmail?.split("_").length)
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
            initiateSignOut()
              .then(() => {
                window.location.reload();
              })
              .catch((e) => console.error(e));
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
