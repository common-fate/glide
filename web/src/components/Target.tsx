import { CheckCircleIcon } from "@chakra-ui/icons";
import {
  Divider,
  Flex,
  FlexProps,
  HStack,
  Stack,
  Text,
  Tooltip,
} from "@chakra-ui/react";
import React from "react";
// @ts-ignore
import { ProviderIcon, ShortTypes } from "../components/icons/providerIcon";
import { Target } from "../utils/backend-client/types";

interface TargetDetailProps extends FlexProps {
  target: Target;
  isChecked?: boolean;
  showIcon?: boolean;
}

export const TargetDetail: React.FC<TargetDetailProps> = ({
  target,
  isChecked,
  showIcon,
  ...rest
}) => {
  return (
    <Flex
      alignContent="flex-start"
      position="relative"
      p={2}
      rounded="md"
      {...rest}
    >
      <CheckCircleIcon
        visibility={isChecked ? "visible" : "hidden"}
        position="absolute"
        top={2}
        right={2}
        boxSize={"12px"}
        color={"brandBlue.300"}
      />
      {showIcon && (
        <Tooltip
          label={`${target.kind.publisher}/${target.kind.name}/${target.kind.kind}`}
          placement="right"
        >
          <Flex p={2} position="relative">
            <ProviderIcon
              boxSize={"24px"}
              shortType={target.kind.icon as ShortTypes}
            />
          </Flex>
        </Tooltip>
      )}
      <HStack>
        {target.fields.map((field, i) => (
          <React.Fragment key={target.id + field.id}>
            <Divider
              orientation="vertical"
              borderColor={"neutrals.300"}
              h="80%"
              hidden={i === 0}
            />

            <Tooltip
              label={
                <Stack>
                  <Text color="white" textStyle={"Body/Small"}>
                    {field.fieldTitle}
                  </Text>
                  <Text color="white" textStyle={"Body/Small"}>
                    {field.fieldDescription}
                  </Text>
                  <Text color="white" textStyle={"Body/Small"}>
                    {field.valueLabel}
                  </Text>
                  <Text color="white" textStyle={"Body/Small"}>
                    {field.value}
                  </Text>
                  <Text color="white" textStyle={"Body/Small"}>
                    {field.valueDescription}
                  </Text>
                </Stack>
              }
              placement="top"
            >
              <Stack h="100%" spacing={0} justify="start" align="start">
                <Text textStyle={"Body/SmallBold"} noOfLines={1}>
                  {field.fieldTitle}
                </Text>
                <Text fontFamily="mono" textStyle={"Body/Small"} noOfLines={1}>
                  {field.valueLabel}
                </Text>
              </Stack>
            </Tooltip>
          </React.Fragment>
        ))}
      </HStack>
    </Flex>
  );
};
