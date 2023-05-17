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
import { ProviderIcon, ShortTypes } from "./icons/providerIcon";
import { access } from "fs";
import {
  useUserGetMe,
  userGetMe,
  userGetUser,
} from "../utils/backend-client/end-user/end-user";

interface ListAccessTemplateProps {
  setChecked: React.Dispatch<React.SetStateAction<Set<string>>>;
}
export const AccessTemplateList: React.FC<ListAccessTemplateProps> = ({
  setChecked,
}) => {
  const { data } = useUserListAccessTemplates();

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
                  _hover={{
                    bg: "neutrals.100",
                    rounded: "lg",
                    textDecoration: "none",
                  }}
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
      </Stack>
    );
  }
  return <></>;
};

const AccessTemplateCard: React.FC<
  {
    accessTemplate: AccessTemplate;
    handleClick: React.MouseEventHandler<HTMLAnchorElement>;
  } & LinkBoxProps
> = ({ accessTemplate, handleClick, ...rest }) => {
  // const user = useUserGetMe();

  // console.log(user);
  const CanUseTemplate = () => {
    // accessTemplate.accessGroups.map((g) => {

    //   if (user && user.data && !user?.data.user.groups.includes(g)) {
    //     console.log(g);
    //     return false;
    //   }
    // });

    return false;
  };
  return (
    <LinkBox {...rest}>
      <Link
        onClick={(e) =>
          CanUseTemplate() ? handleClick(e) : e.preventDefault()
        }
        textDecoration="none"
        _hover={{
          textDecoration: "none",
          cursor: CanUseTemplate() ? "pointer" : "default",
        }}
      >
        <LinkOverlay>
          <Box rounded="lg" w="100%" h="50px">
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
          </Box>
          {/* <Divider borderColor="neutrals.300" /> */}
        </LinkOverlay>
      </Link>
    </LinkBox>
  );
};
