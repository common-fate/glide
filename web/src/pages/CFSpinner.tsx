import { Flex, Spinner } from "@chakra-ui/react";

export const CFSpinner = () => {
  return (
    <Flex
      height="100vh"
      padding="0"
      alignItems="center"
      justifyContent="center"
    >
      <Spinner />
    </Flex>
  );
};

export default CFSpinner;
