import { Center } from "@chakra-ui/react";
import { useCognito } from "../utils/context/cognitoContext";

const Logout = () => {
  const auth = useCognito();
  auth.initiateSignOut();

  return <Center h="80vh">Logging you out 👋</Center>;
};

export default Logout;
