import { SmallAddIcon } from "@chakra-ui/icons";
import {
  Box,
  Button,
  Flex,
  Text,
  Tooltip,
  useDisclosure,
} from "@chakra-ui/react";
import { useMemo } from "react";
import { Column } from "react-table";
import {
  useAdminGetIdentityConfiguration,
  useAdminListGroups,
} from "../../utils/backend-client/admin/admin";

import { MakeGenerics, useNavigate, useSearch } from "react-location";
import { Group, AdminListGroupsSource } from "../../utils/backend-client/types";
import { IDPLogo } from "../../utils/idp-logo";
import { usePaginatorApi } from "../../utils/usePaginatorApi";
import CreateGroupModal from "../modals/CreateGroupModal";
import { SyncUsersAndGroupsButton } from "../SyncUsersAndGroupsButton";
import { GroupsFilterMenu } from "./GroupsFilterMenu";
import { TableRenderer } from "./TableRenderer";

type MyLocationGenerics = MakeGenerics<{
  Search: {
    source?: Lowercase<AdminListGroupsSource>;
  };
}>;

export const GroupsTable = () => {
  const search = useSearch<MyLocationGenerics>();
  const navigate = useNavigate<MyLocationGenerics>();
  const { source } = search;

  const { onOpen, isOpen, onClose } = useDisclosure();
  const paginator = usePaginatorApi<typeof useAdminListGroups>({
    swrHook: useAdminListGroups,
    hookProps: {
      source: source
        ? (source.toUpperCase() as AdminListGroupsSource)
        : undefined,
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
        Cell: ({ cell }) => <IDPLogo idpType={cell.value} boxSize={30} />,
      },
    ],
    []
  );
  const { data } = useAdminGetIdentityConfiguration();
  const AddGroupButton = () => {
    return (
      <Tooltip
        label={
          "Internal groups allows you to configure more granular access policies that may not be possible with your existing identity provider groups."
        }
      >
        <Button
          isLoading={data?.identityProvider === undefined}
          size="sm"
          variant="ghost"
          leftIcon={<SmallAddIcon />}
          onClick={onOpen}
          data-testid={"create-group-button"}
        >
          Add Internal Group
        </Button>
      </Tooltip>
    );
  };
  return (
    <>
      <Flex justify="space-between" my={5}>
        <AddGroupButton />
        <SyncUsersAndGroupsButton
          onSync={() => {
            void paginator.mutate();
          }}
        />
        <GroupsFilterMenu
          onChange={(s) =>
            navigate({
              search: (old) => ({
                ...old,
                source: s?.toLowerCase() as Lowercase<AdminListGroupsSource>,
              }),
            })
          }
          source={source?.toUpperCase() as AdminListGroupsSource}
        />
      </Flex>
      {TableRenderer<Group>({
        columns: cols,
        data: paginator?.data?.groups,
        emptyText: "No groups",
        apiPaginator: paginator,
        linkTo: true,
        rowProps: (row) => ({
          // in our test cases we use the group name for the unique key
          "data-testid": row.original.name,
          "alignItems": "center",
        }),
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
