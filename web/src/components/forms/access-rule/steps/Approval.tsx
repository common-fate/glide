import {
  Flex,
  FormControl,
  FormErrorMessage,
  FormHelperText,
  FormLabel,
  HStack,
  Switch,
  Text,
  VStack,
  Wrap,
  WrapItem,
} from "@chakra-ui/react";
import React from "react";
import { useFormContext } from "react-hook-form";
import { useGetGroup } from "../../../../utils/backend-client/admin/admin";
import { UserAvatarDetails } from "../../../UserAvatar";
import { GroupSelect, UserSelect } from "../components/Select";
import { CreateAccessRuleFormData } from "../CreateForm";

import { FormStep } from "./FormStep";

export const ApprovalStep: React.FC = () => {
  const methods = useFormContext<CreateAccessRuleFormData>();
  const approval = methods.watch("approval");
  // If approval is required, then at least one user or one group needs to be set
  const approverRequired =
    !!approval?.required &&
    !(approval?.groups?.length > 0 || approval?.users?.length > 0);
  return (
    <FormStep
      heading="Approvers"
      subHeading="Who can approve access to the principal?"
      fields={[]}
      hideNext={true}
      preview={<ApprovalPreview />}
    >
      <>
        <FormControl>
          <FormLabel htmlFor="approval.required">
            <HStack>
              <Switch
                id="requires-approval-button"
                bg="neutrals.0"
                {...methods.register("approval.required", {})}
              />
              <Text textStyle={"Body/Medium"}>Approval required</Text>
            </HStack>
          </FormLabel>
        </FormControl>

        <FormControl
          isInvalid={!!methods.formState.errors.approval?.groups}
          isDisabled={!approval?.required}
        >
          <FormLabel htmlFor="approval.groups">
            <Text textStyle={"Body/Medium"}>Add or remove approval groups</Text>
          </FormLabel>

          <GroupSelect
            data-testid="approval-group-select"
            fieldName="approval.groups"
            isDisabled={!approval?.required}
            rules={{ required: approverRequired }}
          />

          <FormErrorMessage>
            At least one approver or one group is required.
          </FormErrorMessage>
        </FormControl>
        <FormControl
          isInvalid={!!methods.formState.errors.approval?.users}
          isDisabled={!approval?.required}
        >
          <FormLabel htmlFor="approval.users">
            <Text textStyle={"Body/Medium"}>
              Add or remove individual approvers
            </Text>
          </FormLabel>
          <UserSelect
            fieldName="approval.users"
            isDisabled={!approval?.required}
            rules={{ required: approverRequired }}
          />

          <FormErrorMessage>
            At least one approver or one group is required.
          </FormErrorMessage>
        </FormControl>
      </>
    </FormStep>
  );
};

export const GroupDisplay: React.FC<{ groupId: string }> = ({ groupId }) => {
  const { data } = useGetGroup(groupId);
  return <Text>{data?.name}</Text>;
};
const ApprovalPreview: React.FC = () => {
  const methods = useFormContext();
  const approval = methods.watch("approval");
  if (!approval?.required) {
    return <Text w="100%">No approval required</Text>;
  }
  return (
    <VStack w="100%" align={"flex-start"}>
      <Flex>
        <Text mr={2}>Users:</Text>
        <Wrap>
          {approval?.users?.map((u: string) => {
            return (
              <WrapItem key={"wrap" + u}>
                <UserAvatarDetails user={u} size="xs" />
              </WrapItem>
            );
          })}
        </Wrap>
      </Flex>
      <Flex>
        <Text mr={2}>Groups:</Text>
        <Wrap>
          {approval?.groups?.map((g: string) => {
            return (
              <WrapItem key={"wrap" + g}>
                <GroupDisplay groupId={g} />
              </WrapItem>
            );
          })}
        </Wrap>
      </Flex>
    </VStack>
  );
};
