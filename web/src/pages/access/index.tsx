import { InfoIcon } from "@chakra-ui/icons";
import { Box, Center, Flex, Spinner, Text, useToast } from "@chakra-ui/react";
import axios, { Axios, AxiosError } from "axios";
import { useEffect, useState } from "react";
import {
  BuildNextOptions,
  Link,
  MakeGenerics,
  useNavigate,
  useSearch,
} from "react-location";
import { CFCodeMultiline } from "../../components/CodeInstruction";
import { UserLayout } from "../../components/Layout";
import { OnboardingCard } from "../../components/OnboardingCard";
import { SelectRuleTable } from "../../components/tables/SelectRuleTable";
import { useAccessRuleLookup } from "../../utils/backend-client/default/default";
import {
  CreateRequestWith,
  LookupAccessRule,
} from "../../utils/backend-client/types";
import { AccessRuleLookupParams } from "../../utils/backend-client/types/accessRuleLookupParams";
import { RequestFormQueryParameters } from "./request/[id]";
// const a: RequestFormQueryParameters = {
//   Search: {
//     reason: getValues("reason"),
//     with: (getValues("with") || [])
//       .filter((fw) => !fw.hidden)
//       .map((fw) => fw.data),
//   },
// };
// const timing: RequestTiming = {
//   durationSeconds: getValues("timing.durationSeconds"),
// };
// if (getValues("when") === "scheduled") {
//   timing.startTime = new Date(getValues("startDateTime")).toISOString();
// }
// a.Search.timing = timing;
// const u = new URL(window.location.href);
// u.search = location.stringifySearch(a.Search);
// setUrlClipboardValue(u.toString());
/**
 * makeLookupAccessRuleRequestLink adds request parameters to the URLso that forms can be prepopulated when there are multiple options for a user
 * @param lookupResult
 * @returns
 */
export const makeLookupAccessRuleRequestLink = (
  lookupResult: LookupAccessRule
): BuildNextOptions<RequestFormQueryParameters> => {
  const w: CreateRequestWith = {};
  lookupResult.selectableWithOptionValues?.forEach(
    (o) => (w[o.key] = [o.value])
  );
  const a: RequestFormQueryParameters = {
    Search: {
      with: [w],
    },
  };
  return {
    to: `/access/request/${lookupResult.accessRule.id}`,
    search: a.Search,
  };
};
const Access = () => {
  type MyLocationGenerics = MakeGenerics<{
    Search: AccessRuleLookupParams;
  }>;

  const search = useSearch<MyLocationGenerics>();

  const navigate = useNavigate();

  const { data, isValidating, error } = useAccessRuleLookup(search);

  const toast = useToast();
  useEffect(() => {
    if (data?.length == 1) {
      navigate(makeLookupAccessRuleRequestLink(data[0]));
    }
  }, [search, data]);

  // navigate away if there was an error
  useEffect(() => {
    // prevent a race condition where the access rule is looked up again after navigating away from the page
    if (error && location.pathname === "access") {
      // if there were search params, then show an error, else just redirect to the requests page
      if (Object.entries(search).length > 0) {
        if (axios.isAxiosError(error)) {
          const e = error as AxiosError<{ error: string }>;
          toast({
            title: "Something went wrong loading access rules for your query",
            description: e?.response?.data.error,
            status: "error",
            variant: "subtle",
            duration: 5000,
            isClosable: true,
          });
        } else {
          toast({
            title: "Unknown error while loading access rules for your query",
            status: "error",
            variant: "subtle",
            duration: 5000,
            isClosable: true,
          });
        }
      }
      navigate({ to: "/requests" });
    }
  }, [error, search]);

  if (!data && isValidating) {
    return <Spinner my={4} pos="absolute" left="50%" top="50vh" />;
  }
  if (error) {
    return null;
  }
  if (!data || data.length == 1) {
    return <Spinner my={4} pos="absolute" left="50%" top="50vh" />;
  }
  return (
    <UserLayout>
      <Center h="80vh">
        <Flex flexDir="column" align="center" minH="60vh" w="100%">
          <br />
          {data.length > 1 && (
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
          {data.length == 0 && (
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
  );
};

export default Access;
