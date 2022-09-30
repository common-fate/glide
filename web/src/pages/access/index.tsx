import {
  Box,
  Button,
  Center,
  chakra,
  Code,
  Flex,
  Spinner,
  Stack,
  Text,
} from "@chakra-ui/react";
import React, { useEffect } from "react";
import { MakeGenerics, useNavigate, useSearch } from "react-location";
import { ProviderIcon } from "../../components/icons/providerIcon";
import { UserLayout } from "../../components/Layout";
import { useAccessRuleLookup } from "../../utils/backend-client/default/default";
import { AccessRuleLookupParams } from "../../utils/backend-client/types/accessRuleLookupParams";
import { Link } from "react-location";
import { SelectRuleTable } from "../../components/tables/SelectRuleTable";
import { CodeInstruction } from "../../components/CodeInstruction";
import { useGetMe } from "../../utils/backend-client/end-user/end-user";
import { OnboardingCard } from "../../components/OnboardingCard";
import { InfoIcon } from "@chakra-ui/icons";

const Access = () => {
  type MyLocationGenerics = MakeGenerics<{
    Search: AccessRuleLookupParams;
  }>;

  const search = useSearch<MyLocationGenerics>();

  const navigate = useNavigate();

  const { data, isValidating } = useAccessRuleLookup(search);

  const me = useGetMe();

  useEffect(() => {
    if (data?.accessRules.length == 1) {
      navigate({
        to: `/access/request/${data.accessRules[0].id}?accountId=${search.accountId}`,
      });
    }
  }, [search, data]);

  return data && !isValidating ? (
    <UserLayout>
      <Center h="80vh">
        <Flex flexDir="column" align="center" minH="60vh" w="100%">
          <br />
          {data && data.accessRules.length > 1 && (
            <Flex flexDir="column" alignItems="center" w="100%">
              <Box w={{ base: "100%", md: "60ch" }}>
                <OnboardingCard
                  // my={4}
                  title="Multiple access rules found"
                  leftIcon={<InfoIcon color="brandBlue.200" />}
                >
                  <Text>
                    If you’re unsure what to choose, choose one based on the
                    name that matches what you’re trying to do
                  </Text>
                </OnboardingCard>
                <br />
              </Box>
              <SelectRuleTable rules={data.accessRules} />
            </Flex>
          )}
          {data && data.accessRules.length == 0 && (
            <>
              We couldn't find any access rules for you
              <CodeInstruction
                mt={2}
                textAlign="left"
                inline={false}
                // @ts-ignore
                children={`Access rule not found, details below:
${JSON.stringify(search, null, 2)}`}
              />
              <Flex _hover={{ textDecor: "underline" }} mt={12}>
                <Link to="/">← Return To Home </Link>
              </Flex>
            </>
          )}
        </Flex>
      </Center>
    </UserLayout>
  ) : (
    <Spinner
      my={4}
      opacity={isValidating ? 1 : 0}
      pos="absolute"
      left="50%"
      top="50vh"
    />
  );
};

export default Access;
