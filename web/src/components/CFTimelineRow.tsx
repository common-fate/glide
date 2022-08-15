import { Box, Flex, Text, useColorModeValue } from "@chakra-ui/react";
import React from "react";

interface Props {
  header: React.ReactNode;
  body?: React.ReactNode;
  /** this should be a 20px react node (use boxSize with chakra elements) */
  marker?: React.ReactNode;
  index: number;
  arrLength: number;
}

export const CFTimelineRow = ({
  header,
  body,
  marker,
  index,
  arrLength,
}: Props) => {
  const textColor = useColorModeValue("gray.700", "white.300");

  return (
    <Flex w="100%">
      <Flex
        flexDir="column"
        flexShrink={0}
        width="24px"
        marginRight={3}
        // webkitBoxAlign="center"
        alignItems="center"
        sx={{
          // This provides an override unless specified in the marker prop
          svg: {
            color: "brandBlue.200",
          },
        }}
      >
        {marker}
        {/* <Icon as={logo} bg={bgIconColor} color={"teal.300"} boxSize="20px" /> */}
        <Box
          w="2px"
          bg="gray.200"
          flexGrow={1}
          boxSizing="border-box"
          margin="8px 0px"
          display={index === arrLength - 1 ? "none" : "block"}
        />
      </Flex>
      <Flex
        w="100%"
        direction="column"
        justifyContent="flex-start"
        h="100%"
        marginBottom={index === arrLength - 1 ? 2 : 6}
      >
        <Text fontSize="sm" color={textColor} fontWeight="bold">
          {header}
        </Text>
        <Text fontSize="sm" color="gray.400" fontWeight="normal">
          {body}
        </Text>
      </Flex>
    </Flex>
  );
};
