import { Button, Flex, useDisclosure } from "@chakra-ui/react";
import format from "date-fns/format";
import { useMemo } from "react";
import { Column } from "react-table";
import { usePaginatorApi } from "../../utils/usePaginatorApi";
import { useGetUsers } from "../../utils/backend-client/admin/admin";
import { User } from "../../utils/backend-client/types";
import { TableRenderer } from "./TableRenderer";
import { SmallAddIcon } from "@chakra-ui/icons";
import CreateUserModal from "../modals/CreateUserModal";
import { SyncUsersAndGroupsButton } from "../SyncUsersAndGroupsButton";

export const UsersTable = () => {
  const { onOpen, isOpen, onClose } = useDisclosure();
  const paginator = usePaginatorApi<typeof useGetUsers>({
    swrHook: useGetUsers,
    hookProps: {},
  });
  const { mutate } = useGetUsers();
  const cols: Column<User>[] = useMemo(
    () => [
      {
        accessor: "firstName",
        Header: "First Name", // blank
        Cell: ({ cell }) => <Flex textStyle="Body/Small">{cell.value}</Flex>,
      },
      {
        accessor: "lastName",
        Header: "Last Name", // blank
        Cell: ({ cell }) => <Flex textStyle="Body/Small">{cell.value}</Flex>,
      },
      {
        accessor: "email",
        Header: "Email", // blank
        Cell: ({ cell }) => <Flex textStyle="Body/Small">{cell.value}</Flex>,
      },
      {
        accessor: "updatedAt",
        Header: "Last Updated",
        Cell: ({ cell }) => (
          <Flex textStyle="Body/Small">
            {format(new Date(cell.value), "p dd/M/yy")}
          </Flex>
        ),
      },
      // {
      //   accessor: "picture",
      //   Header: "", // blank
      //   Cell: ({ cell }) => (
      //     <Button variant="outline" rounded="full" size="xs">
      //       Preview Access
      //     </Button>
      //   ),
      // },
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
          Add User
        </Button>
        <SyncUsersAndGroupsButton
          onSync={() => {
            mutate();
          }}
        />
      </Flex>
      {TableRenderer<User>({
        columns: cols,
        data: paginator?.data?.users,
        emptyText: "No users",
        apiPaginator: paginator,
      })}

      <CreateUserModal isOpen={isOpen} onClose={onClose} />
    </>
  );
};
