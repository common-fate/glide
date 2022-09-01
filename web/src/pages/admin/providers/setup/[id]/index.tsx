import {
  ArrowBackIcon,
  CheckCircleIcon,
  CheckIcon,
  ChevronDownIcon,
  ChevronRightIcon,
  WarningIcon,
} from "@chakra-ui/icons";
import {
  Accordion,
  AccordionButton,
  AccordionIcon,
  AccordionItem,
  AccordionPanel,
  Badge,
  Box,
  Button,
  Center,
  Circle,
  CircularProgress,
  Code,
  Container,
  Flex,
  FormControl,
  FormHelperText,
  FormLabel,
  Grid,
  GridItem,
  HStack,
  IconButton,
  Input,
  InputGroup,
  InputRightElement,
  Popover,
  PopoverArrow,
  PopoverBody,
  PopoverCloseButton,
  PopoverContent,
  PopoverHeader,
  PopoverTrigger,
  useClipboard,
  Spinner,
  Stack,
  Text,
} from "@chakra-ui/react";
import { useEffect, useState } from "react";
import Confetti from "react-confetti";
import { Controller, useForm } from "react-hook-form";
import { Link, useMatch } from "react-location";
import ReactMarkdown from "react-markdown";
import { Sticky, StickyContainer } from "react-sticky";
import useWindowSize from "react-use/lib/useWindowSize";
import { CodeInstruction } from "../../../../../components/CodeInstruction";
import { ConnectorArrow } from "../../../../../components/ConnectorArrow";
import { ApprovalsLogo } from "../../../../../components/icons/Logos";
import { ProviderIcon } from "../../../../../components/icons/providerIcon";
import { AdminLayout } from "../../../../../components/Layout";
import {
  submitProvidersetupStep,
  useGetProvidersetup,
  useGetProvidersetupInstructions,
  validateProvidersetup,
} from "../../../../../utils/backend-client/admin/admin";
import { ProviderSetupStepDetails } from "../../../../../utils/backend-client/types";
import { ProviderConfigValidation } from "../../../../../utils/backend-client/types/accesshandler-openapi.yml/providerConfigValidation";
import { registeredProviders } from "../../../../../utils/providerRegistry";
import { formatValidationErrorToText } from "./copyToClipboard";

const Page = () => {
  const {
    params: { id },
  } = useMatch();

  const { width, height } = useWindowSize();

  const [showConfetti, setShowConfetti] = useState(false);
  const [accordionIndex, setAccordionIndex] = useState([0]);
  const [validationErrorMsg, setValidationErrorMsg] = useState("");
  const { data, mutate } = useGetProvidersetup(id);
  const { data: instructions } = useGetProvidersetupInstructions(id);

  const { hasCopied, onCopy } = useClipboard(validationErrorMsg);

  // used to look up extra details like the name
  const registeredProvider = registeredProviders.find(
    (rp) => rp.type === data?.type
  );

  useEffect(() => {
    if (data !== undefined) {
      if (data.status !== "INITIAL_CONFIGURATION_IN_PROGRESS") {
        setAccordionIndex([]);
        return;
      }

      const initialIndex = data.steps.findIndex((s) => !s.complete);
      if (initialIndex > 0) {
        setAccordionIndex([initialIndex]);
      }
    }
  }, [data]);

  const stepsOverview = data?.steps ?? [];

  const completedSteps = stepsOverview.filter((s) => s.complete).length;

  const completedPercentage =
    stepsOverview.length ?? 0 > 0
      ? (completedSteps / stepsOverview.length) * 100
      : 0;

  const setupComplete = completedSteps === stepsOverview.length;

  const handleStepComplete = (index: number) => {
    setAccordionIndex([index + 1]);
  };

  const handleStartTests = async () => {
    if (data != null) {
      await mutate({ ...data, status: "VALIDATING" }, { revalidate: false });
      const res = await validateProvidersetup(id);

      if (res?.configValidation.length > 0) {
        setValidationErrorMsg(
          formatValidationErrorToText(res?.configValidation || [])
        );
      }
      if (res.status === "VALIDATION_SUCEEDED") {
        setShowConfetti(true);
      }
      await mutate(res);
    }
  };

  if (data === undefined) {
    return (
      <AdminLayout>
        <Center borderBottom="1px solid" borderColor="neutrals.200" h="80px">
          <IconButton
            as={Link}
            aria-label="Go back"
            pos="absolute"
            left={4}
            icon={<ArrowBackIcon />}
            rounded="full"
            variant="ghost"
            to="/admin/providers"
          />
          <Text as="h4" textStyle="Heading/H4"></Text>
        </Center>
        <Container
          my={12}
          // This prevents unbounded widths for small screen widths
          minW={{ base: "100%", xl: "container.xl" }}
          overflowX="auto"
        ></Container>
      </AdminLayout>
    );
  }

  return (
    <AdminLayout>
      {showConfetti && (
        <Confetti
          width={width}
          height={height}
          recycle={false}
          colors={["#ed77c0", "#619eff", "#30d15d"]}
        />
      )}
      <Stack
        justifyContent={"center"}
        alignItems={"center"}
        spacing={{ base: 1, md: 0 }}
        borderBottom="1px solid"
        borderColor="neutrals.200"
        h="80px"
        py={{ base: 4, md: 0 }}
        flexDirection={{ base: "column", md: "row" }}
      >
        <IconButton
          as={Link}
          aria-label="Go back"
          pos="absolute"
          left={4}
          icon={<ArrowBackIcon />}
          rounded="full"
          variant="ghost"
          to="/admin/providers"
        />
        <Text as="h4" textStyle="Heading/H4">
          {registeredProvider !== undefined &&
            `Setting up the ${registeredProvider.name} provider`}
        </Text>
        {data && (
          <HStack
            spacing={3}
            position={{ md: "absolute", base: "relative" }}
            right={{ md: 4, base: 0 }}
          >
            <Text>
              {completedSteps} of {data.steps.length} steps complete
            </Text>
            <CircularProgress value={completedPercentage} color="#449157" />
          </HStack>
        )}
      </Stack>
      <Container
        my={12}
        // This prevents unbounded widths for small screen widths
        minW={{ base: "100%", xl: "container.xl" }}
        overflowX="auto"
      >
        <Stack borderRadius="md" p={0} spacing={8}>
          <Accordion
            bg="neutrals.100"
            index={accordionIndex}
            allowMultiple
            onChange={(e) => Array.isArray(e) && setAccordionIndex(e)}
          >
            {(instructions?.stepDetails ?? []).map((step, index) => (
              <StepDisplay
                key={index}
                readOnly={data.status === "VALIDATION_SUCEEDED"}
                onComplete={() => handleStepComplete(index)}
                configValues={data.configValues}
                setupId={id}
                step={step}
                index={index}
                complete={data.steps[index]?.complete}
              />
            ))}
          </Accordion>
          {setupComplete &&
            data.status === "INITIAL_CONFIGURATION_IN_PROGRESS" && (
              <Stack justifyContent={"center"}>
                <Center>
                  <Stack display="block" flexGrow={0}>
                    <Text textStyle="Body/Small">
                      We'll verify your configuration by making a test
                      connection to the Access Provider.
                    </Text>
                    (
                    <Button w="100%" onClick={handleStartTests}>
                      Complete setup
                    </Button>
                  </Stack>
                </Center>
              </Stack>
            )}
          {data.status !== "INITIAL_CONFIGURATION_IN_PROGRESS" && (
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
                <HStack justifyContent={"center"} alignItems={"center"}>
                  {data.status === "VALIDATING" ? (
                    <Circle size="3" borderWidth={"1px"} />
                  ) : data.status === "VALIDATION_SUCEEDED" ? (
                    <CheckCircleIcon boxSize="3" color="actionSuccess.200" />
                  ) : (
                    <WarningIcon boxSize="3" color="actionWarning.200" />
                  )}
                  <Text textStyle="Body/Medium">Connection Test</Text>
                </HStack>
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
                  >
                    <ApprovalsLogo h="20px" w="auto" />
                    <ConnectorArrow animate={data.status === "VALIDATING"} />
                    <ProviderIcon type={data.type} h="24px" w="auto" />
                  </HStack>
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
                {data.configValidation.map((validation) => (
                  <ValidationResults
                    key={validation.id}
                    loading={data.status === "VALIDATING"}
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
                  <Button
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
                  </Button>
                </Flex>
              </Stack>
            </Stack>
          )}
          {data.status === "VALIDATION_FAILED" && (
            <Stack justifyContent={"center"}>
              <Center>
                <Stack display="block" flexGrow={0}>
                  <Text textStyle="Body/Small">
                    The connection test failed. Fix up your values above and
                    then retry the connection.
                  </Text>
                  (
                  <Button w="100%" onClick={handleStartTests}>
                    Retry connection
                  </Button>
                </Stack>
              </Center>
            </Stack>
          )}
          {data.status === "VALIDATION_SUCEEDED" && (
            <Stack justifyContent={"center"}>
              <Center>
                <Stack display="block" flexGrow={0}>
                  <Text textStyle="Body/Large">
                    Looking good! The connection tests are passing.
                  </Text>
                  (
                  <Button w="100%" as={Link} to="./finish">
                    Next
                  </Button>
                </Stack>
              </Center>
            </Stack>
          )}
        </Stack>
      </Container>
    </AdminLayout>
  );
};

interface ValidationResultsProps {
  loading?: boolean;
  validation: ProviderConfigValidation;
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
          ) : validation.status === "SUCCESS" ? (
            <CheckCircleIcon boxSize={"12px"} color="neutrals.500" />
          ) : (
            <WarningIcon boxSize={"12px"} color="neutrals.500" />
          )}
        </Flex>
        <Text color="#d0d7de">{validation.name}</Text>
      </Flex>
      {expanded && (
        <Stack pl={"60px"} spacing={1}>
          {validation.logs.map((log, index) => (
            <Text
              key={index}
              color="#d0d7de"
              fontSize={"12px"}
              fontFamily="mono"
            >
              {log.level}: {log.msg}
            </Text>
          ))}
        </Stack>
      )}
    </Stack>
  );
};

interface StepDisplayProps {
  setupId: string;
  configValues: Record<string, string>;
  step: ProviderSetupStepDetails;
  index: number;
  readOnly?: boolean;
  onComplete?: () => void;
  complete: boolean;
}

const StepDisplay: React.FC<StepDisplayProps> = ({
  step,
  readOnly,
  index,
  complete,
  onComplete,
  setupId,
  configValues,
}) => {
  const [loading, setLoading] = useState(false);
  const { mutate } = useGetProvidersetup(setupId);
  const { control, handleSubmit } = useForm<Record<string, string>>({
    defaultValues: configValues,
  });
  const onSubmit = async (data: Record<string, string>) => {
    const filteredData: Record<string, string> = {};
    step.configFields.forEach((f) => {
      filteredData[f.id] = data[f.id];
    });
    
    setLoading(true);
    const res = await submitProvidersetupStep(setupId, index, {
      complete: true,
      configValues: filteredData,
    });

    await mutate(res);
    setLoading(false);
    onComplete?.();
  };

  return (
    <AccordionItem>
      <h2>
        <AccordionButton>
          <Flex flex="1" textAlign="left">
            <Box
              display="inline-flex"
              alignItems={"center"}
              justifyContent="center"
              as="span"
              mr="2"
              bg={complete ? "actionSuccess.200" : "white"}
              borderRadius={"50%"}
              w="24px"
              h="24px"
              borderWidth={"1px"}
            >
              <CheckIcon boxSize={"13px"} color="white" />
            </Box>
            {index + 1}: {step.title}
          </Flex>
          <AccordionIcon />
        </AccordionButton>
      </h2>
      <AccordionPanel pb={4}>
        <Grid
          templateColumns={{ md: "repeat(3, 1fr)", base: "repeat(1, 1fr)" }}
          gap={4}
        >
          <GridItem colSpan={2}>
            <Stack pt={2}>
              <Text>Instructions</Text>
              <ReactMarkdown
                components={{
                  a: (props) => (
                    <Link target="_blank" rel="noreferrer" {...props} />
                  ),
                  p: (props) => (
                    <Text
                      pb={3}
                      as="span"
                      color="neutrals.600"
                      textStyle={"Body/Small"}
                    >
                      {props.children}
                    </Text>
                  ),
                  code: CodeInstruction as any,
                }}
              >
                {step.instructions}
              </ReactMarkdown>
            </Stack>
          </GridItem>
          <GridItem position="relative" as={StickyContainer}>
            <Sticky>
              {({ style }) => (
                <Stack
                  style={style}
                  pt={2}
                  as="form"
                  onSubmit={handleSubmit(onSubmit)}
                  autoComplete="off"
                  spacing={5}
                >
                  {step.configFields.length > 0 && (
                    <Stack spacing={5}>
                      {readOnly !== true && <Text>Enter your values</Text>}
                      {step.configFields.map((f) => (
                        <FormControl
                          isReadOnly={readOnly}
                          isRequired={!f.isOptional}
                          key={f.id}
                        >
                          <FormLabel>
                            {f.name}{" "}
                            {f.isSecret && (
                              <Popover>
                                <PopoverTrigger>
                                  <Badge as="button">Secret</Badge>
                                </PopoverTrigger>
                                <PopoverContent>
                                  <PopoverArrow />
                                  <PopoverCloseButton />
                                  <PopoverHeader>Sensitive value</PopoverHeader>
                                  <PopoverBody>
                                    This value will be written to{" "}
                                    <Code>{f.secretPath}</Code>
                                  </PopoverBody>
                                </PopoverContent>
                              </Popover>
                            )}
                          </FormLabel>
                          <Controller
                            control={control}
                            name={f.id}
                            render={({ field }) => {
                              return (
                                <ConfigValueInput
                                  isSecret={f.isSecret}
                                  {...field}
                                />
                              );
                            }}
                          />
                          <FormHelperText>{f.description}</FormHelperText>
                        </FormControl>
                      ))}
                    </Stack>
                  )}
                  <Flex justifyContent={"flex-end"} pt={3}>
                    {readOnly !== true && (
                      <Button
                        flexGrow={0}
                        type="submit"
                        isLoading={loading}
                        variant={complete ? "brandSecondary" : "brandPrimary"}
                      >
                        {complete
                          ? "Update values"
                          : `I've completed step ${index + 1}`}
                      </Button>
                    )}
                  </Flex>
                </Stack>
              )}
            </Sticky>
          </GridItem>
        </Grid>
      </AccordionPanel>
    </AccordionItem>
  );
};

interface ConfigValueInputProps {
  onChange: (n: string) => void;
  value?: string;
  defaultValue?: string;
  isSecret: boolean;
  ref: React.LegacyRef<HTMLInputElement>;
}

export const ConfigValueInput: React.FC<ConfigValueInputProps> = ({
  onChange,
  value,
  defaultValue,
  ref,
  isSecret,
}) => {
  const [locked, setLocked] = useState(
    // we assume anything that starts with awsssm:// is a reference to a secret, rather than a secret itself.
    isSecret && value?.startsWith("awsssm://")
  );
  useEffect(() => {
    if (!locked && isSecret) {
      onChange("");
    }
  }, [locked, isSecret]);

  return (
    <InputGroup size="md">
      <Input
        ref={ref}
        pr={locked ? "4.5rem" : undefined}
        bg="white"
        defaultValue={defaultValue}
        isReadOnly={locked}
        value={value}
        onChange={(e) => onChange(e.target.value)}
      />
      {locked && (
        <InputRightElement width="4.5rem">
          <Button
            variant={"solid"}
            h="1.75rem"
            size="sm"
            onClick={() => setLocked(false)}
          >
            Reset
          </Button>
        </InputRightElement>
      )}
    </InputGroup>
  );
};

export default Page;
