import {
  ArrowBackIcon,
  CheckCircleIcon,
  ChevronDownIcon,
  ChevronRightIcon,
  InfoIcon,
  WarningIcon,
} from "@chakra-ui/icons";
import {
  Box,
  Button,
  Center,
  Collapse,
  Container,
  Flex,
  FormControl,
  FormErrorMessage,
  FormHelperText,
  FormLabel,
  HStack,
  IconButton,
  Input,
  Skeleton,
  SkeletonCircle,
  SkeletonText,
  Spinner,
  Stack,
  Text,
  Textarea,
  useRadioGroup,
  UseRadioGroupProps,
  useToast,
  Wrap,
} from "@chakra-ui/react";
import { format } from "date-fns";
import React, { useEffect, useMemo, useState } from "react";
import { Controller, SubmitHandler, useForm } from "react-hook-form";
import { Link, useMatch, useNavigate } from "react-location";
import Select from "react-select";
import { CFRadioBox } from "../../../components/CFRadioBox";
import {
  DurationInput,
  Hours,
  Minutes,
} from "../../../components/DurationInput";
import { ProviderIcon } from "../../../components/icons/providerIcon";
import { ConnectorArrow } from "../../../components/ConnectorArrow";
import { ApprovalsLogo } from "../../../components/icons/Logos";
import { UserLayout } from "../../../components/Layout";
import { UserAvatarDetails } from "../../../components/UserAvatar";
import {
  getUserGetAccessRuleApproversKey,
  userCreateRequest,
  useUserGetAccessRule,
  useUserGetAccessRuleApprovers,
} from "../../../utils/backend-client/end-user/end-user";
import { CreateRequestRequestBody } from "../../../utils/backend-client/types";
import { durationString } from "../../../utils/durationString";
import { data } from "msw/lib/types/context";
export type When = "asap" | "scheduled";

interface NewRequestFormData extends CreateRequestRequestBody {
  startDateTime: string;
  when: When;
}

interface FieldError {
  error: string;
  field: string;
}

/**
 * returns helper text to be used below form fields for selecting when
 * access should be activated.
 */
export const getWhenHelperText = (
  when: When,
  requiresApproval: boolean
): string => {
  if (when === "asap" && requiresApproval)
    return "Access will be activated immediately after approval";
  if (when === "asap") return "Access will be activated immediately";

  return "Choose a time in future for the access to be activated";
};

const Home = () => {
  const [loading, setLoading] = useState(false);
  const {
    params: { id: ruleId },
  } = useMatch();
  const { data: rule } = useUserGetAccessRule(ruleId);
  const navigate = useNavigate();
  const now = useMemo(() => {
    const d = new Date();
    d.setSeconds(0, 0);
    return format(d, "yyyy-MM-dd'T'HH:mm");
  }, []);

  const {
    register,
    handleSubmit,
    formState: { errors },
    setValue,
    control,
    watch,
    reset,
  } = useForm<NewRequestFormData>({
    shouldUnregister: true,
    defaultValues: {
      when: "asap",
      startDateTime: now,
      timing: {
        durationSeconds: 60,
      },
    },
  });
  const toast = useToast();

  const [validationErrors, setValidationErrors] = useState<FieldError[]>();

  // This use effect sets the duration to either 1 hour or max duration if it is less than one hour
  // it does then when the rule loads for the first time
  useEffect(() => {
    if (rule != undefined) {
      setValue(
        "timing.durationSeconds",
        rule.timeConstraints.maxDurationSeconds > 3600
          ? 3600
          : rule.timeConstraints.maxDurationSeconds
      );
    }
  }, [rule]);

  const when = watch("when");
  const startTimeDate = watch("startDateTime");
  // Don't refetch the approvers
  const {
    data: approvers,
    isValidating: isValidatingApprovers,
  } = useUserGetAccessRuleApprovers(ruleId, {
    swr: {
      swrKey: getUserGetAccessRuleApproversKey(ruleId),
      refreshInterval: 0,
      revalidateOnFocus: false,
    },
  });
  const requiresApproval = !!approvers && approvers.users.length > 0;

  const onSubmit: SubmitHandler<NewRequestFormData> = async (data) => {
    setLoading(true);

    const r: CreateRequestRequestBody = {
      accessRuleId: ruleId,
      timing: {
        durationSeconds: data.timing.durationSeconds,
      },
      reason: data.reason,
      with: data.with,
    };
    if (data.when === "scheduled") {
      r.timing.startTime = new Date(data.startDateTime).toISOString();
    }
    await userCreateRequest(r)
      .then(() => {
        toast({
          title: "Request created",
          status: "success",
          duration: 2200,
          isClosable: true,
        });
        navigate({ to: "/requests" });
      })
      .catch((e) => {
        setLoading(false);

        setValidationErrors(e.response.data.fields);
        toast({
          title: "Request failed",
          status: "error",
          duration: 2200,
          description: e.message,
          isClosable: true,
        });
      });
  };

  return (
    <>
      <UserLayout>
        <Center borderBottom="1px solid" borderColor="neutrals.200" h="80px">
          <IconButton
            as={Link}
            to="/requests"
            aria-label="Go back"
            pos="absolute"
            left={4}
            icon={<ArrowBackIcon />}
            rounded="full"
            variant="ghost"
          />

          <Text as="h4" textStyle="Heading/H4">
            New Access Request
          </Text>
        </Center>
        <Container minW="864px">
          <Box
            p={8}
            bg="neutrals.100"
            mt={12}
            borderRadius="6px"
            as="form"
            onSubmit={handleSubmit(onSubmit)}
          >
            <Text as="h3" textStyle="Heading/H3">
              You are requesting access to
            </Text>

            <Stack
              spacing={2}
              mt={6}
              minH="52px" // prevents layout shift
            >
              {rule ? (
                <>
                  <Flex data-testid="rule-name" align="center" mr="auto">
                    <ProviderIcon shortType={rule?.target.provider.type} />
                    <Text ml={2} textStyle="Body/Medium" color="neutrals.600">
                      {rule?.name}
                    </Text>
                  </Flex>
                  <Text textStyle="Body/Medium">{rule?.description}</Text>
                </>
              ) : (
                <>
                  <Flex align="center">
                    <SkeletonCircle h={8} w={8} mr={2} />
                    <SkeletonText w="14ch" noOfLines={1} />
                  </Flex>
                  <SkeletonText w="10ch" noOfLines={1} />
                </>
              )}
            </Stack>

            <Box mt={12}>
              <Stack spacing={10}>
                {rule &&
                  Object.entries(rule.target.withSelectable).map(
                    ([k, v], i) => {
                      const name = "with." + k;
                      return (
                        <FormControl
                          key={"selectable-" + k}
                          pos="relative"
                          id={name}
                          isInvalid={
                            errors.with && errors.with[k] !== undefined
                          }
                        >
                          <FormLabel
                            textStyle="Body/Medium"
                            fontWeight="normal"
                          >
                            {k}
                          </FormLabel>

                          <Controller
                            name={`with.${k}`}
                            control={control}
                            rules={{ required: true }}
                            render={({
                              field: { value, onChange, ...rest },
                            }) => (
                              <Select
                                isMulti={false}
                                options={v
                                  // exclude invalid options
                                  .filter((op) => op.valid)
                                  .map((op) => {
                                    return op.option;
                                  })}
                                value={
                                  v.find((op) => value === op.option.value)
                                    ?.option
                                }
                                onChange={(val) => {
                                  onChange(val?.value);
                                }}
                                {...rest}
                              />
                            )}
                          />
                          <FormErrorMessage>
                            This field is required
                          </FormErrorMessage>
                        </FormControl>
                      );
                    }
                  )}

                <FormControl
                  pos="relative"
                  id="when"
                  isInvalid={errors.when !== undefined}
                >
                  <FormLabel textStyle="Body/Medium" fontWeight="normal">
                    When do you need access?
                  </FormLabel>

                  <Controller
                    name="when"
                    control={control}
                    render={({ field }) => <WhenRadioGroup {...field} />}
                  />
                  <FormHelperText color="neutrals.600" minH="17px">
                    {isValidatingApprovers ? (
                      <SkeletonText w="24ch" noOfLines={1} />
                    ) : (
                      getWhenHelperText(when, requiresApproval)
                    )}
                  </FormHelperText>
                </FormControl>

                {/* use a Flex here to avoid the Collapse animation jumping due to being nested within a <Stack /> */}
                <Flex direction={"column"}>
                  <Collapse in={when === "scheduled"} animateOpacity>
                    <FormControl mb={10}>
                      <FormLabel textStyle="Body/Medium" fontWeight="normal">
                        Start Time
                      </FormLabel>

                      <Input
                        {...register("startDateTime")}
                        bg="white"
                        type="datetime-local"
                        min={now}
                        defaultValue={now}
                      />

                      {startTimeDate && (
                        <FormHelperText color="neutrals.600">
                          {new Date(startTimeDate).toString()}
                        </FormHelperText>
                      )}
                    </FormControl>
                  </Collapse>

                  <FormControl
                    pos="relative"
                    isInvalid={errors.timing?.durationSeconds !== undefined}
                  >
                    <FormLabel textStyle="Body/Medium" fontWeight="normal">
                      How long do you need access for?
                    </FormLabel>

                    <Controller
                      name="timing.durationSeconds"
                      control={control}
                      rules={{
                        required: "Duration is required.",
                        max: rule?.timeConstraints.maxDurationSeconds,
                        min: 60,
                      }}
                      render={({ field: { ref, ...rest } }) => {
                        return (
                          <DurationInput
                            {...rest}
                            max={rule?.timeConstraints.maxDurationSeconds}
                            min={60}
                          >
                            <Hours />
                            <Minutes />
                            {
                              <Text textStyle={"Body/ExtraSmall"}>
                                Max{" "}
                                {durationString(
                                  rule?.timeConstraints.maxDurationSeconds
                                )}
                                <br />
                                Min 1 minute
                              </Text>
                            }
                          </DurationInput>
                        );
                      }}
                    />

                    {errors.timing?.durationSeconds !== undefined && (
                      <FormErrorMessage>
                        {errors.timing?.durationSeconds.message}
                      </FormErrorMessage>
                    )}
                  </FormControl>
                </Flex>

                <FormControl>
                  <FormLabel textStyle="Body/Medium" fontWeight="normal">
                    Why do you need access?
                  </FormLabel>
                  <Textarea
                    bg="white"
                    id="reasonField"
                    placeholder="Deploying initial Terraform infrastructure for CF-123"
                    {...register("reason")}
                  />
                </FormControl>

                {validationErrors && (
                  <Stack
                    borderWidth={"1px"}
                    justifyContent="center"
                    alignItems={"center"}
                    borderRadius="8px"
                  >
                    <Stack
                      justifyContent={"center"}
                      alignItems={"center"}
                      spacing={{ base: 5, md: 0 }}
                      position="relative"
                      w="100%"
                      p={4}
                      flexDirection={{ base: "column", md: "row" }}
                    >
                      <Text textStyle="Body/Medium">Grant Validation Test</Text>

                      <Flex
                        position={{ md: "absolute", base: "relative" }}
                        right={{ md: 3, base: 0 }}
                        w={{ base: "100%", md: "unset" }}
                        h="100%"
                        justifyContent={"center"}
                        alignItems={"center"}
                      >
                        <HStack
                          spacing={5}
                          justifyContent={"center"}
                          alignItems={"center"}
                        ></HStack>
                      </Flex>
                    </Stack>
                    <Stack
                      spacing={4}
                      bg="neutrals.700"
                      pt={3}
                      w="100%"
                      borderBottomLeftRadius={"8px"}
                      borderBottomRightRadius={"8px"}
                      shadow="inset 0 0 10px #0f0f0f"
                    >
                      {validationErrors.map((validation) => (
                        <ValidationResults
                          key={validation.field}
                          loading={!validationErrors}
                          validation={validation}
                        />
                      ))}
                      <Flex
                        borderTopWidth="1px"
                        borderTopColor={"#383a3c"}
                        w="100%"
                        justifyContent="flex-end"
                        px={3}
                      >
                        {/* <Button
                          onClick={() => onCopy()}
                          borderLeftWidth="1px"
                          borderLeftColor={"#383a3c"}
                          borderLeftRadius={0}
                          pl={3}
                          variant="unstyled"
                          color={hasCopied ? "#22c55e" : "#d0d7de"}
                          textTransform={"uppercase"}
                          fontSize="11px"
                          size="xs"
                        >
                          Copy Diagnostics
                        </Button> */}
                      </Flex>
                    </Stack>
                  </Stack>
                  // <Flex direction="column">
                  //   <Text textStyle="Body/Medium" fontWeight="normal">
                  //     Request failed due to the following errors:
                  //   </Text>
                  //   {validationErrors.map((item) => {
                  //     return (
                  //       <>
                  //         <Text textStyle="Body/Medium">
                  //           Validation: {item.field}
                  //         </Text>
                  //         <Text textStyle={"Body/ExtraSmall"}>
                  //           Error: {item.error}
                  //         </Text>
                  //       </>
                  //     );
                  //   })}
                  // </Flex>
                )}
                {/* Don't show approval section if approvers are still loading */}
                <Approvers approvers={approvers?.users} />
                <Box>
                  <Button type="submit" isLoading={loading} mr={3}>
                    Submit
                  </Button>
                </Box>
              </Stack>
            </Box>
          </Box>
        </Container>
      </UserLayout>
    </>
  );
};

interface ValidationResultsProps {
  loading?: boolean;
  validation: FieldError;
}

const ValidationResults: React.FC<ValidationResultsProps> = ({
  loading,
  validation,
}) => {
  const [expanded, setExpanded] = useState(true);

  return (
    <Stack spacing={0}>
      <Flex color="#d0d7de" alignItems={"center"} py={1} px={3}>
        <Flex w={6} alignItems="center">
          <IconButton
            color="neutrals.500"
            size="s"
            variant={"unstyled"}
            aria-label="expand"
            onClick={() => setExpanded(!expanded)}
            icon={expanded ? <ChevronDownIcon /> : <ChevronRightIcon />}
          />
        </Flex>
        <Flex w={6} alignItems="center">
          {loading ? (
            <Spinner size="xs" color="neutrals.500" />
          ) : (
            <WarningIcon boxSize={"12px"} color="neutrals.500" />
          )}
        </Flex>
        <Text color="#d0d7de">{validation.field}</Text>
      </Flex>
      {expanded && (
        <Stack pl={"60px"} spacing={1}>
          <Text
            key={validation.field}
            color="#d0d7de"
            fontSize={"12px"}
            fontFamily="mono"
          >
            {validation.error}
          </Text>
        </Stack>
      )}
    </Stack>
  );
};

export const WhenRadioGroup: React.FC<UseRadioGroupProps> = (props) => {
  const { getRootProps, getRadioProps } = useRadioGroup(props);
  const group = getRootProps();

  return (
    <HStack {...group}>
      <CFRadioBox {...getRadioProps({ value: "asap" })}>
        <Text textStyle="Body/Medium">ASAP</Text>
      </CFRadioBox>
      <CFRadioBox {...getRadioProps({ value: "scheduled" })}>
        <Text textStyle="Body/Medium">Scheduled</Text>
      </CFRadioBox>
    </HStack>
  );
};

export default Home;

const Approvers: React.FC<{ approvers?: string[] }> = ({ approvers }) => {
  if (approvers === undefined) {
    return <Skeleton w="50%" h={10} />;
  }
  if (approvers.length > 0) {
    return (
      <Box textStyle="Body/Medium" maxW="470px">
        Approvers
        <Wrap spacing={2}>
          {approvers?.map((approver) => (
            // Using style props, we're able to more closely match the figma designs
            <UserAvatarDetails
              key={approver}
              user={approver}
              size="xs"
              textProps={{
                textStyle: "Body/Small",
                color: "neutrals.500",
              }}
            />
          ))}
        </Wrap>
      </Box>
    );
  }
  return (
    <Text color="neutrals.600" display="flex" alignItems="center">
      <InfoIcon mr={2} />
      Approval is not required for this role, so you&apos;ll get access
      immediately
    </Text>
  );
};
