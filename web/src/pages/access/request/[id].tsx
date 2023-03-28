import {
  ArrowBackIcon,
  CheckIcon,
  DeleteIcon,
  InfoIcon,
  LinkIcon,
  SmallAddIcon,
  StarIcon,
} from "@chakra-ui/icons";
import {
  Box,
  Button,
  ButtonGroup,
  Center,
  Collapse,
  Container,
  Flex,
  FormControl,
  FormErrorMessage,
  FormHelperText,
  FormLabel,
  Heading,
  HStack,
  IconButton,
  Input,
  Popover,
  PopoverArrow,
  PopoverBody,
  PopoverCloseButton,
  PopoverContent,
  PopoverFooter,
  PopoverHeader,
  PopoverTrigger,
  Skeleton,
  SkeletonCircle,
  SkeletonText,
  Stack,
  Text,
  Textarea,
  Tooltip,
  useClipboard,
  useDisclosure,
  useRadioGroup,
  UseRadioGroupProps,
  useToast,
  VStack,
  Wrap,
  WrapItem,
} from "@chakra-ui/react";
import axios, { AxiosError } from "axios";
import { format } from "date-fns";
import React, { useEffect, useMemo, useState } from "react";
import { Helmet } from "react-helmet";
import {
  Controller,
  FormProvider,
  SubmitHandler,
  useForm,
  useFormContext,
} from "react-hook-form";
import {
  Link,
  MakeGenerics,
  useLocation,
  useMatch,
  useNavigate,
  useSearch,
} from "react-location";
import { CFRadioBox } from "../../../components/CFRadioBox";
import {
  Days,
  DurationInput,
  Hours,
  Minutes,
  Weeks,
} from "../../../components/DurationInput";
import {
  MultiSelect,
  SelectWithArrayAsValue,
} from "../../../components/forms/access-rule/components/Select";
import { ProviderIcon } from "../../../components/icons/providerIcon";
import { InfoOption } from "../../../components/InfoOption";
import { UserLayout } from "../../../components/Layout";
import {
  getUserGetAccessRuleApproversKey,
  userCreateFavorite,
  userCreateRequest,
  userDeleteFavorite,
  userGetFavorite,
  userUpdateFavorite,
  useUserGetAccessRule,
  useUserGetAccessRuleApprovers,
} from "../../../utils/backend-client/end-user/end-user";
import {
  CreateFavoriteRequestBody,
  CreateRequestRequestBody,
  CreateRequestWith,
  CreateRequestWithSubRequest,
  FavoriteDetail,
  RequestAccessRule,
  RequestAccessRuleTarget,
  RequestArgumentFormElement,
  RequestTiming,
} from "../../../utils/backend-client/types";
import { durationString } from "../../../utils/durationString";
import { colors } from "../../../utils/theme/colors";
import { CFAvatar } from "../../../components/CFAvatar";
export type When = "asap" | "scheduled";

/**
 * The reason I added this type was because I was having trouble being able to remove an array element in the context of the form.
 * Instead, elements are marked as hidden when the remove button is pressed.
 * So when processing the form values, be sure to filter out the hidden elements  first.
 */
interface FormCreateRequestWith {
  hidden?: boolean;
  data: CreateRequestWith;
}
interface NewRequestFormData extends Omit<CreateRequestRequestBody, "with"> {
  startDateTime: string;
  when: When;
  with: FormCreateRequestWith[];
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

type RequestFormFields = {
  with?: CreateRequestWithSubRequest;
  timing?: RequestTiming;
  reason?: string;
};

export type RequestFormQueryParameters = MakeGenerics<{
  Search: {
    favorite?: string;
  } & RequestFormFields;
}>;

// Wrapper that checks if the user can view
// the provided request before renderining AccessRequestComponent.
const WithAccessRequestForm = () => {
  const {
    params: { id: ruleId },
  } = useMatch();

  // prevent the form resetting unexpectedly
  const { data: rule, error } = useUserGetAccessRule(ruleId, {
    swr: {
      refreshInterval: 0,
      revalidateIfStale: false,
      revalidateOnFocus: false,
      revalidateOnReconnect: false,
    },
  });
  const navigate = useNavigate();

  if (
    error &&
    (error?.response?.status === 403 || error.response?.status === 404)
  ) {
    return <ErrorNoAccess />;
  }

  // Users who are approval to the access rule but are not part of request group
  // should not be able to view the request form
  if (rule && !rule?.canRequest) {
    return <ErrorNoAccess />;
  }

  return <AccessRequestForm rule={rule} ruleId={ruleId} />;
};

interface AccessRequestProps {
  rule: RequestAccessRule | undefined;
  ruleId: string;
}

const AccessRequestForm = (props: AccessRequestProps) => {
  const { rule, ruleId } = props;
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();
  const now = useMemo(() => {
    const d = new Date();
    d.setSeconds(0, 0);
    return format(d, "yyyy-MM-dd'T'HH:mm");
  }, []);

  const toast = useToast();
  const search = useSearch<RequestFormQueryParameters>();
  const [favorite, setFavorite] = useState<FavoriteDetail>();

  const methods = useForm<NewRequestFormData>({
    defaultValues: {
      when: "asap",
      startDateTime: now,
      timing: {
        durationSeconds: 60,
      },
    },
  });
  const {
    register,
    handleSubmit,
    formState: { errors },
    setValue,
    control,
    watch,
    reset,
    getValues,
  } = methods;

  const resetForm = (fields: RequestFormFields) => {
    if (fields.timing) {
      setValue("timing.durationSeconds", fields.timing.durationSeconds);
      if (fields.timing.startTime) {
        setValue("startDateTime", fields.timing.startTime);
        setValue("when", "scheduled");
      }
    }
    fields.reason && setValue("reason", fields.reason);
    fields.with &&
      setValue(
        "with",
        fields.with.map((w) => {
          return { data: w };
        })
      );
  };
  // When the rule loads for the first time, this use effect:
  // 1: sets the duration to either 1 hour or max duration if it is less than one hour
  // 2: handles favouriting when the search & rule queries load
  // 3: hydrates the form with rule targets and arguments
  useEffect(() => {
    if (rule != undefined) {
      setValue(
        "timing.durationSeconds",
        rule.timeConstraints.maxDurationSeconds > 3600
          ? 3600
          : rule.timeConstraints.maxDurationSeconds
      );

      if (search.favorite) {
        userGetFavorite(search.favorite)
          .then((favorite) => {
            resetForm(favorite);
            setFavorite(favorite);
          })
          .catch((e) => {
            let description: string | undefined;
            if (axios.isAxiosError(e)) {
              description = (e as AxiosError<{ error: string }>)?.response?.data
                .error;
            }
            toast({
              title: "Failed to load favorite",
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
      } else {
        // The following will attempt to match any query params to withSelectable fields for this rule.
        // If the field matches and the value is a valid option, it will be set in the form values.
        // if it is not a valid value it is ignored.
        // this prevents being able to submit the form with bad options, or being able to submit arbitrary values for the with fields via the UI
        // resetForm(favorite);
        const filteredSearchWith = search.with?.map((w) => {
          const filteredWith: CreateRequestWith = {};
          Object.entries(w).map(([k, v]) => {
            if (rule.target.arguments[k]) {
              filteredWith[k] = v.filter((element) => {
                return !!rule.target.arguments[k].options.find(
                  (s) => s.value === element
                );
              });
            }
          });
          return filteredWith;
        });
        // default value if there is no favorite is an empty selection
        const fields: RequestFormFields = {
          with:
            filteredSearchWith === undefined || filteredSearchWith?.length == 0
              ? [{}]
              : filteredSearchWith,
          reason: search.reason,
          timing: search.timing,
        };

        // This hydrates the form with fields
        resetForm(fields);
      }
    }
  }, [rule, search]);

  const when = watch("when");
  const startTimeDate = watch("startDateTime");
  // Don't refetch the approvers
  const { data: approvers, isValidating: isValidatingApprovers } =
    useUserGetAccessRuleApprovers(ruleId, {
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
      with: data.with.filter((fw) => !fw.hidden).map((fw) => fw.data),
    };
    if (data.when === "scheduled") {
      r.timing.startTime = new Date(data.startDateTime).toISOString();
    }
    await userCreateRequest(r)
      .then(() => {
        toast({
          id: "user-request-created",
          title: "Request created",
          status: "success",
          duration: 2200,
          isClosable: true,
        });
        navigate({ to: "/requests" });
      })
      .catch((e: any) => {
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
  const [urlClipboardValue, setUrlClipboardValue] = useState("");
  const clipboard = useClipboard(urlClipboardValue);
  const location = useLocation();
  const formData = methods.watch();

  // case 0: rule is loading, disabled
  // case 1: has with fields but they're invalid, not filled out, disabled
  // case 1: no `with` fields, enabled
  const isDisabled =
    !rule ||
    !!formData?.with
      ?.filter((fw) => !fw.hidden)
      .find(
        (fw) =>
          !!Object.entries(fw.data).find(
            (o, k) => o[1] === undefined || o[1].length == 0
          )
      ) ||
    false;

  useEffect(() => {
    const a: RequestFormQueryParameters = {
      Search: {
        reason: getValues("reason"),
        with: (getValues("with") || [])
          .filter((fw) => !fw.hidden)
          .map((fw) => fw.data),
      },
    };
    const timing: RequestTiming = {
      durationSeconds: getValues("timing.durationSeconds"),
    };
    if (getValues("when") === "scheduled") {
      timing.startTime = new Date(getValues("startDateTime")).toISOString();
    }
    a.Search.timing = timing;
    const u = new URL(window.location.href);
    u.search = location.stringifySearch(a.Search);
    setUrlClipboardValue(u.toString());
    /** this is needed as redundancy bc. `urlClipboardValue` is not always stateful when used by useClipboard */
    clipboard.setValue(u.toString());
  }, [formData]);

  return (
    <>
      <UserLayout>
        <Helmet>
          <title>New Request</title>
        </Helmet>
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
          <FormProvider {...methods}>
            <Box
              p={8}
              bg="neutrals.100"
              mt={12}
              borderRadius="6px"
              as="form"
              onSubmit={handleSubmit(onSubmit)}
            >
              <Flex justify={"space-between"}>
                <Text as="h3" textStyle="Heading/H3">
                  You are requesting access to
                </Text>
                <ButtonGroup>
                  <FavoriteRequestButton
                    favorite={favorite}
                    ruleId={ruleId}
                    parentFormData={getValues()}
                    onUpdate={(f) => setFavorite(f)}
                  />
                  <Tooltip label="Copy a shareable link for this request">
                    <IconButton
                      variant={"ghost"}
                      aria-label="Copy link"
                      onClick={clipboard.onCopy}
                      icon={clipboard.hasCopied ? <CheckIcon /> : <LinkIcon />}
                    />
                  </Tooltip>
                </ButtonGroup>
              </Flex>
              <Stack
                spacing={2}
                mt={6}
                minH="52px" // prevents layout shift
              >
                {rule ? (
                  <>
                    <Flex data-testid="rule-name" align="center" mr="auto">
                      <ProviderIcon
                        shortType={rule?.target.provider.type}
                        id={rule?.target.provider.id}
                      />
                      <Text ml={2} textStyle="Body/Medium" color="neutrals.600">
                        {rule?.name}
                      </Text>
                    </Flex>
                    <Text textStyle="Body/Medium">{rule?.description}</Text>
                    <AccessRuleArguments target={rule.target} />
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
                        render={({ field: { ref, ...rest } }) => (
                          <DurationInput
                            {...rest}
                            max={rule?.timeConstraints.maxDurationSeconds}
                            min={60}
                            hideUnusedElements
                          >
                            <Weeks />
                            <Days />
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
                        )}
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
                      {...register("reason", { maxLength: 2048 })}
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
                    <Button
                      data-testid="request-submit-button"
                      type="submit"
                      disabled={isDisabled}
                      isLoading={loading}
                      mr={3}
                    >
                      Submit
                    </Button>
                  </Box>
                </Stack>
              </Box>
            </Box>
          </FormProvider>
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

export const AccessRuleArguments: React.FC<{
  target?: RequestAccessRuleTarget;
}> = ({ target }) => {
  const {
    control,
    getValues,
    formState: { errors },
    watch,
    setValue,
  } = useFormContext<NewRequestFormData>();

  if (target === undefined) {
    return <Skeleton minW="30ch" minH="6" mr="auto" />;
  }
  const subRequests = watch("with");
  return (
    <VStack align={"left"}>
      <VStack w="100%" spacing={4}>
        {subRequests?.map((sr, subRequestIndex) => {
          if (sr.hidden) {
            return null;
          }
          return (
            <Box position="relative" w="100%">
              {subRequests?.filter((sr) => !sr.hidden).length > 1 && (
                <IconButton
                  top={0}
                  right={0}
                  position={"absolute"}
                  type="button"
                  size="sm"
                  variant="ghost"
                  aria-label="remove"
                  icon={<DeleteIcon />}
                  onClick={() => {
                    const newSr = [...subRequests];
                    sr.hidden = true;
                    newSr[subRequestIndex] = sr;
                    setValue("with", newSr);
                  }}
                />
              )}
              <VStack
                w="100%"
                key={`subrequest-${subRequestIndex}`}
                id={`subrequest-${subRequestIndex}`}
                border="1px solid"
                borderColor="gray.300"
                rounded="md"
                px={4}
                py={4}
                spacing={4}
                align={"left"}
              >
                {Object.entries(target.arguments).filter(([k, v]) => {
                  return !v.requiresSelection;
                }).length > 0 && (
                  <Wrap spacing={4}>
                    {Object.entries(target.arguments)
                      .filter(([k, v]) => {
                        return !v.requiresSelection;
                      })
                      .map(([k, argument]) => {
                        return (
                          <WrapItem>
                            <VStack align={"left"}>
                              <Text>{argument.title}</Text>
                              <InfoOption
                                label={argument.options[0].label}
                                value={argument.options[0].value}
                              />
                            </VStack>
                          </WrapItem>
                        );
                      })}
                  </Wrap>
                )}
                {Object.entries(target.arguments)
                  .filter(([k, v]) => {
                    return v.requiresSelection;
                  })
                  .map(([k, v], i) => {
                    const name = `with.${subRequestIndex}.data.${k}`;
                    return (
                      <FormControl
                        key={"selectable-" + k}
                        pos="relative"
                        id={name}
                        isInvalid={
                          errors.with &&
                          errors.with?.[subRequestIndex]?.data?.[k] !==
                            undefined
                        }
                      >
                        <FormLabel
                          textStyle="Body/Medium"
                          color="neutrals.600"
                          fontWeight="normal"
                        >
                          {v.title}
                        </FormLabel>
                        {v.formElement === RequestArgumentFormElement.SELECT ? (
                          <SelectWithArrayAsValue
                            fieldName={`with.${subRequestIndex}.data.${k}`}
                            options={v.options
                              // exclude invalid options
                              .filter((op) => op.valid)
                              .map((op) => {
                                return op;
                              })}
                            rules={{
                              required: true,
                              validate: (value) => {
                                // @TODO validate that there is no overlap with other fields
                                return undefined;
                              },
                            }}
                          />
                        ) : (
                          <MultiSelect
                            fieldName={`with.${subRequestIndex}.data.${k}`}
                            options={v.options
                              // exclude invalid options
                              .filter((op) => op.valid)
                              .map((op) => {
                                return op;
                              })}
                            rules={{
                              required: true,
                              minLength: 1,
                              validate: (value) => {
                                // @TODO validate that there is no overlap with other fields
                                return undefined;
                              },
                            }}
                            id="user-request-access"
                          />
                        )}
                        <FormErrorMessage>
                          This field is required
                        </FormErrorMessage>
                      </FormControl>
                    );
                  })}
              </VStack>
            </Box>
          );
        })}
      </VStack>
      {/* Only render the add permissions button if the rule has fields which require selection */}
      {Object.entries(target.arguments).find(
        ([k, v]) => v.requiresSelection
      ) !== undefined && (
        <ButtonGroup>
          <Button
            pl={0}
            type="button"
            size="sm"
            variant="ghost"
            aria-label="add"
            leftIcon={<SmallAddIcon />}
            onClick={() => {
              setValue("with", [...(subRequests || []), { data: {} }]);
            }}
          >
            Add permissions
          </Button>
        </ButtonGroup>
      )}
    </VStack>
  );
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
            <CFAvatar
              // key={approver}
              userId={approver}
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

export default WithAccessRequestForm;

interface FavoriteRequestButtonProps {
  ruleId: string;
  parentFormData: NewRequestFormData;
  // if the page is currently loaded with a favorite
  favorite?: FavoriteDetail;
  onUpdate?: (favorite?: FavoriteDetail) => void;
}
const FavoriteRequestButton: React.FC<FavoriteRequestButtonProps> = ({
  ruleId,
  parentFormData,
  favorite,
  onUpdate,
}) => {
  const [isSubmitting, setIsSubmitting] = useState<boolean>(false);
  const methods = useForm<{ name: string }>({
    defaultValues: { name: favorite?.name },
  });
  useEffect(() => {
    if (favorite) {
      methods.reset({ name: favorite.name });
    }
  }, [favorite]);
  // the state of the parent form
  const popoverDisclosure = useDisclosure();
  const toast = useToast();

  const onSubmit: SubmitHandler<{ name: string }> = async (data) => {
    const r: CreateFavoriteRequestBody = {
      name: data.name,
      accessRuleId: ruleId,
      timing: {
        durationSeconds: parentFormData.timing.durationSeconds,
      },
      reason: parentFormData.reason ? parentFormData.reason : "",
      with: parentFormData.with.filter((fw) => !fw.hidden).map((fw) => fw.data),
    };
    if (parentFormData.when === "scheduled") {
      r.timing.startTime = new Date(parentFormData.startDateTime).toISOString();
    }
    setIsSubmitting(true);

    if (favorite) {
      userUpdateFavorite(favorite.id, r)
        .then((favorite) => {
          toast({
            id: "favourite-updated",
            title: "Favorite updated",
            status: "success",
            duration: 2200,
            isClosable: true,
          });
          popoverDisclosure.onClose();
          methods.reset();
          onUpdate?.(favorite);
        })
        .catch((e: any) => {
          let description: string | undefined;
          if (axios.isAxiosError(e)) {
            description = (e as AxiosError<{ error: string }>)?.response?.data
              .error;
          }
          toast({
            title: "Favorite failed to update",
            status: "error",
            duration: 5000,
            description: (
              <Text color={"white"} whiteSpace={"pre"}>
                {description}
              </Text>
            ),
            isClosable: true,
          });
        })
        .finally(() => {
          setIsSubmitting(false);
        });
    } else {
      userCreateFavorite(r)
        .then((favorite) => {
          toast({
            id: "favourite-created",
            title: "Favorite created",
            status: "success",
            duration: 2200,
            isClosable: true,
          });
          popoverDisclosure.onClose();
          methods.reset();
          onUpdate?.(favorite);
        })
        .catch((e: any) => {
          let description: string | undefined;
          if (axios.isAxiosError(e)) {
            description = (e as AxiosError<{ error: string }>)?.response?.data
              .error;
          }
          toast({
            title: "Favorite failed",
            status: "error",
            duration: 5000,
            description: (
              <Text color={"white"} whiteSpace={"pre"}>
                {description}
              </Text>
            ),
            isClosable: true,
          });
        })
        .finally(() => {
          setIsSubmitting(false);
        });
    }
  };

  const handleDeleteFavorite = () => {
    if (favorite) {
      setIsSubmitting(true);
      userDeleteFavorite(favorite?.id)
        .then(() => {
          toast({
            id: "favourite-removed",
            title: "Favorite removed",
            status: "success",
            duration: 2200,
            isClosable: true,
          });
          popoverDisclosure.onClose();
          methods.reset();
          onUpdate?.();
        })
        .catch((e: any) => {
          let description: string | undefined;
          if (axios.isAxiosError(e)) {
            description = (e as AxiosError<{ error: string }>)?.response?.data
              .error;
          }
          toast({
            title: "Failed to remove favorite",
            status: "error",
            duration: 5000,
            description: (
              <Text color={"white"} whiteSpace={"pre"}>
                {description}
              </Text>
            ),
            isClosable: true,
          });
        })
        .finally(() => {
          setIsSubmitting(false);
        });
    }
  };

  return (
    <Popover
      closeOnBlur={false}
      isOpen={popoverDisclosure.isOpen}
      onOpen={popoverDisclosure.onOpen}
      onClose={popoverDisclosure.onClose}
    >
      <Tooltip
        label={
          favorite
            ? "Update or remove this favorite"
            : "Add this request to your favorites"
        }
      >
        {/* additional element */}
        <Box display="inline-block">
          <PopoverTrigger>
            <IconButton
              data-testid="fav-icon-btn"
              color={favorite ? colors.actionWarning[200] : undefined}
              onClick={popoverDisclosure.onOpen}
              variant={"ghost"}
              aria-label="Favorite"
              icon={<StarIcon />}
            />
          </PopoverTrigger>
        </Box>
      </Tooltip>
      <PopoverContent>
        <PopoverArrow />
        <PopoverCloseButton />
        <PopoverHeader>
          {favorite ? "Update Favorite" : "Add to Favorites"}
        </PopoverHeader>

        {/* I have chosen not to use a native form element wrapper because it can't be easily nested in this popover inside the base request form

I experimented with using a <Portal/> to wrap the popover however this form submitting still triggered the parent form to submit

So I have just submitted the form directly using the submit button*/}
        <PopoverBody>
          <FormControl isInvalid={!!methods.formState.errors?.name}>
            <FormLabel textStyle="Body/Medium" fontWeight="normal">
              Name
            </FormLabel>

            <Input
              bg="white"
              id="nameField"
              data-testid="favourite-request-button"
              placeholder="Daily Development Access"
              {...methods.register("name", {
                required: true,
                minLength: 1,
                maxLength: 128,
                validate: (value) => {
                  const res: string[] = [];
                  [/[^a-zA-Z0-9,.;:()[\]?!\-_`~&/\n\s]/].every((pattern) =>
                    pattern.test(value as string)
                  ) &&
                    res.push(
                      "Invalid characters (only letters, numbers, and punctuation allowed)"
                    );
                  if (value && value.length > 128) {
                    res.push("Maximum length is 128 characters");
                  }
                  return res.length > 0 ? res.join(", ") : undefined;
                },
              })}
              onBlur={() => methods.trigger("name")}
            />
            <FormHelperText>
              Access favorites from your dashboard
            </FormHelperText>
            {methods.formState.errors?.name && (
              <FormErrorMessage>
                {methods.formState.errors?.name.message}
              </FormErrorMessage>
            )}
          </FormControl>
        </PopoverBody>
        <PopoverFooter>
          <Flex justify={"right"}>
            <Button
              size={"sm"}
              onClick={methods.handleSubmit(onSubmit)}
              mr={3}
              isLoading={isSubmitting}
            >
              {favorite ? "Update" : "Save"}
            </Button>
            {favorite && (
              <Button
                data-testid="del-fav-btn"
                variant={"danger"}
                size={"sm"}
                onClick={handleDeleteFavorite}
                mr={3}
                isLoading={isSubmitting}
              >
                Remove
              </Button>
            )}
          </Flex>
        </PopoverFooter>
      </PopoverContent>
    </Popover>
  );
};

const ErrorNoAccess = () => {
  return (
    <Flex
      height="100vh"
      padding="0"
      alignItems="center"
      justifyContent="center"
    >
      <Stack textAlign="center" w="70%" spacing={5}>
        <Heading>You don't have permission to access this</Heading>
        <Text>
          This access rule may no longer exist, or you don't have permission to
          view it.
        </Text>
        <Flex alignItems="center" justifyContent="center">
          <Button as={Link} to={"/"} h="42px" w="auto">
            Go back
          </Button>
        </Flex>
      </Stack>
    </Flex>
  );
};
