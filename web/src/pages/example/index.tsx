import {
  Box,
  Container,
  FormControl,
  FormLabel,
  Heading,
  Select,
  VStack,
} from "@chakra-ui/react";
import { Helmet } from "react-helmet";
import { UserLayout } from "../../components/Layout";
import { UserReviewsTable } from "../../components/tables/UserReviewsTable";
import {
  useUserListEntitlementResources,
  useUserListEntitlements,
} from "../../utils/backend-client/default/default";
import { useState } from "react";

type EntitlementContext = {
  Kind: string;
  Publisher: string;
  Name: string;
  Version: string;
  ResourceType: string;
};

const Home = () => {
  const [Entitlement, setEntitlement] = useState<EntitlementContext>({
    Publisher: "common-fate",
    Name: "AWS",
    Version: "v0.1.0",
    Kind: "Account",
    ResourceType: "accountId",
  });

  const { data } = useUserListEntitlements();

  // const select = document.getElementById("entitlement-select");
  // if (select) {
  //   //@ts-ignore
  //   if (select.options.length > 0 && select.selectedIndex) {
  //     //@ts-ignore
  //     const value = select.options[select.selectedIndex].value;
  //     const splitEnt = value.split("-");
  //     setEntitlement({
  //       Publisher: splitEnt[0],
  //       Name: splitEnt[1],
  //       Version: splitEnt[2],
  //       Kind: splitEnt[3],
  //       ResourceType: "accountId",
  //     });
  //   }
  // }

  const { data: resources } = useUserListEntitlementResources({
    kind: Entitlement.Kind,
    name: Entitlement.Name,
    publisher: Entitlement.Publisher,
    resourceType: Entitlement.ResourceType,
    version: Entitlement.Version,
  });

  return (
    <div>
      <Helmet>
        <title>Example</title>
      </Helmet>
      <UserLayout>
        <Box overflow="auto">
          <Container minW="864px" maxW="container.xl">
            <VStack>
              <FormControl>
                <FormLabel>Select an entitlement</FormLabel>
                <Select>
                  {data?.map((e) => {
                    return (
                      <option
                        value={
                          e.Kind.publisher +
                          "-" +
                          e.Kind.name +
                          "-" +
                          e.Kind.version +
                          "-" +
                          e.Kind.kind
                        }
                      >{`${e.Kind.name} ${e.Kind.version}`}</option>
                    );
                  })}
                </Select>
              </FormControl>
            </VStack>

            {/* <VStack>
              <Heading>Select an option for Account</Heading>
              {resources?.resources && (
                <Select>
                  {resources?.resources.map((e) => {
                    return <option>{`${e.name}: ${e.value} `}</option>;
                  })}
                </Select>
              )}
            </VStack> */}
          </Container>
        </Box>
      </UserLayout>
    </div>
  );
};

export default Home;
