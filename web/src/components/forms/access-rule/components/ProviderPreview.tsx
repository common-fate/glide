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
import {
  AccessRuleTarget,
  Provider,
} from "../../../../utils/backend-client/types";
import { CopyableOption } from "../../../CopyableOption";
import { ProviderIcon } from "../../../icons/providerIcon";
import { AccessRuleFormDataTarget } from "../CreateForm";

// TODO: Update ProviderPreview component based on new arg schema response object.
export const ProviderPreview: React.FC<{
  target: AccessRuleFormDataTarget;
  provider: Provider;
}> = ({ target, provider }) => {
  const { data } = useGetProviderArgs(provider?.id || "");

  console.log({ target, useGetProviderArgs: data });

  if (provider?.id === undefined || provider?.id === "" || data === undefined) {
    return null;
  }
  const ProviderPreviewWithDisplay: React.FC<FieldProps> = (props) => {
    const { data } = useListProviderArgOptions(provider.id, props.name);

    // Handling for target.with params
    if (props.name in target.with) {
      const value = target.with[props.name];
      let string = Array.isArray(value) ? value.join(",") : value;
      return (
        <VStack w="100%" align={"flex-start"} spacing={0}>
          <Text>
            ok
            {/* {props.schema.title}: {string} */}
          </Text>
        </VStack>
      );
    }
    // Handling for target.withFilter params
    // if (props.name in target.withFilter) {
    //   const value = target.withFilter[props.name];
    //   return (
    //     <VStack w="100%" align={"flex-start"} spacing={0}>
    //       <Text>{props.schema.title}</Text>
    //       <Wrap>
    //         {value.map((opt: any) => {
    //           return (
    //             <CopyableOption
    //               key={"cp-" + opt}
    //               label={
    //                 data?.options.find((d) => d.value === opt)?.label ?? ""
    //               }
    //               value={opt}
    //             />
    //           );
    //         })}
    //       </Wrap>
    //     </VStack>
    //   );
    // }
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
        <ProviderIcon shortType={provider.type} />
        <Text>{provider.id}</Text>
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
          // showErrorList={false}
          showErrorList={true}
          schema={data}
          fields={{
            StringField: ProviderPreviewWithDisplay,
          }}
          // This field template override removes all form elements wrapping fields, this is so that no form ui is rendered other than the stringField override above
          FieldTemplate={(props) => props.children}
        />
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

      {/* @TODO resolve typing issue once above is compelte  */}
      {/* <ProviderPreview target={target} provider={target.provider} /> */}
    </VStack>
  );
};
