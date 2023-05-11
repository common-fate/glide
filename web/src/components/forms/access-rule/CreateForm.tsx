import { Container, useToast, VStack } from "@chakra-ui/react";
import axios from "axios";
import { FormProvider, useForm } from "react-hook-form";
import { useNavigate } from "react-location";
import { adminCreateAccessRule } from "../../../utils/backend-client/admin/admin";

import {
  CreateAccessRuleRequestBody,
  CreateAccessRuleTarget,
  Operation,
  ResourceFilterOperationTypeEnum,
} from "../../../utils/backend-client/types";
import { ApprovalStep } from "./steps/Approval";
import { GeneralStep } from "./steps/General";
import { TargetStep } from "./steps/Provider";
import { RequestsStep } from "./steps/Request";
import { TimeStep } from "./steps/Time";
import { StepsProvider } from "./StepsContext";
import { FieldStep } from "./steps/Field";

interface TargetGroupFilterOperation {
  attribute: string;
  values: string[];
  value: string;
  operationType: string;
}

type TargetGroups = {
  [key: string]: Record<string, TargetGroupFilterOperation>;
};

export interface AccessRuleFormData extends CreateAccessRuleRequestBody {
  approval: { required: boolean; users: string[]; groups: string[] };
  targetgroups: TargetGroups;
}

export const accessRuleFormDataToApi = (
  formData: AccessRuleFormData
): CreateAccessRuleRequestBody => {
  const { approval, targetgroups, ...d } = formData;

  const ruleData: CreateAccessRuleRequestBody = {
    approval: { users: [], groups: [] },
    ...d,
  };

  const targets: Map<string, CreateAccessRuleTarget> = new Map();

  if (targetgroups) {
    Object.entries(targetgroups).map(([k, v]) => {
      let selectedTarget = d.targets.find((t) => t.targetGroupId === k);

      if (selectedTarget) {
        Object.entries(v).map(([k, v]) => {
          if (selectedTarget) {
            if (v.value == "" && v.values.length === 0) {
              selectedTarget.fieldFilterExpessions[k] = [];
            } else {
              const filterOperation: Operation = {
                operationType:
                  v.operationType as ResourceFilterOperationTypeEnum,
                attribute: v.attribute,
                ...(v.value != ""
                  ? {
                      value: v.value,
                    }
                  : null),
                ...(v.values.length !== 0
                  ? {
                      values: v.values,
                    }
                  : null),
              };

              // NOTE: once multi-operations feature is supported. we need to update this
              // to add operations value as well.
              selectedTarget.fieldFilterExpessions[k] = [filterOperation];
            }

            targets.set(selectedTarget.targetGroupId, selectedTarget);
          }
        });
      }
    });
  }

  // mutate the targets field with updated structure.
  ruleData.targets = Array.from(targets.values());

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
  const methods = useForm<AccessRuleFormData>({
    shouldUnregister: true,
    defaultValues: {
      targets: [],
      priority: 1,
    },
  });
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
              <TargetStep />
              {/* <FieldStep /> */}
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
