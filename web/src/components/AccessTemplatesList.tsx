import {
  Box,
  Button,
  Divider,
  Flex,
  Grid,
  HStack,
  Link,
  LinkBox,
  LinkBoxProps,
  LinkOverlay,
  Popover,
  PopoverArrow,
  PopoverBody,
  PopoverCloseButton,
  PopoverContent,
  PopoverHeader,
  PopoverTrigger,
  Stack,
  Text,
  Tooltip,
  VStack,
  chakra,
} from "@chakra-ui/react";
import { useUserListAccessTemplates } from "../utils/backend-client/default/default";
import {
  AccessTemplate,
  AuthUserResponseResponse,
  User,
} from "src/utils/backend-client/types";
import { ProviderIcon, ShortTypes } from "./icons/providerIcon";
import { access } from "fs";
import {
  useUserGetMe,
  userGetMe,
  userGetUser,
} from "../utils/backend-client/end-user/end-user";
import { useUser } from "../utils/context/userContext";

interface ListAccessTemplateProps {
  setChecked: React.Dispatch<React.SetStateAction<Set<string>>>;
}
export const AccessTemplateList: React.FC<ListAccessTemplateProps> = ({
  setChecked,
}) => {
  const { data } = useUserListAccessTemplates();

  if (!data || data.accessTemplates.length === 0) {
    return (
      <Stack>
        <Flex
          p={1}
          rounded="lg"
          bg="white"
          // columns={2}
          borderWidth={1}
          borderColor="neutrals.300"
          direction="column"
          w="350px"
          h="70vh"
        >
          <Text as="h4" textStyle="Heading/H4" my="10px" pl="5px">
            Access Templates
          </Text>
          <Text>No Access Templates made yet...</Text>
        </Flex>
      </Stack>
    );
  }

  if (data && data?.accessTemplates) {
    return (
      <Stack>
        <Flex
          p={1}
          rounded="lg"
          bg="white"
          // columns={2}
          borderWidth={1}
          borderColor="neutrals.300"
          direction="column"
          w="350px"
          h="70vh"
        >
          <Text as="h4" textStyle="Heading/H4" my="10px" pl="5px">
            Access Templates
          </Text>
          <Grid templateColumns="repeat(1, 1fr)" gap={2}>
            {data.accessTemplates.map((template) => {
              return (
                <AccessTemplateCard
                  accessTemplate={template}
                  handleClick={() => {
                    const set = new Set();

                    template.accessGroups.forEach((group) => {
                      group.targets.forEach((target) => {
                        set.add(target.id.toLowerCase());
                      });
                    });
                    //@ts-ignore
                    setChecked(() => set);
                  }}
                />
              );
            })}
          </Grid>
        </Flex>
      </Stack>
    );
  }
  return <></>;
};

const CanUseTemplate = (user: User, accessTemplate: AccessTemplate) => {
  const mergedArr = user.groups.concat(accessTemplate.groupsWithAccess); // Merge the arrays into a single array

  // Create an object to store the occurrence count of each item
  //@ts-ignore
  const count: Record<T, number> = {};

  // Iterate over the merged array and count the occurrences of each item
  for (let i = 0; i < mergedArr.length; i++) {
    const item = mergedArr[i];
    count[item] = (count[item] || 0) + 1;
  }

  // Check if any item has occurred more than once
  for (const item in count) {
    if (count[item] > 1) {
      return true; // Found an overlapping item
    }
  }

  return false; // No overlapping item found
};

const AccessTemplateCard: React.FC<
  {
    accessTemplate: AccessTemplate;
    handleClick: React.MouseEventHandler<HTMLDivElement>;
  } & LinkBoxProps
> = ({ accessTemplate, handleClick, ...rest }) => {
  const user = useUser();

  if (!user || !user.user) {
    return <Text>Loading...</Text>;
  }

  const canUse = CanUseTemplate(user.user, accessTemplate);

  return (
    <>
      {!canUse && (
        <Box
          {...rest}
          _hover={{
            cursor: "default",
          }}
          rounded="lg"
          w="100%"
          h="50px"
        >
          <Tooltip
            hasArrow
            label={
              "You do not have access to all the resources in this Access Template"
            }
          >
            <Flex px={3} py={2}>
              <HStack>
                <Text
                  textStyle="Body/medium"
                  color="neutrals.400"
                  decoration="none"
                >
                  {accessTemplate.name}
                </Text>
                {accessTemplate.accessGroups.map((group) => {
                  return (
                    <ProviderIcon
                      h="18px"
                      w="18px"
                      shortType={group.targets[0].kind.name as ShortTypes}
                      mr={2}
                      color="neutrals.400"
                    />
                  );
                })}
              </HStack>
            </Flex>
          </Tooltip>
        </Box>
      )}
      {canUse && (
        <Box
          {...rest}
          _hover={{
            bg: "neutrals.100",
            rounded: "lg",
            textDecoration: "none",
            cursor: "pointer",
          }}
          rounded="lg"
          w="100%"
          h="50px"
          onClick={(e) => handleClick(e)}
        >
          <Tooltip hasArrow label={accessTemplate.description}>
            <Flex px={3} py={2}>
              <HStack>
                <Text
                  textStyle="Body/medium"
                  color="neutrals.700"
                  decoration="none"
                >
                  {accessTemplate.name}
                </Text>
                {accessTemplate.accessGroups.map((group) => {
                  return (
                    <ProviderIcon
                      h="18px"
                      w="18px"
                      shortType={group.targets[0].kind.name as ShortTypes}
                      mr={2}
                    />
                  );
                })}
              </HStack>
            </Flex>
          </Tooltip>
        </Box>
      )}
    </>
  );
};
