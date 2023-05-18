import {
  Center,
  Circle,
  Flex,
  FlexProps,
  SkeletonCircle,
  SkeletonText,
  Text,
} from "@chakra-ui/react";
import React from "react";
import {
  RequestAccessGroupApprovalMethod,
  RequestAccessGroupTargetStatus,
  RequestStatus,
} from "../utils/backend-client/types";
interface Props<T> extends FlexProps {
  replaceValue?: string;
  value?: T;
  success?: T | T[];
  danger?: T | T[];
  warning?: T | T[];
  info?: T | T[];
}
// StatusCell displays a status icon and can be configured to map success,warning or danger to some strings
// this will match the value to the configured color groups
export const StatusCell = <T,>({
  replaceValue,
  value,
  success,
  danger,
  warning,
  info,
  ...rest
}: Props<T>) => {
  // We may want to handle loading/null states separately, for now this is to serve as a skeleton component
  if (!value || typeof value != "string")
    return (
      <Flex align="center" h="21px" {...rest}>
        <SkeletonCircle size="8px" mr={2} />{" "}
        <SkeletonText noOfLines={1} lineHeight="12px" w="10ch" />
      </Flex>
    );

  // default color is warning
  let statusColor = "actionWarning.200";

  // Pulsing status added to match Figma Designs
  // https://www.figma.com/file/ziXjEufb8v3FVDZQo55ZK2/CF-UI-Designs?type=design&node-id=3149%3A18483&t=bfWcublL9Zq6bwPX-1
  let statusColorforPulse = "actionSuccess.100";

  if (Array.isArray(success) ? success.includes(value) : success === value) {
    statusColor = "actionSuccess.200";
    statusColorforPulse = "actionSuccess.100";
  }
  if (Array.isArray(danger) ? danger.includes(value) : danger === value) {
    statusColor = "actionDanger.200";
    statusColorforPulse = "actionDanger.100";
  }
  if (Array.isArray(warning) ? warning.includes(value) : warning === value) {
    statusColor = "actionWarning.200";
    statusColorforPulse = "actionWarning.100";
  }
  if (Array.isArray(info) ? info.includes(value) : info === value) {
    statusColor = "actionInfo.200";
    statusColorforPulse = "actionInfo.100";
  }

  //Dont pulse status for warning and error status's
  if (
    statusColor === "actionWarning.200" ||
    statusColor === "actionDanger.200"
  ) {
    return (
      <Flex minW="75px" align="center" {...rest}>
        <Circle bg={statusColor} size="8px" mr={2} />{" "}
        <Text
          as="span"
          css={{ ":first-letter": { textTransform: "uppercase" } }}
        >
          {replaceValue ? replaceValue : (value as string).toLowerCase()}
        </Text>
      </Flex>
    );
  }

  return (
    <Flex minW="75px" align="center" {...rest}>
      <Center
        sx={{
          "@keyframes pulse": {
            "0%": {
              outlineColor: statusColorforPulse,
              transform: "scale(1)",
            },
            "50%": {
              outlineColor: "transparent",
              transform: "scale(0.95)",
            },
            "100%": {
              outlineColor: statusColorforPulse,
              transform: "scale(1)",
            },
          },
        }}
      >
        <Circle
          bg={statusColor}
          size="8px"
          mr={2}
          animation="pulse 2s infinite"
          outline="4px solid transparent"
          // border on outside
        />{" "}
        {/* <SkeletonCircle endColor={statusColor} size="8px" mr={2} />{" "} */}
      </Center>
      <Text as="span" css={{ ":first-letter": { textTransform: "uppercase" } }}>
        {replaceValue ? replaceValue : (value as string).toLowerCase()}
      </Text>
    </Flex>
  );
};

interface RequestStatusCellProps extends FlexProps {
  value: string | undefined;
  approvalMethod: RequestAccessGroupApprovalMethod | undefined;
}

// RequestStatusCell providers a slim wrapper to remove boilerplate for request statuses
export const RequestStatusCell: React.FC<RequestStatusCellProps> = ({
  value,
  approvalMethod,
  ...rest
}) => {
  const isAuto = approvalMethod === "AUTOMATIC" && value;

  return (
    <StatusCell
      value={isAuto ? "Automatically approved" : value}
      // success={[RequestStatus.APPROVED, "Automatically approved"]}
      // danger={[RequestStatus.DECLINED, RequestStatus.CANCELLED]}
      warning={RequestStatus.PENDING}
      textStyle="Body/Small"
      {...rest}
    ></StatusCell>
  );
};

interface GrantStatusCellProps extends FlexProps {
  targetStatus: RequestAccessGroupTargetStatus;
}

// RequestStatusCell providers a slim wrapper to remove boilerplate for request statuses
export const GrantStatusCell: React.FC<GrantStatusCellProps> = ({
  targetStatus,
  ...rest
}) => {
  switch (targetStatus) {
    case "ACTIVE":
      return (
        <StatusCell
          {...rest}
          success="ACTIVE"
          value={targetStatus}
          replaceValue={"Active"}
          fontSize="12px"
        />
      );
    case "AWAITING_START":
      return (
        <StatusCell
          {...rest}
          info={targetStatus}
          value={targetStatus}
          replaceValue={"Awaiting Start"}
          fontSize="12px"
        />
      );
    case "ERROR":
      return (
        <StatusCell
          danger={targetStatus}
          value={targetStatus}
          replaceValue={"Error"}
          fontSize="12px"
        />
      );
    case "EXPIRED":
      return (
        <StatusCell
          {...rest}
          success="ACTIVE"
          value={targetStatus}
          replaceValue={"Expired"}
          fontSize="12px"
        />
      );
    case "PENDING_PROVISIONING":
      return (
        <StatusCell
          {...rest}
          info={targetStatus}
          value={targetStatus}
          replaceValue={"Pending"}
          fontSize="12px"
        />
      );
    case "REVOKED":
      return (
        <StatusCell
          {...rest}
          success="ACTIVE"
          value={targetStatus}
          replaceValue={"Revoked"}
          fontSize="12px"
        />
      );
    default:
      return (
        <StatusCell
          {...rest}
          success="ACTIVE"
          value={undefined}
          replaceValue={"test"}
          fontSize="12px"
        />
      );
  }
};
