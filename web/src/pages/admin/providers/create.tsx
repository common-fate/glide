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
  SimpleGrid,
  Stack,
  Text,
} from "@chakra-ui/react";
import { StickyContainer, Sticky } from "react-sticky";
import { Link, MakeGenerics, useNavigate, useSearch } from "react-location";
import ReactMarkdown from "react-markdown";
import { CodeInstruction } from "../../../components/CodeInstruction";
import { ProviderIcon } from "../../../components/icons/providerIcon";
import { AdminLayout } from "../../../components/Layout";
import { registeredProviders } from "../../../utils/providerRegistry";

type MyLocationGenerics = MakeGenerics<{
  Search: {
    type?: string;
  };
}>;

const CreateProvider = () => {
  const search = useSearch<MyLocationGenerics>();
  const navigate = useNavigate<MyLocationGenerics>();

  const { type } = search;

  if (type === undefined) {
    // provider selection UI
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
            New Access Provider
          </Text>
        </Center>
        <Container
          my={12}
          // This prevents unbounded widths for small screen widths
          minW={{ base: "100%", xl: "container.xl" }}
          overflowX="auto"
        >
          <SimpleGrid columns={2} spacing={4} p={1}>
            {registeredProviders.map((provider) => (
              <Box
                as="button"
                className="group"
                textAlign="center"
                bg="neutrals.100"
                p={6}
                rounded="md"
                data-testid={"provider_" + provider.type}
                onClick={() => navigate({ search: { type: provider.type } })}
              >
                <ProviderIcon provider={provider.type} mb={3} h="8" w="8" />

                <Text textStyle="Body/SmallBold" color="neutrals.700">
                  {provider.name}
                </Text>
              </Box>
            ))}
          </SimpleGrid>
        </Container>
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
          to="/admin/providers/create"
        />
        <Text as="h4" textStyle="Heading/H4">
          Setting up the AWS SSO provider
        </Text>
        <HStack spacing={3} position="absolute" right={4}>
          <Text>2 of 5 steps complete</Text>
          <CircularProgress value={80} color="#449157" />
        </HStack>
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
          <Accordion defaultIndex={[0]} allowMultiple>
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
                      bg="brandGreen.300"
                      borderRadius={"50%"}
                      w="24px"
                      h="24px"
                      borderWidth={"1px"}
                    >
                      <CheckIcon boxSize={"13px"} color="white" />
                    </Box>
                    1: Create an IAM role
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
                        {STEP_1}
                      </ReactMarkdown>
                    </Stack>
                  </GridItem>
                  <GridItem position="relative" as={StickyContainer}>
                    <Sticky>
                      {({ style }) => (
                        <Stack style={style} pt={2}>
                          <Flex justifyContent={"flex-end"} mt={3}>
                            <Button flexGrow={0}>I've completed step 1</Button>
                          </Flex>
                        </Stack>
                      )}
                    </Sticky>
                  </GridItem>
                </Grid>
              </AccordionPanel>
            </AccordionItem>

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
                      bg="brandGreen.300"
                      borderRadius={"50%"}
                      w="24px"
                      h="24px"
                      borderWidth={"1px"}
                    >
                      <CheckIcon boxSize={"13px"} color="white" />
                    </Box>
                    2: Enter your AWS SSO instance details
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
                        {STEP_2}
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
                          autoComplete="off"
                          spacing={5}
                        >
                          <Stack>
                            <Text>Enter your values</Text>
                            <FormControl>
                              <FormLabel>Instance ARN</FormLabel>
                              <Input bg="white" />
                            </FormControl>
                            <FormControl>
                              <FormLabel>Instance API key</FormLabel>
                              <InputGroup>
                                <Input
                                  bg="white"
                                  type="password"
                                  autoComplete="off"
                                />
                                {/* <InputRightElement width="4rem">
                                <Button
                                  mr={1}
                                  h="1.75rem"
                                  size="sm"
                                  variant={"solid"}
                                >
                                  {"Save"}
                                </Button>
                              </InputRightElement> */}
                              </InputGroup>
                              <FormHelperText>
                                Will be written to AWS SSM:
                                <Code
                                  fontSize={"xs"}
                                  wordBreak="break-all"
                                  overflowWrap="break-word"
                                >
                                  /granted/providers/aws/testing/secret:1
                                </Code>
                              </FormHelperText>
                            </FormControl>
                          </Stack>
                          <Flex justifyContent={"flex-end"} mt={3}>
                            <Button flexGrow={0}>I've completed step 2</Button>
                          </Flex>
                        </Stack>
                      )}
                    </Sticky>
                  </GridItem>
                </Grid>
              </AccordionPanel>
            </AccordionItem>
          </Accordion>
        </Stack>
      </Container>
    </AdminLayout>
  );
};

const STEP_1 = `
Create an IAM role in the AWS console:

\`\`\`json
{
  "Version": "2012-10-17",
  "Statement": [
      {
      "Action": [
        "sso:CreateAccountAssignment",
        "sso:DeleteAccountAssignment",
        "sso:ListAccountAssignments",
        "sso:ListTagsForResource",
        "identitystore:ListUsers",
        "organizations:DescribeAccount",
        "sso:CreatePermissionSet",
        "sso:PutInlinePolicyToPermissionSet",
        "sso:ListPermissionSets",
        "sso:DescribePermissionSet",
        "sso:DeletePermissionSet",
        "iam:ListRoles"
      ],
      "Resource": "*",
      "Effect": "Allow"
    }
  ]
}
\`\`\`
`;

const STEP_2 = `
Open the AWS console in the account that your AWS SSO instance is deployed to. If your company is using AWS Control Tower, this will be the root account in your AWS organisation.

Visit the Settings tab. The information about your SSO instance will be shown here, including the Instance ARN (as the “ARN” field) and the Identity Store ID.

Open the AWS console in the account that your AWS SSO instance is deployed to. If your company is using AWS Control Tower, this will be the root account in your AWS organisation.

Visit the Settings tab. The information about your SSO instance will be shown here, including the Instance ARN (as the “ARN” field) and the Identity Store ID.

Open the AWS console in the account that your AWS SSO instance is deployed to. If your company is using AWS Control Tower, this will be the root account in your AWS organisation.

Visit the Settings tab. The information about your SSO instance will be shown here, including the Instance ARN (as the “ARN” field) and the Identity Store ID.

Open the AWS console in the account that your AWS SSO instance is deployed to. If your company is using AWS Control Tower, this will be the root account in your AWS organisation.

Visit the Settings tab. The information about your SSO instance will be shown here, including the Instance ARN (as the “ARN” field) and the Identity Store ID.

Open the AWS console in the account that your AWS SSO instance is deployed to. If your company is using AWS Control Tower, this will be the root account in your AWS organisation.

Visit the Settings tab. The information about your SSO instance will be shown here, including the Instance ARN (as the “ARN” field) and the Identity Store ID.

Open the AWS console in the account that your AWS SSO instance is deployed to. If your company is using AWS Control Tower, this will be the root account in your AWS organisation.

Visit the Settings tab. The information about your SSO instance will be shown here, including the Instance ARN (as the “ARN” field) and the Identity Store ID.

Open the AWS console in the account that your AWS SSO instance is deployed to. If your company is using AWS Control Tower, this will be the root account in your AWS organisation.

Visit the Settings tab. The information about your SSO instance will be shown here, including the Instance ARN (as the “ARN” field) and the Identity Store ID.

Open the AWS console in the account that your AWS SSO instance is deployed to. If your company is using AWS Control Tower, this will be the root account in your AWS organisation.

Visit the Settings tab. The information about your SSO instance will be shown here, including the Instance ARN (as the “ARN” field) and the Identity Store ID.

Open the AWS console in the account that your AWS SSO instance is deployed to. If your company is using AWS Control Tower, this will be the root account in your AWS organisation.

Visit the Settings tab. The information about your SSO instance will be shown here, including the Instance ARN (as the “ARN” field) and the Identity Store ID.

Open the AWS console in the account that your AWS SSO instance is deployed to. If your company is using AWS Control Tower, this will be the root account in your AWS organisation.

Visit the Settings tab. The information about your SSO instance will be shown here, including the Instance ARN (as the “ARN” field) and the Identity Store ID.

Open the AWS console in the account that your AWS SSO instance is deployed to. If your company is using AWS Control Tower, this will be the root account in your AWS organisation.

Visit the Settings tab. The information about your SSO instance will be shown here, including the Instance ARN (as the “ARN” field) and the Identity Store ID.

Open the AWS console in the account that your AWS SSO instance is deployed to. If your company is using AWS Control Tower, this will be the root account in your AWS organisation.

Visit the Settings tab. The information about your SSO instance will be shown here, including the Instance ARN (as the “ARN” field) and the Identity Store ID.

Open the AWS console in the account that your AWS SSO instance is deployed to. If your company is using AWS Control Tower, this will be the root account in your AWS organisation.

Visit the Settings tab. The information about your SSO instance will be shown here, including the Instance ARN (as the “ARN” field) and the Identity Store ID.

Open the AWS console in the account that your AWS SSO instance is deployed to. If your company is using AWS Control Tower, this will be the root account in your AWS organisation.

Visit the Settings tab. The information about your SSO instance will be shown here, including the Instance ARN (as the “ARN” field) and the Identity Store ID.

Open the AWS console in the account that your AWS SSO instance is deployed to. If your company is using AWS Control Tower, this will be the root account in your AWS organisation.

Visit the Settings tab. The information about your SSO instance will be shown here, including the Instance ARN (as the “ARN” field) and the Identity Store ID.

Open the AWS console in the account that your AWS SSO instance is deployed to. If your company is using AWS Control Tower, this will be the root account in your AWS organisation.

Visit the Settings tab. The information about your SSO instance will be shown here, including the Instance ARN (as the “ARN” field) and the Identity Store ID.

Open the AWS console in the account that your AWS SSO instance is deployed to. If your company is using AWS Control Tower, this will be the root account in your AWS organisation.

Visit the Settings tab. The information about your SSO instance will be shown here, including the Instance ARN (as the “ARN” field) and the Identity Store ID.

Open the AWS console in the account that your AWS SSO instance is deployed to. If your company is using AWS Control Tower, this will be the root account in your AWS organisation.

Visit the Settings tab. The information about your SSO instance will be shown here, including the Instance ARN (as the “ARN” field) and the Identity Store ID.
`;

export default CreateProvider;
