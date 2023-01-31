import { AvatarProps, Avatar, HStack, Text, TextProps } from "@chakra-ui/react";
import React from "react";
import { TerraformIcon } from "./icons/Icons";

interface BotAvatarProps extends AvatarProps {
  botType: string;
  textProps?: TextProps;
  tooltip?: boolean;
}

// UserAvatar loads a user avatar from a user ID
export const BotAvatarDetails: React.FC<BotAvatarProps> = ({
  botType,
  textProps,
  ...rest
}) => {
  switch (botType) {
    case "bot_governance_api":
      return (
        <HStack>
          <Avatar
            icon={<TerraformIcon fontSize="0.9rem" />}
            {...rest}
            bg="white"
          />
          <Text {...textProps}>Governance API</Text>
        </HStack>
      );

    default:
      return (
        <HStack>
          <Avatar {...rest} bg="white" />
          <Text {...textProps}>{botType}</Text>
        </HStack>
      );
      break;
  }
};
