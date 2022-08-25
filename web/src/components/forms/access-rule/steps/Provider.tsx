import {
  Box,
  Flex,
  FormControl,
  FormErrorMessage,
  FormHelperText,
  FormLabel,
  HStack,
  IconButton,
  Input,
  Skeleton,
  SkeletonText,
  Spinner,
  Text,
} from "@chakra-ui/react";
import Form from "@rjsf/chakra-ui";
import { FieldProps } from "@rjsf/core";
import React, { useEffect, useState } from "react";
import { Controller, useFormContext } from "react-hook-form";
import RSelect from "react-select";
import {
  useGetProvider,
  useGetProviderArgs,
  useListProviderArgOptions,
  listProviderArgOptions,
} from "../../../../utils/backend-client/admin/admin";
import { colors } from "../../../../utils/theme/colors";
import ProviderSetupNotice from "../../../ProviderSetupNotice";
import { ProviderPreview } from "../components/ProviderPreview";
import { ProviderRadioSelector } from "../components/ProviderRadio";
import { CustomOption } from "../components/Select";
import { CreateAccessRuleFormData } from "../CreateForm";
import { FormStep } from "./FormStep";
import { JSONSchema7 } from "json-schema";
import { RefreshIcon } from "../../../icons/Icons";

export const ProviderStep: React.FC = () => {
  const methods = useFormContext<CreateAccessRuleFormData>();
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
    if (!target || !provider || !target?.with) {
      return null;
    }
    return (
      <ProviderPreview
        target={{
          provider: provider,
          with: target.with,
        }}
      />
    );
  };
  return (
    <FormStep
      heading="Provider"
      subHeading="The group or role that the rule gives access to."
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
    <Form
      // tagname is a prop that allows us to prevent this using a <form> element to wrap this, this avoids a nested form error
      tagName={"div"}
      uiSchema={{
        "ui:submitButtonOptions": {
          props: {
            disabled: true,
            className: "btn btn-info",
          },
          norender: true,
          submitText: "",
        },
      }}
      showErrorList={false}
      schema={data}
      fields={{
        StringField: SelectField,
        // I would have overridden the DescriptionField to make it formatted nicer but its broken in RJSF :(
        // using a FieldTemplate does allow you to overide the whole thing sort of, but then we may as well write our own library
      }}
    ></Form>
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

// SelectField is used to render the select input for a provider args field, the data is saved to target.with.<fieldName> in the formdata
const SelectField: React.FC<FieldProps> = (props) => {
  const {
    control,
    watch,
    formState,
    trigger,
  } = useFormContext<CreateAccessRuleFormData>();
  const providerId = watch("target.providerId");
  const { data } = useListProviderArgOptions(providerId, props.name);
  const withError = formState.errors.target?.with;
  if (data === undefined) {
    return (
      <FormControl
        isInvalid={withError && withError[props.name] !== undefined}
        w="100%"
      >
        <FormLabel htmlFor="target.providerId">
          <Text textStyle={"Body/Medium"}>{props.schema.title}</Text>
        </FormLabel>
        <Skeleton h={8} />
      </FormControl>
    );
  }
  return (
    <FormControl
      isInvalid={withError && withError[props.name] !== undefined}
      w="100%"
    >
      <FormLabel htmlFor="target.providerId">
        <Text textStyle={"Body/Medium"}>{props.schema.title}</Text>
      </FormLabel>
      <Controller
        control={control}
        rules={{ required: props.required }}
        name={`target.with.${props.name}`}
        render={({ field: { onChange, ref, value } }) => {
          return data.hasOptions ? (
            <HStack minW={{ base: "200px", md: "500px" }}>
              <Box>
                <HStack>
                  <RSelect
                    options={data.options}
                    components={{ Option: CustomOption }}
                    ref={ref}
                    value={data.options.find((o) => o.value === value)}
                    onChange={(val) => {
                      // TS improperly infers this as MultiValue<Option>, when Option works fine?
                      // @ts-ignore
                      onChange(val?.value);
                      void trigger(`target.with.${props.name}`);
                    }}
                    styles={{
                      multiValue: (provided, state) => {
                        return {
                          minWidth: "100%",
                          borderRadius: "20px",
                          background: colors.neutrals[100],
                        };
                      },
                      container: (provided, state) => {
                        return {
                          minWidth: "300px",
                          width: "100%",
                        };
                      },
                    }}
                    // data-testid={rest.testId}
                  />
                  <RefreshButton providerId={providerId} argId={props.name} />
                </HStack>
                <FormHelperText>{value}</FormHelperText>
              </Box>
            </HStack>
          ) : (
            <>
              <Input
                id="provider-vault"
                bg="white"
                ref={ref}
                onChange={(e) => {
                  onChange(e);
                  void trigger(`target.with.${props.name}`); // this triggers the form to revalidate
                }}
                value={value}
                placeholder={props.schema.default?.toString() ?? ""}
              />
            </>
          );
        }}
      />
      <FormErrorMessage>{props.schema.title} is required</FormErrorMessage>
    </FormControl>
  );
};
