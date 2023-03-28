import { Box, BoxProps, Flex, HStack, Text } from "@chakra-ui/react";
import React from "react";
import { Link } from "react-location";
import { useAdminGetAccessRule } from "../utils/backend-client/admin/admin";
import { useUserGetAccessRule } from "../utils/backend-client/end-user/end-user";
import { ProviderIcon } from "./icons/providerIcon";

type Props = {
  reason?: string;
  accessRuleId: string;
  adminRoute: boolean;
} & Omit<BoxProps, "children">;

export const RuleNameCell: React.FC<Props> = ({
  accessRuleId,
  reason,
  adminRoute,
  ...rest
}) => {
  const { data } = adminRoute
    ? useUserGetAccessRule(accessRuleId)
    : useAdminGetAccessRule(accessRuleId);

  // For now we're disabling linking/click-through
  const isAdmin = false; // window.location.pathname.includes("admin");

  return (
    <Link to={isAdmin ? "/admin/access-rules/" + accessRuleId : "#"}>
      <Flex
        direction="column"
        className="group"
        textStyle="Body/Small"
        minW="200px"
        maxW="600px"
        as="a"
        {...rest}
      >
        <HStack>
          <ProviderIcon
            shortType={data?.target.provider.type}
            id={data?.target.provider.id}
          />
          <Text
            _groupHover={{
              textDecor: isAdmin ? "underline" : undefined,
            }}
            textStyle="Body/SmallBold"
            color="neutrals.700"
          >
            {data?.name}
          </Text>
        </HStack>
        {reason && <Text color="neutrals.500">{reason}</Text>}
      </Flex>
    </Link>
  );
};
