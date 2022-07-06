import { RequestTiming } from "./backend-client/types";
import { durationString } from "./durationString";

export const renderTiming = (timing: RequestTiming | undefined): string => {
  if (timing === undefined) return "";
  const duration = durationString(timing?.durationSeconds);
  const startTime = timing.startTime
    ? new Date(timing.startTime).toString()
    : undefined;

  return startTime
    ? `${duration}, starting at ${startTime}`
    : `${duration}, starting ASAP`;
};
