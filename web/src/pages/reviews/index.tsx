import { Box, Container } from "@chakra-ui/react";
import { UserLayout } from "../../components/Layout";
import { UserReviewsTable } from "../../components/tables/UserReviewsTable";

const Home = () => {
  return (
    <div>
      <UserLayout>
        <Box overflow="auto">
          <Container minW="864px" maxW="container.xl">
            <UserReviewsTable />
          </Container>
        </Box>
      </UserLayout>
    </div>
  );
};

export default Home;
