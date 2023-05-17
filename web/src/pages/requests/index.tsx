import { Box, Flex, Stack, Text } from "@chakra-ui/react";
import { Helmet } from "react-helmet";
import { EntitlementCheckout } from "../../components/EntitlementCheckout";
import { UserLayout } from "../../components/Layout";
import { RecentRequests } from "../../components/RecentRequests";
import { AccessTemplateList } from "../../components/AccessTemplatesList";
import { useState } from "react";

const Home = () => {
  const [checked, setChecked] = useState<Set<string>>(new Set());

  return (
    <>
      <UserLayout>
        <Helmet>
          <title>Common Fate</title>
        </Helmet>
        <Box overflowY={"auto"}>
          <Stack
            direction={{ lg: "row", md: "column", sm: "column" }}
            pt={{ base: 12, lg: 32 }}
            spacing="50px"
            justify={"center"}
            align={{ lg: "flex-start", md: "center", sm: "center" }}
          >
            <Flex w={["350px"]}>
              <AccessTemplateList setChecked={setChecked} />
            </Flex>
            <Flex w={["770px"]}>
              <EntitlementCheckout checked={checked} setChecked={setChecked} />
            </Flex>
            <Flex w={["550px"]}>
              <RecentRequests />
            </Flex>
          </Stack>
        </Box>
      </UserLayout>
    </>
  );
};

export default Home;
