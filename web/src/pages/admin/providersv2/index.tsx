import { CopyIcon } from "@chakra-ui/icons";
import {
  ButtonGroup,
  Circle,
  Code,
  Container,
  Flex,
  IconButton,
  Text,
  useClipboard,
} from "@chakra-ui/react";
import { useMemo } from "react";
import { Helmet } from "react-helmet";
import { Column } from "react-table";
import { AdminLayout } from "../../../components/Layout";
import { TabsStyledButton } from "../../../components/nav/Navbar";
import { TableRenderer } from "../../../components/tables/TableRenderer";

import {
  Diagnostic,
  TGHandler,
  TargetGroup,
} from "../../../utils/backend-client/types";
import { usePaginatorApi } from "../../../utils/usePaginatorApi";
import {
  useAdminListHandlers,
  useAdminListTargetGroups,
} from "../../../utils/backend-client/admin/admin";

// using a chakra tab component and links, link to /admin/providers and /admin/providersv2
export const ProvidersV2Tabs = () => {
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
          return (
            <Code
              rounded="md"
              fontSize="sm"
              p={2}
              noOfLines={3}
              position="relative"
            >
              {strippedCode}
              <IconButton
                aria-label="Copy"
                variant="ghost"
                icon={<CopyIcon />}
                size="xs"
                position="absolute"
                bottom={0}
                right={0}
                opacity={0.5}
                onClick={() => {
                  clippy.setValue(strippedCode);
                  clippy.onCopy();
                  console.log("copied", strippedCode);
                }}
              />
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
    data: data?.targetGroups,
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
