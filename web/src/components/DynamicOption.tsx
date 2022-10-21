import React from "react";
import {
  Flex,
  Popover,
  PopoverArrow,
  PopoverBody,
  PopoverContent,
  PopoverTrigger,
  Text,
  WrapItem,
} from "@chakra-ui/react";
import { BoltIcon } from "./icons/Icons";

export const DynamicOption: React.FC<{
  label: string;
  value: string;
  isParentGroup?: boolean;
}> = ({ label, value, isParentGroup }) => {
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
            {isParentGroup && (
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
              {isParentGroup && (
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
