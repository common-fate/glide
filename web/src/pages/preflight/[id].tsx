import { ArrowBackIcon } from "@chakra-ui/icons";
import {
  Center,
  Container,
  Flex,
  IconButton,
  Text,
  Textarea,
  useToast,
} from "@chakra-ui/react";
import { Helmet } from "react-helmet";
import { Link, useMatch } from "react-location";

import { UserLayout } from "../../components/Layout";
// import { useUserGetPreflight } from "../../utils/backend-client/default/default";

const Home = () => {
  const {
    params: { id: preflightId },
  } = useMatch();

  const toast = useToast();
  return (
    <div>
      <UserLayout>
        <Helmet>
          <title>Preflight</title>
        </Helmet>

        <Center>
          <Flex
            background="neutrals.100"
            px={"200px"}
            pt={"100px"}
            pb={"150px"}
          >
            <Textarea></Textarea>
          </Flex>
        </Center>
      </UserLayout>
    </div>
  );
};

export default Home;
