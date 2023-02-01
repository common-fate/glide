import React from "react";
import { BoxProps, Circle, Flex, Text } from "@chakra-ui/react";
import { ProviderV2 } from "../utils/local-client/types/openapi.yml";

type Props = {
  provider: ProviderV2;
};

const ProviderStatus = ({ provider, ...boxProps }: Props & BoxProps) => {
  // default color is warning
  let statusColor = "actionWarning.200";
  if (provider.status === "DEPLOYED") {
    statusColor = "actionSuccess.200";
  }
  if (provider.status === "UPDATING" || provider.status === "CREATING") {
    statusColor = "actionWarning.200";
  }
  if (provider.status === "DELETED") {
    statusColor = "actionDanger.200";
  }

  return (
    <Flex minW="75px" align="center" {...boxProps}>
      <Circle bg={statusColor} size="8px" mr={2} />{" "}
      <Text as="span" css={{ ":first-letter": { textTransform: "uppercase" } }}>
        {provider.status.toLowerCase()}
      </Text>
    </Flex>
  );
};

export default ProviderStatus;
