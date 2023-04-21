import { durationString } from "./durationString";
// import { RequestTiming } from "./backend-client/types";

// export const renderTiming = (timing: RequestTiming | undefined): string => {
//   if (timing === undefined) return "";
//   const duration = durationString(timing?.durationSeconds);
//   const startTime = timing.startTime
//     ? new Date(timing.startTime).toString()
//     : undefined;

//   return startTime
//     ? `${duration}, starting at ${startTime}`
//     : `${duration}, starting ASAP`;
// };

export const renderTiming = (timing: any): string => "deprecated";
