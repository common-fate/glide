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
    setHours(Math.floor(value / HOUR));
    setMinutes(Math.floor((value % HOUR) / MINUTE));
  }, [value]);

  const setValue = (d: DurationInterval, v: number) => {
    switch (d) {
      case "HOUR":
        let newTime = v * HOUR + minutes * MINUTE;
        if (max && newTime > max) {
          onChange(
            v * HOUR +
              Math.min(Math.floor((max - v * HOUR) / MINUTE), 59) * MINUTE
          );
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
  const { maxHours, hours, setValue, register } = useContext(Context);
  const [defaultValue] = useState(hours);
  useEffect(() => {
    register("HOUR");
  });
  return (
    <NumberInput
      // variant="reveal"
      precision={0}
      defaultValue={defaultValue}
      min={0}
      step={1}
      role="group"
      max={maxHours}
      w="100px"
      value={hours}
      onChange={(s: string, n: number) => setValue("HOUR", n)}
      className="peer"
      // prevent chackra component from controlling the value on blur because we fully control the values via the context
      onBlur={undefined}
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
        hrs
      </InputRightElement>
      <NumberInputStepper>
        <NumberIncrementStepper />
        <NumberDecrementStepper />
      </NumberInputStepper>
    </NumberInput>
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
    <NumberInput
      // variant="reveal"
      precision={0}
      defaultValue={defaultValue}
      max={maxMinutes}
      min={minMinutes}
      step={1}
      role="group"
      w="100px"
      value={minutes}
      onChange={(s: string, n: number) => setValue("MINUTE", n)}
      className="peer"
      // prevent chackra component from controlling the value on blur because we fully control the values via the context
      onBlur={undefined}
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
        mins
      </InputRightElement>
      <NumberInputStepper>
        <NumberIncrementStepper />
        <NumberDecrementStepper />
      </NumberInputStepper>
    </NumberInput>
  );
};
