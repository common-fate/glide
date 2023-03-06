import {
  ButtonGroup,
  Circle,
  Code,
  Container,
  Flex,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalHeader,
  ModalOverlay,
  Text,
  useClipboard,
  useDisclosure,
} from "@chakra-ui/react";
import { useMemo, useState } from "react";
import { Helmet } from "react-helmet";
import { Column } from "react-table";
import { AdminLayout } from "../../../components/Layout";
import { TabsStyledButton } from "../../../components/nav/Navbar";
import { TableRenderer } from "../../../components/tables/TableRenderer";
import {
  useAdminListHandlers,
  useAdminListTargetGroups,
} from "../../../utils/backend-client/admin/admin";
import {
  Diagnostic,
  TargetGroup,
  TGHandler,
} from "../../../utils/backend-client/types";
import { usePaginatorApi } from "../../../utils/usePaginatorApi";

// using a chakra tab component and links, link to /admin/providers and /admin/providersv2
export const ProvidersV2Tabs = () => {
  return (
    <ButtonGroup variant="ghost" spacing="0" mb={"-32px !important;"} my={4}>
      <TabsStyledButton href="/admin/providers">V1</TabsStyledButton>
      <TabsStyledButton href="/admin/providersv2">V2</TabsStyledButton>
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
        // @ts-ignore this is required because ts cannot infer the nexted object types correctly
        accessor: "targetGroupAssignment.TargetGroupId",
        Header: "Target Group Id",
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
  const { data } = useAdminListTargetGroups();
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
    data: [],
    emptyText: "No Target Groups have been set up yet.",
    linkTo: false,
  });
};

const Providers = () => {
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
          <ProvidersV2Tabs />
        </Flex>

        <Container
          pb={9}
          // This prevents unbounded widths for small screen widths
          minW={{ base: "100%", xl: "container.xl" }}
          overflowX="auto"
        >
          Target Groups
          <AdminTargetGroupsTable />
        </Container>

        <Container
          pb={9}
          // This prevents unbounded widths for small screen widths
          minW={{ base: "100%", xl: "container.xl" }}
          overflowX="auto"
        >
          Handlers
          <AdminProvidersTable />
        </Container>
      </Container>
    </AdminLayout>
  );
};

export default Providers;
