import {
  Box,
  Flex,
  Popover,
  PopoverArrow,
  PopoverBody,
  PopoverCloseButton,
  PopoverContent,
  PopoverTrigger,
  Text,
  useColorModeValue,
} from "@chakra-ui/react";
import { differenceInSeconds, formatDistanceToNowStrict } from "date-fns";
import React, { useEffect, useState } from "react";

interface Props {
  header: React.ReactNode;
  timestamp: Date;
  /** this should be a 20px react node (use boxSize with chakra elements) */
  marker?: React.ReactNode;
  index: number;
  arrLength: number;
}

const TimestampLine = React.forwardRef<React.ReactNode, { timestamp: Date }>(
  ({ timestamp }, ref) => (
    <Text fontSize="sm" color="gray.400" fontWeight="normal">
      {formatDistanceToNowStrict(timestamp, { addSuffix: true })}
    </Text>
  )
);

export const CFTimelineRow = ({
  header,
  timestamp,
  marker,
  index,
  arrLength,
}: Props) => {
  const textColor = useColorModeValue("gray.700", "white.300");
  const [timeSince, setTimeSince] = useState("");
  const renderTimeSince = () => {
    if (differenceInSeconds(new Date(), timestamp) < 60) {
      setTimeSince("a few seconds ago");
    } else {
      setTimeSince(formatDistanceToNowStrict(timestamp, { addSuffix: true }));
    }
  };
  useEffect(() => {
    renderTimeSince();
    const interval = setInterval(renderTimeSince, 60000);
    return () => clearInterval(interval);
  }, [timestamp]);
  return (
    <Flex>
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
        direction="column"
        justifyContent="flex-start"
        h="100%"
        marginBottom={index === arrLength - 1 ? 2 : 6}
      >
        <Text fontSize="sm" color={textColor} fontWeight="bold">
          {header}
        </Text>
        <Box>
          <Popover>
            <PopoverTrigger>
              <Text
                flexGrow={0}
                tabIndex={0}
                as="button"
                display={"inline-block"}
                textAlign="left"
                fontSize="sm"
                color="gray.400"
                fontWeight="normal"
              >
                {timeSince}
              </Text>
            </PopoverTrigger>
            <PopoverContent>
              <PopoverArrow />
              <PopoverCloseButton />
              <PopoverBody>{timestamp.toString()}</PopoverBody>
            </PopoverContent>
          </Popover>
        </Box>
      </Flex>
    </Flex>
  );
};
