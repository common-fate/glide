import { AvatarProps, Avatar, HStack, Text, TextProps } from "@chakra-ui/react";
import React from "react";
import { TerraformIcon } from "./icons/Icons";

export enum BotType {
  Terraform,
}

interface BotAvatarProps extends AvatarProps {
  botType: BotType;
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
    case BotType.Terraform:
      return (
        <HStack>
          <Avatar
            icon={<TerraformIcon fontSize="0.9rem" />}
            {...rest}
            bg="white"
          />
          <Text {...textProps}>Terraform</Text>
        </HStack>
      );

    default:
      return (
        <HStack>
          <Avatar
            icon={<TerraformIcon fontSize="0.9rem" />}
            {...rest}
            bg="white"
          />
          <Text {...textProps}>Terraform</Text>
        </HStack>
      );
      break;
  }
};
