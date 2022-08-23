import { Center } from "@chakra-ui/react";
import { useCognito } from "../utils/context/cognitoContext";

const Logout = () => {
  const auth = useCognito();
  auth.initiateSignOut().catch((e) => console.error(e));

  return <Center h="80vh">Logging you out 👋</Center>;
};

export default Logout;
