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
}

export const TargetDetail: React.FC<TargetDetailProps> = ({
  target,
  isChecked,
  ...rest
}) => {
  return (
    <Flex alignContent="flex-start" p={2} rounded="md" {...rest}>
      <Tooltip
        label={`${target.kind.publisher}/${target.kind.name}/${target.kind.kind}`}
        placement="right"
      >
        <Flex p={6} position="relative">
          <CheckCircleIcon
            visibility={isChecked ? "visible" : "hidden"}
            position="absolute"
            top={0}
            left={0}
            boxSize={"12px"}
            color={"brandBlue.300"}
          />
          <ProviderIcon
            boxSize={"24px"}
            shortType={target.kind.icon as ShortTypes}
          />
        </Flex>
      </Tooltip>

      <HStack>
        {target.fields.map((field, i) => (
          <React.Fragment key={target.id + field.id}>
            <Divider
              orientation="vertical"
              borderColor={"black"}
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
              <Stack>
                <Text textStyle={"Body/SmallBold"} noOfLines={1}>
                  {field.fieldTitle}
                </Text>
                <Text textStyle={"Body/Small"} noOfLines={1}>
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
