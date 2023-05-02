import { Box, chakra, Code, CodeProps, Tooltip } from "@chakra-ui/react";
import React from "react";
import { TargetField } from "../utils/backend-client/types";

type Props = {
  fields: TargetField[];
  showTooltips?: boolean;
};

const FieldsCodeBlock = ({
  fields,
  showTooltips,
  ...props
}: Props & CodeProps) => {
  return (
    <Code bg="none" whiteSpace="pre-wrap" {...props}>
      {fields.map((f) =>
        showTooltips ? (
          <Tooltip
            key={f.id}
            label={
              <>
                <Box fontWeight="bold">{f.fieldTitle}</Box>
                {f.fieldDescription && <Box>{f.fieldDescription}</Box>}
                {f.valueDescription && <Box mt={2}>{f.valueDescription}</Box>}
              </>
            }
            placement="top"
          >
            <chakra.span
              lineHeight="1.05rem"
              sx={{
                _hover: {
                  // bg: "#CCF8FE99",
                  // color: "actionInfo.200",
                  span: {
                    textDecor: "underline",
                  },
                  rounded: "md",
                },
              }}
              w="min-content"
              // display="inline-flex"
              // display="flex"
              // flexDir="row"
              // verticalAlign="top"
              // noOfLines={1}
              // wordBreak="break-all"
              // textOverflow="ellipsis"
              display="block"
              overflow="hidden"
              whiteSpace="nowrap"
              // h="12px"
              // display="block"
              // d="flex"
            >
              <chakra.span
                noOfLines={1}
                textOverflow="ellipsis"
                wordBreak="break-all"
                // bg="pink"
                display="inline-block"
                maxW="400px"
                // w="20%"
              >
                {f.valueLabel}
              </chakra.span>

              <chakra.span
                noOfLines={1}
                textOverflow="ellipsis"
                wordBreak="break-all"
                display="inline-block"
                // bg="tomato"
                // w="80%"
              >
                :{f.value}
              </chakra.span>
            </chakra.span>
          </Tooltip>
        ) : (
          <chakra.span
            noOfLines={1}
            wordBreak="break-all"
            textOverflow="ellipsis"
          >
            {f.value}
          </chakra.span>
        )
      )}
    </Code>
  );
};

export default FieldsCodeBlock;
