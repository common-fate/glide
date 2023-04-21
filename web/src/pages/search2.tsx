import {
  Box,
  Button,
  Center,
  chakra,
  Code,
  Container,
  Flex,
  Text,
  Heading,
  Input,
  Stack,
  useBoolean,
  useDisclosure,
  useEventListener,
} from "@chakra-ui/react";
import { Command } from "cmdk";
import { Command as CommandNew } from "../utils/cmdk";
import React from "react";
import { ProviderIcon, ShortTypes } from "../components/icons/providerIcon";
import { UserLayout } from "../components/Layout";
import {
  userPostRequests,
  useUserListEntitlements,
  useUserListEntitlementTargets,
} from "../utils/backend-client/default/default";
import { ArrowBackIcon, CheckCircleIcon } from "@chakra-ui/icons";
import { Link, useRouter } from "react-location";
import Counter from "../components/Counter";

const search2 = () => {
  const targets = useUserListEntitlementTargets(
    {},
    {
      swr: { refreshInterval: 10000 },
      request: {
        baseURL: "http://127.0.0.1:3100",
        headers: {
          Prefer: "code=200, example=example_targets",
        },
      },
    }
  );

  return (
    <UserLayout>
      <Container mt={24}>
        {/* main */}
        <Text textStyle="Body/Medium">Access</Text>
        <Stack spacing={1} borderColor="neutrals.300" rounded="md">
          {[].map((target) => {
            return <>ok</>;
          })}
        </Stack>

        {/* buttons */}
        <Flex w="100%" mt={4}>
          <Button
            ml="auto"
            // disabled={checked.length == 0}
            // onClick={handleSubmit}
            variant="brandSecondary"
            leftIcon={<ArrowBackIcon />}
            to="/search"
            as={Link}
          >
            Go back
          </Button>
          <Button
            ml="auto"
            // disabled={checked.length == 0}
            // onClick={handleSubmit}
            // isLoading={submitLoading}
            loadingText="Processing request..."
          >
            Next (⌘+Enter)
          </Button>
        </Flex>
      </Container>
    </UserLayout>
  );
};

export default search2;
