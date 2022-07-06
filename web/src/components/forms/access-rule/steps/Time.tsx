import {
  FormControl,
  FormErrorMessage,
  FormLabel,
  NumberDecrementStepper,
  NumberIncrementStepper,
  NumberInput,
  NumberInputField,
  NumberInputStepper,
  Text,
  VStack,
} from "@chakra-ui/react";
import React from "react";
import { Controller, useFormContext } from "react-hook-form";
import { FormStep } from "./FormStep";

export const TimeStep: React.FC = () => {
  const methods = useFormContext();
  const time = methods.watch("timeConstraints");
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
              ? time.maxDurationSeconds / 60 / 60 + " hours"
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
          <Text textStyle={"Body/Medium"}>Maximum Duration (hours)</Text>
        </FormLabel>
        <Controller
          control={methods.control}
          rules={{
            required: true,
            // these need to be in seconds
            min: 15 * 60,
            max: 12 * 60 * 60,
            // ensure value is divisible by 15 minutes
            validate: (v) => {
              return v % (15 * 60) === 0;
            },
          }}
          defaultValue={0}
          name={"timeConstraints.maxDurationSeconds"}
          render={({ field: { ref, onChange, name, value } }) => (
            <NumberInput
              step={0.25}
              w="200px"
              value={value / 60 / 60}
              name={name}
              ref={ref}
              onChange={(s, n) => {
                onChange(n * 60 * 60);
              }}
              onBlur={() =>
                methods.trigger("timeConstraints.maxDurationSeconds")
              }
            >
              <NumberInputField bg="neutrals.0" />
              <NumberInputStepper>
                <NumberIncrementStepper />
                <NumberDecrementStepper />
              </NumberInputStepper>
            </NumberInput>
          )}
        />

        <FormErrorMessage>
          Duration must be in 0.25 hour increments. Minimum duration 0.25 hours.
          Max duration 12 hours.
        </FormErrorMessage>
      </FormControl>
    </FormStep>
  );
};
