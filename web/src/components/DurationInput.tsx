import { HStack } from "@chakra-ui/layout";
import {
  forwardRef,
  InputRightElement,
  NumberDecrementStepper,
  NumberIncrementStepper,
  NumberInput,
  NumberInputField,
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

type DurationInterval = "MINUTE" | "HOUR";
interface DurationInputContext {
  maxHours?: number;
  maxMinutes?: number;
  minHours: number;
  minMinutes: number;
  hours: number;
  minutes: number;
  setValue: (d: DurationInterval, v: number) => void;
  // Register should be called once on mount of the child duration intervals hours or minutes etc
  register: (d: DurationInterval) => void;
}

const Context = createContext<DurationInputContext>({
  setValue: (a, b) => {},
  register: (d) => {},
  minHours: 0,
  minMinutes: 0,
  hours: 0,
  minutes: 0,
});
const HOUR = 3600;
const MINUTE = 60;

const maxMinutesFn = (
  hasHours: boolean,
  hours: number,
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
            Math.floor((maxDurationSeconds - hours * HOUR) / MINUTE),
            59
          );
    }
  } else if (maxDurationSeconds != undefined) {
    // if there is no hours component, and max is defined, then get the minutes component of the max
    return Math.floor(maxDurationSeconds / MINUTE);
  }
  return undefined;
};
const minMinutesFn = (duration: number, minDurationSeconds: number) =>
  duration < HOUR ? Math.floor((minDurationSeconds % HOUR) / MINUTE) : 0;

/*
  DurationInput is intended to be a composable duration input element, it can be used with either hour minute or both hours and minutes.
  In future we may wish to add Days as well.

usage example 
  <DurationInput>
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
  const defaultValue = dv ?? 0;
  const value = v ?? 0;
  const min = minv || 0;
  const [hours, setHours] = useState<number>(Math.floor(defaultValue / HOUR));
  const [minutes, setMinutes] = useState<number>(
    Math.floor((defaultValue % HOUR) / MINUTE)
  );

  // The children components can register which means you can use this duration input with hours, minutes or both
  const [hasHours, setHasHours] = useState(false);
  const [hasMinutes, setHasMinutes] = useState(false);

  useEffect(() => {
    // The following effect updates the hours and minutes values when the external value changes after a call to onChange
    // it supports having eitehr hours and minutes or just hours or just minutes component
    if (hasHours) {
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
  }, [value, hasHours]);

  const setValue = (d: DurationInterval, v: number) => {
    switch (d) {
      case "HOUR":
        let newTime = v * HOUR + minutes * MINUTE;
        if (max && newTime > max) {
          onChange(
            v * HOUR +
              Math.min(Math.floor((max - v * HOUR) / MINUTE), 59) * MINUTE
          );
        } else if (newTime < min) {
          onChange(v * HOUR + minMinutesFn(newTime, min) * MINUTE);
        } else {
          onChange(v * HOUR + minutes * MINUTE);
        }

        break;
      case "MINUTE":
        onChange(hours * HOUR + v * MINUTE);
        break;
    }
  };
  const register = (d: DurationInterval) => {
    switch (d) {
      case "HOUR":
        setHasHours(true);
        break;
      case "MINUTE":
        setHasMinutes(true);
        break;
    }
  };

  const maxHours =
    hasHours && max != undefined ? Math.floor(max / HOUR) : undefined;
  const maxMinutes = hasMinutes
    ? maxMinutesFn(hasHours, hours, max)
    : undefined;
  const minHours = Math.floor(min / HOUR);
  const minMinutes = minMinutesFn(value, min);
  return (
    <Context.Provider
      value={{
        setValue,
        register,
        minHours,
        minMinutes,
        maxHours,
        maxMinutes,
        hours,
        minutes,
      }}
    >
      <HStack>{children}</HStack>
    </Context.Provider>
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
  max?: number;
  min?: number;
  defaultValue: number;
  value: number;
  onChange: (n: number) => void;
  rightElement?: React.ReactNode;
}
const InputElement: React.FC<InputElementProps> = ({
  defaultValue,
  onChange,
  value,
  max,
  min,
  rightElement,
}) => {
  const [v, setV] = useState<string | number>(value);
  useEffect(() => {
    if (typeof v === "string" || v != value) {
      setV(value);
    }
  }, [value]);
  return (
    <NumberInput
      // variant="reveal"
      precision={0}
      defaultValue={defaultValue}
      max={max}
      min={min}
      step={1}
      role="group"
      w="100px"
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
    >
      <NumberInputField bg="white" />
      <InputRightElement
        pos="absolute"
        right={10}
        w="8px"
        color="neutrals.500"
        userSelect="none"
        textAlign="left"
      >
        {rightElement}
      </InputRightElement>
      <NumberInputStepper>
        <NumberIncrementStepper />
        <NumberDecrementStepper />
      </NumberInputStepper>
    </NumberInput>
  );
};
