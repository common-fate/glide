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
import { ArgumentFormElement } from "../../../../utils/backend-client/types/accesshandler-openapi.yml";
import {
  listProviderArgOptions,
  useGetProvider,
  useGetProviderArgs,
  useListProviderArgOptions,
} from "../../../../utils/backend-client/admin/admin";

import { RefreshIcon } from "../../../icons/Icons";
import ProviderSetupNotice from "../../../ProviderSetupNotice";
import ArgField from "../components/ArgField";
import { ProviderPreview } from "../components/ProviderPreview";
import { ProviderRadioSelector } from "../components/ProviderRadio";
import { AccessRuleFormData } from "../CreateForm";
import { FormStep } from "./FormStep";

export const ProviderStep: React.FC = () => {
  const methods = useFormContext<AccessRuleFormData>();
  const target = methods.watch("target");

  const { data: provider } = useGetProvider(target?.providerId);
  const { data: providerArgs } = useGetProviderArgs(target?.providerId ?? "");

  // trigger a refresh of all provider arg options in the background when the provider is selected.
  // this helps to keep the cached options fresh.
  useEffect(() => {
    if (providerArgs != null) {
      const args = Object.values(providerArgs);

      // TODO: Currenly, we have only multi-select and input form element defined.
      // If in future, we have other form element that doesn't have options then we need to change
      // the if condition here.
      args.forEach((arg) => {
        if (arg.formElement != ArgumentFormElement.INPUT) {
          void listProviderArgOptions(target.providerId, arg.id, {
            refresh: true,
          });
        }
      });
    }
  }, [providerArgs, target?.providerId]);

  const Preview = () => {
    if (!target || !provider || !(target?.inputs || target?.multiSelects)) {
      return null;
    }
    return <ProviderPreview />;
  };
  return (
    <FormStep
      heading="Provider"
      subHeading="The permissions that the rule gives access to"
      fields={["target.with", "target.providerId"]}
      preview={<Preview />}
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
        <ProviderWithQuestions />
      </>
    </FormStep>
  );
};

// // Enable chakra styling of the react json schema form component!!!!
// // https://chakra-ui.com/docs/styled-system/chakra-factory
// const StyledForm = chakra(Form);
const ProviderWithQuestions: React.FC = () => {
  const { watch } = useFormContext();
  const providerId = watch("target.providerId");
  const { data } = useGetProviderArgs(providerId ?? "");

  if (providerId === undefined || providerId === "") {
    return null;
  }

  if (data === undefined) {
    return <Spinner />;
  }

  return (
    <>
      {Object.values(data).map((v) => (
        <ArgField argument={v} providerId={providerId} />
      ))}
    </>
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
  const { mutate } = useListProviderArgOptions(providerId, argId);

  const onClick = async () => {
    setLoading(true);
    const res = await listProviderArgOptions(providerId, argId, {
      refresh: true,
    });
    await mutate(res);
    setLoading(false);
  };

  return (
    <Tooltip>
      <IconButton
        {...props}
        onClick={onClick}
        isLoading={loading}
        icon={<RefreshIcon boxSize="24px" />}
        aria-label="Refresh"
        variant={"ghost"}
      />
    </Tooltip>
  );
};
