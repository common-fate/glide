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

import { Entitlement, Resource } from "../../utils/backend-client/types";
import { useEffect, useState } from "react";

interface EntitlementStore {
  [key: string]: Entitlement;
}

interface ResourceStore {
  [key: string]: Resource[];
}

const Home = () => {
  const { data } = useUserListEntitlements();

  //@ts-ignore
  const [ent, setEnt] = useState<Entitlement>({});

  const [filters, setFilters] = useState<string[]>([]);

  const entitlementStore: EntitlementStore = {};
  const resourceStore: ResourceStore = {};

  useEffect(() => {
    if (data) {
      setEnt(data[1]);
    }
  }, [data]);
  // //key value pair to figure out which entitlement we have selected
  if (data) {
    data.forEach((e) => {
      const name =
        e.Kind.publisher +
        "-" +
        e.Kind.name +
        "-" +
        e.Kind.version +
        "-" +
        e.Kind.kind;
      entitlementStore[name] = e;

      //initialise resource store
      Object.entries(ent.Schema).map(([key, val]) => {
        resourceStore[key] = [];
      });
    });
  }
  const { data: resources } = useUserListEntitlementResources({
    kind: data && ent.Kind ? ent.Kind.kind : "",
    publisher: data && ent.Kind ? ent.Kind.publisher : "",
    name: data && ent.Kind ? ent.Kind.name : "",
    version: data && ent.Kind ? ent.Kind.version : "",
    resourceType: data && ent.Kind ? "accountId" : "accountId",
    filters: filters.length > 0 ? filters[0] : undefined,
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
            {!resources ||
              (!resources.resources && (
                <Text>You dont seem to have access to this entitlement</Text>
              ))}
            {ent.Kind &&
              resources &&
              resources.resources &&
              Object.entries(ent.Schema).map(([key, obj]) => {
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
