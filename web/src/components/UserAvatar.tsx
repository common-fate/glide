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
import { useGetUser } from "../utils/backend-client/end-user/end-user";
import { TerraformIcon } from "./icons/Icons";

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
  const { data, isValidating } = useGetUser(user);
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
  const Component = tooltip ? TooltipAvatar : Avatar;

  //skip lookup user step if creating user was terraform
  if (user === "bot_governance_api") {
    // Loading/loaded states
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
  }
  // @ts-ignore
  const { data, isValidating } = useGetUser(user);

  if (!data) {
    return <Avatar {...rest} />;
  }

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
