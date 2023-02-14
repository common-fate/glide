import { CloseIcon, SmallAddIcon } from "@chakra-ui/icons";
import {
  Button,
  ButtonGroup,
  Circle,
  CircularProgress,
  Code,
  Container,
  Flex,
  HStack,
  IconButton,
  LinkBox,
  LinkOverlay,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  Stack,
  Text,
  useDisclosure,
} from "@chakra-ui/react";
import { useMemo, useState } from "react";
import { Helmet } from "react-helmet";
import { Link } from "react-location";
import { Column } from "react-table";
import { AdminLayout } from "../../../components/Layout";
import { TabsStyledButton } from "../../../components/nav/Navbar";
import { StatusCell } from "../../../components/StatusCell";
import { TableRenderer } from "../../../components/tables/TableRenderer";
import {
  adminDeleteProvidersetup,
  useAdminListProvidersetups,
} from "../../../utils/backend-client/admin/admin";
import { useListTargetGroupDeployments } from "../../../utils/backend-client/target-groups/target-groups";
import {
  Provider,
  ProviderSetup,
  TargetGroupDeployment,
} from "../../../utils/backend-client/types";
import { usePaginatorApi } from "../../../utils/usePaginatorApi";

// using a chakra tab component and links, link to /admin/providers and /admin/providersv2
export const ProvidersV2Tabs = () => {
  return (
    <ButtonGroup variant="ghost" spacing="0" mb={"-32px !important;"} my={4}>
      <TabsStyledButton href="/admin/providers">Legacy</TabsStyledButton>
      <TabsStyledButton href="/admin/providersv2">V2</TabsStyledButton>
    </ButtonGroup>
  );
};

const AdminProvidersTable = () => {
  const paginator = usePaginatorApi<typeof useListTargetGroupDeployments>({
    swrHook: useListTargetGroupDeployments,
    hookProps: {},
  });

  const cols: Column<TargetGroupDeployment>[] = useMemo(
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
    ],
    []
  );

  return TableRenderer<TargetGroupDeployment>({
    columns: cols,
    data: paginator?.data?.res,
    emptyText: "No providers have been set up yet.",
    linkTo: false,
    apiPaginator: paginator,
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
        <div style={{ height: "32px" }} />
        <ProvidersV2Tabs />
        <AdminProvidersTable />
        {/* <HStack mt={2} spacing={1} w="100%" justify={"center"}>
          <Text textStyle={"Body/ExtraSmall"}>
            View the full configuration of each access provider in your{" "}
          </Text>
          <Code fontSize={"12px"}>deployment.yml</Code>
          <Text textStyle={"Body/ExtraSmall"}>file.</Text>
        </HStack> */}
      </Container>
    </AdminLayout>
  );
};

export default Providers;
