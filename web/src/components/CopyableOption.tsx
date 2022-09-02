import React from "react";
import { CheckIcon, CopyIcon } from "@chakra-ui/icons";
import {
  Flex,
  IconButton,
  Tooltip,
  useClipboard,
  WrapItem,
} from "@chakra-ui/react";

export const CopyableOption: React.FC<{ label: string; value: string }> = ({
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
            icon={hasCopied ? <CheckIcon /> : <CopyIcon />}
            onClick={onCopy}
            aria-label={"Copy"}
          />
        </Flex>
      </Tooltip>
    </WrapItem>
  );
};
