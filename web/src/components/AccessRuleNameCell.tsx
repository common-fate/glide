import { Box, BoxProps, HStack, Text } from "@chakra-ui/react";
import React from "react";
import { Link } from "react-location";
import { useUserGetAccessRule } from "../utils/backend-client/end-user/end-user";
import { ProviderIcon } from "./icons/providerIcon";

type Props = {
  reason?: string;
  accessRuleId: string;
} & Omit<BoxProps, "children">;

export const RuleNameCell: React.FC<Props> = ({
  accessRuleId,
  reason,
  ...rest
}) => {
  const { data } = useUserGetAccessRule(accessRuleId);

  // For now we're disabling linking/click-through
  const isAdmin = false; // window.location.pathname.includes("admin");

  return (
    <Link to={isAdmin ? "/admin/access-rules/" + accessRuleId : "#"}>
      <Box
        className="group"
        textStyle="Body/Small"
        minW="200px"
        as="a"
        {...rest}
      >
        <HStack>
          <ProviderIcon provider={data?.target.provider} />
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
      </Box>
    </Link>
  );
};
