import {
  Flex,
  FormControl,
  FormErrorMessage,
  FormLabel,
  InputRightElement,
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
          <Text textStyle={"Body/Medium"}>Maximum Duration </Text>
        </FormLabel>
        <Controller
          control={methods.control}
          /**
           * @TODO:
           * after lunch
           * resolve the final issue witih the minute section, it should pass through nicely
           */
          rules={{
            required: true,
            // these need to be in seconds
            min: 15 * 60,
            max: 12 * 60 * 60,
            // ensure value is divisible by 15 minutes
            // validate: (v) => {
            //   return v % (15 * 60) === 0;
            // },
          }}
          defaultValue={0}
          name={"timeConstraints.maxDurationSeconds"}
          render={({ field: { ref, onChange, name, value } }) => {
            const NaN1 = Number.isNaN(value / 60 / 60);
            const NaN2 = Number.isNaN(value / 60 / 60);

            // let hours = (value / 60 / 60).toFixed(0).toString();
            let hours = Math.floor(value / 3600);
            let mins = Math.floor((value % 3600) / 60);

            return (
              <Flex>
                <NumberInput
                  step={1}
                  w="130px"
                  min={0}
                  mr={2}
                  // get only the hours from a milisecond value, 0 decimal places
                  value={NaN1 ? 0 : hours}
                  name={name}
                  ref={ref}
                  onChange={(s, n) => {
                    console.log({ n, value });
                    onChange(n * 60 * 60);
                  }}
                  onBlur={() =>
                    methods.trigger("timeConstraints.maxDurationSeconds")
                  }
                  sx={{
                    "#step": {
                      opacity: 0,
                    },
                    "_focusWithin": {
                      "#step": {
                        opacity: 1,
                      },
                    },
                  }}
                >
                  <NumberInputField bg="neutrals.0" id="rule-max-duration" />
                  <NumberInputStepper id="step">
                    <NumberIncrementStepper />
                    <NumberDecrementStepper />
                  </NumberInputStepper>
                  <InputRightElement
                    pos="absolute"
                    right={10}
                    color="neutrals.500"
                    userSelect="none"
                  >
                    hrs
                  </InputRightElement>
                </NumberInput>
                <NumberInput
                  step={10}
                  min={0}
                  w="130px"
                  value={(value / 60) % 60}
                  name={name}
                  ref={ref}
                  onChange={(s, n) => {
                    console.log({ n, value });
                    // if (n < 0) {
                    //   onChange(value - 10 * 60);
                    // } else {
                    onChange(value + 10 * 60);
                    //}
                  }}
                  onBlur={() =>
                    methods.trigger("timeConstraints.maxDurationSeconds")
                  }
                  sx={{
                    "#step": {
                      opacity: 0,
                    },
                    "_focusWithin": {
                      "#step": {
                        opacity: 1,
                      },
                    },
                  }}
                >
                  <NumberInputField bg="neutrals.0" id="rule-max-duration" />
                  <NumberInputStepper id="step">
                    <NumberIncrementStepper />
                    <NumberDecrementStepper />
                  </NumberInputStepper>
                  <InputRightElement
                    pos="absolute"
                    right={10}
                    color="neutrals.500"
                    userSelect="none"
                  >
                    mins
                  </InputRightElement>
                </NumberInput>
              </Flex>
            );
          }}
        />

        <FormErrorMessage>
          Duration must be in 0.25 hour increments. Minimum duration 0.25 hours.
          Max duration 12 hours.
        </FormErrorMessage>
      </FormControl>
    </FormStep>
  );
};
