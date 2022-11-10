import {
  HStack,
  InputGroup,
  InputRightElement,
  NumberDecrementStepper,
  NumberIncrementStepper,
  NumberInput,
  NumberInputField,
  NumberInputProps,
  NumberInputStepper,
} from "@chakra-ui/react";
import React, { createContext, useContext, useEffect, useState } from "react";

interface DurationInputProps {
  onChange: (n: number) => void;
  /**  maximum duration in seconds*/
  max?: number;
  /**  minimum duration in seconds, defaults to 0s when not provided*/
  min?: number;
  /** value, provide this to control the component */
  value?: number;
  defaultValue?: number;
  children?: React.ReactNode;
}

type DurationInterval = "MINUTE" | "HOUR" | "DAY" | "WEEK" | "MONTH";
interface DurationInputContext {
  maxHours?: number;
  maxMinutes?: number;
  maxDays?: number;
  maxWeeks?: number;
  maxMonths?: number;
  minHours: number;
  minMinutes: number;
  minDays: number;
  minWeeks: number;
  minMonths: number;
  minutes: number;
  hours: number;
  days: number;
  weeks: number;
  months: number;
  setValue: (d: DurationInterval, v: number) => void;
  // Register should be called once on mount of the child duration intervals hours or minutes etc
  register: (d: DurationInterval) => void;
}

const Context = createContext<DurationInputContext>({
  setValue: (a, b) => {
    undefined;
  },
  register: (d) => {
    undefined;
  },
  minMinutes: 0,
  minHours: 0,
  minDays: 0,
  minWeeks: 0,
  minMonths: 0,
  hours: 0,
  minutes: 0,
  days: 0,
  weeks: 0,
  months: 0,
});
const MINUTE = 60;
const HOUR = 3600;
const DAY = 86400;
const WEEK = 7 * DAY;
const MONTH = 30 * DAY;

const minMinutesFn = (duration: number, minDurationSeconds: number) =>
  duration < HOUR ? Math.floor((minDurationSeconds % HOUR) / MINUTE) : 0;

/*
  DurationInput is intended to be a composable duration input element, it can be used with either hour minute or both hours and minutes.
  In future we may wish to add Days as well.

usage example 
  <DurationInput>
    <Weeks />
    <Days />
    <Hour>
    <Minute>
    <Text>
      some text on the right of the inputs
    </Text>
  </DurationInput>
  */

export const DurationInput: React.FC<DurationInputProps> = ({
  children,
  onChange,
  value: v,
  defaultValue: dv,
  max,
  min: minv,
}) => {
  const defaultValue = dv ?? minv ?? 0;
  const value = v ?? defaultValue;
  const min = minv || 0;
  const [minutes, setMinutes] = useState<number>(
    Math.floor((value % HOUR) / MINUTE)
  );
  const [hours, setHours] = useState<number>(Math.floor(value / HOUR));
  const [days, setDays] = useState<number>(Math.floor(value / DAY));
  const [weeks, setWeeks] = useState<number>(Math.floor(value / WEEK));
  const [months, setMonths] = useState<number>(Math.floor(value / MONTH));

  // The children components can register which means you can use this duration input with hours, minutes or both
  const [hasMinutes, setHasMinutes] = useState(false);
  const [hasHours, setHasHours] = useState(false);
  const [hasDays, setHasDays] = useState(false);
  const [hasWeeks, setHasWeeks] = useState(false);
  const [hasMonths, setHasMonths] = useState(false);

  // on first load, if v is undefined, call onChange with the default value to update the form
  useEffect(() => {
    if (v == undefined) {
      onChange(value);
    }
  }, [v, value]);

  useEffect(() => {
    // The following effect updates the hours and minutes values when the external value changes after a call to onChange
    // it supports having eitehr hours and minutes or just hours or just minutes components,
    // we prioritise the larger units (months, weeks, days, hours) and then the smaller units (minutes)
    if (hasWeeks) {
      setWeeks(Math.floor(value / WEEK));
      if (hasDays) {
        setDays(Math.floor((value % WEEK) / DAY));
        if (hasHours) {
          setHours(Math.floor(((value % WEEK) % DAY) / HOUR));
        }
        if (hasMinutes) {
          setMinutes(Math.floor((((value % WEEK) % DAY) % HOUR) / MINUTE));
        }
      } else if (hasHours) {
        setHours(Math.floor(value / HOUR));
        if (hasMinutes) {
          setMinutes(Math.floor((value % HOUR) / MINUTE));
        } else {
          setMinutes(0);
        }
      } else if (hasMinutes) {
        setHours(0);
        setMinutes(Math.floor(value / MINUTE));
      }
    }
  }, [value, hasHours, hasDays, hasWeeks]);

  // setValue checks whether the change to one field needs to affect the other field
  // e.g if reducing an hour to 0 does the minute field need to be increased
  // the validation logic on the input components themselves handle "most" of the actual validation
  // however they are not aware of each other, so edge cases are handled in here
  const setValue = (d: DurationInterval, v: number) => {
    switch (d) {
      case "MINUTE": {
        const newTime = weeks * WEEK + days * DAY + hours * HOUR + v * MINUTE;
        // should also do min/max handling  here.....
        if (max && newTime > max) {
          onChange(max);
        } else if (min && newTime < min) {
          onChange(min);
        } else {
          onChange(newTime);
        }
        break;
      }
      case "HOUR": {
        const newTime = weeks * WEEK + days * DAY + v * HOUR + minutes * MINUTE;
        if (max && newTime > max) {
          onChange(max);
          // onChange(
          //   weeks * WEEK +
          //     days * DAY +
          //     v * HOUR +
          //     Math.min(
          //       Math.floor((max - (days * DAY + v * HOUR)) / MINUTE),
          //       59
          //     ) *
          //       MINUTE
          // );
        } else if (newTime < min) {
          onChange(
            weeks * WEEK +
              days * DAY +
              v * HOUR +
              minMinutesFn(newTime, min) * MINUTE
          );
        } else {
          // onChange(weeks * WEEK + days * DAY + v * HOUR + minutes * MINUTE);
          onChange(newTime);
        }

        break;
      }
      case "DAY": {
        const newTime =
          weeks * WEEK + v * DAY + hours * HOUR + minutes * MINUTE;
        if (max && newTime > max) {
          onChange(max);
        } else if (newTime < min) {
          // onChange(
          //   weeks * WEEK + v * DAY + minMinutesFn(newTime, min) * MINUTE
          // );
          onChange(min);
        } else {
          // onChange(weeks * WEEK + v * DAY + hours * HOUR + minutes * MINUTE);
          onChange(newTime);
        }
        break;
      }
      case "WEEK": {
        const newTime = v * WEEK + days * DAY + hours * HOUR + minutes * MINUTE;
        if (max && newTime > max) {
          onChange(max);
        } else if (newTime < min) {
          onChange(v * WEEK + days * DAY + minMinutesFn(newTime, min) * MINUTE);
        } else {
          onChange(newTime);
        }
        break;
      }
    }
  };

  // Register is meant to register capability of the component, there may be a better way to work out if the minutes or hours components are present
  const register = (d: DurationInterval) => {
    switch (d) {
      case "MINUTE":
        setHasMinutes(true);
        break;
      case "HOUR":
        setHasHours(true);
        break;
      case "DAY":
        setHasDays(true);
        break;
      case "WEEK":
        setHasWeeks(true);
        break;
      case "MONTH":
        setHasMonths(true);
        break;
    }
  };

  const maxMinutesFn = (
    hasHours: boolean,
    hours: number,
    days: number,
    maxDurationSeconds?: number
  ) => {
    if (hasHours) {
      if (maxDurationSeconds == undefined) {
        // if the hours component is available, but no max is set, then 59 minutes is the maximum
        return 59;
      } else {
        // if a max is set and the hours component available, then get the minimum of 59 or the remainder of minutes from (the max - the current value) after removing hours
        return maxDurationSeconds < HOUR
          ? Math.floor(maxDurationSeconds / MINUTE)
          : Math.min(
              Math.floor(
                (maxDurationSeconds -
                  weeks * WEEK -
                  days * DAY -
                  hours * HOUR -
                  days * DAY) /
                  MINUTE
              ),
              59
            );
      }
    } else if (maxDurationSeconds != undefined) {
      // if there is no hours component, and max is defined, then get the minutes component of the max
      return Math.floor(maxDurationSeconds / MINUTE);
    }
    return undefined;
  };

  // max constraints
  const maxMinutes = hasMinutes
    ? maxMinutesFn(hasHours, hours, days, max)
    : undefined;

  const maxHoursFn = (
    hasDays: boolean,
    hours: number,
    maxDurationSeconds?: number
  ) => {
    if (hasDays) {
      if (maxDurationSeconds == undefined) {
        // if the hours component is available, but no max is set, then 23 hours is the maximum
        return 23;
      } else {
        // if a max is set and the hours component available, then get the minimum of 23 or the remainder of hours from (the max - the current value) after removing days
        return maxDurationSeconds < DAY
          ? Math.floor(maxDurationSeconds / HOUR)
          : Math.min(
              Math.floor(
                (maxDurationSeconds - weeks * WEEK - days * DAY) / HOUR
              ),
              23
            );
      }
    } else if (maxDurationSeconds != undefined) {
      // if there is no hours component, and max is defined, then get the minutes component of the max
      return Math.floor(maxDurationSeconds / HOUR);
    }
    return undefined;
  };
  const maxDaysFn = (
    hasWeeks: boolean,
    days: number,
    maxDurationSeconds?: number
  ) => {
    if (hasWeeks) {
      if (maxDurationSeconds == undefined) {
        // if the hours component is available, but no max is set, then 23 hours is the maximum
        return 6;
      } else {
        // if a max is set and the hours component available, then get the minimum of 23 or the remainder of hours from (the max - the current value) after removing days
        return maxDurationSeconds < WEEK
          ? Math.floor(maxDurationSeconds / DAY)
          : Math.min(Math.floor((maxDurationSeconds - weeks * WEEK) / DAY), 7);
      }
    } else if (maxDurationSeconds != undefined) {
      // if there is no hours component, and max is defined, then get the minutes component of the max
      return Math.floor(maxDurationSeconds / DAY);
    }
    return undefined;
  };

  const maxHours = hasHours ? maxHoursFn(hasDays, hours, max) : undefined;
  // todo: this should work, but we may want to be careful with handling cross over UX for incrementing superior units
  const maxDays = hasDays ? maxDaysFn(hasWeeks, days, max) : undefined;
  const maxWeeks =
    hasWeeks && max != undefined ? Math.floor(max / WEEK) : undefined;
  const maxMonths =
    hasMonths && max != undefined ? Math.floor(max / MONTH) : undefined;
  // min constraints
  const minMinutes = minMinutesFn(value, min);
  const minHours = hasMinutes
    ? Math.floor(min / HOUR)
    : min < HOUR
    ? 1
    : Math.floor(min / HOUR);
  const minDays = 0;
  const minWeeks = 0;
  const minMonths = 0;
  return (
    <Context.Provider
      value={{
        setValue,
        register,
        minMinutes,
        minHours,
        minDays,
        minWeeks,
        minMonths,
        maxMinutes,
        maxHours,
        maxDays,
        maxWeeks,
        maxMonths,
        minutes,
        hours,
        days,
        weeks,
        months,
      }}
    >
      <HStack>{children}</HStack>
    </Context.Provider>
  );
};

export const Months: React.FC = () => {
  const { maxMonths, minMonths, months, setValue, register } = useContext(
    Context
  );
  const [defaultValue] = useState(months);
  useEffect(() => {
    register("MONTH");
  });
  return (
    <InputElement
      inputId="month-duration-input"
      defaultValue={defaultValue}
      onChange={(n: number) => setValue("MONTH", n)}
      value={months}
      min={minMonths}
      max={maxMonths}
      rightElement="months"
    />
  );
};

export const Weeks: React.FC = () => {
  const { maxWeeks, minWeeks, weeks, setValue, register } = useContext(Context);
  const [defaultValue] = useState(weeks);
  useEffect(() => {
    register("WEEK");
  });
  return (
    <InputElement
      inputId="week-duration-input"
      defaultValue={defaultValue}
      onChange={(n: number) => setValue("WEEK", n)}
      value={weeks}
      min={minWeeks}
      max={maxWeeks}
      rightElement="weeks"
      w="122px"
    />
  );
};

export const Days: React.FC = () => {
  const { maxDays, minDays, days, setValue, register } = useContext(Context);
  const [defaultValue] = useState(days);
  useEffect(() => {
    register("DAY");
  });
  console.log({ maxDays });
  return (
    <InputElement
      inputId="day-duration-input"
      defaultValue={defaultValue}
      onChange={(n: number) => setValue("DAY", n)}
      value={days}
      max={maxDays}
      min={minDays}
      rightElement="days"
      // w="112px"
    />
  );
};
export const Hours: React.FC = () => {
  const { maxHours, minHours, hours, setValue, register } = useContext(Context);
  const [defaultValue] = useState(hours);
  useEffect(() => {
    register("HOUR");
  });
  return (
    <InputElement
      inputId="hour-duration-input"
      defaultValue={defaultValue}
      onChange={(n: number) => setValue("HOUR", n)}
      value={hours}
      max={maxHours}
      min={minHours}
      rightElement="hrs"
    />
  );
};
export const Minutes: React.FC = () => {
  const { maxMinutes, minMinutes, minutes, setValue, register } = useContext(
    Context
  );
  const [defaultValue] = useState(minutes);
  useEffect(() => {
    register("MINUTE");
  });

  return (
    <InputElement
      inputId="minute-duration-input"
      defaultValue={defaultValue}
      onChange={(n: number) => setValue("MINUTE", n)}
      value={minutes}
      max={maxMinutes}
      min={minMinutes}
      rightElement="mins"
    />
  );
};
interface InputElementProps {
  // input id is set on the input element if present
  inputId?: string;
  max?: number;
  min?: number;
  defaultValue: number;
  value: number;
  onChange: (n: number) => void;
  rightElement?: React.ReactNode;
  w?: string;
}
const InputElement: React.FC<InputElementProps> = ({
  inputId,
  defaultValue,
  onChange,
  value,
  max,
  min,
  rightElement,
  w,
}) => {
  const [v, setV] = useState<string | number>(value);
  useEffect(() => {
    if (typeof v === "string" || v != value) {
      setV(value);
    }
  }, [value]);
  return (
    <InputGroup w="unset">
      <NumberInput
        // variant="reveal"
        precision={0}
        id="minute-duration-input"
        defaultValue={defaultValue}
        max={max}
        min={min}
        step={1}
        role="group"
        width={w ?? "100px"}
        value={v}
        // if you backspace the value then click out, this resets the value to the current value
        onBlur={() => {
          if (typeof v === "string" || isNaN(v)) {
            setV(value);
          }
        }}
        onChange={(s: string, n: number) => {
          if (isNaN(n)) {
            setV(s);
          } else if (max && n > max) {
            // don't allow typed inputs greater than max
            setV(max);
            onChange(max);
          } else {
            setV(n);
            onChange(n);
          }
        }}
        className="peer"
        pos="relative"
      >
        <NumberInputField bg="white" id={inputId} />
        <InputRightElement
          pos="absolute"
          right={"40%"}
          w="8px"
          color="neutrals.500"
          userSelect="none"
          textAlign="left"
        >
          {rightElement}
        </InputRightElement>
        <NumberInputStepper>
          <NumberIncrementStepper id="increment" />
          <NumberDecrementStepper id="decrement" />
        </NumberInputStepper>
      </NumberInput>
    </InputGroup>
  );
};
