import { ArrowBackIcon, CheckIcon } from "@chakra-ui/icons";
import {
  Accordion,
  AccordionButton,
  AccordionIcon,
  AccordionItem,
  AccordionPanel,
  Box,
  Button,
  Center,
  CircularProgress,
  Container,
  Flex,
  FormControl,
  FormLabel,
  Grid,
  GridItem,
  HStack,
  IconButton,
  Input,
  Spinner,
  useAccordionContext,
  Stack,
  Text,
} from "@chakra-ui/react";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { Link, useMatch } from "react-location";
import ReactMarkdown from "react-markdown";
import { Sticky, StickyContainer } from "react-sticky";
import { CodeInstruction } from "../../../../components/CodeInstruction";
import { AdminLayout } from "../../../../components/Layout";
import {
  submitProvidersetupStep,
  useGetProvidersetup,
  useGetProvidersetupInstructions,
} from "../../../../utils/backend-client/admin/admin";
import { ProviderSetupStepDetails } from "../../../../utils/backend-client/types";
import { registeredProviders } from "../../../../utils/providerRegistry";

const Page = () => {
  const {
    params: { id },
  } = useMatch();

  const [accordionIndex, setAccordionIndex] = useState([0]);
  const { data } = useGetProvidersetup(id);
  const { data: instructions } = useGetProvidersetupInstructions(id);

  // used to look up extra details like the name
  const registeredProvider = registeredProviders.find(
    (rp) => rp.type === data?.type
  );

  useEffect(() => {
    if (data !== undefined) {
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

  const handleStepComplete = (index: number) => {
    setAccordionIndex([index + 1]);
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
        <Text as="h4" textStyle="Heading/H4">
          {registeredProvider !== undefined &&
            `Setting up the ${registeredProvider.name} provider`}
        </Text>
        {data && (
          <HStack spacing={3} position="absolute" right={4}>
            <Text>
              {completedSteps} of {data.steps.length} steps complete
            </Text>
            <CircularProgress value={completedPercentage} color="#449157" />
          </HStack>
        )}
        {/* <Button
          pos="absolute"
          right={0}
          size="sm"
          variant="ghost"
          leftIcon={<DeleteIcon />}
        >
          Cancel setup
        </Button> */}
      </Center>
      <Container
        my={12}
        // This prevents unbounded widths for small screen widths
        minW={{ base: "100%", xl: "container.xl" }}
        overflowX="auto"
      >
        <Stack bg="neutrals.100" borderRadius="md" p={0}>
          <Accordion
            index={accordionIndex}
            allowMultiple
            onChange={(e) => Array.isArray(e) && setAccordionIndex(e)}
          >
            {(instructions?.stepDetails ?? []).map((step, index) => (
              <StepDisplay
                onComplete={() => handleStepComplete(index)}
                configValues={data.configValues}
                setupId={id}
                step={step}
                index={index}
                complete={data.steps[index].complete}
              />
            ))}
          </Accordion>
        </Stack>
      </Container>
    </AdminLayout>
  );
};

interface StepDisplayProps {
  setupId: string;
  configValues: Record<string, string>;
  step: ProviderSetupStepDetails;
  index: number;
  onComplete?: () => void;
  complete: boolean;
}

const StepDisplay: React.FC<StepDisplayProps> = ({
  step,
  index,
  complete,
  onComplete,
  setupId,
  configValues,
}) => {
  const [loading, setLoading] = useState(false);
  const { mutate } = useGetProvidersetup(setupId);
  const { register, handleSubmit } = useForm<Record<string, string>>({
    defaultValues: configValues,
  });

  const onSubmit = async (data: Record<string, string>) => {
    setLoading(true);
    const res = await submitProvidersetupStep(setupId, index, {
      complete: true,
      configValues: data,
    });
    void mutate(res);
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
              bg={complete ? "brandGreen.300" : "white"}
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
        <Grid templateColumns="repeat(3, 1fr)" gap={4}>
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
                    <Stack>
                      <Text>Enter your values</Text>
                      {step.configFields.map((field) => (
                        <FormControl
                          isRequired={!field.isOptional}
                          key={field.id}
                        >
                          <FormLabel>{field.name}</FormLabel>
                          <Input bg="white" {...register(field.id)} />
                        </FormControl>
                      ))}
                    </Stack>
                  )}
                  <Flex justifyContent={"flex-end"} mt={3}>
                    <Button flexGrow={0} type="submit" isLoading={loading}>
                      I've completed step {index + 1}
                    </Button>
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

export default Page;
function useAccordionItemContext(): { isOpen: any; isDisabled: any } {
  throw new Error("Function not implemented.");
}
