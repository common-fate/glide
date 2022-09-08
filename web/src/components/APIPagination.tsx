import {
  ButtonProps,
  chakra,
  SystemStyleObject,
  HStack,
} from "@chakra-ui/react";
import { TableInstance, UsePaginationInstanceProps } from "react-table";
import { PaginationProps } from "../utils/usePaginatorApi";

export const APIPagination = ({
  paginator,
  ...props
}: {
  paginator?: PaginationProps<any>;
} & UsePaginationInstanceProps<Record<string, never>> &
  TableInstance<any>) => {
  return (
    <HStack spacing={2}>
      <PaginationButton
        disabled={!paginator?.canPrevPage}
        cursor={paginator?.canPrevPage ? "pointer" : "not-allowed"}
        onClick={() => {
          props.previousPage();
          paginator?.decrementPage();
        }}
      >
        {"<"}
      </PaginationButton>
      {paginator?.pageOptions.map((el) => (
        <PaginationButton
          key={el}
          // @ts-ignore, unfortunately these props aren't carrying through too well
          aria-current={paginator.pageIndex == el ? "step" : "false"}
          onClick={() => {
            paginator?.setPageIndex(el);
            paginator?.selectPage(el);
          }}
        >
          {el + 1}
        </PaginationButton>
      ))}
      {paginator?.canNextPage && (
        <chakra.span color="neutrals.600">...</chakra.span>
      )}
      <PaginationButton
        disabled={!paginator?.canNextPage}
        cursor={paginator?.canNextPage ? "pointer" : "not-allowed"}
        onClick={() => {
          props.nextPage();
          paginator?.incrementPage();
        }}
      >
        {">"}
      </PaginationButton>
    </HStack>
  );
};

export const PaginationButton = (props: ButtonProps) => {
  const commonStyle: SystemStyleObject = {
    bg: "neutrals.200",
  };

  const activeStyle: SystemStyleObject = {
    bg: "brandGreen.100",
    color: "white",
  };

  return (
    <chakra.button
      sx={{
        px: 2,
        py: 1,
        rounded: "4px",
        textStyle: "Body/Small",
        color: "neutrals.600",
        _hover: commonStyle,
        _focus: {
          boxShadow: "outline",
          border: "none",
          outline: "none",
        },
        _activeStep: activeStyle,
      }}
      {...props}
    />
  );
};
