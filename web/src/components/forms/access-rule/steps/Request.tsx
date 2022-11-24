import {
  Flex,
  FormControl,
  FormErrorMessage,
  FormLabel,
  Text,
  VStack,
  Wrap,
  WrapItem,
} from "@chakra-ui/react";
import React from "react";
import { useFormContext } from "react-hook-form";
import { GroupSelect } from "../components/Select";
import { AccessRuleFormData } from "../CreateForm";
import { GroupDisplay } from "./Approval";
import { FormStep } from "./FormStep";

export const RequestsStep: React.FC = () => {
  const methods = useFormContext<AccessRuleFormData>();
  const groups = methods.watch("groups");

  return (
    <FormStep
      heading="Request"
      subHeading="Who can request access to the permissions?"
      fields={["groups"]}
      preview={
        <VStack width={"100%"} align="flex-start">
          <Flex>
            <Text mr={2}>Groups:</Text>
            <Wrap>
              {groups?.map((g) => (
                <WrapItem key={"rwrap" + g}>
                  <GroupDisplay groupId={g} />
                </WrapItem>
              ))}
            </Wrap>
          </Flex>
        </VStack>
      }
    >
      <FormControl isInvalid={!!methods.formState.errors.groups}>
        <FormLabel htmlFor="groups">
          <Text textStyle={"Body/Medium"}>Add or remove groups</Text>
        </FormLabel>

        <GroupSelect
          testId="group-select"
          fieldName="groups"
          rules={{ required: true, minLength: 1 }}
          onBlurSecondaryAction={() => {
            void methods.trigger("groups");
          }}
        />
        <FormErrorMessage>At least one group is required.</FormErrorMessage>
      </FormControl>
    </FormStep>
  );
};
