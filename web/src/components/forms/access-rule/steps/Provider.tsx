import {
  FormControl,
  FormErrorMessage,
  FormLabel,
  IconButton,
  IconButtonProps,
  Spinner,
  Text,
  Tooltip,
} from "@chakra-ui/react";
import React, { useEffect, useState } from "react";
import { Controller, useFormContext } from "react-hook-form";

import { RefreshIcon } from "../../../icons/Icons";
import ProviderArgumentField from "../components/ProviderArgumentField";
import { ProviderPreview } from "../components/ProviderPreview";
import { ProviderRadioSelector } from "../components/ProviderRadio";
import { AccessRuleFormData, AccessRuleFormDataTarget } from "../CreateForm";
import { FormStep } from "./FormStep";
import { Provider, TargetGroup } from "../../../../utils/backend-client/types";
import {
  adminListTargetGroupArgOptions,
  useAdminGetTargetGroup,
  useAdminGetTargetGroupArgs,
  useAdminListTargetGroupArgOptions,
} from "src/utils/backend-client/admin/admin";

interface PreviewProps {
  target: AccessRuleFormDataTarget;
  provider?: TargetGroup;
}

const Preview = (props: PreviewProps) => {
  if (!props.provider) return null;

  return <ProviderPreview provider={props.provider} />;
};

export const ProviderStep: React.FC = () => {
  const methods = useFormContext<AccessRuleFormData>();
  const target = methods.watch("target");

  const { data: provider, isValidating: ivp } = useAdminGetTargetGroup(
    target?.providerId
  );
  const { data: providerArgs, isValidating: ivpa } = useAdminGetTargetGroupArgs(
    target?.providerId ?? ""
  );

  const isFieldLoading = (!provider && ivp) || (!providerArgs && ivpa);

  return (
    <FormStep
      heading="Provider"
      subHeading="The permissions that the rule gives access to"
      fields={["target", "target.providerId"]}
      preview={<Preview target={target} provider={provider} />}
      isFieldLoading={isFieldLoading}
    >
      <>
        <FormControl isInvalid={!!methods.formState.errors.target?.providerId}>
          <FormLabel htmlFor="target.providerId">
            <Text textStyle={"Body/Medium"}>Provider</Text>
          </FormLabel>
          {/* <ProviderSetupNotice /> */}
          <Controller
            control={methods.control}
            rules={{ required: true }}
            name={"target.providerId"}
            render={({ field: { ref, onChange, ...rest } }) => (
              <ProviderRadioSelector
                onChange={async (t) => {
                  onChange(t);
                  await methods.trigger("target.providerId");
                }}
                {...rest}
              />
            )}
          />

          <FormErrorMessage>Provider is required</FormErrorMessage>
        </FormControl>
        {providerArgs &&
          target?.providerId &&
          Object.values(providerArgs).map((v) => (
            <ProviderArgumentField
              argument={v}
              providerId={target?.providerId}
            />
          ))}
      </>
    </FormStep>
  );
};

type RefreshButtonProps = {
  providerId: string;
  argId: string;
} & Omit<IconButtonProps, "aria-label">;

export const RefreshButton: React.FC<RefreshButtonProps> = ({
  argId,
  providerId,
  ...props
}) => {
  const [loading, setLoading] = useState(false);
  const { data, mutate, isValidating } = useAdminListTargetGroupArgOptions(
    providerId,
    argId
  );
  const onClick = async () => {
    setLoading(true);
    await mutate(
      adminListTargetGroupArgOptions(providerId, argId, {
        refresh: true,
      })
    );
    setLoading(false);
  };

  return (
    <Tooltip>
      <IconButton
        {...props}
        onClick={onClick}
        isLoading={(!data && isValidating) || loading}
        icon={<RefreshIcon boxSize="24px" />}
        aria-label="Refresh"
        variant={"ghost"}
      />
    </Tooltip>
  );
};
