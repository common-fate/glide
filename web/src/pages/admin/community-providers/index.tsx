import {
  Button,
  ButtonGroup,
  Circle,
  Code,
  Container,
  Flex,
  HStack,
  Heading,
  Link,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalHeader,
  ModalOverlay,
  Text,
  VStack,
  useClipboard,
  useDisclosure,
  useToast,
  Tooltip,
} from "@chakra-ui/react";

import { useMemo, useState } from "react";
import { Helmet } from "react-helmet";
import { Column } from "react-table";
import { AdminLayout } from "../../../components/Layout";
import { TabsStyledButton } from "../../../components/nav/Navbar";
import { TableRenderer } from "../../../components/tables/TableRenderer";
import {
  adminHealthcheckHandlers,
  useAdminListHandlers,
  useAdminListTargetGroups,
} from "../../../utils/backend-client/admin/admin";
import {
  Diagnostic,
  TargetGroup,
  TGHandler,
} from "../../../utils/backend-client/types";
import { usePaginatorApi } from "../../../utils/usePaginatorApi";
import { HealthCheckIcon, RefreshIcon } from "../../../components/icons/Icons";

import axios from "axios";

// using a chakra tab component and links, link to /admin/providers and /admin/providersv2
export const CommunityProvidersTabs = () => {
  return (
    <ButtonGroup variant="ghost" spacing="0" mb={"-32px !important;"} my={4}>
      <TabsStyledButton href="/admin/providers">
        Built-In Providers
      </TabsStyledButton>
      <TabsStyledButton href="/admin/providersv2">
        PDK Providers
      </TabsStyledButton>
    </ButtonGroup>
  );
};

const AdminProvidersTable = () => {
  const paginator = usePaginatorApi<typeof useAdminListHandlers>({
    swrHook: useAdminListHandlers,
    hookProps: {},
  });

  const clippy = useClipboard("");

  const diagnosticModal = useDisclosure();
  const [diagnosticText, setDiagnosticText] = useState("");

  const cols: Column<TGHandler>[] = useMemo(
    () => [
      {
        accessor: "id",
        Header: "ID",
      },
      {
        accessor: "awsRegion",
        Header: "Region",
      },
      {
        accessor: "awsAccount",
        Header: "Account",
      },
      {
        accessor: "healthy",
        Header: "Health",
        Cell: ({ value }) => (
          <Flex minW="75px" align="center">
            <Circle
              bg={value ? "actionSuccess.200" : "actionWarning.200"}
              size="8px"
              mr={2}
            />
            <Text as="span">{value ? "Healthy" : "Unhealthy"}</Text>
          </Flex>
        ),
      },
      {
        accessor: "diagnostics",
        Header: "Diagnostics",
        Cell: ({ value }) => {
          // Strip out the code from the diagnostics, it's currently an empty field
          const strippedCode = JSON.stringify(
            (value as Partial<Diagnostic>[]).map((v) => {
              delete v["code"];
              return v;
            })
          );

          const maxDiagnosticChars = 200;
          let expandCode = false;
          if (strippedCode.length > maxDiagnosticChars) {
            expandCode = true;
          }

          const handleClick = () => {
            if (expandCode) {
              diagnosticModal.onOpen();
              setDiagnosticText(strippedCode);
            }
          };

          return (
            <Code
              rounded="md"
              fontSize="sm"
              userSelect={expandCode ? "none" : "auto"}
              p={2}
              noOfLines={3}
              onClick={handleClick}
              position="relative"
              _hover={{
                "backgroundColor": expandCode ? "gray.600" : "gray.200",
                "cursor": expandCode ? "pointer" : "default",
                "#expandCode": {
                  display: "block",
                },
              }}
            >
              {expandCode && (
                <Text
                  id="expandCode"
                  display="none"
                  position="absolute"
                  left="50%"
                  top="50%"
                  transform="translate(-50%, -50%)"
                  zIndex={2}
                  size="md"
                  color="gray.50"
                >
                  Expand code
                </Text>
              )}
              {strippedCode}
            </Code>
          );
        },
      },
    ],
    []
  );

  return TableRenderer<TGHandler>({
    columns: cols,
    data: paginator?.data?.res,
    emptyText: "No Handlers have been set up yet.",
    linkTo: false,
    apiPaginator: paginator,
    additionalChildren: (
      <Modal isOpen={diagnosticModal.isOpen} onClose={diagnosticModal.onClose}>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>Diagnostics</ModalHeader>
          <ModalCloseButton />
          <ModalBody pb={4}>
            <Code rounded="md" minH="200px" fontSize="sm" p={2}>
              {diagnosticText}
            </Code>
          </ModalBody>
        </ModalContent>
      </Modal>
    ),
  });
};

const AdminTargetGroupsTable = () => {
  const paginator = usePaginatorApi<typeof useAdminListTargetGroups>({
    swrHook: useAdminListTargetGroups,
    hookProps: {},
  });
  // @ts-ignore this is required because ts cannot infer the nexted object types correctly
  const cols: Column<TargetGroup>[] = useMemo(
    () => [
      {
        accessor: "id",
        Header: "ID",
      },
      {
        accessor: "targetSchema.From",
        Header: "From",
      },
    ],
    []
  );

  return TableRenderer<TargetGroup>({
    columns: cols,
    data: paginator?.data?.targetGroups,
    apiPaginator: paginator,
    emptyText: "No Target Groups have been set up yet.",
    linkTo: true,
  });
};

const Providers = () => {
  const [loading, setLoading] = useState(false);
  const toast = useToast();
  const onClick = async () => {
    setLoading(true);

    await adminHealthcheckHandlers()
      .then(() => {
        toast({
          title: "Health check run",
          status: "success",
          variant: "subtle",
          duration: 2200,
          isClosable: true,
        });
      })
      .catch((err) => {
        let description: string | undefined;
        if (axios.isAxiosError(err)) {
          // @ts-ignore
          description = err?.response?.data.error;
        }
        toast({
          title: "Error running health check",
          description,
          status: "error",
          variant: "subtle",
          duration: 2200,
          isClosable: true,
        });
      });
    setLoading(false);
  };
  return (
    <AdminLayout>
      <Helmet>
        <title>Providers</title>
      </Helmet>
      <Container
        my={12}
        // This prevents unbounded widths for small screen widths
        minW={{ base: "100%", xl: "container.xl" }}
        overflowX="auto"
      >
        {/* spacer of 32px to acccount for un-needed UI/CLS */}
        <Flex justify="space-between" align="center">
          <CommunityProvidersTabs />
          <HStack spacing="1px">
            <Button
              isLoading={loading}
              onClick={() => onClick()}
              size="s"
              padding="5px"
              variant="ghost"
              iconSpacing="0px"
              leftIcon={<HealthCheckIcon boxSize="24px" />}
              data-testid="create-access-rule-button"
            ></Button>
            <Tooltip
              hasArrow
              label="Healthcheck will poll any deployed lambda providers to validate health and validity"
            >
              <Text textStyle={"Body/Small"}>Run Health Check</Text>
            </Tooltip>
          </HStack>
        </Flex>

        <VStack pb={9} align={"left"}>
          <Text textStyle="Heading/H4">Target Groups</Text>
          <AdminTargetGroupsTable />
        </VStack>

        <VStack pb={9} align={"left"}>
          <Text textStyle="Heading/H4">Handlers</Text>
          <AdminProvidersTable />
        </VStack>
      </Container>
    </AdminLayout>
  );
};

export default Providers;
