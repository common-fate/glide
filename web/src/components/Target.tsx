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
import { F } from "dist/assets/Layout-e4945437";

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
      w="100%"
      alignContent="flex-start"
      p={2}
      rounded="md"
      id="targetParent"
      {...rest}
    >
      {showIcon && (
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
      )}
      <Grid autoColumns="minmax(0, 1fr)" gridAutoFlow="column">
        {target.fields.map((field, i) => (
          <Flex alignItems="center" w="100%">
            <Divider
              orientation="vertical"
              borderColor={"black"}
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
                <Stack wrap="wrap">
                  <Text textStyle={"Body/SmallBold"}>{field.fieldTitle}</Text>
                  <Text textStyle={"Body/Small"}>{field.valueLabel}</Text>
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
