import { CheckCircleIcon } from "@chakra-ui/icons";
import {
  Box,
  Divider,
  Flex,
  FlexProps,
  Grid,
  GridItem,
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
      alignItems="start"
      position="relative"
      p={2}
      rounded="md"
      w="100%"
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
          <Flex pr={2} position="relative">
            <ProviderIcon
              boxSize={"24px"}
              shortType={target.kind.name as ShortTypes}
            />
          </Flex>
        </Tooltip>
      )}
      <Grid autoColumns="minmax(0, 1fr)" gridAutoFlow="column">
        {target.fields.map((field, i) => (
          <Flex alignItems="center" w="100%">
            <Divider
              orientation="vertical"
              borderColor={"neutrals.300"}
              h="80%"
              hidden={i === 0}
              mx="20px"
            />
            <GridItem alignItems="center">
              {/* <React.Fragment key={target.id + field.id}> */}

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
                <Stack wrap="wrap" spacing={0}>
                  <Text textStyle={"Body/SmallBold"}>{field.fieldTitle}</Text>
                  <Text fontFamily="mono" textStyle={"Body/Small"}>
                    {field.valueLabel}
                  </Text>
                </Stack>
              </Tooltip>
              {/* </React.Fragment> */}
            </GridItem>
          </Flex>
        ))}
      </Grid>
    </Flex>
  );
};
