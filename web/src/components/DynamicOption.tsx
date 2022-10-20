import React from "react";
import { CheckIcon, CopyIcon } from "@chakra-ui/icons";
import {
  Flex,
  IconButton,
  Popover,
  PopoverArrow,
  PopoverBody,
  PopoverCloseButton,
  PopoverContent,
  PopoverHeader,
  PopoverTrigger,
  Text,
  Tooltip,
  useClipboard,
  WrapItem,
} from "@chakra-ui/react";
import { BoltIcon } from "./icons/Icons";
import { GroupOption } from "../utils/backend-client/types/accesshandler-openapi.yml";

export const DynamicOption: React.FC<{
  label: string;
  value: string;
  parentGroup?: GroupOption;
}> = ({ label, value, parentGroup }) => {
  const { hasCopied, onCopy } = useClipboard(value);

  const colArr = [""];

  return (
    <WrapItem>
      <Popover trigger="hover">
        <PopoverTrigger>
          <Flex
            textStyle={"Body/Small"}
            rounded="full"
            bg="neutrals.300"
            py={1}
            px={4}
          >
            {label}{" "}
            {parentGroup && (
              <BoltIcon
                transition="all .2s ease"
                color="neutrals.400"
                h="20px"
                ml={2}
              />
            )}
            {/* <IconButton
            variant="ghost"
            h="20px"
            size="xs"
            // color
            icon={(parentGroup && (hasCopied ? <CheckIcon /> : <BoltIcon />)}
            onClick={onCopy}
            aria-label={"Copy"}
          /> */}
          </Flex>
        </PopoverTrigger>
        <PopoverContent maxW="180px">
          <PopoverArrow />
          <PopoverBody>
            <Text fontWeight="semibold" textStyle={"Body/Medium"}>
              {label}
              {parentGroup && (
                <BoltIcon
                  // filter="grayscale(1);"
                  color="neutrals.400"
                  h="12px"
                  ml={2}
                />
              )}
            </Text>
            {value}
          </PopoverBody>
        </PopoverContent>
      </Popover>
    </WrapItem>
  );
};
