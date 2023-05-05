import { ArrowBackIcon } from "@chakra-ui/icons";
import {
  Box,
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  Button,
  Center,
  Container,
  Flex,
  FormControl,
  FormErrorMessage,
  FormLabel,
  IconButton,
  Input,
  Spacer,
  Stack,
  Text,
  Textarea,
  useToast,
} from "@chakra-ui/react";
import axios from "axios";
import { Helmet } from "react-helmet";
import { Controller, FormProvider, useForm } from "react-hook-form";
import { Link, useMatch, useNavigate } from "react-location";
import { ProviderIcon, ShortTypes } from "../../components/icons/providerIcon";

import { UserLayout } from "../../components/Layout";
import { TargetDetail } from "../../components/Target";
import {
  useUserGetPreflight,
  userPostRequests,
} from "../../utils/backend-client/default/default";
import { CreateAccessRequestRequestBody } from "../../utils/backend-client/types";

const Home = () => {
  const {
    params: { id: preflightId },
  } = useMatch();
  const navigate = useNavigate();
  const toast = useToast();
  const { data: preflight } = useUserGetPreflight(preflightId);

  const methods = useForm<CreateAccessRequestRequestBody>({
    defaultValues: {
      preflightId: preflightId,
    },
  });
  const onSubmit = async (data: CreateAccessRequestRequestBody) => {
    console.debug("submit form data", { data });

    try {
      const request = await userPostRequests(data);
      navigate({ to: `/requests/${request.id}` });
    } catch (err) {
      let description: string | undefined;
      if (axios.isAxiosError(err)) {
        // @ts-ignore
        description = err?.response?.data.error;
      }
      toast({
        title: "Error submitting request",
        description,
        status: "error",
        variant: "subtle",
        duration: 2200,
        isClosable: true,
      });
    }
  };

  return (
    <div>
      <UserLayout>
        <Helmet>
          <title>Preflight</title>
        </Helmet>
        <FormProvider {...methods}>
          <form onSubmit={methods.handleSubmit(onSubmit)}>
            <Center pt={{ base: 12, lg: 32 }}>
              <Stack
                spacing={7}
                background="neutrals.100"
                px={"200px"}
                pt={"100px"}
                pb={"150px"}
              >
                <Breadcrumb>
                  <BreadcrumbItem>
                    <BreadcrumbLink href="/requests">
                      New Request
                    </BreadcrumbLink>
                  </BreadcrumbItem>
                  <BreadcrumbItem>
                    <BreadcrumbLink href="#">Review</BreadcrumbLink>
                  </BreadcrumbItem>
                </Breadcrumb>
                <Stack spacing={2} w="100%">
                  {preflight?.accessGroups.map((group, i) => {
                    return (
                      <FormControl
                        isInvalid={
                          !!methods.formState.errors?.groupOptions?.[i]
                        }
                      >
                        <Controller
                          control={methods.control}
                          name={`groupOptions.${i}`}
                          // sets the field up with default value of max request duration
                          defaultValue={{
                            id: group.id,
                            timing: {
                              durationSeconds:
                                group.timeConstraints.maxDurationSeconds,
                            },
                          }}
                          render={({
                            field: { onChange, ref, value, onBlur },
                          }) => {
                            return (
                              <Box
                                p={2}
                                w="100%"
                                borderColor="neutrals.300"
                                borderWidth="1px"
                                rounded="lg"
                              >
                                {/* <HeaderStatusCell group={group} /> */}
                                <Stack spacing={2}>
                                  <Flex>
                                    <Text>
                                      {group.requiresApproval
                                        ? "Requires Approval"
                                        : "No Approval Required"}
                                    </Text>
                                    <Spacer />
                                  </Flex>
                                  {group.targets.map((target) => {
                                    return (
                                      <Flex
                                        p={2}
                                        borderColor="neutrals.300"
                                        borderWidth="1px"
                                        rounded="lg"
                                        flexDir="row"
                                        background="white"
                                      >
                                        <TargetDetail
                                          showIcon
                                          target={target}
                                        />
                                      </Flex>
                                    );
                                  })}
                                </Stack>
                                {methods.formState.errors?.reason?.message && (
                                  <FormErrorMessage>
                                    {methods.formState.errors.reason?.message?.toString()}
                                  </FormErrorMessage>
                                )}
                              </Box>
                            );
                          }}
                        />
                      </FormControl>
                    );
                  })}
                </Stack>

                <FormControl isInvalid={!!methods.formState.errors.reason}>
                  <FormLabel htmlFor="reason">
                    <Text textStyle={"Body/Medium"}>
                      Why do you need access?
                    </Text>
                  </FormLabel>
                  <Textarea
                    placeholder="Deploying initial Terraform infrastructure..."
                    bg="neutrals.0"
                    {...methods.register("reason")}
                    onBlur={() => void methods.trigger("reason")}
                  />
                  {methods.formState.errors?.reason?.message && (
                    <FormErrorMessage>
                      {methods.formState.errors.reason?.message?.toString()}
                    </FormErrorMessage>
                  )}
                </FormControl>
                {/* buttons */}
                <Flex w="100%" mt={4}>
                  <Button
                    type="button"
                    isDisabled={methods.formState.isSubmitting}
                    variant="brandSecondary"
                    leftIcon={<ArrowBackIcon />}
                    to="/requests"
                    as={Link}
                  >
                    Go back
                  </Button>
                  <Button
                    type="submit"
                    ml="auto"
                    isDisabled={
                      methods.formState.isValidating ||
                      !methods.formState.isValid
                    }
                    isLoading={methods.formState.isSubmitting}
                    loadingText="Processing request..."
                  >
                    Next (⌘+Enter)
                  </Button>
                </Flex>
              </Stack>
            </Center>
          </form>
        </FormProvider>
      </UserLayout>
    </div>
  );
};

export default Home;
