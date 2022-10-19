import { Button, ButtonGroup, Flex } from "@chakra-ui/react";
import format from "date-fns/format";
import { useMemo } from "react";
import { Link, MakeGenerics, useNavigate, useSearch } from "react-location";
import { Column } from "react-table";
import { usePaginatorApi } from "../../utils/usePaginatorApi";
import { useUserListRequests } from "../../utils/backend-client/end-user/end-user";
import { Request, RequestStatus } from "../../utils/backend-client/types";
import { durationString } from "../../utils/durationString";
import { RuleNameCell } from "../AccessRuleNameCell";
import { RequestStatusDisplay } from "../Request";
import { UserAvatarDetails } from "../UserAvatar";
import { RequestsFilterMenu } from "./RequestsFilterMenu";
import { TableRenderer } from "./TableRenderer";

type MyLocationGenerics = MakeGenerics<{
  Search: {
    status?: Lowercase<RequestStatus>;
  };
}>;

export const UserReviewsTable = () => {
  const search = useSearch<MyLocationGenerics>();
  const navigate = useNavigate<MyLocationGenerics>();

  const { status } = search;

  // const [status, setStatus] = useState<RequestStatus | undefined>();

  const paginator = usePaginatorApi<typeof useUserListRequests>({
    swrHook: useUserListRequests,
    hookProps: {
      reviewer: true,
      status: status ? (status.toUpperCase() as RequestStatus) : undefined,
    },
  });

  const cols: Column<Request>[] = useMemo(
    () => [
      {
        accessor: "reason",
        Header: "Request",
        Cell: (props) => (
          <Link to={"/requests/" + props.row.original.id}>
            <RuleNameCell
              accessRuleId={props.row.original.accessRuleId}
              reason={props.value ?? ""}
              as="a"
              _hover={{
                textDecor: "underline",
              }}
              adminRoute={false}
            />
          </Link>
        ),
      },
      {
        accessor: "timing",
        Header: "Duration",
        Cell: ({ cell }) => (
          <Flex textStyle="Body/Small">
            {durationString(cell.value.durationSeconds)}
          </Flex>
        ),
      },
      {
        accessor: "requestor",
        Header: "Requested by",
        Cell: ({ cell }) => (
          <Flex textStyle="Body/Small">
            <UserAvatarDetails
              textProps={{
                maxW: "20ch",
                noOfLines: 1,
              }}
              tooltip={true}
              variant="withBorder"
              mr={0}
              size="xs"
              user={cell.value}
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
      {
        accessor: "status",
        Header: "Status",
        Cell: (props) => {
          return <RequestStatusDisplay request={props.row.original} />;
        },
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
                status: s?.toLowerCase() as Lowercase<RequestStatus>,
              }),
            })
          }
          status={status?.toUpperCase() as RequestStatus}
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
          "data-testid": row.original.reason,
          "alignItems": "center",
          "onClick": () => {
            navigate({ to: "/requests/" + row.original.id });
          },
        }),
      })}
    </>
  );
};
