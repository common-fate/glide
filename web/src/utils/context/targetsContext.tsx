import { useBoolean } from "@chakra-ui/react";
import React, { useEffect, useState } from "react";
import { userListEntitlementTargets } from "../backend-client/default/default";
import { Target } from "../backend-client/types";
import { createCtx } from "./createCtx";
export interface TargetsContextProps {
  targets: Target[];
  validating: boolean;
}
const CACHE_KEY = "cachedTargets";
const [useTargets, TargetsContextProvider] = createCtx<TargetsContextProps>();

interface Props {
  children: React.ReactNode;
}

const TargetsProvider: React.FC<Props> = ({ children }) => {
  const [targets, setTargets] = useState<Target[]>([]);
  const [validating, setValidating] = useBoolean(true);

  useEffect(() => {
    const c = localStorage.getItem(CACHE_KEY);
    setTargets(JSON.parse(c || "[]"));
    const fetchData = async (
      nextToken: string | undefined
    ): Promise<Target[]> => {
      const result = await userListEntitlementTargets({ nextToken });
      const targets = result.targets;
      if (result.next) {
        const remainingTargets = await fetchData(result.next);
        targets.push(...remainingTargets);
      }
      return targets;
    };
    fetchData(undefined)
      .then((t) => {
        setTargets(t);
        localStorage.setItem(CACHE_KEY, JSON.stringify(t));
      })
      .finally(() => {
        setValidating.off();
      });
  }, []);

  return (
    <TargetsContextProvider value={{ targets, validating }}>
      {children}
    </TargetsContextProvider>
  );
};
export { useTargets, TargetsProvider };
