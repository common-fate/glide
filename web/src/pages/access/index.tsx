import { InfoIcon } from "@chakra-ui/icons";
import { Box, Center, Flex, Spinner, Text } from "@chakra-ui/react";
import { useEffect } from "react";
import { Link, MakeGenerics, useNavigate, useSearch } from "react-location";
import { CFCodeMultiline } from "../../components/CodeInstruction";
import { UserLayout } from "../../components/Layout";
import { OnboardingCard } from "../../components/OnboardingCard";
import { SelectRuleTable } from "../../components/tables/SelectRuleTable";
import { useAccessRuleLookup } from "../../utils/backend-client/default/default";
import { useGetMe } from "../../utils/backend-client/end-user/end-user";
import { LookupAccessRule } from "../../utils/backend-client/types";
import { AccessRuleLookupParams } from "../../utils/backend-client/types/accessRuleLookupParams";

/**
 * makeLookupAccessRuleRequestLink adds request parameters to the URLso that forms can be prepopulated when there are multiple options for a user
 * @param lookupResult
 * @returns
 */
export const makeLookupAccessRuleRequestLink = (
  lookupResult: LookupAccessRule
) => {
  const query = lookupResult.selectableWithOptionValues
    ?.map((o) => `${o.key}=${o.value}`)
    .join("&");
  return `/access/request/${lookupResult.accessRule.id}${
    query ? `?${query}` : ""
  }`;
};
const Access = () => {
  type MyLocationGenerics = MakeGenerics<{
    Search: AccessRuleLookupParams;
  }>;

  const search = useSearch<MyLocationGenerics>();

  const navigate = useNavigate();

  const { data, isValidating } = useAccessRuleLookup(search);

  const me = useGetMe();

  useEffect(() => {
    if (data?.length == 1) {
      navigate({
        to: makeLookupAccessRuleRequestLink(data[0]),
      });
    }
  }, [search, data]);

  return data && !isValidating ? (
    <UserLayout>
      <Center h="80vh">
        <Flex flexDir="column" align="center" minH="60vh" w="100%">
          <br />
          {data && data.length > 1 && (
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
              <SelectRuleTable rules={data} />
            </Flex>
          )}
          {data && data.length == 0 && (
            <>
              <Text mb={2}>We couldn't find any access rules for you</Text>
              <CFCodeMultiline
                text={`Access rule not found, details below:
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
