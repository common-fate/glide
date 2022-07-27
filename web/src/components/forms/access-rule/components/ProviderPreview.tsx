import { Box, Flex, HStack, Spacer, Text, VStack } from "@chakra-ui/react";
import Form, { FieldProps } from "@rjsf/core";
import React, { useEffect, useState } from "react";
import {
  useGetProviderArgs,
  useListProviderArgOptions,
} from "../../../../utils/backend-client/default/default";
import { AccessRuleTarget } from "../../../../utils/backend-client/types";
import { ProviderIcon } from "../../../icons/providerIcon";

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

    const value = target.with[props.name];
    const [label, setLabel] = useState<string>();
    useEffect(() => {
      const l = data?.options.find((d) => d.value === value);
      setLabel(l?.label ?? "");
    }, [data, value]);
    return (
      <VStack w="100%" align={"flex-start"} spacing={0}>
        <Text>
          {props.schema.title}: {label}
        </Text>
        <Text textStyle={"Body/SmallBold"}>{value}</Text>
      </VStack>
    );
  };
  // Using a schema form here to do the heavy lifting of parsing the schema
  //  so we can get field names
  return (
    <VStack w="100%" align="flex-start">
      <HStack>
        <ProviderIcon provider={target.provider} />

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
