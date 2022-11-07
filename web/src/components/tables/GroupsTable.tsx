import { SmallAddIcon } from "@chakra-ui/icons";
import { Box, Button, Flex, Text, useDisclosure } from "@chakra-ui/react";
import { useMemo } from "react";
import { Column } from "react-table";
import {
  useIdentityConfiguration,
  useListGroups,
} from "../../utils/backend-client/admin/admin";

import {
  Group,
  ListGroupsSource,
  RequestStatus,
} from "../../utils/backend-client/types";
import { usePaginatorApi } from "../../utils/usePaginatorApi";
import CreateGroupModal from "../modals/CreateGroupModal";
import { TableRenderer } from "./TableRenderer";
import { MakeGenerics, useSearch, useNavigate } from "react-location";
import { GroupsFilterMenu } from "./GroupsFilterMenu";
import {
  ApprovalsLogo,
  CognitoLogo,
  GoogleLogo,
  OneLoginLogo,
} from "../icons/Logos";
import { SyncUsersAndGroupsButton } from "../SyncUsersAndGroupsButton";
import { AWSIcon, AzureIcon, GrantedKeysIcon, OktaIcon } from "../icons/Icons";
import { GetIDPLogo } from "../../utils/idp-logo";

type MyLocationGenerics = MakeGenerics<{
  Search: {
    source?: Lowercase<ListGroupsSource>;
  };
}>;

export const GroupsTable = () => {
  const search = useSearch<MyLocationGenerics>();
  const navigate = useNavigate<MyLocationGenerics>();
  const { source } = search;

  console.log(source);

  const { onOpen, isOpen, onClose } = useDisclosure();
  const paginator = usePaginatorApi<typeof useListGroups>({
    swrHook: useListGroups,
    hookProps: {
      source: source ? (source.toUpperCase() as ListGroupsSource) : undefined,
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
          <Box>{GetIDPLogo({ idpType: cell.value, size: 30 })}</Box>
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
                source: s?.toLowerCase() as Lowercase<ListGroupsSource>,
              }),
            })
          }
          source={source?.toUpperCase() as ListGroupsSource}
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
