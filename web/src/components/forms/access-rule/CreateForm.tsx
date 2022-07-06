import { Container, useToast, VStack } from "@chakra-ui/react";
import axios from "axios";
import { FormProvider, useForm } from "react-hook-form";
import { useNavigate } from "react-location";
import { adminCreateAccessRule } from "../../../utils/backend-client/admin/admin";
import {
  AccessRuleTarget,
  CreateAccessRuleRequestBody,
} from "../../../utils/backend-client/types";
import { ApprovalStep } from "./steps/Approval";
import { GeneralStep } from "./steps/General";
import { ProviderStep } from "./steps/Provider";
import { RequestsStep } from "./steps/Request";
import { TimeStep } from "./steps/Time";
import { StepsProvider } from "./StepsContext";

export interface CreateAccessRuleFormData extends CreateAccessRuleRequestBody {
  approval: { required: boolean; users: string[]; groups: string[] };
}

const CreateAccessRuleForm = () => {
  const navigate = useNavigate();

  const toast = useToast();
  //  Should unregister controls how the form will persist data if a component is unmounted
  // we use this to ensure that data for selected and then deselected providers is not included.
  const methods = useForm<CreateAccessRuleFormData>({ shouldUnregister: true });
  const onSubmit = async (data: CreateAccessRuleFormData) => {
    console.debug("submit form data", { data });

    const { approval, timeConstraints, ...d } = data;
    const ruleData: CreateAccessRuleRequestBody = {
      approval: { users: [], groups: [] },
      timeConstraints: {
        maxDurationSeconds: timeConstraints.maxDurationSeconds,
      },
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
