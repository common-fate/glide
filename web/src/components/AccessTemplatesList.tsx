import {
  Box,
  Divider,
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

interface ListAccessTemplateProps {
  setChecked: React.Dispatch<React.SetStateAction<Set<string>>>;
}
export const AccessTemplateList: React.FC<ListAccessTemplateProps> = ({
  setChecked,
}) => {
  const { data } = useUserListAccessTemplates();

  if (data && data?.accessTemplates) {
    return (
      <Flex
        p={1}
        rounded="lg"
        bg="white"
        // columns={2}
        borderWidth={1}
        borderColor="neutrals.300"
        direction="column"
        w="400px"
        h="70vh"
      >
        <Text pb="20px" textStyle="Heading/H4">
          Access Templates
        </Text>

        <Grid templateColumns="repeat(1, 1fr)" gap={2}>
          {data.accessTemplates.map((template) => {
            return (
              <AccessTemplateCard
                _hover={{ bg: "neutrals.100" }}
                accessTemplate={template}
                handleClick={() => {
                  template.accessGroups.forEach((group) => {
                    group.targets.forEach((target) => {
                      setChecked((old) => {
                        const newSet = new Set(old);
                        newSet.add(target.id.toLowerCase());
                        return newSet;
                      });
                    });
                  });
                }}
              />
            );
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
    handleClick: React.MouseEventHandler<HTMLAnchorElement> | undefined;
  } & LinkBoxProps
> = ({ accessTemplate, handleClick, ...rest }) => {
  return (
    <LinkBox {...rest}>
      <Link onClick={handleClick}>
        <LinkOverlay>
          <Box rounded="lg" w="100%" h="50px">
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
          {/* <Divider borderColor="neutrals.300" /> */}
        </LinkOverlay>
      </Link>
    </LinkBox>
  );
};
