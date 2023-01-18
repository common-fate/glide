import {
  AvatarProps,
  HStack,
  TextProps,
  SkeletonCircle,
  SkeletonText,
} from "@chakra-ui/react";
import React from "react";
import { useUserGetUser } from "../utils/backend-client/end-user/end-user";

import { UserAvatarDetails } from "./UserAvatar";
import { BotAvatarDetails } from "./BotAvatar";

enum AvatarType {
  User,
  Bot,
}

interface CFAvatarProps extends AvatarProps {
  //   type: AvatarType;
  userId: string;
  textProps?: TextProps;
  tooltip?: boolean;
}

// UserAvatar loads a user avatar from a user ID
export const CFAvatar: React.FC<CFAvatarProps> = ({ userId }) => {
  const { data, isValidating } = useUserGetUser(userId);

  if (!data && isValidating) {
    return (
      <HStack>
        <SkeletonCircle size="6" />
        <SkeletonText noOfLines={1} w="12ch" />
      </HStack>
    );
  }

  if (userId != "tf_bot" && data) {
    return (
      <UserAvatarDetails
        tooltip
        user={data}
        size="xs"
        variant="withBorder"
        textProps={{
          textStyle: "Body/Small",
          maxW: "20ch",
          noOfLines: 1,
          color: "neutrals.700",
        }}
      />
    );
  } else {
    return (
      <BotAvatarDetails
        tooltip
        size="xs"
        variant="withBorder"
        textProps={{
          textStyle: "Body/Small",
          maxW: "20ch",
          noOfLines: 1,
          color: "neutrals.700",
        }}
        botType={"Terraform"}
      />
    );
  }
};
