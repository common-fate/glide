import React from "react";

/**
 * Gives the inner height of the browser allowing elements to be sized vertically for mobile.
 * see: https://github.com/chakra-ui/chakra-ui/issues/2468
 * @returns the inner height value in pixels
 */
export const useInnerHeight = (): number => {
  const [innerHeight, setInnerHeight] = React.useState<number>(
    typeof window !== "undefined" ? window.innerHeight : 100
  );

  function windowResizeHandler() {
    if (window !== undefined) {
      setInnerHeight(window.innerHeight);
    }
  }

  React.useEffect(() => {
    if (window !== undefined) {
      window.addEventListener("resize", windowResizeHandler);
      return () => {
        window.removeEventListener("resize", windowResizeHandler);
      };
    }
  }, []);

  return innerHeight;
};
