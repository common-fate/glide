import {
  chakra,
  FormControl,
  FormErrorMessage,
  FormLabel,
  Input,
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
import { MultiSelectOptions } from "../components/Select";
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
          with: {},
          // with: target.with,
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

// Enable chakra styling of the react json schema form component!!!!
// https://chakra-ui.com/docs/styled-system/chakra-factory
const StyledForm = chakra(Form);
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
    <StyledForm
      // using chakra styling props to set the width to 100%
      w="100%"
      // tagname is a prop that allows us to prevent this using a <form> element to wrap this, this avoids a nested form error
      tagName={"div"}
      uiSchema={{
        "ui:options": {
          chakra: {
            w: "100%",
          },
        },
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
        StringField: WithField,
      }}
    ></StyledForm>
  );
};

// WithField is used to render the select input for a provider args field, the data is saved to target.with.<fieldName> in the formdata
const WithField: React.FC<FieldProps> = (props) => {
  const { watch, formState, register } = useFormContext();
  const providerId = watch("target.providerId");
  const { data } = useListProviderArgOptions(providerId, props.name);
  const withError = formState.errors.target?.with;

  // @TODO: find a way to ensure validation chill bruh

  if (data === undefined) {
    return <Spinner />;
  }
  return (
    <FormControl isInvalid={withError && withError[props.name]} w="100%">
      <FormLabel htmlFor="target.providerId">
        <Text textStyle={"Body/Medium"}>{props.schema.title}</Text>
      </FormLabel>
      {data.hasOptions ? (
        <MultiSelectOptions
          fieldName={`target.with.${props.name}`}
          options={data.options}
        />
      ) : (
        <Input
          id="provider-vault"
          bg="white"
          placeholder={props.schema.default?.toString() ?? ""}
          {...register(`target.with.${props.name}`)}
        />
      )}
      <FormErrorMessage>{props.schema.title} is required</FormErrorMessage>
    </FormControl>
  );
};
