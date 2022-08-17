import { Dispatch, SetStateAction, useEffect, useMemo, useState } from "react";

type InputParams<T extends (...args: any[]) => any> = {
  /** This will be called within useInfiniteScrollApi, the resulting nextToken will be extracted */
  swrHook: T;
  hookProps?: Parameters<T>[0];
  /** by specifying the listObjKey we know what array value to mutate from the returned data type (helps keep this method generic!)  */
  listObjKey: keyof ReturnType<T>["data"];
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
  listObjKey,
}: InputParams<T>): PaginationProps<T> => {
  const [nextToken, setNextToken] = useState<string | undefined>();

  const { data, mutate, isValidating } = swrHook({
    ...hookProps,
    nextToken: nextToken,
  });

  const [virtualData, setVirtualData] = useState();

  // only set virtual data on first load
  useEffect(() => {
    if (data?.[listObjKey].length > 0) {
      if (data?.[listObjKey] && !virtualData) {
        setVirtualData(data);
      } else if (
        data?.[listObjKey] &&
        virtualData &&
        virtualData.next != data?.next
      ) {
        setVirtualData((curr) => {
          let prevListItems =
            curr?.[listObjKey].length > 0 ? curr?.[listObjKey] : [];
          return {
            ...curr,
            [listObjKey]: [...prevListItems, ...data?.[listObjKey]],
            next: data?.next,
          };
        });
      }
    }
  }, [data, isValidating]);

  const incrementPage = () => {
    data?.next && setNextToken(data.next);
  };

  const canNextPage = useMemo(() => !!virtualData?.next, [virtualData]);

  return {
    data: virtualData,
    isValidating,
    canNextPage,
    incrementPage,
  };
};
