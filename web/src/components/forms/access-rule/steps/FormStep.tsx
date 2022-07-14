import {
  Box,
  Button,
  ButtonGroup,
  Collapse,
  Flex,
  Spacer,
  Text,
  VStack,
} from "@chakra-ui/react";
import React from "react";
import { useFormContext } from "react-hook-form";
import { useFormStep } from "../FormStepContext";

interface FormStepProps {
  heading: string;
  subHeading: string;
  // the fields that this form component contains, used for validating the step before enabling next
  fields: string[];
  hideNext?: boolean;
  preview?: React.ReactElement;
  children?: React.ReactElement;
}

export const FormStep: React.FC<FormStepProps> = ({
  heading,
  children,
  subHeading,
  fields,
  preview,
}) => {
  const { trigger, getFieldState } = useFormContext();
  const { active } = useFormStep();
  let hasFieldErrors = false;
  fields.forEach((f) => {
    const s = getFieldState(f);
    if (s.error) {
      hasFieldErrors = true;
    }
  });
  const runValidation = () => trigger(fields);
  return (
    <VStack px={8} py={8} bg="neutrals.100" rounded="md" w="100%">
      <Flex w="100%">
        <Text textStyle="Heading/H3" opacity={active ? 1 : 0.6}>
          {heading}
        </Text>
        <Spacer />
        <TopActionButtons validate={runValidation} hasErrors={hasFieldErrors} />
      </Flex>
      {preview && <Preview>{preview}</Preview>}
      <SubHeading subHeading={subHeading} />
      <Box
        w="100%"
        sx={{
          ".chakra-collapse": {
            overflow: active ? "visible !important" : "visible !important",
          },
        }}
      >
        <Collapse in={active}>
          <VStack spacing={10} align={"flex-start"} w="100%">
            <VStack spacing={6} align={"flex-start"} w="100%" p={1}>
              {children}
            </VStack>
            <BottomActionButtons
              validate={runValidation}
              hasErrors={hasFieldErrors}
            />
          </VStack>
        </Collapse>
      </Box>
    </VStack>
  );
};
const SubHeading: React.FC<{ subHeading: string }> = ({ subHeading }) => {
  const { active } = useFormStep();
  if (!active) {
    return null;
  }
  return (
    <Text w="100%" textStyle={"Body/Medium"} color="neutrals.600">
      {subHeading}
    </Text>
  );
};
const BottomActionButtons: React.FC<{
  validate: () => Promise<boolean>;
  hasErrors: boolean;
}> = ({ validate, hasErrors }) => {
  const { showNext, showSubmit, next } = useFormStep();
  const { formState } = useFormContext();
  return (
    <Flex w="100%" justify={"left"}>
      {showNext && (
        <Button
          id="form-step-next-button"
          isDisabled={hasErrors}
          onClick={async () => {
            // only go to next if there are no field errors
            (await validate()) && next();
          }}
        >
          Next
        </Button>
      )}
      {showSubmit && (
        <Button isLoading={formState.isSubmitting} type="submit">
          Create
        </Button>
      )}
    </Flex>
  );
};

const Preview: React.FC<{
  children: React.ReactElement;
}> = ({ children }) => {
  const { showPreview } = useFormStep();
  return showPreview && children ? children : null;
};
const TopActionButtons: React.FC<{
  validate: () => Promise<boolean>;
  hasErrors: boolean;
}> = ({ hasErrors, validate }) => {
  const { showEdit, showClose, edit, close } = useFormStep();
  return (
    <ButtonGroup>
      {showEdit && (
        <Button variant="brandSecondary" size="sm" onClick={edit}>
          Edit
        </Button>
      )}
      {showClose && (
        <Button
          variant="brandSecondary"
          isDisabled={hasErrors}
          size="sm"
          onClick={async () => {
            (await validate()) && close();
          }}
        >
          Close
        </Button>
      )}
    </ButtonGroup>
  );
};
