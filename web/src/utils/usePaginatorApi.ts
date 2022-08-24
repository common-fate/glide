import { Dispatch, SetStateAction, useEffect, useMemo, useState } from "react";

type InputParams<T extends (...args: any[]) => any> = {
  pageSize?: number | undefined;
  /** This will be called within usePaginatorAPI, the resulting nextToken will be extracted */
  swrHook: T;
  hookProps: Parameters<T>[0];
};

export type PaginationProps<T extends (...args: any[]) => any> = {
  /** The pageSize that is passed to usePaginator (default value is 5) */
  pageSize: number;
  /** Numerical index representing the page i.e. 1 of 1,2,3 */
  pageIndex: number;
  /** a number representing the last item in a list of rows, this should be the pageSize for the first query */
  /** stateful array of all page indexes */
  pageOptions: number[];
  /** boolean representing whether or not the next page is available */
  canNextPage: boolean;
  /** boolean representing whether or not the prev page is available */
  canPrevPage: boolean;

  //   mutating methods
  setPageIndex: Dispatch<SetStateAction<number>>;
  selectPage: (pageIndex: number) => void;
  incrementPage: () => void;
  decrementPage: () => void;

  // Could also extend this to return isValidate etc.
  data: ReturnType<T>["data"] | undefined;
};

/**
 * The goal for this paginator hook is to pass state between the
 * SWR hook and the TableRender component.
 */
export const usePaginatorApi = <T extends (...args: any[]) => any>({
  pageSize,
  swrHook,
  hookProps,
}: InputParams<T>): PaginationProps<T> => {
  const [pageIndex, setPageIndex] = useState(0);
  const pageSizeOrDefault = pageSize || 5;

  const [nextToken, setNextToken] = useState<string | undefined>();

  const [tokenStack, setTokenStack] = useState<(string | undefined)[]>([
    undefined,
  ]);

  const { data } = swrHook({
    ...hookProps,
    nextToken: nextToken,
  });

  const [pageOptions, setPageOptions] = useState<number[]>([pageIndex]);

  useEffect(() => {
    if (data?.next && pageOptions?.length === 1) {
      setPageOptions((curr) => [...curr, pageIndex + 1]);
      setTokenStack((curr) => [...curr, data.next]);
    }
  }, [data, pageOptions, pageIndex]);

  const incrementPage = () => {
    // For adding a new page to the stack
    if (!pageOptions.includes(pageIndex + 1)) {
      console.log("new page adding to stack", pageIndex);
      setPageOptions((curr) => [...curr, pageIndex + 1]);
      setPageIndex((curr) => curr + 1);
      setTokenStack((curr) => [...curr, data.next]);
      setNextToken(data.next);
      // Otherwise incrementing internally between pages
    } else {
      const tokenIfPresent = tokenStack[pageIndex + 1];
      if (tokenIfPresent) {
        setNextToken(tokenStack[pageIndex + 1]);
      } else {
        setNextToken(data.next);
      }
      setPageIndex((curr) => curr + 1);
    }
  };

  const decrementPage = () => {
    setPageIndex((curr) => curr - 1);
    if (pageIndex === 1) {
      setNextToken(undefined);
    } else {
      setNextToken(
        tokenStack[tokenStack.findIndex((t) => t === nextToken) - 1]
      );
    }
  };

  const selectPage = (pageIndex: number) => {
    setPageIndex(pageIndex);
    setNextToken(tokenStack[pageIndex]);
    if (pageIndex == 0) {
      setNextToken(undefined);
    }
  };

  //  canNextPage support
  const canNextPage = useMemo(() => !(pageOptions?.length > 0 && !data?.next), [
    data,
  ]);

  //  canPrevPage support
  const canPrevPage = useMemo(() => pageIndex != 0, [pageIndex]);

  return {
    data,
    pageSize: pageSizeOrDefault,
    pageIndex,
    pageOptions,
    setPageIndex,
    canPrevPage,
    canNextPage,
    selectPage,
    incrementPage,
    decrementPage,
  };
};
