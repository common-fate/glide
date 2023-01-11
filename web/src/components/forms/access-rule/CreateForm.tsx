import { Container, useToast, VStack } from "@chakra-ui/react";
import axios from "axios";
import { FormProvider, useForm } from "react-hook-form";
import { useNavigate } from "react-location";
import { adminCreateAccessRule } from "../../../utils/backend-client/admin/admin";

import {
  CreateAccessRuleRequestBody,
  AccessRuleTarget,
  Provider,
  CreateAccessRuleTarget,
} from "../../../utils/backend-client/types";
import { ApprovalStep } from "./steps/Approval";
import { GeneralStep } from "./steps/General";
import { ProviderStep } from "./steps/Provider";
import { RequestsStep } from "./steps/Request";
import { TimeStep } from "./steps/Time";
import { StepsProvider } from "./StepsContext";

export type AccessRuleFormDataTarget = {
  providerId: string;
  multiSelects: { [key: string]: string[] };
  argumentGroups: { [key: string]: { [key: string]: string[] } };
  inputs: { [key: string]: string };
};
export interface AccessRuleFormData
  extends Omit<CreateAccessRuleRequestBody, "target"> {
  approval: { required: boolean; users: string[]; groups: string[] };
  // with text is used for single text fields
  target: AccessRuleFormDataTarget;
}

export const accessRuleFormDataTargetToApi = (
  target: AccessRuleFormDataTarget
): CreateAccessRuleTarget => {
  const t: CreateAccessRuleTarget = {
    providerId: target.providerId,
    with: {},
  };
  for (const k in target.inputs) {
    t.with[k] = {
      groupings: {},
      values: [target.inputs[k]],
    };
  }
  for (const k in target.multiSelects) {
    t.with[k] = {
      groupings: target.argumentGroups?.[k] || {},
      values: target.multiSelects[k],
    };
  }
  return t;
};

export const accessRuleFormDataToApi = (
  formData: AccessRuleFormData
): CreateAccessRuleRequestBody => {
  const { approval, target, ...d } = formData;

  const ruleData: CreateAccessRuleRequestBody = {
    approval: { users: [], groups: [] },
    target: accessRuleFormDataTargetToApi(target),
    ...d,
  };
  // only apply these fields if approval is enabled
  if (approval.required) {
    ruleData["approval"].users = approval.users;
    ruleData["approval"].groups = approval.groups;
  } else {
    ruleData["approval"].users = [];
  }
  return ruleData;
};

const CreateAccessRuleForm = () => {
  const navigate = useNavigate();

  const toast = useToast();
  //  Should unregister controls how the form will persist data if a component is unmounted
  // we use this to ensure that data for selected and then deselected providers is not included.
  const methods = useForm<AccessRuleFormData>({ shouldUnregister: true });
  const onSubmit = async (data: AccessRuleFormData) => {
    console.debug("submit form data", { data });

    try {
      await adminCreateAccessRule(accessRuleFormDataToApi(data));
      toast({
        id: "access-rule-created",
        title: "Access rule created",
        status: "success",
        variant: "subtle",
        duration: 2200,
        isClosable: true,
      });
      navigate({ to: "/admin/access-rules" });
    } catch (err) {
      let description: string | undefined;
      if (axios.isAxiosError(err)) {
        // @ts-ignore
        description = err?.response?.data.error;
      }
      toast({
        title: "Error creating access rule",
        description,
        status: "error",
        variant: "subtle",
        duration: 2200,
        isClosable: true,
      });
    }
  };
  return (
    <Container pt={12} maxW="container.md">
      <FormProvider {...methods}>
        <form onSubmit={methods.handleSubmit(onSubmit)}>
          <VStack w="100%" spacing={6}>
            <StepsProvider>
              <GeneralStep />
              <ProviderStep />
              <TimeStep />
              <RequestsStep />
              <ApprovalStep />
            </StepsProvider>
          </VStack>
        </form>
      </FormProvider>
    </Container>
  );
};

export default CreateAccessRuleForm;
