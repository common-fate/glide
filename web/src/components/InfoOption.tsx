import { InfoIcon } from "@chakra-ui/icons";
import {
  Center,
  HStack,
  Popover,
  PopoverArrow,
  PopoverBody,
  PopoverCloseButton,
  PopoverContent,
  PopoverTrigger,
  Text,
  WrapItem,
} from "@chakra-ui/react";
import React from "react";
import { colors } from "../utils/theme/colors";

export const InfoOption: React.FC<{ label: string; value: string }> = ({
  label,
  value,
}) => {
  return (
    <WrapItem>
      <Popover>
        <HStack
          textStyle={"Body/Small"}
          rounded="full"
          bg="neutrals.300"
          py={1}
          px={4}
        >
          <Text>{label}</Text>{" "}
          <PopoverTrigger>
            <InfoIcon
              cursor={"pointer"}
              h="10px"
              color={colors.neutrals[600]}
            />
          </PopoverTrigger>
        </HStack>
        <PopoverContent>
          <PopoverArrow />
          <PopoverCloseButton />
          <PopoverBody>
            <Center>{value}</Center>
          </PopoverBody>
        </PopoverContent>
      </Popover>
    </WrapItem>
  );
};
