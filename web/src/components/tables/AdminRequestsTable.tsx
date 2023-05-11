import { Flex } from "@chakra-ui/react";
import format from "date-fns/format";
import { useMemo } from "react";
import { MakeGenerics, useNavigate, useSearch } from "react-location";
import { Column } from "react-table";
import { useAdminListRequests } from "../../utils/backend-client/admin/admin";
import {
  RequestStatus,
  Request,
  AdminListRequestsStatus,
} from "../../utils/backend-client/types";
import { usePaginatorApi } from "../../utils/usePaginatorApi";
import { CFAvatar } from "../CFAvatar";
import { RequestsFilterMenu } from "./RequestsFilterMenu";
import { TableRenderer } from "./TableRenderer";

type MyLocationGenerics = MakeGenerics<{
  Search: {
    status?: Lowercase<AdminListRequestsStatus>;
  };
}>;

export const AdminRequestsTable = () => {
  const search = useSearch<MyLocationGenerics>();
  const navigate = useNavigate<MyLocationGenerics>();

  const { status } = search;

  const paginator = usePaginatorApi<typeof useAdminListRequests>({
    swrHook: useAdminListRequests,
    hookProps: {
      status: status
        ? (status.toUpperCase() as AdminListRequestsStatus)
        : undefined,
    },
    swrProps: { swr: { refreshInterval: 10000 } },
  });

  const cols: Column<Request>[] = useMemo(
    () => [
      // {
      //   accessor: "context",
      //   Header: "Request",
      //   // Cell: (props) => (
      //   //   <Link to={"/requests/" + props.row.original.id}>
      //   //     <RuleNameCell
      //   //       accessRuleId={props.row.original.accessRuleId}
      //   //       reason={props.value ?? ""}
      //   //       as="a"
      //   //       _hover={{
      //   //         textDecor: "underline",
      //   //       }}
      //   //       adminRoute={false}
      //   //     />
      //   //   </Link>
      //   // ),
      // },
      // {
      //   accessor: "timing",
      //   Header: "Duration",
      //   Cell: ({ cell }) => (
      //     <Flex textStyle="Body/Small">
      //       {durationString(cell.value.durationSeconds)}
      //     </Flex>
      //   ),
      // },
      // {
      //   accessor: "user",
      //   Header: "Requested by",
      //   Cell: ({ cell }) => (
      //     <Flex textStyle="Body/Small">
      //       <CFAvatar
      //         textProps={{
      //           maxW: "20ch",
      //           noOfLines: 1,
      //         }}
      //         tooltip={true}
      //         variant="withBorder"
      //         mr={0}
      //         size="xs"
      //         userId={cell.value.id}
      //       />
      //     </Flex>
      //   ),
      // },
      // {
      //   accessor: "createdAt",
      //   Header: "Date Requested",
      //   Cell: ({ cell }) => (
      //     <Flex textStyle="Body/Small">
      //       {format(new Date(Date.parse(cell.value)), "p dd/M/yy")}
      //     </Flex>
      //   ),
      // },
      // {
      //   accessor: "status",
      //   Header: "Status",
      //   Cell: (props) => {
      //     return <RequestStatusDisplay request={props.row.original} />;
      //   },
      // },
    ],
    []
  );

  return (
    <>
      <Flex justify="flex-end" my={5}>
        <RequestsFilterMenu
          onChange={(s) =>
            navigate({
              search: (old) => ({
                ...old,
                status: s?.toLowerCase() as Lowercase<AdminListRequestsStatus>,
              }),
            })
          }
          status={status?.toUpperCase() as RequestStatus}
        />
      </Flex>
      {TableRenderer<Request>({
        columns: cols,
        data: paginator?.data?.requests,
        emptyText: "No requests",
        apiPaginator: paginator,
        linkTo: true,
      })}
    </>
  );
};
