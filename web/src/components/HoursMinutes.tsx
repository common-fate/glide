import { HStack } from "@chakra-ui/layout";
import {
  NumberInput,
  NumberInputField,
  InputRightElement,
  NumberInputStepper,
  NumberIncrementStepper,
  NumberDecrementStepper,
} from "@chakra-ui/react";
import React, { useEffect, useState } from "react";

type Props = {
  setValue: (n: number) => void;
  /** maxDurationSeconds in seconds, this must be set for onBlurFn validation to occur */
  maxDurationSeconds?: number;
  /** intitialValue that hydrates the component with hours/mins */
  initialValue?: number;
  rightElement?: React.ReactNode;
};

const HoursMinutes = ({
  maxDurationSeconds,
  setValue,
  initialValue,
  rightElement,
}: Props) => {
  const [hours, setHours] = useState<number>();
  const [mins, setMins] = useState<number>();

  useEffect(() => {
    if (
      initialValue != undefined &&
      hours === undefined &&
      mins === undefined
    ) {
      setHours(Math.trunc((initialValue / 60 / 60)));
      setMins((initialValue / 60) % 60);
    }
  }, [initialValue]);

  let maxH = maxDurationSeconds ? maxDurationSeconds / 3600 : 24;

  const onBlurFn = () => {
    if (hours != undefined && mins != undefined) {
      const duration = hours * 60 * 60 + mins * 60;
      // DE = when an out of bounds value is adjusted to maxSeconds, we need to update the hours and mins to match
      if (maxDurationSeconds && duration > maxDurationSeconds) {
        setValue(maxDurationSeconds);
        let h = Math.trunc(maxDurationSeconds / 60 / 60).;
        let m = (maxDurationSeconds / 60) % 60;

        setHours(h);
        setMins(m);
        // DE =  If the duration is under 60s, set the time to one minute
      } else if (duration < 60 || hours + mins == 0) {
        setValue(60);
        setHours(0);
        setMins(1);
      } else {
        setValue(duration);
      }
    } else {
      console.log("cannot update values before initialValue has been set");
    }
  };

  return (
    <HStack>
      <NumberInput
        // variant="reveal"
        precision={0}
        defaultValue={1}
        min={0}
        step={1}
        role="group"
        max={maxH}
        w="100px"
        value={hours}
        onChange={(s, n) => setHours(n)}
        className="peer"
        onBlur={onBlurFn}
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
      <NumberInput
        // variant="reveal"
        role="group"
        precision={0}
        defaultValue={1}
        min={hours == 0 ? 1 : 0}
        step={1}
        max={59}
        w="100px"
        value={mins}
        onChange={(s, n) => {
          if (hours != undefined && mins != undefined) {
            if (hours * 3600 + mins * 60 >= maxDurationSeconds) {
              return;
            } else setMins(n);
          }
        }}
        onBlur={onBlurFn}
        onKeyDown={(e) => {
          if (hours != undefined) {
            // allow stepping up from 59 to 0
            if (e.key === "ArrowUp") {
              if (mins === 59 && hours < maxH) {
                setMins(0);
                setHours((h) => h && h + 1);
              }
            } else if (e.key === "ArrowDown") {
              if (mins === 0 && hours > 0) {
                setMins(59);
                setHours((h) => h && h - 1);
              }
            }
          }
        }}
      >
        <NumberInputField bg="white" />
        <NumberInputStepper>
          <NumberIncrementStepper />
          <NumberDecrementStepper />
        </NumberInputStepper>
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
      </NumberInput>
      {rightElement}
    </HStack>
  );
};

export default HoursMinutes;
