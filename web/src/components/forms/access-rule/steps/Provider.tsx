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
import { ArgumentRuleFormElement } from "../../../../utils/backend-client/types/accesshandler-openapi.yml";
import {
  adminListProviderArgOptions,
  useAdminGetProvider,
  useAdminGetProviderArgs,
  useAdminListProviderArgOptions,
} from "../../../../utils/backend-client/admin/admin";

import { RefreshIcon } from "../../../icons/Icons";
import ProviderSetupNotice from "../../../ProviderSetupNotice";
import ProviderArgumentField from "../components/ProviderArgumentField";
import { ProviderPreview } from "../components/ProviderPreview";
import { ProviderRadioSelector } from "../components/ProviderRadio";
import { AccessRuleFormData } from "../CreateForm";
import { FormStep } from "./FormStep";

export const ProviderStep: React.FC = () => {
  const methods = useFormContext<AccessRuleFormData>();
  const target = methods.watch("target");

  const { data: provider, isValidating: ivp } = useAdminGetProvider(
    target?.providerId
  );
  const { data: providerArgs, isValidating: ivpa } = useAdminGetProviderArgs(
    target?.providerId ?? ""
  );

  const Preview = () => {
    if (!target || !provider || !(target?.inputs || target?.multiSelects)) {
      return null;
    }
    return <ProviderPreview provider={provider} />;
  };
  const isFieldLoading = (!provider && ivp) || (!providerArgs && ivpa);

  return (
    <FormStep
      heading="Provider"
      subHeading="The permissions that the rule gives access to"
      fields={["target", "target.providerId"]}
      preview={<Preview />}
      isFieldLoading={isFieldLoading}
    >
      <>
        <FormControl isInvalid={!!methods.formState.errors.target?.providerId}>
          <FormLabel htmlFor="target.providerId">
            <Text textStyle={"Body/Medium"}>Provider</Text>
          </FormLabel>
          <ProviderSetupNotice />
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
  const { data, mutate, isValidating } = useAdminListProviderArgOptions(
    providerId,
    argId
  );
  const onClick = async () => {
    setLoading(true);
    await mutate(
      adminListProviderArgOptions(providerId, argId, {
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
