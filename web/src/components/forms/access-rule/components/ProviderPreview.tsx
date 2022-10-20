import { Flex, HStack, Spacer, Text, VStack, Wrap } from "@chakra-ui/react";
import React from "react";
import { useFormContext } from "react-hook-form";

import {
  useGetProvider,
  useGetProviderArgs,
  useListProviderArgOptions,
} from "../../../../utils/backend-client/admin/admin";
import { CopyableOption } from "../../../CopyableOption";
import { ProviderIcon } from "../../../icons/providerIcon";
import { AccessRuleFormData } from "../CreateForm";

export const ProviderPreview: React.FC = () => {
  const { watch } = useFormContext<AccessRuleFormData>();
  const target = watch("target");
  const { data } = useGetProviderArgs(target?.providerId || "");
  const { data: provider } = useGetProvider(target?.providerId);

  if (
    target?.providerId === undefined ||
    target?.providerId === "" ||
    data === undefined ||
    provider === undefined
  ) {
    return null;
  }

  return (
    <VStack w="100%" align="flex-start">
      <HStack>
        <ProviderIcon shortType={provider.type} />
        <Text>{provider.id}</Text>
      </HStack>
      <VStack w="100%" align={"flex-start"} spacing={2}>
        {data &&
          target.multiSelects &&
          Object.entries(target.multiSelects).map(([k, v]) => {
            const arg = data[k];

            // This will now fetch all arg options i.e.
            // { label: 'AWSReadOnlyAccess', value: 'arn:aws...' }
            // This can make our flat values copyable
            const { data: argOptions } = useListProviderArgOptions(
              provider.id,
              k
            );
            if (v.length === 0) return null;

            return (
              <VStack w="100%" align={"flex-start"} spacing={0} key={k}>
                <Text>{arg.title}</Text>
                <Wrap>
                  {v?.map((opt) => {
                    return (
                      <CopyableOption
                        key={"cp-" + opt}
                        label={
                          argOptions?.options?.find((d) => d.value === opt)
                            ?.label ?? ""
                        }
                        value={opt}
                      />
                    );
                  })}
                </Wrap>
                {target.argumentGroups &&
                  target.argumentGroups[k] &&
                  arg.groups &&
                  Object.entries(target.argumentGroups[k]).map(
                    ([groupId, groupValues]) => {
                      if (!arg.groups || groupValues.length === 0) {
                        return null;
                      }
                      const group = arg.groups[groupId];
                      return (
                        <VStack>
                          <Text>{group.title}</Text>
                          {groupValues.map((groupValue) => {
                            if (!argOptions?.groups) {
                              return null;
                            }
                            const groupOptions = argOptions.groups[groupId];
                            return (
                              <CopyableOption
                                key={"cp-" + groupValue}
                                label={
                                  groupOptions.find(
                                    (d) => d.value === groupValue
                                  )?.label ?? ""
                                }
                                value={groupValue}
                              />
                            );
                          })}
                        </VStack>
                      );
                    }
                  )}
              </VStack>
            );
          })}
      </VStack>
    </VStack>
  );
};

export const ProviderPreviewOnlyStep: React.FC = () => {
  return (
    <VStack px={8} py={8} bg="neutrals.100" rounded="md" w="100%">
      <Flex w="100%">
        <Text textStyle="Heading/H3" opacity={0.6}>
          Provider
        </Text>
        <Spacer />
      </Flex>
      <ProviderPreview />
    </VStack>
  );
};
