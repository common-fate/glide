import { Box, Flex, Stack } from "@chakra-ui/react";
import { Helmet } from "react-helmet";
import { EntitlementCheckout } from "../../components/EntitlementCheckout";
import { UserLayout } from "../../components/Layout";
import { RecentRequests } from "../../components/RecentRequests";

const Home = () => {
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
            <Flex w={["770px"]}>
              <EntitlementCheckout />
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
