import { SmallAddIcon } from "@chakra-ui/icons";
import { Box, Button, Flex, Text, useDisclosure } from "@chakra-ui/react";
import { useMemo } from "react";
import { Column } from "react-table";
import {
  useGetGroups,
  useIdentityConfiguration,
} from "../../utils/backend-client/admin/admin";

import { Group, IdpStatus } from "../../utils/backend-client/types";
import { usePaginatorApi } from "../../utils/usePaginatorApi";
import CreateGroupModal from "../modals/CreateGroupModal";
import { TableRenderer } from "./TableRenderer";

export const ManagementGroupsTable = () => {
  const { onOpen, isOpen, onClose } = useDisclosure();
  const paginator = usePaginatorApi<typeof useGetGroups>({
    swrHook: useGetGroups,
    hookProps: { status: IdpStatus.INTERNAL },
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
      </Flex>
      {TableRenderer<Group>({
        columns: cols,
        data: paginator?.data?.groups,
        emptyText: "No groups",
        apiPaginator: paginator,
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
