import { FormControl, FormLabel, Text, VStack } from "@chakra-ui/react";
import React from "react";
import { Controller, useFormContext } from "react-hook-form";
import { durationString } from "../../../../utils/durationString";
import {
  Days,
  DurationInput,
  Hours,
  Minutes,
  Weeks,
} from "../../../DurationInput";
import { AccessRuleFormData } from "../CreateForm";
import { FormStep } from "./FormStep";

export const TimeStep: React.FC = () => {
  const methods = useFormContext<AccessRuleFormData>();
  const time = methods.watch("timeConstraints");
  const sixMonthsMaxInSeconds =
    26 * // number of weeks in 6 months
    7 * // number of days in a week
    24 * // hours
    3600; // seconds in an hour

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
            max: sixMonthsMaxInSeconds,
            min: 60,
          }}
          name="timeConstraints.maxDurationSeconds"
          render={({ field: { ref, ...rest } }) => {
            console.log(sixMonthsMaxInSeconds, "maxDurationSeconds");
            return (
              <>
                <DurationInput
                  {...rest}
                  max={sixMonthsMaxInSeconds}
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
