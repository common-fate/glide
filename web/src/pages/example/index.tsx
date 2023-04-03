import {
  Box,
  Container,
  FormControl,
  FormLabel,
  Heading,
  Select,
  VStack,
  Text,
} from "@chakra-ui/react";
import { Helmet } from "react-helmet";
import { UserLayout } from "../../components/Layout";
import { UserReviewsTable } from "../../components/tables/UserReviewsTable";
import {
  useUserListEntitlementResources,
  useUserListEntitlements,
} from "../../utils/backend-client/default/default";

import { TargetGroup, Resource } from "../../utils/backend-client/types";
import { useEffect, useState } from "react";

interface EntitlementStore {
  [key: string]: TargetGroup;
}

interface ResourceStore {
  [key: string]: Resource[];
}

const Home = () => {
  const { data } = useUserListEntitlements();

  //@ts-ignore
  const [ent, setEnt] = useState<TargetGroup>({});

  const [filters, setFilters] = useState<string[]>([]);

  const entitlementStore: EntitlementStore = {};
  const resourceStore: ResourceStore = {};

  useEffect(() => {
    if (data) {
      setEnt(data.targetGroups[1]);
    }
  }, [data]);
  // //key value pair to figure out which entitlement we have selected
  if (data) {
    data.targetGroups.forEach((e) => {
      const name =
        e.from.publisher +
        "-" +
        e.from.name +
        "-" +
        e.from.version +
        "-" +
        e.from.kind;
      entitlementStore[name] = e;

      //initialise resource store
      Object.entries(ent.schema).map(([key, val]) => {
        resourceStore[key] = [];
      });
    });
  }
  const { data: resources } = useUserListEntitlementResources({
    kind: data && ent.from ? ent.from.kind : "",
    publisher: data && ent.from ? ent.from.publisher : "",
    name: data && ent.from ? ent.from.name : "",
    version: data && ent.from ? ent.from.version : "",
    resourceType: data && ent.from ? "accountId" : "accountId",
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
                <Select
                  onChange={(e) =>
                    // console.log(entitlementStore[e.target.value])
                    setEnt(entitlementStore[e.target.value])
                  }
                >
                  {data?.targetGroups?.map((e) => {
                    return (
                      <option
                        value={
                          e.from.publisher +
                          "-" +
                          e.from.name +
                          "-" +
                          e.from.version +
                          "-" +
                          e.from.kind
                        }
                      >{`${e.from.name} ${e.from.version}`}</option>
                    );
                  })}
                </Select>
              </FormControl>
            </VStack>
            {!resources ||
              (!resources.resources && (
                <Text>You dont seem to have access to this entitlement</Text>
              ))}
            {ent.from &&
              resources &&
              resources.resources &&
              Object.entries(ent.schema).map(([key, obj]) => {
                return (
                  <VStack>
                    <Heading>
                      {`Select an option for
                       ${key}`}
                    </Heading>
                    {resources && resources?.resources && (
                      <Select
                        onChange={(e) =>
                          // console.log(entitlementStore[e.target.value])
                          setFilters((f) => f.concat(e.target.value))
                        }
                      >
                        {resourceStore[key].map((e: Resource) => {
                          return (
                            <option
                              value={e.value}
                            >{`${e.name}: ${e.value} `}</option>
                          );
                        })}
                      </Select>
                    )}
                  </VStack>
                );
              })}
          </Container>
        </Box>
      </UserLayout>
    </div>
  );
};

export default Home;
