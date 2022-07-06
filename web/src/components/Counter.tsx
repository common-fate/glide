import { Tag, TagProps } from "@chakra-ui/react";
import React from "react";

interface Props extends TagProps {
  count?: number;
}

const Counter = ({ count, children, ...rest }: Props) => {
  // If the number > 2 digits, good to have less padding
  let p = count && count >= 10 ? 1 : 2;

  return (
    <Tag pl={p} pr={p} fontWeight="500" rounded="full" size="sm" {...rest}>
      {count ?? 0}
      {children}
    </Tag>
  );
};

export default Counter;
