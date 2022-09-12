import React, { useEffect } from "react";
import { MakeGenerics, useNavigate, useSearch } from "react-location";
import { UserLayout } from "../components/Layout";
import { AccessRuleLookupParams } from "src/utils/backend-client/types/accessRuleLookupParams";
import { Box, Center, Flex, Spinner } from "@chakra-ui/react";

const assume = () => {
  type MyLocationGenerics = MakeGenerics<{
    Search: AccessRuleLookupParams;
  }>;

  const search = useSearch<MyLocationGenerics>();

  const navigate = useNavigate();

  useEffect(() => {
    if (search.accountId && search.roleName && search.type) {
      // Run account lookup
      //   wait 3 seconds then redirect to http://localhost:3000/access/request/rul_2Eenkkrj4ka2uChX9BIoALuo3Ws
      setTimeout(
        () =>
          navigate({ to: "/access/request/rul_2Eenkkrj4ka2uChX9BIoALuo3Ws" }),
        2000
      );
    }
  }, [search]);

  return (
    <UserLayout>
      <Center h="80vh">
        <Flex flexDir="column" align="center" textAlign="center">
          <Spinner my={4} />
          Finding your access request now...
          <br />
          {/* {JSON.stringify(search)} */}
        </Flex>
      </Center>
    </UserLayout>
  );
};

export default assume;
