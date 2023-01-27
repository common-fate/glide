import { useEffect, useMemo, useState } from "react";

type InputParams<T extends (...args: any[]) => any> = {
  /** This will be called within useInfiniteScrollApi, the resulting nextToken will be extracted */
  swrHook: T;
  hookProps?: Parameters<T>[0];
  swrProps?: Parameters<T>[1];
  /** by specifying the listObjKey we know what array value to mutate from the returned data type (helps keep this method generic!)  */
  listObjKey: keyof Exclude<ReturnType<T>["data"], undefined>;
};

export type PaginationProps<T extends (...args: any[]) => any> = {
  /** boolean representing whether or not the next page is available */
  canNextPage: boolean;
  incrementPage: () => void;
  data: ReturnType<T>["data"] | undefined;
  isValidating: ReturnType<T>["isValidating"];
};

/**
 * useInfiniteScrollApi
 */
export const useInfiniteScrollApi = <T extends (...args: any[]) => any>({
  swrHook,
  hookProps,
  swrProps,
  listObjKey,
}: InputParams<T>): PaginationProps<T> => {
  const [nextToken, setNextToken] = useState<string | undefined>();

  const { data, mutate, isValidating } = swrHook(
    {
      nextToken: nextToken,
      ...hookProps,
    },
    swrProps
  );

  const [virtualData, setVirtualData] = useState<any | undefined>();

  // only set virtual data on first load
  useEffect(() => {
    // case 1: no virtual data yet, set it
    if (data?.[listObjKey] && !virtualData) {
      setVirtualData(data);
      // case 2: virtual data already exists, but it's just been updated, set it
    } else if (data?.[listObjKey] && !virtualData.next) {
      setVirtualData(data);
    }
    // case 3: virtual data already exists, but they've scrolled down, append it
    else if (
      data?.[listObjKey]?.length > 0 &&
      virtualData &&
      virtualData.next != data?.next
    ) {
      setVirtualData((curr: any) => {
        const prevListItems =
          curr?.[listObjKey].length > 0 ? curr?.[listObjKey] : [];
        return {
          ...curr,
          [listObjKey]: [...prevListItems, ...(data?.[listObjKey] ?? [])],
          next: data?.next,
        };
      });
    }
  }, [data, isValidating]);

  const incrementPage = () => {
    data?.next && !isValidating && setNextToken(data.next);
  };

  const canNextPage = useMemo(() => !!virtualData?.next, [virtualData]);

  return {
    data: virtualData,
    isValidating,
    canNextPage,
    incrementPage,
  };
};
