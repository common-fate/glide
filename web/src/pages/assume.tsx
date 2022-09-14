import {
  Box,
  Button,
  Center,
  Flex,
  Spinner,
  Stack,
  Text,
} from "@chakra-ui/react";
import React, { useEffect } from "react";
import { MakeGenerics, useNavigate, useSearch } from "react-location";
import { ProviderIcon } from "../components/icons/providerIcon";
import { UserLayout } from "../components/Layout";
import { useAccessRuleLookup } from "../utils/backend-client/default/default";
import { AccessRuleLookupParams } from "../utils/backend-client/types/accessRuleLookupParams";
import { Link } from "react-location";

const assume = () => {
  type MyLocationGenerics = MakeGenerics<{
    Search: AccessRuleLookupParams;
  }>;

  const search = useSearch<MyLocationGenerics>();

  const [loadText, setLoadText] = React.useState(
    "Finding your access request now..."
  );

  const navigate = useNavigate();

  const { data, isValidating } = useAccessRuleLookup(search);

  useEffect(() => {
    // Run account lookup
    if (search.accountId && search.roleName && search.type) {
      if (data && data.accessRules.length > 0) {
        if (data.accessRules.length == 1) {
          setLoadText("Access rule found 🚀 Redirecting now...");
          setTimeout(() => {
            // navigate({ to: "/access/request/" + data.accessRules[0].id });
          }, 300);
        } else {
          // add handling for multi rule resolution...
          setLoadText("Multiple access rules found, choose one to continue");
        }
      }
    }
  }, [search, data?.accessRules]);

  return (
    <UserLayout>
      <Center h="80vh">
        <Flex flexDir="column" align="center" textAlign="center" minH="400px">
          <Spinner my={4} opacity={isValidating ? 1 : 0} />
          {loadText}
          <br />
          <Stack mt={4} direction="row" spacing={4}>
            {data &&
              data.accessRules.length > 1 &&
              data.accessRules.map((r) => (
                <Link
                  style={{ display: "flex" }}
                  to={"/access/request/" + r.id}
                  key={r.id}
                >
                  <Box
                    className="group"
                    textAlign="center"
                    bg="neutrals.100"
                    p={6}
                    h="172px"
                    w="232px"
                    rounded="md"
                  >
                    <ProviderIcon
                      shortType={r.target.provider.type}
                      mb={3}
                      h="8"
                      w="8"
                    />

                    <Text textStyle="Body/SmallBold" color="neutrals.700">
                      {r.name}
                    </Text>

                    <Button
                      mt={4}
                      variant="brandSecondary"
                      size="sm"
                      opacity={0}
                      sx={{
                        // This media query ensure always visible for touch screens
                        "@media (hover: none)": {
                          opacity: 1,
                        },
                      }}
                      transition="all .2s ease-in-out"
                      transform="translateY(8px)"
                      _groupHover={{
                        bg: "white",
                        opacity: 1,
                        transform: "translateY(0px)",
                      }}
                    >
                      Request
                    </Button>
                  </Box>
                </Link>
              ))}
          </Stack>
        </Flex>
      </Center>
    </UserLayout>
  );
};

export default assume;
