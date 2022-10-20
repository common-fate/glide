import {
  Box,
  Flex,
  HStack,
  Spacer,
  Text,
  Tooltip,
  VStack,
  Wrap,
  Circle,
} from "@chakra-ui/react";
import React from "react";
import { useFormContext } from "react-hook-form";

import {
  useGetProvider,
  useGetProviderArgs,
  useListProviderArgOptions,
} from "../../../../utils/backend-client/admin/admin";
import { CopyableOption } from "../../../CopyableOption";
import { DynamicOption } from "../../../DynamicOption";
import { BoltIcon } from "../../../icons/Icons";
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
      <VStack w="100%" align={"flex-start"} spacing={4}>
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
              <VStack
                w="100%"
                align={"flex-start"}
                spacing={4}
                key={k}
                p={4}
                rounded="md"
                border="1px solid"
                borderColor="gray.300"
              >
                <Box>
                  <Text textStyle={"Body/Medium"} color="neutrals.500">
                    {arg.title}s
                  </Text>
                  {/* {arg.description && (
                    <Text textStyle={"Body/Medium"} color="neutrals.500">
                      {arg.description}
                    </Text>
                  )} */}
                  <Wrap>
                    {v?.map((opt) => {
                      return (
                        <DynamicOption
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
                </Box>
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
                        <Box>
                          <Flex>
                            <Text
                              textStyle={"Body/Medium"}
                              color="neutrals.500"
                            >
                              {group.title}s
                            </Text>
                            <Tooltip label="Dynamic Field" hasArrow={true}>
                              <Circle
                                display="inline-flex"
                                size="24px"
                                px={1}
                                rounded="full"
                              >
                                <BoltIcon boxSize="12px" color="neutrals.400" />
                              </Circle>
                            </Tooltip>
                          </Flex>
                          {/* {group.description && (
                            <Text
                              textStyle={"Body/Medium"}
                              color="neutrals.500"
                            >
                              {group.description}
                            </Text>
                          )} */}
                          <Wrap>
                            {groupValues.map((groupValue) => {
                              if (!argOptions?.groups) {
                                return null;
                              }
                              const groupOptions = argOptions.groups[groupId];
                              return (
                                <DynamicOption
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
                          </Wrap>
                        </Box>
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
