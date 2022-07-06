import { BoxProps, HStack, Stack, Text } from "@chakra-ui/layout";
import React from "react";
import { CFCard } from "./CFCard";

interface Props extends BoxProps {
  leftIcon?: React.ReactElement;
  title: string;
}

export const OnboardingCard: React.FC<Props> = ({
  leftIcon,
  title,
  children,
  ...rest
}) => {
  return (
    <CFCard {...rest}>
      <HStack align="flex-start" spacing={4} flex={1}>
        {leftIcon}
        <Stack spacing={1} flex={1}>
          <Text color="black" fontWeight="medium">
            {title}
          </Text>
          {children}
        </Stack>
      </HStack>
    </CFCard>
  );
};
