import { Box, VStack } from "@chakra-ui/layout";
import { Skeleton } from "@chakra-ui/react";
import React from "react";
import { RequestDetail } from "../utils/backend-client/types";
import { CFTimelineRow } from "./CFTimelineRow";
export const AuditLog: React.FC<{ request?: RequestDetail }> = ({
  request,
}) => {
  if (!request) {
    return (
      <VStack flex={1} align="left">
        <Box textStyle="Heading/H4" as="h4" mb={8}>
          Audit Log
        </Box>
        <Skeleton h={30} w="100%" />
      </VStack>
    );
  }

  const GrantItems = () => {
    const items = [];
    switch (request?.grant?.status) {
      case "PENDING":
        items.push(
          <CFTimelineRow
            arrLength={2}
            header={"Grant Created "}
            index={2}
            body={new Date(request.updatedAt).toString()}
          />
        );
      case "ACTIVE":
        items.push(
          <CFTimelineRow
            arrLength={2}
            header={"Grant Activated "}
            index={2}
            body={new Date(request.updatedAt).toString()}
          />
        );
    }
    return <>{items}</>;
  };

  return (
    <VStack flex={1} align="left">
      <Box textStyle="Heading/H4" as="h4" mb={8}>
        Audit Log
      </Box>
      {request?.grant && <GrantItems />}
      {request?.status !== "PENDING" && (
        <CFTimelineRow
          arrLength={2}
          header={
            "Request " +
            request.status.toLowerCase().charAt(0).toUpperCase() +
            request.status.toLowerCase().slice(1)
          }
          index={2}
          body={new Date(request.updatedAt).toString()}
        />
      )}
      <CFTimelineRow
        arrLength={2}
        header={"Request Created"}
        index={1}
        body={new Date(request.requestedAt).toString()}
      />
    </VStack>
  );
};
