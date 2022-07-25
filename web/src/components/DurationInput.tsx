import { HStack, Stack, VStack } from "@chakra-ui/layout";
import {
  NumberInput,
  NumberInputField,
  InputRightElement,
  NumberInputStepper,
  NumberIncrementStepper,
  NumberDecrementStepper,
  Text,
  Skeleton,
} from "@chakra-ui/react";
import { setHours } from "date-fns";
import React, { createContext, useContext, useEffect, useState } from "react";

interface DurationInputProps {
  onChange: (n: number) => void;
  /**  maximum duration in seconds*/
  max?: number;
  /**  minimum duration in seconds, defaults to 0s when not provided*/
  min?: number;
  /** value, provide this to control the component */
  initialValue?: number;
  isLoading?: boolean;
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
  loading: boolean;
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
  loading: true,
});
const HOUR = 3600;
const MINUTE = 60;

export const DurationInput: React.FC<DurationInputProps> = ({
  children,
  onChange,
  initialValue,
  max,
  min: minv,
  isLoading,
}) => {
  console.log({ minv, max });
  const initial = initialValue || 0;
  const min = minv || 0;
  const [hours, setHours] = useState<number>(Math.floor(initial / HOUR));
  const [minutes, setMinutes] = useState<number>(
    Math.floor((initial % HOUR) / MINUTE)
  );
  const [hasHours, setHasHours] = useState(false);
  const [hasMinutes, setHasMinutes] = useState(false);
  const v = hours * HOUR + minutes * MINUTE;
  useEffect(() => {
    console.log({ h: initial / HOUR });
    setHours(Math.floor(initial / HOUR));
    setMinutes(Math.floor((initial % HOUR) / MINUTE));
  }, [initialValue]);

  let maxHours = undefined;
  let maxMinutes = undefined;
  if (hasHours && max != undefined) {
    maxHours = Math.floor(max / HOUR);
  }
  if (hasMinutes) {
    if (hasHours) {
      if (max == undefined) {
        // if the hours component is available, but no max is set, then 59 minutes is the maximum
        maxMinutes = 59;
      } else {
        // if a max is set and the hours component available, then get the minimum of 59 or the remainder of minutes from (the max - the current value) after removing hours
        maxMinutes = Math.min(Math.floor((max - v) % HOUR), 59);
      }
    } else if (max != undefined) {
      // if there is no hours component, and max is defined, then get the minutes component of the max
      maxMinutes = Math.floor(max / MINUTE);
    }
  }
  const minHours = Math.floor(min / HOUR);
  const minMinutesFn = (duration: number, minDurationSeconds: number) =>
    duration < HOUR ? Math.floor((minDurationSeconds % HOUR) / MINUTE) : 0; // the minute component of the min , e.g if min is 3540 then min minutes it 59, if min is 3600 then min minutes is 0
  const minMinutes = minMinutesFn(v, min);

  const setValue = (d: DurationInterval, v: number) => {
    switch (d) {
      case "HOUR":
        setHours(v);
        if (v + minutes * MINUTE < min) {
          setMinutes(minMinutesFn(v + minutes * MINUTE, min));
        }
        break;
      case "MINUTE":
        setMinutes(v);
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
        loading: isLoading || false,
      }}
    >
      <VStack>
        <Text>
          other max {max} min {min} value {initialValue} v {v}
        </Text>
        <Text>
          hours {hours} {maxHours} {minHours}
        </Text>
        <Text>
          minutes {minutes} {maxMinutes} {minMinutes}
        </Text>
      </VStack>
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
      // onBlur={onBlurFn}
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

  console.log({ minMinutes, minutes });
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
      // onBlur={onBlurFn}
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
