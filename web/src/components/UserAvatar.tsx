import {
  AvatarProps,
  Avatar,
  HStack,
  Text,
  TextProps,
  Tooltip,
} from "@chakra-ui/react";
import React from "react";

import { User } from "src/utils/backend-client/types";

interface UserAvatarProps extends AvatarProps {
  user: User;
  textProps?: TextProps;
  tooltip?: boolean;
}

const TooltipAvatar: typeof Avatar = (props: any) => (
  <Tooltip hasArrow={true} label={props.name}>
    <Avatar {...props} />
  </Tooltip>
);

// UserAvatar loads a user avatar from a user ID
export const UserAvatarDetails: React.FC<UserAvatarProps> = ({
  user,
  textProps,
  tooltip,
  ...rest
}) => {
  const Component = tooltip ? TooltipAvatar : Avatar;

  // Loading/loaded states
  return (
    <HStack>
      <Component name={user?.email} {...rest} />
      <Text {...textProps}>{user?.email}</Text>
    </HStack>
  );
};
