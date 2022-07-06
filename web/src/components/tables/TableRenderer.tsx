import {
  Box,
  Center,
  Flex,
  HStack,
  Skeleton,
  SkeletonText,
  Stack,
  Table,
  TableRowProps,
  Tbody,
  Td,
  Th,
  Thead,
  Tr,
} from "@chakra-ui/react";
import React, { useCallback } from "react";
import { Link, useNavigate } from "react-location";
import {
  Column,
  Row,
  TableInstance,
  TableOptions,
  usePagination,
  UsePaginationInstanceProps,
  UsePaginationOptions,
  useTable,
} from "react-table";
import { PaginationProps } from "../../utils/usePaginatorApi";

import { APIPagination } from "../APIPagination";
import { Pagination } from "../Pagination";
interface TableRendererProps<T extends object> {
  data: T[] | undefined;
  columns: Column<T>[];
  emptyText: string;
  rowProps?: (row: Row<T>) => TableRowProps;
  apiPaginator?: PaginationProps<any>;
  linkTo?: boolean;
}

export function TableRenderer<T extends object>(
  props: React.PropsWithChildren<TableRendererProps<T>>
) {
  if (!props.data) {
    return (
      <Table roundedTop="lg">
        <Thead h="45px">
          <Tr bg="neutrals.100" h="45px">
            {props.columns.map((column) =>
              column?.Header ? (
                <Th
                  h="45px"
                  fontFamily="Rubik"
                  color="neutrals.700"
                  fontWeight={400}
                  textTransform="none"
                  _first={{
                    roundedTopLeft: "xl",
                  }}
                  _last={{
                    roundedTopRight: "xl",
                  }}
                  key={"column-" + column.id}
                >
                  {column?.Header.toString()}
                </Th>
              ) : (
                <Th key={column.id} h="45px" />
              )
            )}
          </Tr>
        </Thead>
        <Tbody>
          {[1, 2, 3, 4, 5, 6, 7, 8].map((s) => (
            <Tr key={s}>
              {props.columns.map((i) =>
                i.Header ? (
                  <Td
                    // maxHeight="32px"
                    py="10px"
                    key={i.id}
                  >
                    <SkeletonText noOfLines={1} w="12ch" h="20px" />
                  </Td>
                ) : (
                  <Td py="10px" key={i.id} />
                )
              )}
            </Tr>
          ))}
        </Tbody>
      </Table>
    );
  }

  return (
    <_TableRenderer
      columns={props.columns}
      data={props.data}
      emptyText={props.emptyText}
      rowProps={props.rowProps}
      paginator={props.apiPaginator}
      linkTo={props.linkTo}
    />
  );
}

interface _TableRendererProps {
  data: any;
  columns: any;
  emptyText: string;
  rowProps?: (row: Row<any>) => TableRowProps;
  paginator?: PaginationProps<any>;
  linkTo?: boolean;
}
export const _TableRenderer: React.FC<_TableRendererProps> = ({
  rowProps,
  paginator,
  ...props
}) => {
  // https://react-table-v7.tanstack.com/docs/api/usePagination#instance-properties
  const paginatorProps: TableOptions<any> = paginator
    ? {
        initialState: {
          pageIndex: paginator?.pageIndex,
          pageSize: paginator?.pageSize,
          pageOptions: paginator?.pageOptions,
        },
        // @ts-ignore, this one type wont resolve
        manualPagination: true,
      }
    : {};

  const instance = useTable<any>(
    {
      ...paginatorProps,
      data: props.data,
      columns: props.columns,
    },
    usePagination
  ) as TableInstance<any> & UsePaginationInstanceProps<any>;

  const navigate = useNavigate();

  return (
    <Box overflowX="auto">
      <Table
        {...instance.getTableProps()}
        // borderCollapse="separate !important"
        roundedTop="lg"
      >
        <Thead h="45px">
          {instance.headerGroups.map((headerGroup) => (
            <Tr
              {...headerGroup.getHeaderGroupProps()}
              bg="neutrals.100"
              h="45px"
              key={"headergroup-" + headerGroup.id}
              // borderCollapse="separate !important"
              // roundedTop="48px"
            >
              {headerGroup.headers.map((column) => (
                <Th
                  {...column.getHeaderProps()}
                  h="45px"
                  fontFamily="Rubik"
                  color="neutrals.700"
                  fontWeight={400}
                  textTransform="none"
                  // borderCollapse="separate !important"
                  // roundedTop="md"
                  _first={{
                    roundedTopLeft: "xl",
                  }}
                  _last={{
                    roundedTopRight: "xl",
                  }}
                  key={"column-" + column.id}
                  // isNumeric={column}
                >
                  {column.render("Header")}

                  {/* <chakra.span pl="4">
                {column.isSorted ? (
                  column.isSortedDesc ? (
                    <TriangleDownIcon aria-label="sorted descending" />
                  ) : (
                    <TriangleUpIcon aria-label="sorted ascending" />
                  )
                ) : null}
              </chakra.span> */}
                </Th>
              ))}
            </Tr>
          ))}
        </Thead>
        <Tbody {...instance.getTableBodyProps()}>
          {instance?.page.map((row, i) => {
            instance.prepareRow(row);

            let extraProps = rowProps ? rowProps(row) : {};

            //optional prop to turn clickable rows on or off
            if (props.linkTo) {
              return (
                <Tr
                  cursor="pointer"
                  {...extraProps}
                  userSelect="all"
                  key={"tablerow-" + i}
                  onClick={(e) => {
                    e.preventDefault();

                    navigate({ to: row.original.id });
                  }}
                  _hover={{
                    backgroundColor: "neutrals.200",
                  }}
                >
                  {row.cells.map((cell, i) => (
                    <Td
                      {...cell.getCellProps()}
                      // maxHeight="32px"
                      py="10px"
                      key={"rowcell-" + i}
                      // isNumeric={cell.column.isNumeric}
                    >
                      {cell.render("Cell")}
                    </Td>
                  ))}
                </Tr>
              );
            } else {
              return (
                <Tr {...extraProps} userSelect="all" key={"tablerow-" + i}>
                  {row.cells.map((cell, i) => (
                    <Td
                      {...cell.getCellProps()}
                      // maxHeight="32px"
                      py="10px"
                      key={"rowcell-" + i}
                      // isNumeric={cell.column.isNumeric}
                    >
                      {cell.render("Cell")}
                    </Td>
                  ))}
                </Tr>
              );
            }
          })}
        </Tbody>
      </Table>

      {instance.page.length === 0 && (
        /* empty state  */
        <Center mt={8} h="200px">
          {props.emptyText}
        </Center>
      )}
      {(instance.data.length > 0 ||
        (paginator && paginator.pageOptions?.length > 0)) &&
      paginator ? (
        <Center my={6}>
          <APIPagination paginator={paginator} {...instance} />
        </Center>
      ) : (
        <Center my={6}>
          <Pagination {...instance} />
        </Center>
      )}
    </Box>
  );
};
