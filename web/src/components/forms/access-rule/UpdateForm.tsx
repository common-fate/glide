import { DeleteIcon } from "@chakra-ui/icons";
import {
  Button,
  ButtonGroup,
  Container,
  Flex,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  useDisclosure,
  useToast,
  VStack,
} from "@chakra-ui/react";
import axios from "axios";
import { useEffect, useState } from "react";
import { FormProvider, useForm, useFormContext } from "react-hook-form";
import { useMatch, useNavigate } from "react-location";
import {
  adminUpdateAccessRule,
  useAdminGetAccessRule,
  useAdminGetTargetGroup,
} from "../../../utils/backend-client/admin/admin";
import { adminDeleteAccessRule } from "../../../utils/backend-client/default/default";
import { AccessRule, TargetGroup } from "../../../utils/backend-client/types";
import { AccessRuleFormData, accessRuleFormDataToApi } from "./CreateForm";

import { ApprovalStep } from "./steps/Approval";
import { GeneralStep } from "./steps/General";
import { TargetStep } from "./steps/Provider";
import { RequestsStep } from "./steps/Request";
import { TimeStep } from "./steps/Time";
import { StepsProvider } from "./StepsContext";
interface Props {
  data: AccessRule;
}

const UpdateAccessRuleForm = ({ data }: Props) => {
  const {
    params: { id: ruleId },
  } = useMatch();

  const { isOpen, onClose, onOpen } = useDisclosure();
  const navigate = useNavigate();
  const toast = useToast();
  // const ruleId = typeof query?.id == "string" ? query.id : "";
  // we use this to ensure that data for selected and then deselected providers is not included.
  const methods = useForm<AccessRuleFormData>({
    shouldUnregister: true,
    defaultValues: {
      targets: [],
    },
  });

  const [isDeleting, setIsDeleting] = useState<boolean>(false);

  const { mutate } = useAdminGetAccessRule(ruleId);

  useEffect(() => {
    // We will only reset form data if it has changed on the backend
    if (data !== undefined) {
      //set accessRuleTargetData from rule details from api

      let approvalRequired = false;
      if (data.approval.users) {
        if (data.approval.users.length > 0) {
          approvalRequired = true;
        }
      }

      if (data.approval.groups) {
        if (data.approval.groups.length > 0) {
          approvalRequired = true;
        }
      }

      const f: AccessRuleFormData = {
        description: data.description,
        groups: data.groups,
        name: data.name,
        timeConstraints: {
          maxDurationSeconds: data.timeConstraints.maxDurationSeconds,
          defaultDurationSeconds: data.timeConstraints.defaultDurationSeconds,
        },
        approval: {
          required: approvalRequired,
          users: data.approval.users ? data.approval.users : [],
          groups: data.approval.groups ? data.approval.groups : [],
        },
        targets: [],
        priority: data.priority,
      };
      methods.reset(f);
    }
  }, [data, methods]);

  const onSubmit = async (data: AccessRuleFormData) => {
    console.debug("submit form data for edit", { data });
    try {
      await adminUpdateAccessRule(ruleId, accessRuleFormDataToApi(data));
      toast({
        id: "access-rule-updated",
        title: "Access rule updated",
        status: "success",
        variant: "subtle",
        duration: 3000,
        isClosable: true,
      });
      void mutate();
      navigate({ to: "/admin/access-rules" });
    } catch (err) {
      let description: string | undefined;
      if (axios.isAxiosError(err)) {
        // @ts-ignore
        description = err?.response?.data.error;
      }
      toast({
        title: "Error updating access rule",
        description,
        status: "error",
        variant: "subtle",
        duration: 2200,
        isClosable: true,
      });
    }
  };

  const handleDelete = async () => {
    try {
      setIsDeleting(true);
      await adminDeleteAccessRule(ruleId);

      toast({
        title: "Access rule deleted",
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
        title: "Error deleting access rule",
        description,
        status: "error",
        variant: "subtle",
        duration: 2200,
        isClosable: true,
      });
    } finally {
      setIsDeleting(false);
    }
  };

  return (
    <Container pt={6} maxW="container.md">
      <Flex justifyContent="flex-end" w="100%" flexGrow={1} mb={4}>
        <Button
          size="sm"
          variant="ghost"
          leftIcon={<DeleteIcon />}
          onClick={onOpen}
        >
          Delete Access Rule
        </Button>
      </Flex>

      <FormProvider {...methods}>
        <form onSubmit={methods.handleSubmit(onSubmit)}>
          <VStack w="100%" spacing={6}>
            <StepsProvider isEditMode={true}>
              <GeneralStep />
              <TargetStep />
              <TimeStep />
              <RequestsStep />
              <ApprovalStep />
            </StepsProvider>
            <BottomActionButtons rule={data} />
          </VStack>
        </form>
      </FormProvider>
      <Modal isOpen={isOpen} onClose={onClose}>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>Delete Access Rule</ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            Are you sure you want to delete this access rule?
          </ModalBody>

          <ModalFooter>
            <Button
              variant={"solid"}
              colorScheme="red"
              rounded="full"
              mr={3}
              onClick={handleDelete}
              isLoading={isDeleting}
            >
              Delete Rule
            </Button>
            <Button
              variant={"brandSecondary"}
              onClick={onClose}
              isDisabled={isDeleting}
            >
              Cancel
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </Container>
  );
};

const BottomActionButtons: React.FC<{ rule: AccessRule }> = ({ rule }) => {
  const { formState } = useFormContext();
  const navigate = useNavigate();

  return (
    <ButtonGroup w="100%">
      <Button isLoading={formState.isSubmitting} type="submit">
        Update
      </Button>
      <Button
        variant="brandSecondary"
        isDisabled={formState.isSubmitting}
        type="button"
        onClick={() => navigate({ to: "/admin/access-rules" })}
      >
        Cancel
      </Button>
    </ButtonGroup>
  );
};
export default UpdateAccessRuleForm;
