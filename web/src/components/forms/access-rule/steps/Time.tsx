import {
  FormControl,
  FormLabel,
  HStack,
  InputRightElement,
  NumberDecrementStepper,
  NumberIncrementStepper,
  NumberInput,
  NumberInputField,
  NumberInputStepper,
  Text,
  VStack,
} from "@chakra-ui/react";
import React, { useState } from "react";
import { Controller, useFormContext } from "react-hook-form";
import { FormStep } from "./FormStep";

export const TimeStep: React.FC = () => {
  const methods = useFormContext();
  const time = methods.watch("timeConstraints");
  const maxDurationSeconds = 24 * 3600;

  const [hours, setHours] = useState(0);
  const [mins, setMins] = useState(0);

  return (
    <FormStep
      heading="Time"
      subHeading="How long and when can access be requested?"
      fields={["timeConstraints.maxDurationSeconds"]}
      preview={
        <VStack width={"100%"} align="flex-start">
          <Text textStyle={"Body/Medium"} color="neutrals.600">
            Max duration:{" "}
            {time?.maxDurationSeconds
              ? time.maxDurationSeconds / 60 / 60 +
                " hours " +
                ((time.maxDurationSeconds / 60) % 60) +
                " minutes"
              : ""}
          </Text>
        </VStack>
      }
    >
      <FormControl
        isInvalid={
          !!methods.formState.errors.timeConstraints?.maxDurationSeconds
        }
      >
        <FormLabel htmlFor="timeConstraints.maxDurationSeconds">
          <Text textStyle={"Body/Medium"}>Maximum Duration </Text>
        </FormLabel>
        <Controller
          control={methods.control}
          rules={{ required: "Duration is required." }}
          defaultValue={0}
          name="timeConstraints.maxDurationSeconds"
          render={({ field, fieldState }) => {
            const onBlurFn = () => {
              const duration = hours * 60 * 60 + mins * 60;

              if (maxDurationSeconds && duration > maxDurationSeconds) {
                methods.setValue("timeConstraints.maxDurationSeconds", 0);

                // DE = when an out of bounds value is adjusted to maxSeconds, we need to update the hours and mins to match
                // Firstly calculate what the hours would be
                let h = maxDurationSeconds / 60 / 60;
                let m = (maxDurationSeconds / 60) % 60;

                setHours(h);
                setMins(m);

                // Invalidate the field
              } else {
                methods.setValue(
                  "timeConstraints.maxDurationSeconds",
                  duration
                );
              }
            };

            let maxH = maxDurationSeconds ? maxDurationSeconds / 3600 : 24;

            return (
              <HStack>
                <NumberInput
                  variant="reveal"
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
                  variant="reveal"
                  role="group"
                  defaultValue={1}
                  min={0}
                  step={1}
                  max={59}
                  w="100px"
                  value={mins}
                  onChange={(s, n) => {
                    if (hours * 3600 + mins * 60 >= maxDurationSeconds) {
                      return;
                    } else setMins(n);
                  }}
                  className="peer"
                  onBlur={onBlurFn}
                  onKeyDown={(e) => {
                    // allow stepping up from 59 to 0
                    if (e.key === "ArrowUp") {
                      if (mins === 59 && hours < maxH) {
                        setMins(0);
                        setHours((h) => h + 1);
                      }
                    } else if (e.key === "ArrowDown") {
                      if (mins === 0 && hours > 0) {
                        setMins(59);
                        setHours((h) => h - 1);
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
              </HStack>
            );
          }}
        />

        {/* <FormErrorMessage>
          Duration must be in 0.25 hour increments. Minimum duration 0.25 hours.
          Max duration 12 hours.
        </FormErrorMessage> */}
      </FormControl>
    </FormStep>
  );
};
