import { FormControl, FormLabel, Text, VStack } from "@chakra-ui/react";
import React from "react";
import { Controller, useFormContext } from "react-hook-form";
import { DurationInput, Hours, Minutes } from "../../../DurationInput";
import { FormStep } from "./FormStep";

export const TimeStep: React.FC = () => {
  const methods = useFormContext();
  const time = methods.watch("timeConstraints");
  const maxDurationSeconds = 24 * 3600;

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
              ? Math.floor(time.maxDurationSeconds / 60 / 60) +
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
          rules={{
            required: "Duration is required.",
            max: maxDurationSeconds,
            min: 60,
          }}
          name="timeConstraints.maxDurationSeconds"
          render={({ field: { ref, ...rest } }) => {
            return (
              <DurationInput
                {...rest}
                max={maxDurationSeconds}
                min={60}
                defaultValue={time?.maxDurationSeconds}
              >
                <Hours />
                <Minutes />
              </DurationInput>
            );
          }}
        />
      </FormControl>
    </FormStep>
  );
};
