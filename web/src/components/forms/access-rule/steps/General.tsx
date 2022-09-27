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
            {...methods.register("name", { required: true })}
            onBlur={() => void methods.trigger("name")}
          />
          <FormErrorMessage>Name is required.</FormErrorMessage>
        </FormControl>
        <FormControl isInvalid={!!methods.formState.errors.description}>
          <FormLabel htmlFor="Description">
            <Text textStyle={"Body/Medium"}>Description</Text>
          </FormLabel>
          <Textarea
            bg="neutrals.0"
            {...methods.register("description", {
              required: true,
            })}
            onBlur={() => void methods.trigger("description")}
          />
          <FormErrorMessage>Description is required.</FormErrorMessage>
        </FormControl>
      </>
    </FormStep>
  );
};
