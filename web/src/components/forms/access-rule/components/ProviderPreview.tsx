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
  // I need to be run per arg... (i should be in a for loop)
  // const { data: argOptions } = useListProviderArgOptions(provider.id, props.name);

  // Using a schema form here to do the heavy lifting of parsing the schema
  //  so we can get field names
  return (
    <VStack w="100%" align="flex-start">
      <HStack>
        <ProviderIcon shortType={provider.type} />
        <Text>{provider.id}</Text>
      </HStack>
      {data &&
        Object.keys(data).map((key) => {
          const arg = data[key];
          // const { data: argOptions } = useListProviderArgOptions(provider.id, arg.id);
          return (
            <VStack w="100%" align={"flex-start"} spacing={0}>
              <Text>{arg.title}</Text>
              <Wrap>
                {/* {value.map((opt: any) => {
            return (
              <CopyableOption
                key={"cp-" + opt}
                label={data?.options.find((d) => d.value === opt)?.label ?? ""}
                value={opt}
              />
            );
          })} */}
              </Wrap>
            </VStack>
          );
        })}
      {/* <Box w="100%">
      </Box> */}
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
