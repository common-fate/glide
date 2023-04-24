import { chakra, Code, CodeProps } from "@chakra-ui/react";
import React from "react";
import { TargetField } from "../utils/backend-client/types";

type Props = {
  fields: TargetField[];
};

const FieldsCodeBlock = ({ fields, ...props }: Props & CodeProps) => {
  return (
    <Code bg="white" whiteSpace="pre-wrap" {...props}>
      {fields.map((f) => (
        <chakra.span
          noOfLines={1}
          wordBreak="break-all"
          textOverflow="ellipsis"
        >
          {f.value}
        </chakra.span>
      ))}
    </Code>
  );
};

export default FieldsCodeBlock;
