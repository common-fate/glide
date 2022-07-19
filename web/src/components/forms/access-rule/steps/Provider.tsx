import {
  FormControl,
  FormErrorMessage,
  FormHelperText,
  FormLabel,
  Input,
  Select,
  Spinner,
  Text,
} from "@chakra-ui/react";
import Form from "@rjsf/chakra-ui";
import { FieldProps } from "@rjsf/core";
import React from "react";
import { Controller, useFormContext } from "react-hook-form";
import {
  useGetProvider,
  useGetProviderArgs,
  useListProviderArgOptions,
} from "../../../../utils/backend-client/default/default";
import ProviderSetupNotice from "../../../ProviderSetupNotice";
import { ProviderPreview } from "../components/ProviderPreview";
import { ProviderRadioSelector } from "../components/ProviderRadio";
import { CreateAccessRuleFormData } from "../CreateForm";

import { FormStep } from "./FormStep";

export const ProviderStep: React.FC = () => {
  const methods = useFormContext<CreateAccessRuleFormData>();
  const target = methods.watch("target");
  const Preview = () => {
    const { data: provider } = useGetProvider(target?.providerId);
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
                onChange={(t) => {
                  onChange(t);
                  methods.trigger("target.providerId");
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
      }}
    ></Form>
  );
};

// SelectField is used to render the select input for a provider args field, the data is saved to target.with.<fieldName> in the formdata
const SelectField: React.FC<FieldProps> = (props) => {
  const { control, watch, formState, unregister, trigger } = useFormContext();
  const providerId = watch("target.providerId");
  const { data } = useListProviderArgOptions(providerId, props.name);
  const withError = formState.errors.target?.with;

  if (data === undefined) {
    return <Spinner />;
  }
  return (
    <FormControl isInvalid={withError && withError[props.name]}>
      <FormLabel htmlFor="target.providerId">
        <Text textStyle={"Body/Medium"}>{props.schema.title}</Text>
      </FormLabel>
      <Controller
        control={control}
        // not sure how much validation we will support in these schemas, is this the best way to pass validation rules through?
        rules={{ required: props.required }}
        name={`target.with.${props.name}`}
        render={({ field: { onChange, ref, value } }) => {
          return data.hasOptions ? (
            <>
              <Select
                bg="white"
                ref={ref}
                value={data.options.find((o) => o.value === value)?.label}
                minW={{ base: "200px", md: "300px" }}
                onChange={(e: any) => {
                  {
                    onChange(
                      data.options.find((o) => o.label === e.target.value)
                        ?.value
                    );
                  }
                  trigger(`target.with.${props.name}`);
                }}
              >
                {<option></option>}
                {data.options.map((o, i) => (
                  <option key={o.label} value={o.label} label={o.label}>
                    {o.label}
                  </option>
                ))}
              </Select>
              <FormHelperText>{value}</FormHelperText>
            </>
          ) : (
            <>
              <Input
                id="provider-vault"
                bg="white"
                ref={ref}
                onChange={onChange}
                value={value}
                placeholder={props.schema.default?.toString() ?? ""}
              />
            </>
          );
        }}
      ></Controller>

      <FormErrorMessage>{props.schema.title} is required</FormErrorMessage>
    </FormControl>
  );
};
