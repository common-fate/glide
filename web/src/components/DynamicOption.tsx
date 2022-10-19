import React from "react";
import { CheckIcon, CopyIcon } from "@chakra-ui/icons";
import {
  Flex,
  IconButton,
  Tooltip,
  useClipboard,
  WrapItem,
} from "@chakra-ui/react";
import { BoltIcon } from "./icons/Icons";

export const DynamicOption: React.FC<{ label: string; value: string }> = ({
  label,
  value,
}) => {
  const { hasCopied, onCopy } = useClipboard(value);
  return (
    <WrapItem>
      <Tooltip label={value}>
        <Flex
          textStyle={"Body/Small"}
          rounded="full"
          bg="neutrals.300"
          py={1}
          px={4}
        >
          {label}{" "}
          <IconButton
            variant="ghost"
            h="20px"
            size="xs"
            icon={hasCopied ? <CheckIcon /> : <BoltIcon />}
            onClick={onCopy}
            aria-label={"Copy"}
          />
        </Flex>
      </Tooltip>
    </WrapItem>
  );
};
