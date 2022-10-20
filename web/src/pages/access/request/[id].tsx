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
  VStack,
  Wrap,
  WrapItem,
} from "@chakra-ui/react";
import { format } from "date-fns";
import React, { useEffect, useMemo, useState } from "react";
import { Controller, SubmitHandler, useForm } from "react-hook-form";
import { Link, useMatch, useNavigate } from "react-location";
import Select, { components, GroupBase, OptionProps } from "react-select";
import { CFRadioBox } from "../../../components/CFRadioBox";
import {
  DurationInput,
  Hours,
  Minutes,
} from "../../../components/DurationInput";
import { ProviderIcon } from "../../../components/icons/providerIcon";
import { ConnectorArrow } from "../../../components/ConnectorArrow";
import { ApprovalsLogo } from "../../../components/icons/Logos";
import { InfoOption } from "../../../components/InfoOption";
import { UserLayout } from "../../../components/Layout";
import { UserAvatarDetails } from "../../../components/UserAvatar";
import {
  getUserGetAccessRuleApproversKey,
  userCreateRequest,
  useUserGetAccessRule,
  useUserGetAccessRuleApprovers,
} from "../../../utils/backend-client/end-user/end-user";
import {
  AccessRuleTargetDetail,
  CreateRequestRequestBody,
  WithOption,
} from "../../../utils/backend-client/types";
import { durationString } from "../../../utils/durationString";
import { data } from "msw/lib/types/context";
import axios, { AxiosError } from "axios";
import { colors } from "../../../utils/theme/colors";
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
  console.log(rule);

  const {
    register,
    handleSubmit,
    formState: { errors },
    setValue,
    control,
    watch,
    reset,
    getValues,
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
      // The following will attempt to match any query params to withSelectable fields for this rule.
      // If the field matches and the value is a valid option, it will be set in the form values.
      // if it is not a valid value it is ignored.
      // this prevents being able to submit the form with bad options, or being able to submit arbitrary values for the with fields via the UI
      Object.entries(rule.target.withSelectable).map(([k, v]) => {
        const queryParamValue = new URLSearchParams(
          location.search.substring(1)
        ).get(k);
        if (
          queryParamValue !== null &&
          v.options.find((s) => {
            return s.value === queryParamValue;
          }) !== undefined
        ) {
          setValue(`with.${k}`, queryParamValue);
        }
      });
    }
  }, [rule, location.search]);

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
      reason: data.reason ? data.reason : "",
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
        let description: string | undefined;
        if (axios.isAxiosError(e)) {
          description = (e as AxiosError<{ error: string }>)?.response?.data
            .error;
        }
        toast({
          title: "Request failed",
          status: "error",
          duration: 5000,
          description: (
            <Text color={"white"} whiteSpace={"pre"}>
              {description}
            </Text>
          ),
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
                  <AccessRuleWithDisplay rule={rule.target} />
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
                              color="neutrals.600"
                              fontWeight="normal"
                            >
                              {v.title}
                            </FormLabel>

                            <Controller
                              name={`with.${k}`}
                              control={control}
                              rules={{ required: true }}
                              render={({
                                field: { value, onChange, ...rest },
                              }) => (
                                <>
                                  <Select
                                    components={{
                                      Option: CustomOption,
                                    }}
                                    styles={{
                                      option: (provided, state) => {
                                        return {
                                          ...provided,
                                          background: state.isSelected
                                            ? colors.blue[200]
                                            : provided.background,
                                          color: state.isSelected
                                            ? colors.neutrals[800]
                                            : provided.color,
                                        };
                                      },
                                    }}
                                    isMulti={false}
                                    options={v.options
                                      // exclude invalid options
                                      .filter((op) => op.valid)
                                      .map((op) => {
                                        return op;
                                      })
                                      .sort((a, b) => {
                                        return a.label < b.label
                                          ? -1
                                          : a.label === b.label
                                          ? 0
                                          : 1;
                                      })}
                                    value={v.options.find(
                                      (op) => value === op.value
                                    )}
                                    onChange={(val) => {
                                      onChange(val?.value);
                                    }}
                                    {...rest}
                                  />
                                  <Text
                                    textStyle={"Body/Small"}
                                    color="neutrals.600"
                                  >
                                    {value}
                                  </Text>
                                </>
                              )}
                            />

                            <FormErrorMessage>
                              This field is required
                            </FormErrorMessage>
                          </FormControl>
                        );
                      }
                    )}
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

                <FormControl isInvalid={!!errors?.reason}>
                  <FormLabel textStyle="Body/Medium" fontWeight="normal">
                    Why do you need access?
                  </FormLabel>
                  <Textarea
                    bg="white"
                    id="reasonField"
                    placeholder="Deploying initial Terraform infrastructure for CF-123"
                    {...register("reason", {
                      validate: (value) => {
                        const res: string[] = [];
                        [
                          /[^a-zA-Z0-9,.;:()[\]?!\-_`~&/\n\s]/,
                        ].every((pattern) => pattern.test(value as string)) &&
                          res.push(
                            "Invalid characters (only letters, numbers, and punctuation allowed)"
                          );
                        if (value && value.length > 2048) {
                          res.push("Maximum length is 2048 characters");
                        }
                        return res.length > 0 ? res.join(", ") : undefined;
                      },
                    })}
                  />
                  {errors?.reason && (
                    <FormErrorMessage>
                      {errors?.reason.message}
                      {JSON.stringify(errors?.reason.types)}
                    </FormErrorMessage>
                  )}
                </FormControl>

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

export const AccessRuleWithDisplay: React.FC<{
  rule?: AccessRuleTargetDetail;
}> = ({ rule }) => {
  if (rule === undefined) {
    return <Skeleton minW="30ch" minH="6" mr="auto" />;
  }
  if (Object.entries(rule.with).length > 0) {
    return (
      <Wrap>
        {rule.with &&
          Object.entries(rule.with).map(([k, v]) => {
            return (
              <WrapItem>
                <VStack align={"left"}>
                  <Text>{v.title}</Text>
                  <InfoOption label={v.label} value={v.value} />
                </VStack>
              </WrapItem>
            );
          })}
      </Wrap>
    );
  }
  return null;
};
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
const CustomOption = ({
  children,
  ...innerProps
}: OptionProps<WithOption, false, GroupBase<WithOption>>) => (
  <div data-testid={innerProps.data.value}>
    <components.Option {...innerProps}>
      <>
        {children}
        {<Text>{innerProps.data.value}</Text>}
      </>
    </components.Option>
  </div>
);
export default Home;
