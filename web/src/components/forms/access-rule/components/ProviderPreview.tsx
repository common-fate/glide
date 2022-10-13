import {
  Box,
  Flex,
  HStack,
  Spacer,
  Text,
  VStack,
  Wrap,
} from "@chakra-ui/react";
import Form, { FieldProps } from "@rjsf/core";
import React from "react";

import {
  useGetProviderArgs,
  useListProviderArgOptions,
} from "../../../../utils/backend-client/admin/admin";
import { AccessRuleTarget } from "../../../../utils/backend-client/types";
import { CopyableOption } from "../../../CopyableOption";
import { ProviderIcon } from "../../../icons/providerIcon";

// TODO: Update ProviderPreview component based on new arg schema response object.
export const ProviderPreview: React.FC<{ target: AccessRuleTarget }> = ({
  target,
}) => {
  const { data } = useGetProviderArgs(target.provider?.id || "");
  if (
    target.provider?.id === undefined ||
    target.provider?.id === "" ||
    data === undefined
  ) {
    return null;
  }
  const ProviderPreviewWithDisplay: React.FC<FieldProps> = (props) => {
    const { data } = useListProviderArgOptions(target.provider.id, props.name);

    if (props.name in target.with) {
      const value = target.with[props.name];
      return (
        <VStack w="100%" align={"flex-start"} spacing={0}>
          <Text>
            {props.schema.title}: {value}
          </Text>
        </VStack>
      );
    }
    if (props.name in target.withSelectable) {
      const value = target.withSelectable[props.name];
      return (
        <VStack w="100%" align={"flex-start"} spacing={0}>
          <Text>{props.schema.title}</Text>
          <Wrap>
            {value.map((opt: any) => {
              return (
                <CopyableOption
                  key={"cp-" + opt}
                  label={
                    data?.options.find((d) => d.value === opt)?.label ?? ""
                  }
                  value={opt}
                />
              );
            })}
          </Wrap>
        </VStack>
      );
    }
    return (
      <VStack w="100%" align={"flex-start"} spacing={0}>
        <Text>{props.schema.title}:</Text>
      </VStack>
    );
  };
  // Using a schema form here to do the heavy lifting of parsing the schema
  //  so we can get field names
  return (
    <VStack w="100%" align="flex-start">
      <HStack>
        <ProviderIcon shortType={target.provider.type} />

        <Text>{target.provider.id}</Text>
      </HStack>

      {/* The purpose of the form here is that is will do the heavy lifting of parsing the json schema for us, we then render each attribute with a display only element and remove all the form components, so the purpose is purly to part jsonschema and provider some rendering hooks */}
      <Box w="100%">
        <Form
          // tagname is a prop that allows us to prevent this using a <form> element to wrap this, this avoids a nested form error
          tagName={"div"}
          uiSchema={{
            "ui:options": { title: false },
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
            StringField: ProviderPreviewWithDisplay,
          }}
          // This field template override removes all form elements wrapping fields, this is so that no form ui is rendered other than the stringField override above
          FieldTemplate={(props) => props.children}
        ></Form>
      </Box>
    </VStack>
  );
};

export const ProviderPreviewOnlyStep: React.FC<{
  target: AccessRuleTarget;
}> = ({ target }) => {
  return (
    <VStack px={8} py={8} bg="neutrals.100" rounded="md" w="100%">
      <Flex w="100%">
        <Text textStyle="Heading/H3" opacity={0.6}>
          Provider
        </Text>
        <Spacer />
      </Flex>

      <ProviderPreview target={target} />
    </VStack>
  );
};
