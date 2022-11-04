import { SmallAddIcon } from "@chakra-ui/icons";
import { Box, Button, Flex, Text, useDisclosure } from "@chakra-ui/react";
import { useMemo } from "react";
import { Column } from "react-table";
import { useGetGroupBySource } from "../../utils/backend-client/default/default";
import { useIdentityConfiguration } from "../../utils/backend-client/admin/admin";

import {
  Group,
  GroupSource,
  RequestStatus,
} from "../../utils/backend-client/types";
import { usePaginatorApi } from "../../utils/usePaginatorApi";
import CreateGroupModal from "../modals/CreateGroupModal";
import { TableRenderer } from "./TableRenderer";
import { MakeGenerics, useSearch, useNavigate } from "react-location";
import { GroupsFilterMenu } from "./GroupsFilterMenu";
import { ApprovalsLogo } from "../../components/icons/Logos";
import { AzureIcon, OktaIcon } from "../icons/Icons";

type MyLocationGenerics = MakeGenerics<{
  Search: {
    source?: Lowercase<GroupSource>;
  };
}>;

export const ManagementGroupsTable = () => {
  const search = useSearch<MyLocationGenerics>();
  const navigate = useNavigate<MyLocationGenerics>();
  const { source } = search;

  console.log(source);

  const { onOpen, isOpen, onClose } = useDisclosure();
  const paginator = usePaginatorApi<typeof useGetGroupBySource>({
    swrHook: useGetGroupBySource,
    hookProps: {
      source: source ? (source.toUpperCase() as GroupSource) : undefined,
    },
    swrProps: {},
  });

  const cols: Column<Group>[] = useMemo(
    () => [
      {
        accessor: "name",
        Header: "Name",
        Cell: ({ cell }) => (
          <Box>
            <Text color="neutrals.900">{cell.value}</Text>
          </Box>
        ),
      },
      {
        accessor: "description",
        Header: "Description",
        Cell: ({ cell }) => (
          <Box>
            <Text color="neutrals.900">{cell.value}</Text>
          </Box>
        ),
      },
      {
        accessor: "memberCount",
        Header: "Members",
        Cell: ({ cell }) => (
          <Box>
            <Text color="neutrals.900">{cell.value}</Text>
          </Box>
        ),
      },
      {
        accessor: "source",
        Header: "",
        Cell: ({ cell }) => (
          <Box>
            {cell.value == "INTERNAL" && <ApprovalsLogo h="20px" w="auto" />}
            {cell.value == "AZURE" && <AzureIcon h="20px" w="auto" />}
            {cell.value == "ONELOGIN" && <OktaIcon h="20px" w="auto" />}
          </Box>
        ),
      },
    ],
    []
  );
  const { data } = useIdentityConfiguration();
  const AddGroupButton = () => {
    return (
      <Button
        isLoading={data?.identityProvider === undefined}
        size="sm"
        variant="ghost"
        leftIcon={<SmallAddIcon />}
        onClick={onOpen}
      >
        Add Group
      </Button>
    );
  };
  return (
    <>
      <Flex justify="space-between" my={5}>
        <AddGroupButton />
        <GroupsFilterMenu
          onChange={(s) =>
            navigate({
              search: (old) => ({
                ...old,
                source: s?.toLowerCase() as Lowercase<GroupSource>,
              }),
            })
          }
          source={source?.toUpperCase() as GroupSource}
        />
      </Flex>
      {TableRenderer<Group>({
        columns: cols,
        data: paginator?.data?.groups,
        emptyText: "No groups",
        apiPaginator: paginator,
        linkTo: true,
      })}

      <CreateGroupModal
        isOpen={isOpen}
        onClose={() => {
          void paginator.mutate();
          onClose();
        }}
      />
    </>
  );
};
