import { Center } from "@chakra-ui/react";
import { useEffect } from "react";
import { useNavigate } from "react-location";
import { useCognito } from "../utils/context/cognitoContext";

const Logout = () => {
  const auth = useCognito();

  //This gets unmounted and remounted once before an actual logout occurs - causing a 'double' logout redirect
  useEffect(() => {
    auth.initiateSignOut().catch((e) => console.error(e));
  }, []);

  return <Center h="80vh">Logging you out 👋</Center>;
};

export default Logout;
