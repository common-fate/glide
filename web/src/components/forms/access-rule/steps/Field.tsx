import { FormControl, FormLabel, HStack, Text, VStack } from "@chakra-ui/react";
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
import { TargetGroupField } from "../components/TargetGroupField";

export const FieldStep: React.FC = () => {
  const methods = useFormContext<AccessRuleFormData>();

  return (
    <FormStep
      heading="TargetGroupField"
      subHeading="abcd"
      fields={["field"]}
      preview={
        <VStack width={"100%"} align="flex-start">
          <Text textStyle={"Body/Medium"} color="neutrals.600">
            This will contain some items + more
          </Text>
        </VStack>
      }
    >
      <FormControl
        isInvalid={
          !!methods.formState.errors.timeConstraints?.maxDurationSeconds
        }
      ></FormControl>
    </FormStep>
  );
};
