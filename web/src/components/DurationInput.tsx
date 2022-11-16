import {
  HStack,
  InputGroup,
  InputRightElement,
  NumberDecrementStepper,
  NumberIncrementStepper,
  NumberInput,
  NumberInputField,
  NumberInputStepper,
} from "@chakra-ui/react";
import React, {
  createContext,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";

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
  /** set this to true to omit rendering input elements which are outside the range of the max value */
  hideUnusedElements?: boolean;
}

type DurationInterval = "MINUTE" | "HOUR" | "DAY" | "WEEK";
interface DurationInputContext {
  maxHours?: number;
  maxMinutes?: number;
  maxDays?: number;
  maxWeeks?: number;
  minHours: number;
  minMinutes: number;
  minDays: number;
  minWeeks: number;
  minutes: number;
  hours: number;
  days: number;
  weeks: number;
  /** true if the component should render */
  shouldRenderMinutesInput: boolean;
  /** true if the component should render */
  shouldRenderHoursInput: boolean;
  /** true if the component should render */
  shouldRenderDaysInput: boolean;
  /** true if the component should render */
  shouldRenderWeeksInput: boolean;
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
  hours: 0,
  minutes: 0,
  days: 0,
  weeks: 0,
  shouldRenderWeeksInput: false,
  shouldRenderDaysInput: false,
  shouldRenderHoursInput: false,
  shouldRenderMinutesInput: true,
});
const MINUTE = 60;
const HOUR = 3600;
const DAY = 86400;
const WEEK = 7 * DAY;

const minMinutesFn = (duration: number, minDurationSeconds: number) =>
  duration < HOUR ? Math.floor((minDurationSeconds % HOUR) / MINUTE) : 0;

const maxMinutesFn = (
  hasHours: boolean,
  days: number,
  hours: number,
  weeks: number,
  maxDurationSeconds?: number
) => {
  if (maxDurationSeconds == undefined) {
    // if the hours component is available, but no max is set, then 59 minutes is the maximum
    return 59;
  } else {
    // if a max is set and the hours component available, then get the minimum of 59 or the remainder of minutes from (the max - the current value) after removing hours
    return maxDurationSeconds < HOUR
      ? Math.floor(maxDurationSeconds / MINUTE)
      : Math.min(
          Math.floor(
            (maxDurationSeconds - weeks * WEEK - days * DAY - hours * HOUR) /
              MINUTE
          ),
          59
        );
  }
  // return undefined;
};

// max constraints

const maxHoursFn = (
  hasDays: boolean,
  days: number,
  hours: number,
  weeks: number,
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
            Math.floor((maxDurationSeconds - weeks * WEEK - days * DAY) / HOUR),
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
  hours: number,
  weeks: number,
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
  hideUnusedElements,
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

  // The children components can register which means you can use this duration input with hours, minutes or both
  const [hasMinutes, setHasMinutes] = useState(false);
  const [hasHours, setHasHours] = useState(false);
  const [hasDays, setHasDays] = useState(false);
  const [hasWeeks, setHasWeeks] = useState(false);

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
  }, [value, hasHours, hasDays, hasWeeks]);

  // setValue checks whether the change to one field needs to affect the other field
  // e.g if reducing an hour to 0 does the minute field need to be increased
  // the validation logic on the input components themselves handle "most" of the actual validation
  // however they are not aware of each other, so edge cases are handled in here
  const setValue = (d: DurationInterval, v: number) => {
    switch (d) {
      case "MINUTE": {
        const newTime = weeks * WEEK + days * DAY + hours * HOUR + v * MINUTE;
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
        } else if (newTime < min) {
          onChange(
            weeks * WEEK +
              days * DAY +
              v * HOUR +
              minMinutesFn(newTime, min) * MINUTE
          );
        } else {
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
          onChange(min);
        } else {
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
    }
  };

  const maxMinutes = useMemo(
    () =>
      hasMinutes ? maxMinutesFn(hasHours, days, hours, weeks, max) : undefined,
    [hasMinutes, hasHours, days, hours, weeks, minutes, max]
  );
  const maxHours = useMemo(
    () => (hasHours ? maxHoursFn(hasDays, days, hours, weeks, max) : undefined),
    [hasHours, hasDays, days, hours, weeks, minutes, max]
  );
  const maxDays = useMemo(
    () => (hasDays ? maxDaysFn(hasWeeks, days, hours, weeks, max) : undefined),
    [hasDays, hasWeeks, days, hours, weeks, max]
  );
  const maxWeeks =
    hasWeeks && max != undefined ? Math.floor(max / WEEK) : undefined;

  // min constraints
  const minMinutes = minMinutesFn(value, min);
  const minHours = value < DAY ? Math.floor((min % DAY) / HOUR) : 0;
  const minDays = value < WEEK ? Math.floor((min % WEEK) / DAY) : 0;
  const minWeeks = Math.floor(min / WEEK);

  return (
    <Context.Provider
      value={{
        setValue,
        register,
        minMinutes,
        minHours,
        minDays,
        minWeeks,
        maxMinutes,
        maxHours,
        maxDays,
        maxWeeks,
        minutes,
        hours,
        days,
        weeks,
        shouldRenderDaysInput: hideUnusedElements
          ? !!(max && max >= DAY)
          : true,
        shouldRenderHoursInput: hideUnusedElements
          ? !!(max && max >= HOUR)
          : true,
        shouldRenderMinutesInput: true,
        shouldRenderWeeksInput: hideUnusedElements
          ? !!(max && max >= WEEK)
          : true,
      }}
    >
      <HStack>{children}</HStack>
    </Context.Provider>
  );
};

export const Weeks: React.FC = () => {
  const {
    maxWeeks,
    minWeeks,
    weeks,
    setValue,
    register,
    shouldRenderWeeksInput,
  } = useContext(Context);
  const [defaultValue] = useState(weeks);
  useEffect(() => {
    shouldRenderWeeksInput && register("WEEK");
  }, [shouldRenderWeeksInput]);
  if (!shouldRenderWeeksInput) {
    return null;
  }
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
  const {
    maxDays,
    minDays,
    days,
    setValue,
    register,
    shouldRenderDaysInput,
  } = useContext(Context);
  const [defaultValue] = useState(days);

  useEffect(() => {
    shouldRenderDaysInput && register("DAY");
  }, [shouldRenderDaysInput]);
  if (!shouldRenderDaysInput) {
    return null;
  }
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
  const {
    maxHours,
    minHours,
    hours,
    setValue,
    register,
    shouldRenderHoursInput,
  } = useContext(Context);
  const [defaultValue] = useState(hours);
  useEffect(() => {
    shouldRenderHoursInput && register("HOUR");
  }, [shouldRenderHoursInput]);
  if (!shouldRenderHoursInput) {
    return null;
  }
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
