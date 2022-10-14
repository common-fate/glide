import {
  chakra,
  FormControl,
  FormErrorMessage,
  FormLabel,
  HStack,
  IconButton,
  Input,
  Skeleton,
  Spinner,
  Text,
} from "@chakra-ui/react";
import Form from "@rjsf/chakra-ui";
import { FieldProps } from "@rjsf/core";
import { JSONSchema7 } from "json-schema";
import React, { useEffect, useState } from "react";
import { Controller, useFormContext } from "react-hook-form";
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
import { MultiSelect } from "../components/Select";
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
      // example schema
      // {"$defs":{"Args":{"properties":{"vault":{"description":"example","title":"Vault","type":"string"}},"required":["vault"],"type":"object"}},"$id":"https://commonfate.io/demo/1password/args","$ref":"#/$defs/Args","$schema":"http://json-schema.org/draft/2020-12/schema"}
      const schema = providerArgs as JSONSchema7;
      const argSchema = schema.$defs?.Args;
      if (argSchema !== undefined && typeof argSchema !== "boolean") {
        const args = Object.keys(argSchema.properties ?? {});
        args.forEach((arg) => {
          void listProviderArgOptions(target.providerId, arg, {
            refresh: true,
          });
        });
      }
    }
  }, [providerArgs, target?.providerId]);

  const Preview = () => {
    if (!target || !provider || !(target?.with || target?.withText)) {
      return null;
    }

    return (
      <ProviderPreview
        target={{
          provider: provider,
          withSelectable: {},
          with: target.withText ?? {},
        }}
      />
    );
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
      {Object.keys(data).map((arg) => (
        <ArgField value={arg} providerId={providerId} />
      ))}
    </>
  );
};

interface RefreshButtonProps {
  providerId: string;
  argId: string;
}

const RefreshButton: React.FC<RefreshButtonProps> = ({ argId, providerId }) => {
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
    <IconButton
      onClick={onClick}
      isLoading={loading}
      icon={<RefreshIcon boxSize="24px" />}
      aria-label="Refresh"
      variant={"ghost"}
    />
  );
};

// TODO:
// Instead of deleting this commenting out for now.
//
//
//
// WithField is used to render the select input for a provider args field, the data is saved to target.with.<fieldName> in the formdata
// const WithField: React.FC<FieldProps> = (props) => {
//   const { watch, formState, register } = useFormContext<AccessRuleFormData>();
//   const providerId = watch("target.providerId");
//   const { data } = useListProviderArgOptions(providerId, props.name);
//   const withError = formState.errors.target?.with;
//   if (data === undefined) {
//     return (
//       <FormControl
//         isInvalid={withError && withError[props.name] !== undefined}
//         w="100%"
//       >
//         <FormLabel htmlFor="target.providerId">
//           <Text textStyle={"Body/Medium"}>{props.schema.title}</Text>
//         </FormLabel>
//         <Skeleton h={8} />
//       </FormControl>
//     );
//   }
//   return (
//     <FormControl
//       isInvalid={withError && withError[props.name] !== undefined}
//       w="100%"
//     >
//       <FormLabel htmlFor="target.providerId">
//         <Text textStyle={"Body/Medium"}>{props.schema.title}</Text>
//       </FormLabel>
//       {data.hasOptions ? (
//         <HStack>
//           <MultiSelect
//             rules={{ required: true, minLength: 1 }}
//             fieldName={`target.with.${props.name}`}
//             options={data.options}
//             shouldAddSelectAllOption={true}
//           />

//           <RefreshButton providerId={providerId} argId={props.name} />
//         </HStack>
//       ) : (
//         <Input
//           id="provider-vault"
//           bg="white"
//           placeholder={props.schema.default?.toString() ?? ""}
//           {...register(`target.withText.${props.name}`)}
//         />
//       )}
//       <FormErrorMessage>{props.schema.title} is required</FormErrorMessage>
//     </FormControl>
//   );
// };
