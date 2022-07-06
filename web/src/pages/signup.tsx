import {
  Box,
  Button,
  Container,
  Divider,
  Flex,
  Heading,
  Stack,
  useBreakpointValue,
} from "@chakra-ui/react";
import { ApprovalsLogo } from "../components/icons/Logos";
import SvgAwsLogo from "../components/icons/SvgAwsLogo";
import { useUser } from "../utils/context/userContext";

export const Signup = () => {
  const { initiateAuth } = useUser();

  return (
    <Box bg="neutrals.100" minH="100vh">
      <Container maxW="md" py={{ base: "12", md: "24" }} minH="90vh">
        <Stack spacing="8">
          <Flex
            flexDir="column"
            // spacing="6"
            align="center"
            mt={{ base: 8, md: 12 }}
          >
            <Flex
              rounded="full"
              // borderColor="neutrals.500"
              borderWidth="1px"
              borderStyle="solid"
              py={3}
              //   px={5}
              bg="white"
              //   w="min-content"
            >
              <ApprovalsLogo h="32px" w="200px" />
            </Flex>
            <Heading
              mt={16}
              fontWeight="bold"
              size={useBreakpointValue({ base: "xs", md: "lg" })}
            >
              Sign in to your account
            </Heading>
          </Flex>
          <Stack spacing="6" alignItems="center">
            <Button
              variant="solid"
              leftIcon={<SvgAwsLogo height="18px" />}
              iconSpacing="3"
              // bg="white"
              w="350px"
              onClick={initiateAuth}
            >
              Continue with AWS Cognito
            </Button>
            <Divider />
            {/* <Stack spacing="4">
          <Input placeholder="Enter your email" />
          <Button variant="primary">Continue with email</Button>
        </Stack> */}
          </Stack>
          {/* <Button variant="link" colorScheme="blue" size="sm">
          Continue using Single Sign-on (SSO)
        </Button> */}
        </Stack>
      </Container>
    </Box>
  );
};

export default Signup;
