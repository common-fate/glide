import { Box, Flex, Text } from "@chakra-ui/react";
import { useMemo, useState } from "react";
import { Column } from "react-table";
import { useGetGroups } from "../../utils/backend-client/admin/admin";
import { Group } from "../../utils/backend-client/types";
import GroupModal from "../modals/GroupModal";
import { TableRenderer } from "./TableRenderer";

export const GroupsTable = () => {
  const { data } = useGetGroups();

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
      // {
      //   accessor: "users",
      //   Header: "Members",
      //   Cell: ({ cell }) => (
      //     <Flex>
      //       <AvatarGroup
      //         max={3}
      //         size="sm"
      //         spacing="-6px"
      //         sx={{
      //           ".chakra-avatar__excess": {
      //             fontWeight: "normal",
      //             bg: "white",
      //             border: "1px solid #E5E5E5",
      //             fontSize: "12px",
      //           },
      //         }}
      //       >
      //         {cell.value.map((approver) => (
      //           <Avatar key={approver.email} name={approver.name} />
      //         ))}
      //       </AvatarGroup>
      //       <Button
      //         variant="outline"
      //         size="sm"
      //         ml={2}
      //         rounded="full"
      //         onClick={() => {
      //           groupModal.onOpen();
      //           setSelectedGroup(cell.row.values);
      //         }}
      //       >
      //         See all members
      //       </Button>
      //     </Flex>
      //   ),
      // },
    ],
    []
  );

  return (
    <>
      <Flex justify="space-between" my={5}></Flex>
      {TableRenderer<Group>({
        columns: cols,
        data: data?.groups,
        emptyText: "No groups",
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
