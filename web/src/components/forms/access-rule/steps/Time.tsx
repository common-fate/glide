import { QuestionIcon } from "@chakra-ui/icons";
import {
  FormControl,
  FormLabel,
  Text,
  Tooltip,
  VStack,
} from "@chakra-ui/react";
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
          <Text textStyle={"Body/Medium"} color="neutrals.600">
            Default duration:{"  "}
            {time && durationString(time?.defaultDurationSeconds)}
          </Text>
          {/* <Text textStyle={"Body/Medium"} color="neutrals.600">
            Suggested duration:{"  "}
            {time && durationString(time?.reccomdenedDurationSeconds)}
          </Text> */}
        </VStack>
      }
    >
      <>
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
        <FormControl
          isInvalid={
            !!methods.formState.errors.timeConstraints?.defaultDurationSeconds
          }
        >
          <FormLabel htmlFor="timeConstraints.defaultDurationSeconds">
            <Text textStyle={"Body/Medium"}>Default Duration </Text>
          </FormLabel>
          <Controller
            control={methods.control}
            rules={{
              required: "Duration is required.",
              max: methods.watch("timeConstraints.maxDurationSeconds"),
              min: 60,
            }}
            name="timeConstraints.defaultDurationSeconds"
            render={({ field: { ref, ...rest } }) => {
              return (
                <>
                  <DurationInput
                    {...rest}
                    max={methods.watch("timeConstraints.maxDurationSeconds")}
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
      </>
      {/* <FormControl
        isInvalid={
          !!methods.formState.errors.timeConstraints?.reccomdenedDurationSeconds
        }
      >
        <FormLabel
          display="flex"
          flexDir="row"
          htmlFor="timeConstraints.reccomdenedDurationSeconds"
          alignItems="center"
        >
          <Text textStyle={"Body/Medium"}>Recommended Duration </Text>
          <Tooltip
            hasArrow={true}
            label="The default session duration displayed when making an Access Request. Requestors can choose to adjust this on a per-request basis"
          >
            <QuestionIcon color="neutrals.800" ml={1} />
          </Tooltip>
        </FormLabel>
        <Controller
          control={methods.control}
          rules={{
            required: "Duration is required.",
            max: sixMonthsMaxInSeconds,
            min: 60,
          }}
          name="timeConstraints.reccomdenedDurationSeconds"
          render={({ field: { ref, ...rest } }) => {
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
      </FormControl> */}
    </FormStep>
  );
};
