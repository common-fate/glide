import { Container, useToast, VStack } from "@chakra-ui/react";
import axios from "axios";
import { FormProvider, useForm } from "react-hook-form";
import { useNavigate } from "react-location";
import { adminCreateAccessRule } from "../../../utils/backend-client/admin/admin";

import {
  CreateAccessRuleRequestBody,
  CreateAccessRuleTarget,
  Operation,
  ResourceFilter,
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
  // contains informatation related to how many targets are selected.
  // this is irrelevant information for CreateAccessRuleRequestBody.
  targetFieldMap: Map<string, ResourceFilter>;
}

export const accessRuleFormDataToApi = (
  formData: AccessRuleFormData
): CreateAccessRuleRequestBody => {
  const { approval, targetgroups, targetFieldMap, ...d } = formData;

  const ruleData: CreateAccessRuleRequestBody = {
    approval: { users: [], groups: [] },
    ...d,
  };

  const targets: Map<string, CreateAccessRuleTarget> = new Map();

  // TODO: Check for condition where there is no targetgroups.
  if (targetgroups) {
    Object.entries(targetgroups).map(([k, v]) => {
      const target: CreateAccessRuleTarget = {
        targetGroupId: k,
        fieldFilterExpessions: {},
      };

      const targetFields = targetgroups[k];

      Object.entries(targetFields).map(([k, v]) => {
        // if no filter selected then we will return empty array which signifies to select everything
        if (
          v.operationType === ResourceFilterOperationTypeEnum.IN &&
          v.values.length === 0
        ) {
          target.fieldFilterExpessions[k] = [];
        } else {
          const filterOperation: Operation = {
            operationType: v.operationType as ResourceFilterOperationTypeEnum,
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
          target.fieldFilterExpessions[k] = [filterOperation];
        }
      });

      targets.set(k, target);
    });

    ruleData.targets = Array.from(targets.values());
  }

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
