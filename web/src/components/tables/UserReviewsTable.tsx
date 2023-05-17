import { Flex } from "@chakra-ui/react";
import format from "date-fns/format";
import { useMemo } from "react";
import { MakeGenerics, useNavigate, useSearch } from "react-location";
import { Column } from "react-table";
import {
  Request,
  UserListReviewsStatus,
} from "../../utils/backend-client/types";
import { usePaginatorApi } from "../../utils/usePaginatorApi";

import { useUserListReviews } from "../../utils/backend-client/default/default";
import { CFAvatar } from "../CFAvatar";
import { RequestsFilterMenu } from "./RequestsFilterMenu";
import { TableRenderer } from "./TableRenderer";

type MyLocationGenerics = MakeGenerics<{
  Search: {
    status?: Lowercase<UserListReviewsStatus>;
  };
}>;

export const UserReviewsTable = () => {
  const search = useSearch<MyLocationGenerics>();
  const navigate = useNavigate<MyLocationGenerics>();

  const { status } = search;

  // const [status, setStatus] = useState<UserListReviewsStatus | undefined>();

  const paginator = usePaginatorApi<typeof useUserListReviews>({
    swrHook: useUserListReviews,
    hookProps: {
      status: status
        ? (status.toUpperCase() as UserListReviewsStatus)
        : undefined,
    },
    swrProps: {
      // @ts-ignore; type discrepancy with latest SWR client
      swr: { refreshInterval: 10000 },
    },
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
        Header: "Date Requested",
        Cell: ({ cell }) => (
          <Flex textStyle="Body/Small">
            {format(new Date(Date.parse(cell.value)), "p dd/M/yy")}
          </Flex>
        ),
      },
      // {
      //   accessor: "status",
      //   Header: "Status",
      //   Cell: (props) => {
      //     return <UserListReviewsStatusDisplay request={props.row.original} />;
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
                status: s?.toLowerCase() as Lowercase<UserListReviewsStatus>,
              }),
            })
          }
          status={status?.toUpperCase() as UserListReviewsStatus}
        />
      </Flex>
      {TableRenderer<Request>({
        columns: cols,
        data: paginator?.data?.requests,
        apiPaginator: paginator,
        emptyText: "🎉 No outstanding reviews",
        rowProps: (row) => ({
          "_hover": { bg: "gray.50" },
          "cursor": "pointer",
          // in our test cases we use reason for the unique key
          "data-testid": row.original.purpose.reason,
          "alignItems": "center",
          "onClick": () => {
            navigate({ to: "/requests/" + row.original.id });
          },
        }),
      })}
    </>
  );
};
