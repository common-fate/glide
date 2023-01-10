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
import { useUserGetUser } from "../utils/backend-client/end-user/end-user";

interface UserAvatarProps extends AvatarProps {
  user: string | undefined;
  textProps?: TextProps;
  tooltip?: boolean;
}

const TooltipAvatar: typeof Avatar = (props: any) => (
  <Tooltip hasArrow={true} label={props.name}>
    <Avatar {...props} />
  </Tooltip>
);

// UserAvatar loads a user avatar from a user ID
export const UserAvatar: React.FC<UserAvatarProps> = ({
  user,
  tooltip,
  ...rest
}) => {
  // @ts-ignore, swr should handle this fine
  const { data, isValidating } = useUserGetUser(user);
  if (!data && isValidating) {
    return <SkeletonCircle size="6" />;
  }
  const Component = tooltip ? TooltipAvatar : Avatar;
  return <Component name={data?.email} {...rest} />;
};

// UserAvatar loads a user avatar from a user ID
export const UserAvatarDetails: React.FC<UserAvatarProps> = ({
  user,
  textProps,
  tooltip,
  ...rest
}) => {
  // @ts-ignore
  const { data, isValidating } = useUserGetUser(user);

  if (!data) {
    return <Avatar {...rest} />;
  }

  const Component = tooltip ? TooltipAvatar : Avatar;

  // Loading/loaded states
  return !data && isValidating ? (
    <HStack>
      <SkeletonCircle size="6" />
      <SkeletonText noOfLines={1} w="12ch" />
    </HStack>
  ) : (
    <HStack>
      <Component name={data?.email} {...rest} />
      <Text {...textProps}>{data?.email}</Text>
    </HStack>
  );
};
