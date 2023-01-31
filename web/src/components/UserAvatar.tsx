import {
  Avatar,
  AvatarProps,
  HStack,
  SkeletonCircle,
  SkeletonText,
  Text,
  TextProps,
  Tooltip,
} from "@chakra-ui/react";
import React from "react";

import { useUserGetUser } from "../utils/backend-client/end-user/end-user";

interface UserAvatarProps extends AvatarProps {
  userId: string;
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
  userId,
  textProps,
  tooltip,
  ...rest
}) => {
  const { data, isValidating } = useUserGetUser(userId);

  if (!data && isValidating) {
    return (
      <HStack>
        <SkeletonCircle size="6" />
        <SkeletonText noOfLines={1} w="12ch" />
      </HStack>
    );
  }

  const Component = tooltip ? TooltipAvatar : Avatar;

  // Loading/loaded states
  return (
    <HStack>
      <Component name={data?.email} {...rest} />
      <Text {...textProps}>{data?.email}</Text>
    </HStack>
  );
};
