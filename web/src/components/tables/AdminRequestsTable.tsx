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
import { StatusCell } from "../StatusCell";

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
    swrProps: {
      // @ts-ignore; type discrepancy with latest SWR client
      swr: { refreshInterval: 10000 },
    },
  });

  const cols: Column<Request>[] = useMemo(
    () => [
      {
        accessor: "targetCount",
        Header: "Number of targets",
        Cell: ({ cell }) => <Flex textStyle="Body/Small">{cell.value}</Flex>,
      },
      {
        accessor: "purpose",
        Header: "Reason",
        Cell: ({ cell }) => (
          <Flex textStyle="Body/Small">{cell.value.reason}</Flex>
        ),
      },

      {
        accessor: "requestedBy",
        Header: "Requested by",
        Cell: ({ cell }) => (
          <Flex textStyle="Body/Small">
            <CFAvatar
              textProps={{
                maxW: "20ch",
                noOfLines: 1,
              }}
              tooltip={true}
              variant="withBorder"
              mr={0}
              size="xs"
              userId={cell.value.id}
            />
          </Flex>
        ),
      },
      {
        accessor: "requestedAt",
        Header: "Requested At",
        Cell: ({ cell }) => (
          // @NOTE: date resolution is currently not working, we can type these in OpenAPI, but ultimately it will come from BE
          <Flex textStyle="Body/Small">
            {format(new Date(Date.parse(cell.value)), "p dd/M/yy")}
          </Flex>
        ),
      },
      {
        accessor: "status",
        Header: "Status",
        Cell: ({ cell }) => (
          <StatusCell
            sx={{
              span: {
                textStyle: "Body/Small",
              },
            }}
            value={cell.value}
            replaceValue={cell.value.toLowerCase()}
            success={["COMPLETE", "ACTIVE"]}
            warning={["PENDING", "REVOKING"]}
            danger={["REVOKED", "CANCELLED"]}
          />
        ),
      },
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

// export const RequestStatusDisplay: React.FC<{
//   request: Request | undefined;
// }> = ({ request }) => {
//   const activeTimeString =
//     request?.grant && request?.grant.status === "ACTIVE"
//       ? "Active for the next " +
//         durationStringHoursMinutes(
//           intervalToDuration({
//             start: new Date(),
//             end: new Date(Date.parse(request.grant.end)),
//           })
//         )
//       : undefined;

//   const status = getStatus(request, activeTimeString);

//   return (
//     <StatusCell
//       value={status}
//       success={[
//         activeTimeString ?? "",
//         GrantStatus.ACTIVE,
//         "Automatically approved",
//       ]}
//       info={[RequestStatus.APPROVED]}
//       danger={[
//         RequestStatus.DECLINED,
//         RequestStatus.CANCELLED,
//         GrantStatus.REVOKED,
//       ]}
//       warning={RequestStatus.PENDING}
//       textStyle="Body/Small"
//     />
//   );
// };
