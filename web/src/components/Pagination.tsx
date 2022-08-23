import { ButtonProps, chakra, CSSObject, HStack } from "@chakra-ui/react";
import React from "react";
import {
  TableInstance,
  TableState,
  UsePaginationInstanceProps,
} from "react-table";

export const Pagination = (
  props: UsePaginationInstanceProps<Record<string, never>> & TableInstance<any>
) => {
  return (
    <HStack spacing={2}>
      <PaginationButton
        disabled={!props.canPreviousPage}
        cursor={props.canPreviousPage ? "pointer" : "not-allowed"}
        onClick={() => props.previousPage()}
      >
        {"<"}
      </PaginationButton>
      {props?.pageOptions.map((el) => (
        <PaginationButton
          key={el}
          // @ts-ignore, unfortunately these props aren't carrying through too well
          aria-current={props?.state?.pageIndex == el ? "step" : "false"}
          onClick={() => props.gotoPage(el)}
        >
          {el + 1}
        </PaginationButton>
      ))}
      <PaginationButton
        disabled={!props.canNextPage}
        cursor={props.canNextPage ? "pointer" : "not-allowed"}
        onClick={() => props.nextPage()}
      >
        {">"}
      </PaginationButton>
    </HStack>
  );
};

export const PaginationButton = (props: ButtonProps) => {
  const commonStyle: CSSObject = {
    bg: "neutrals.200",
  };

  const activeStyle: CSSObject = {
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
