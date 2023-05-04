import {
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
interface Props extends FlexProps {
  replaceValue?: string;
  value?: string;
  success?: string | string[];
  danger?: string | string[];
  warning?: string | string[];
  info?: string | string[];
}
// StatusCell displays a status icon and can be configured to map success,warning or danger to some strings
// this will match the value to the configured color groups
export const StatusCell: React.FC<Props> = ({
  replaceValue,
  value,
  success,
  danger,
  warning,
  info,
  ...rest
}) => {
  // We may want to handle loading/null states separately, for now this is to serve as a skeleton component
  if (!value)
    return (
      <Flex align="center" h="21px" {...rest}>
        <SkeletonCircle size="8px" mr={2} />{" "}
        <SkeletonText noOfLines={1} lineHeight="12px" w="10ch" />
      </Flex>
    );

  // default color is warning
  let statusColor = "actionWarning.200";

  if (Array.isArray(success) ? success.includes(value) : success === value) {
    statusColor = "actionSuccess.200";
  }
  if (Array.isArray(danger) ? danger.includes(value) : danger === value) {
    statusColor = "actionDanger.200";
  }
  if (Array.isArray(warning) ? warning.includes(value) : warning === value) {
    statusColor = "actionWarning.200";
  }
  if (Array.isArray(info) ? info.includes(value) : info === value) {
    statusColor = "actionInfo.200";
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
          {replaceValue ? replaceValue : value.toLowerCase()}
        </Text>
      </Flex>
    );
  }
  return (
    <Flex minW="75px" align="center" {...rest}>
      <SkeletonCircle endColor={statusColor} size="8px" mr={2} />{" "}
      <Text as="span" css={{ ":first-letter": { textTransform: "uppercase" } }}>
        {replaceValue ? replaceValue : value.toLowerCase()}
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
