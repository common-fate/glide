import {
  Avatar,
  Box,
  Button,
  Checkbox,
  CheckboxGroup,
  Flex,
  FormControl,
  FormLabel,
  Input,
  Stack,
  Tab,
  TabPanel,
  TabPanels,
  Tabs,
  Text,
  useBoolean,
} from "@chakra-ui/react";
import React, { useState } from "react";
import { CommonFateLogo } from "../../components/icons/Logos";
import { BlurShapesBox } from "../../components/icons/BlurShapes";
import { useForm } from "react-hook-form";
import { Helmet } from "react-helmet";
import StaticProviderSelect from "../../components/StaticProviderSelect";

const Index = () => {
  const [loading, setLoading] = useBoolean(false);

  const [tabIndex, setTabIndex] = useState(0);

  const { register, control, handleSubmit } = useForm<{
    name: string;
  }>();

  //   const router = useRouter();

  const onSubmit = handleSubmit((data) => {
    setLoading.on();

    //   createAccount(data)
    //       .then(async (r) => {
    //           submitState.setLoading(false);
    //           router.push(`/hm/${r.data.id}`);
    //           console.log(data);
    //       })
    //       .catch((err) => {
    //           submitState.setLoading(false);
    //           console.error(err);
    //       });
  });

  // const NO_OF_QS = Object.keys({ name: "test" }).length - 1;
  const NO_OF_QS = 20;

  const isFirstPage = tabIndex === 0;
  const isLastPage = tabIndex === NO_OF_QS;

  // hotkey listen for CMD + ENTER to submit

  return (
    <Box bg="neutrals.700" h="100vh" w="100vw">
      <Helmet>
        <title>Signup</title>
      </Helmet>
      <Flex>
        {/* LHS TABBED FORM INPUTS */}
        <Flex
          h="100vh"
          w={{ base: "100%", lg: "66%" }}
          borderRight="1px solid"
          borderRightColor="neutrals.600"
          justifyContent="center"
          alignItems="start"
        >
          <BlurShapesBox
            display={{ base: "block", lg: "none" }}
            pos="absolute"
            top={0}
            left={0}
            zIndex={0}
          />
          <Stack
            mt={12}
            rounded="md"
            p={8}
            w="387px"
            bg="neutrals.100"
            zIndex={1}
          >
            <CommonFateLogo h="24px" width="auto" mr="auto" />
            <Tabs index={tabIndex}>
              <TabPanels>
                <TabPanel px={0}>
                  <Stack spacing={4}>
                    <Box>
                      <Text textStyle="Body/Large" color="neutrals.800">
                        Choose a name for your Common Fate Account
                      </Text>
                      <Text textStyle="Body/Small" color="neutrals.500">
                        This is what will be displayed to your team
                      </Text>
                    </Box>
                    <Input type="text" placeholder="Acme Corp" />
                  </Stack>
                </TabPanel>
                <TabPanel px={0}>
                  <Stack spacing={4}>
                    <Box>
                      <Text textStyle="Body/Large" color="neutrals.800">
                        How did you hear about us?
                      </Text>
                      <Text textStyle="Body/Small" color="neutrals.500">
                        This is what will be displayed to your team
                      </Text>
                    </Box>
                    <CheckboxGroup
                      // colorScheme="green"
                      defaultValue={["naruto", "kakashi"]}
                    >
                      <Stack spacing={[1]} direction={{ base: "column" }}>
                        {[
                          "Google",
                          "Hacker News",
                          "Twitter",
                          "OSS Tools (like Granted CLI)",
                          "Other",
                        ].map((val) => (
                          <Checkbox value={val}>{val}</Checkbox>
                        ))}
                      </Stack>
                    </CheckboxGroup>
                    {/* field with textarea */}
                    <FormControl>
                      <FormLabel>Other</FormLabel>
                      <Input
                        as="textarea"
                        minH="24px"
                        placeholder="..."
                        rows={3}
                      />
                    </FormControl>
                  </Stack>
                </TabPanel>
                <TabPanel px={0}>
                  <Stack spacing={4}>
                    <Box>
                      <Text textStyle="Body/Large" color="neutrals.800">
                        What Cloud Providers/SaaS Applications do you use?
                      </Text>
                      <Text textStyle="Body/Small" color="neutrals.500">
                        Common Fate is actively building integrations
                      </Text>
                    </Box>
                    {/* add in provider select here */}
                    <StaticProviderSelect />
                    {/* field with textarea */}
                    <FormControl>
                      <FormLabel>Other</FormLabel>
                      <Input
                        as="textarea"
                        minH="24px"
                        placeholder="..."
                        rows={3}
                      />
                    </FormControl>
                  </Stack>
                </TabPanel>
              </TabPanels>
            </Tabs>
            <Flex>
              <Button
                variant="link"
                color="neutrals.500"
                disabled={isFirstPage}
                onClick={
                  isFirstPage ? () => null : () => setTabIndex((t) => t - 1)
                }
              >
                Back
              </Button>
              <Button
                ml="auto"
                isLoading={loading}
                onClick={() =>
                  isLastPage ? onSubmit() : setTabIndex((t) => t + 1)
                }
              >
                Continue
              </Button>
            </Flex>
          </Stack>
        </Flex>
        {/* RHS BLURRED SHAPES */}
        <Flex
          w="33%"
          h="100%"
          // blurred
          flexDir="column"
          alignItems="center"
          display={{ base: "none", lg: "block" }}
        >
          <BlurShapesBox pos="absolute" top={0} />
          {/* Quote Block: needs to be made into component */}
          <Box
            mx="auto"
            mt={12}
            zIndex={2}
            pos="relative"
            rounded="md"
            p={6}
            w="387px"
            bg="neutrals.100"
          >
            <Flex>
              <Avatar
                variant="withBorder"
                src={"https://avatars.githubusercontent.com/u/810438?v=4"}
                mr={4}
              />
              <Box>
                <Text textStyle="Body/Large" color="neutrals.800">
                  Dan Abramov
                </Text>
                <Text textStyle="Body/Small" color="neutrals.500">
                  Company
                </Text>
              </Box>
            </Flex>
            <Text mt={2} textStyle="Body/Small" color="neutrals.600">
              Managed to integrate AppSync with @ClerkDev over the weekend.
              Implementing a Lambda authorizer with Clerk's Go SDK was super
              easy.
            </Text>
          </Box>
        </Flex>
      </Flex>
    </Box>
  );
};

export default Index;
