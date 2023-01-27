import { intervalToDuration, formatDuration } from "date-fns";

export const durationString = (durationSeconds?: number): string => {
  if (durationSeconds) {
    const d = intervalToDuration({ start: 0, end: durationSeconds * 1000 });

    return formatDuration(d);
    // In odd occasions where the duration is nullish,
    // we prefer to display an empty string than a NaN/invalid value
  } else return "";
};

export const durationStringHoursMinutes = (d?: Duration): string => {
  if (d) {
    if (
      !d.years &&
      !d.months &&
      !d.weeks &&
      !d.days &&
      !d.hours &&
      !d.minutes
    ) {
      return "few seconds";
    }
    return formatDuration(d, {
      format: ["days", "hours", "minutes"],
    });
  } else return "";
};
