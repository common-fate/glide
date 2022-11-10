import { FormControl, FormLabel, Text, VStack } from "@chakra-ui/react";
import React from "react";
import { Controller, useFormContext } from "react-hook-form";
import { durationString } from "../../../../utils/durationString";
import {
  Days,
  DurationInput,
  Hours,
  Minutes,
  Months,
  Weeks,
} from "../../../DurationInput";
import { AccessRuleFormData } from "../CreateForm";
import { FormStep } from "./FormStep";

export const TimeStep: React.FC = () => {
  const methods = useFormContext<AccessRuleFormData>();
  const time = methods.watch("timeConstraints");
  const maxDurationSeconds = 6 * 4 * 7 * 24 * 3600; // 6 months in seconds

  return (
    <FormStep
      heading="Time"
      subHeading="How long can access be requested for?"
      fields={["timeConstraints.maxDurationSeconds"]}
      preview={
        <VStack width={"100%"} align="flex-start">
          <Text textStyle={"Body/Medium"} color="neutrals.600">
            Max duration:{"  "}
            {time && durationString(time?.maxDurationSeconds)}
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
              <>
                <DurationInput
                  {...rest}
                  max={maxDurationSeconds}
                  min={60}
                  defaultValue={3600}
                >
                  <Weeks />
                  <Days />
                  <Hours />
                  <Minutes />
                </DurationInput>
              </>
            );
          }}
        />
      </FormControl>
    </FormStep>
  );
};
