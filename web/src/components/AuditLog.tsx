import { Box, VStack } from "@chakra-ui/layout";
import { Skeleton } from "@chakra-ui/react";
import React, { useMemo } from "react";
import { useListRequestEvents } from "../utils/backend-client/end-user/end-user";
import { RequestDetail } from "../utils/backend-client/types";
import { CFTimelineRow } from "./CFTimelineRow";
export const AuditLog: React.FC<{ request?: RequestDetail }> = ({
  request,
}) => {
  const { data } = useListRequestEvents(request?.id || "");
  const events = useMemo(() => {
    const items: {
      timestamp: string;
      react: JSX.Element;
    }[] = [];
    data?.events.forEach((e) => {
      if (e.grantCreated) {
        items.push({
          timestamp: e.createdAt,
          react: (
            <CFTimelineRow
              arrLength={2}
              header={"Grant created"}
              index={2}
              body={new Date(e.createdAt).toString()}
            />
          ),
        });
      }
      if (e.fromGrantStatus) {
        items.push({
          timestamp: e.createdAt,
          react: (
            <CFTimelineRow
              arrLength={2}
              header={`Grant status changed from ${e.fromGrantStatus} to ${e.toGrantStatus}`}
              index={2}
              body={new Date(e.createdAt).toString()}
            />
          ),
        });
      }
      if (e.fromStatus) {
        items.push({
          timestamp: e.createdAt,
          react: (
            <CFTimelineRow
              arrLength={2}
              header={`${
                e.actor || "Granted Approvals"
              } changed request status from ${e.fromStatus} to ${e.toStatus}`}
              index={2}
              body={new Date(e.createdAt).toString()}
            />
          ),
        });
      }
      if (e.requestCreated) {
        items.push({
          timestamp: e.createdAt,
          react: (
            <CFTimelineRow
              arrLength={2}
              header={`Request created by ${e.actor}`}
              index={2}
              body={new Date(e.createdAt).toString()}
            />
          ),
        });
      }
    });
    items.sort();
    return items.map((i) => i.react);
  }, [data]);
  console.log(events);
  if (!request || data === undefined) {
    return (
      <VStack flex={1} align="left">
        <Box textStyle="Heading/H4" as="h4" mb={8}>
          Audit Log
        </Box>
        <Skeleton h={30} w="100%" />
      </VStack>
    );
  }

  return (
    <VStack flex={1} align="left">
      <Box textStyle="Heading/H4" as="h4" mb={8}>
        Audit Log
      </Box>
      {events}
    </VStack>
  );
};
