import { Button, Flex, useDisclosure } from "@chakra-ui/react";
import format from "date-fns/format";
import { useMemo } from "react";
import { Column } from "react-table";
import { usePaginatorApi } from "../../utils/usePaginatorApi";
import {
  useAdminListUsers,
  useAdminGetIdentityConfiguration,
} from "../../utils/backend-client/admin/admin";
import { User } from "../../utils/backend-client/types";
import { TableRenderer } from "./TableRenderer";
import { SmallAddIcon } from "@chakra-ui/icons";
import CreateUserModal from "../modals/CreateUserModal";
import { SyncUsersAndGroupsButton } from "../SyncUsersAndGroupsButton";

export const UsersTable = () => {
  const { onOpen, isOpen, onClose } = useDisclosure();

  const paginator = usePaginatorApi<typeof useAdminListUsers>({
    swrHook: useAdminListUsers,
    hookProps: {},
  });
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
  const { data } = useAdminGetIdentityConfiguration();
  const AddUsersButton = () => {
    if (data?.identityProvider !== "cognito") {
      return <div />;
    }
    return (
      <Button
        isLoading={data?.identityProvider === undefined}
        size="sm"
        variant="ghost"
        leftIcon={<SmallAddIcon />}
        onClick={onOpen}
      >
        Add User
      </Button>
    );
  };

  return (
    <>
      <Flex justify="space-between" my={5}>
        <AddUsersButton />
        <SyncUsersAndGroupsButton
          onSync={() => {
            void paginator.mutate();
          }}
        />
      </Flex>
      {TableRenderer<User>({
        columns: cols,
        data: paginator?.data?.users,
        emptyText: "No users",
        apiPaginator: paginator,
        linkTo: true,
        rowProps: (row) => ({
          // in our test cases we use the email for the unique key
          "data-testid": row.original.email,
          "alignItems": "center",
        }),
      })}

      <CreateUserModal
        isOpen={isOpen}
        onClose={() => {
          void paginator.mutate();
          onClose();
        }}
      />
    </>
  );
};
