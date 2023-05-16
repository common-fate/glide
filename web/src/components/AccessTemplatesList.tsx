import {
  Box,
  Flex,
  Grid,
  HStack,
  Link,
  LinkBox,
  LinkBoxProps,
  LinkOverlay,
  Stack,
  Text,
  VStack,
  chakra,
} from "@chakra-ui/react";
import { useUserListAccessTemplates } from "../utils/backend-client/default/default";
import { AccessTemplate } from "src/utils/backend-client/types";
import { StatusCell } from "./StatusCell";
import { ProviderIcon, ShortTypes } from "./icons/providerIcon";
import { access } from "fs";

// interface ListAccessTemplateProps {
//   test: string;
// }
export const AccessTemplateList: React.FC = () => {
  const { data } = useUserListAccessTemplates();

  if (data && data?.accessTemplates) {
    return (
      <Flex direction="column">
        <Text pb="20px" textStyle="Heading/H4">
          Access Templates
        </Text>

        <Grid templateColumns="repeat(2, 1fr)" gap={2}>
          {data.accessTemplates.map((template) => {
            return <AccessTemplateCard accessTemplate={template} />;
          })}
        </Grid>
      </Flex>
    );
  }
  return <></>;
};

const AccessTemplateCard: React.FC<
  {
    accessTemplate: AccessTemplate;
  } & LinkBoxProps
> = ({ accessTemplate, ...rest }) => {
  return (
    <LinkBox {...rest}>
      <Link /*to={"/requests/" + request.id}*/>
        <LinkOverlay>
          <Box
            rounded="lg"
            bg="neutrals.100"
            // columns={2}
            borderWidth={1}
            borderColor="neutrals.300"
            w="100%"
            h="100px"
          >
            <Flex px={1} align="center">
              <Text
                textStyle="Body/Small"
                color="neutrals.800"
                decoration="none"
              >
                {accessTemplate.name}
              </Text>
            </Flex>
          </Box>
        </LinkOverlay>
      </Link>
    </LinkBox>
  );
};
