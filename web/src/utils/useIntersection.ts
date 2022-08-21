import React from "react";

/**
 * useIntersection
 * Uses the IntersectionObserver API to determine when a component is in view.
 *
 * NOTE: it's behavior is not 100% reliable, can do with further testing
 */
export const useIntersection = (
  /** Pass in the ref of a given element i.e.
   *   const ref = useRef(); ... <div ref={ref}>
   */
  element: React.MutableRefObject<undefined>,
  /** rootMargin
   * Specify a rootMargin value as in a plain css string i.e.
   * '100px' will trigger 100px from the top of it's container.
   * '-100px' will trigger once 100px down
   */
  rootMargin: string,
  /**
   * Beware of closure state here (variables may have unexpected values if passed into callbackFn)
   */
  callbackFn?: () => void
) => {
  const [isVisible, setState] = React.useState(false);

  React.useEffect(() => {
    const observer = new IntersectionObserver(
      ([entry]) => {
        setState(entry.isIntersecting);
        callbackFn && callbackFn();
      },
      { rootMargin }
    );

    element.current && observer.observe(element.current);

    // return () => element.current && observer.unobserve(element.current);
  }, [element]);

  return isVisible;
};
