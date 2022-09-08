import { Box, Button, Flex, Text, useDisclosure } from "@chakra-ui/react";
import { useMemo, useState } from "react";
import { Column } from "react-table";
import { usePaginatorApi } from "../../utils/usePaginatorApi";
import { useGetGroups } from "../../utils/backend-client/admin/admin";
import { Group } from "../../utils/backend-client/types";
import GroupModal from "../modals/GroupModal";
import { TableRenderer } from "./TableRenderer";
import { SmallAddIcon } from "@chakra-ui/icons";

export const GroupsTable = () => {
  const { onOpen, isOpen, onClose } = useDisclosure();
  const paginator = usePaginatorApi<typeof useGetGroups>({
    swrHook: useGetGroups,
    hookProps: {},
  });

  const [selectedGroup, setSelectedGroup] = useState<Group>();

  const cols: Column<Group>[] = useMemo(
    () => [
      {
        accessor: "name",
        Header: "Name", // blank
        Cell: ({ cell }) => (
          <Box>
            <Text color="neutrals.900">{cell.value}</Text>
          </Box>
        ),
      },
    ],
    []
  );

  return (
    <>
      <Flex justify="space-between" my={5}>
        <Button
          size="sm"
          variant="ghost"
          leftIcon={<SmallAddIcon />}
          onClick={onOpen}
        >
          Add Group
        </Button>
      </Flex>
      {TableRenderer<Group>({
        columns: cols,
        data: paginator?.data?.groups,
        emptyText: "No groups",
        apiPaginator: paginator,
      })}

      <GroupModal
        isOpen={selectedGroup !== undefined}
        onClose={() => setSelectedGroup(undefined)}
        group={selectedGroup}
        members={[]}
      />
    </>
  );
};
