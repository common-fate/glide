import { Container, useToast, VStack } from "@chakra-ui/react";
import axios from "axios";
import { FormProvider, useForm } from "react-hook-form";
import { useNavigate } from "react-location";
import { adminCreateAccessRule } from "../../../utils/backend-client/admin/admin";

import {
  CreateAccessRuleRequestBody,
  CreateAccessRuleTarget,
  CreateAccessRuleTargetWith,
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

const CreateAccessRuleForm = () => {
  const navigate = useNavigate();

  const toast = useToast();
  //  Should unregister controls how the form will persist data if a component is unmounted
  // we use this to ensure that data for selected and then deselected providers is not included.
  const methods = useForm<AccessRuleFormData>({ shouldUnregister: true });
  const onSubmit = async (data: AccessRuleFormData) => {
    console.debug("submit form data", { data });

    const { approval, timeConstraints, target, ...d } = data;
    const t: {
      providerId: string;
      with: CreateAccessRuleTargetWith;
    } = {
      providerId: target.providerId,
      with: {},
    };

    // // For fields with text i.e input type add the values to
    // // with.values for the API.
    // for (const k in target.withText) {
    //   t.with[k] = {
    //     values: [target.withText[k]],
    //     groupings: {},
    //   };
    // }

    // // First add everything in `target.with` to values.
    // for (const arg in target.with) {
    //   t.with[arg] = {
    //     ...t.with[arg],
    //     values: target.with[arg] as any,
    //     groupings: {},
    //   };
    // }

    // // TODO: Grouping can be made an optional value.
    // for (const arg in target.withFilter) {
    //   // Loop over any withFilter key for that arg
    //   for (const key of Object.keys(target.withFilter[arg])) {
    //     t.with[arg] = {
    //       ...t.with[arg],
    //       groupings: {
    //         ...t.with[arg].groupings,
    //         [key]: target.withFilter[arg][key],
    //       },
    //     };
    //   }
    // }

    // for (const k in target.withText) {
    //   t.with[k] = {
    //     ...t.with[k],
    //     values: [target.withText[k]],
    //   };
    // }

    const ruleData: CreateAccessRuleRequestBody = {
      approval: { users: [], groups: [] },
      timeConstraints: {
        maxDurationSeconds: timeConstraints.maxDurationSeconds,
      },
      target: t,
      ...d,
    };
    // only apply these fields if approval is enabled
    if (approval.required) {
      ruleData["approval"].users = data.approval.users;
      ruleData["approval"].groups = data.approval.groups;
    } else {
      ruleData["approval"].users = [];
    }
    try {
      await adminCreateAccessRule(ruleData);
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
