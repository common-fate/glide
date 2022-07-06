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
} from "../../../utils/backend-client/admin/admin";
import { adminArchiveAccessRule } from "../../../utils/backend-client/default/default";
import {
  AccessRuleDetail,
  UpdateAccessRuleRequestBody,
} from "../../../utils/backend-client/types";
import { ProviderPreviewOnlyStep } from "./components/ProviderPreview";

import { ApprovalStep } from "./steps/Approval";
import { GeneralStep } from "./steps/General";
import { RequestsStep } from "./steps/Request";
import { TimeStep } from "./steps/Time";
import { StepsProvider } from "./StepsContext";

interface FormData extends UpdateAccessRuleRequestBody {
  approval: { required: boolean; users: string[]; groups: string[] };
}

interface Props {
  data: AccessRuleDetail;
  readOnly?: boolean;
}

const EditAccessRuleForm = ({ data, readOnly }: Props) => {
  const {
    params: { id: ruleId },
  } = useMatch();

  const { isOpen, onClose, onOpen } = useDisclosure();
  const navigate = useNavigate();
  const toast = useToast();
  // const ruleId = typeof query?.id == "string" ? query.id : "";
  // we use this to ensure that data for selected and then deselected providers is not included.
  const methods = useForm<FormData>({
    shouldUnregister: true,
  });

  const [isArchiving, setIsArchiving] = useState<boolean>(false);

  const [cachedRule, setCachedRule] = useState<AccessRuleDetail | undefined>();
  const { mutate } = useAdminGetAccessRule(ruleId);

  useEffect(() => {
    // We will only reset form data if it has changed on the backend
    if (data && (!cachedRule || cachedRule != data)) {
      const f: FormData = {
        description: data.description,
        groups: data.groups,
        name: data.name,
        timeConstraints: {
          maxDurationSeconds: data.timeConstraints.maxDurationSeconds,
        },
        approval: {
          required:
            data.approval.users.length > 0 || data.approval.groups?.length > 0,
          users: data.approval.users,
          groups: data.approval.groups,
        },
      };
      methods.reset(f);
      setCachedRule(data);
    }
    return () => {
      setCachedRule(undefined);
    };
  }, [data, methods]);

  const onSubmit = async (data: FormData) => {
    console.debug("submit form data for edit", { data });

    const { approval, ...d } = data;
    const ruleData: UpdateAccessRuleRequestBody = {
      approval: { users: [], groups: [] },
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
      await adminUpdateAccessRule(ruleId, ruleData);
      toast({
        title: "Access rule updated",
        status: "success",
        variant: "subtle",
        duration: 2200,
        isClosable: true,
      });
      mutate();
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

  const handleArchive = async () => {
    try {
      setIsArchiving(true);
      await adminArchiveAccessRule(ruleId);

      toast({
        title: "Access rule archived",
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
        title: "Error archiving access rule",
        description,
        status: "error",
        variant: "subtle",
        duration: 2200,
        isClosable: true,
      });
    } finally {
      setIsArchiving(false);
    }
  };

  return (
    <Container pt={6} maxW="container.md">
      {!readOnly && (
        <Flex justifyContent="flex-end" w="100%" flexGrow={1} mb={4}>
          <Button
            size="sm"
            variant="ghost"
            leftIcon={<DeleteIcon />}
            onClick={onOpen}
          >
            Archive Access Rule
          </Button>
        </Flex>
      )}
      <FormProvider {...methods}>
        <form onSubmit={methods.handleSubmit(onSubmit)}>
          <VStack w="100%" spacing={6}>
            <ProviderPreviewOnlyStep target={data.target} />
            <StepsProvider isEditMode={!readOnly} isReadOnly={readOnly}>
              <GeneralStep />
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
          <ModalHeader>Archive Access Rule</ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            Are you sure you want to archive this access rule?
          </ModalBody>

          <ModalFooter>
            <Button
              variant={"solid"}
              colorScheme="red"
              rounded="full"
              mr={3}
              onClick={handleArchive}
              isLoading={isArchiving}
            >
              Archive Rule
            </Button>
            <Button
              variant={"brandSecondary"}
              onClick={onClose}
              isDisabled={isArchiving}
            >
              Cancel
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </Container>
  );
};

const BottomActionButtons: React.FC<{ rule: AccessRuleDetail }> = ({
  rule,
}) => {
  const { formState } = useFormContext();
  const navigate = useNavigate();
  const toast = useToast();
  const [isArchiving, setIsArchiving] = useState<boolean>(false);
  const handleArchive = async () => {
    try {
      setIsArchiving(true);
      const res = await adminArchiveAccessRule(rule.id, {});

      toast({
        title: "Access rule archived",
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
        title: "Error archiving access rule",
        description,
        status: "error",
        variant: "subtle",
        duration: 2200,
        isClosable: true,
      });
    } finally {
      setIsArchiving(false);
    }
  };

  // No available actions for archived rules
  if (rule.status === "ARCHIVED") {
    return <ButtonGroup w="100%"></ButtonGroup>;
  }

  return (
    <ButtonGroup w="100%">
      <Button
        isLoading={formState.isSubmitting}
        disabled={isArchiving}
        type="submit"
      >
        Update
      </Button>
      <Button
        disabled={formState.isSubmitting}
        variant={"solid"}
        colorScheme="red"
        isLoading={isArchiving}
        onClick={handleArchive}
        type="button"
      >
        Archive
      </Button>
      <Button
        variant="brandSecondary"
        disabled={isArchiving || formState.isSubmitting}
        type="button"
        onClick={() => navigate({ to: "/admin/access-rules" })}
      >
        Cancel
      </Button>
    </ButtonGroup>
  );
};
export default EditAccessRuleForm;
