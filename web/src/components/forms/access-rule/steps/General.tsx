import {
  FormControl,
  FormErrorMessage,
  FormLabel,
  Input,
  Text,
  Textarea,
  VStack,
} from "@chakra-ui/react";
import React from "react";
import { useFormContext } from "react-hook-form";
import { FormStep } from "./FormStep";

export const GeneralStep: React.FC = () => {
  const methods = useFormContext();
  const name = methods.watch("name");
  const description = methods.watch("description");
  return (
    <FormStep
      heading="General"
      subHeading="General information about this rule."
      fields={["name", "description"]}
      preview={
        <VStack width={"100%"} align="flex-start">
          <Text textStyle={"Body/Medium"} color="neutrals.600">
            Name: {name}
          </Text>
          <Text
            textStyle={"Body/Medium"}
            color="neutrals.600"
            wordBreak={"break-word"}
            flexWrap="wrap"
          >
            Description: {description}
          </Text>
        </VStack>
      }
    >
      <>
        <FormControl isInvalid={!!methods.formState.errors.name}>
          <FormLabel htmlFor="name">
            <Text textStyle={"Body/Medium"}>Name</Text>
          </FormLabel>
          <Input
            bg="neutrals.0"
            {...methods.register("name", {
              // required: true,
              validate: (value) => {
                const res: string[] = [];
                if (!value || value.length == 0) {
                  res.push("Field is required");
                }
                [/[^a-zA-Z0-9,.;:()[\]?!\-_`~&/\n\s]/].every((pattern) =>
                  pattern.test(value as string)
                ) &&
                  res.push(
                    "Invalid characters (only letters, numbers, and punctuation allowed)"
                  );
                if (value && value.length > 400) {
                  res.push("Maximum length is 400 characters");
                }
                return res.length > 0 ? res.join(", ") : undefined;
              },
            })}
            onBlur={() => void methods.trigger("name")}
          />
          {methods.formState.errors?.name?.message && (
            <FormErrorMessage>
              {methods.formState.errors.name?.message?.toString()}
            </FormErrorMessage>
          )}
        </FormControl>
        <FormControl isInvalid={!!methods.formState.errors.description}>
          <FormLabel htmlFor="Description">
            <Text textStyle={"Body/Medium"}>Description</Text>
          </FormLabel>
          <Textarea
            bg="neutrals.0"
            {...methods.register("description", {
              // required: true,
              validate: (value) => {
                const res: string[] = [];
                if (!value || value.length == 0) {
                  res.push("Field is required");
                }
                [/[^a-zA-Z0-9,.;:()[\]?!\-_`~&/\n\s]/].every((pattern) =>
                  pattern.test(value as string)
                ) &&
                  res.push(
                    "Invalid characters (only letters, numbers, and punctuation allowed)"
                  );
                if (value && value.length > 2048) {
                  res.push("Maximum length is 2048 characters");
                }
                return res.length > 0 ? res.join(", ") : undefined;
              },
            })}
            onBlur={() => void methods.trigger("description")}
          />
          {methods.formState.errors?.description && (
            <FormErrorMessage>
              {methods.formState.errors.description?.message?.toString()}
            </FormErrorMessage>
          )}
        </FormControl>
      </>
    </FormStep>
  );
};
