import {
  AvatarProps,
  Avatar,
  HStack,
  Text,
  TextProps,
  Tooltip,
  SkeletonCircle,
  SkeletonText,
} from "@chakra-ui/react";
import React from "react";
import { TerraformIcon } from "./icons/Icons";

interface BotAvatarProps extends AvatarProps {
  botType: string;
  textProps?: TextProps;
  tooltip?: boolean;
}

const TooltipAvatar: typeof Avatar = (props: any) => (
  <Tooltip hasArrow={true} label={props.name}>
    <Avatar {...props} />
  </Tooltip>
);

// UserAvatar loads a user avatar from a user ID
export const BotAvatarDetails: React.FC<BotAvatarProps> = ({
  botType,
  textProps,
  tooltip,
  ...rest
}) => {
  const Component = tooltip ? TooltipAvatar : Avatar;

  // Loading/loaded states
  return (
    <HStack>
      <Avatar icon={<TerraformIcon fontSize="0.9rem" />} {...rest} bg="white" />
      <Text {...textProps}>Terraform</Text>
    </HStack>
  );
};
