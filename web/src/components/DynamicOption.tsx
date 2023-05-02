import React from "react";
import {
  Flex,
  HStack,
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
}> = ({ label, value }) => {
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
          </Flex>
        </PopoverTrigger>
        <PopoverContent maxW="180px">
          <PopoverArrow />
          <PopoverBody>
            <Text fontWeight="semibold" textStyle={"Body/Medium"}>
              {label}
            </Text>
            {value}
          </PopoverBody>
        </PopoverContent>
      </Popover>
    </WrapItem>
  );
};
