import { AvatarProps, TextProps } from "@chakra-ui/react";
import React from "react";

import { BotAvatarDetails } from "./BotAvatar";
import { UserAvatarDetails } from "./UserAvatar";

interface CFAvatarProps extends AvatarProps {
  //   type: AvatarType;
  userId: string;
  textProps?: TextProps;
  tooltip?: boolean;
}

// UserAvatar loads a user avatar from a user ID
export const CFAvatar: React.FC<CFAvatarProps> = ({ userId }) => {
  if (userId != "bot_governance_api") {
    return (
      <UserAvatarDetails
        tooltip
        userId={userId}
        size="xs"
        variant="withBorder"
        textProps={{
          textStyle: "Body/Small",
          // maxW: "20ch",
          // noOfLines: 1,
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
        botType={userId}
      />
    );
  }
};
