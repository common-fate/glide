import {
  Box,
  Flex,
  HStack,
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
import {
  Argument,
  GroupOption,
  Option,
  Provider,
} from "../../../../utils/backend-client/types/accesshandler-openapi.yml";
import { DynamicOption } from "../../../DynamicOption";
import { BoltIcon } from "../../../icons/Icons";
import { ProviderIcon } from "../../../icons/providerIcon";
import { AccessRuleFormData } from "../CreateForm";

interface ProviderPreviewProps {
  provider: Provider;
}

export const ProviderPreview: React.FC<ProviderPreviewProps> = ({
  provider,
}) => {
  const { data: providerArgs } = useGetProviderArgs(provider.id ?? "");

  if (!provider) return null;

  return (
    <VStack w="100%" align="flex-start">
      <HStack>
        <ProviderIcon shortType={provider.type} />
        <Text>{provider.id}</Text>
      </HStack>
      {providerArgs &&
        Object.values(providerArgs).map((v) => (
          <PreviewArgument argument={v} providerId={provider.id} />
        ))}
    </VStack>
  );
};

interface ProviderArgFieldProps {
  argument: Argument;
  providerId: string;
}

export const PreviewArgument: React.FC<ProviderArgFieldProps> = ({
  argument,
  providerId,
}) => {
  const { formState, watch } = useFormContext<AccessRuleFormData>();

  const { data: argOptions } = useListProviderArgOptions(
    providerId,
    argument.id
  );
  const multiSelectsError = formState.errors.target?.multiSelects;

  const [argumentGroups, multiSelects] = watch([
    `target.argumentGroups.${argument.id}`,
    `target.multiSelects.${argument.id}`,
  ]);

  /** get all the group children (for aws these are the accounts for the OUs) */
  const effectiveViaGroups =
    (argumentGroups &&
      Object.entries(argumentGroups || {}).flatMap(
        ([groupId, selectedGroupValues]) => {
          // get all the accounts for the selected group value
          const allOptionsForSelectedGroup = argOptions?.groups
            ? argOptions?.groups[groupId]
            : [];

          return selectedGroupValues.flatMap((groupValue) => {
            const selectedGroupDetails = allOptionsForSelectedGroup.filter(
              (group) => selectedGroupValues.includes(group.value)
            );

            return selectedGroupDetails.flatMap((g) => g.children || []);
            // return group.find((g) => g.value === groupValue)?.children || [];
          });
        }
      )) ||
    [];

  const selectedGroups =
    (argumentGroups &&
      Object.entries(argumentGroups || {}).flatMap(
        ([groupId, selectedGroupValues]) => {
          // get all the accounts for the selected group value
          const groupDetails = argOptions?.groups
            ? argOptions?.groups[groupId]
            : ([] as GroupOption[]);

          const group = groupDetails.filter((group) =>
            selectedGroupValues.includes(group.value)
          );
          return group;
        }
      )) ??
    [];

  type Obj = {
    option: Option;
    parentGroup?: GroupOption;
  };

  // Desired output: an array of Options (value and label) with ID(s) that can lookup into the richer map object
  // Step 1: get all the options from the groups
  // Step 2: get all the options from the multi-selects
  // Step 3: store the options in an array of type Obj
  // Step 4: remove any duplicate Option.key Option.value paires

  const effectiveGroups =
    argumentGroups &&
    Object.entries(argumentGroups || {}).flatMap(
      ([groupId, selectedGroupValues]) => {
        // get all the accounts for the selected group value
        const group = argOptions?.groups ? argOptions?.groups?.[groupId] : [];
        return (
          selectedGroupValues
            .flatMap((groupValue) => {
              return group.find((g) => g.value === groupValue) ?? null;
            })
            // Now remove any null values
            .filter((g) => g)
        );
      }
    );

  const res: Obj[] = [];

  effectiveGroups?.forEach((g) => {
    g?.children?.forEach((c) => {
      const option = argOptions?.options?.find((o) => o.value === c);
      option && res.push({ option, parentGroup: g });
    });
  });
  multiSelects?.forEach((ms) => {
    const option = argOptions?.options?.find((o) => o.value === ms);
    option && res.push({ option });
  });

  // Now remove any duplicate Option.label Option.value paires
  // With a preference for the parentGroup option
  const uniqueRes = res.reduce((acc, cur) => {
    const existing = acc.find(
      (a) =>
        a.option.value === cur.option.value &&
        a.option.label === cur.option.label
    );
    if (existing) {
      if (cur.parentGroup) {
        acc.splice(acc.indexOf(existing), 1, cur);
      }
    } else {
      acc.push(cur);
    }
    return acc;
  }, [] as Obj[]);

  // effectiveViaGroups.filter(g => {
  // })

  // Filter to remove duplicates
  const effectiveAccountIds = [
    ...(multiSelects || []),
    ...effectiveViaGroups,
  ].filter((v, i, a) => a.indexOf(v) === i);

  // DE = we want to associate each 'effectiveOption' with either a group (dynamic) or a multiSelect (single)
  // We can do this by checking if the effectiveOption is in the multiSelects array
  // Or, avoid casting it to a string to begin with
  const effectiveOptions =
    argOptions?.options.filter((option) => {
      return effectiveAccountIds.includes(option.value);
    }) || [];
  const required = effectiveOptions.length === 0;

  return (
    <VStack
      w="100%"
      align={"flex-start"}
      spacing={4}
      // key={k}
      p={4}
      rounded="md"
      border="1px solid"
      borderColor="gray.300"
    >
      <Box>
        <Text textStyle={"Body/Medium"} color="neutrals.500">
          {argument.title}s
        </Text>
        {/* {arg.description && (
                    <Text textStyle={"Body/Medium"} color="neutrals.500">
                      {arg.description}
                    </Text>
                  )} */}
        <Wrap>
          {[...effectiveAccountIds].map((opt) => {
            return (
              <DynamicOption
                key={"cp-" + opt}
                label={
                  argOptions?.options?.find((d) => d.value === opt)?.label ?? ""
                }
                value={opt}
                parentGroup={multiSelects.find(c => c === opt)? undefined : ["random"] as any}
              />
            );
          })}
        </Wrap>
      </Box>
      <>
        {argumentGroups && !!selectedGroups.length &&
          Object.entries(argumentGroups).map(([groupName, groupId]) => {
            return (
              <Box>
                <Flex>
                  <Text textStyle={"Body/Medium"} color="neutrals.500">
                    {groupName}s
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

                {selectedGroups && (
                  <Box>
                    <Wrap>
                      {selectedGroups.map((group) => (
                        <DynamicOption
                          label={group.label}
                          value={group.value}
                          parentGroup={group}
                        />
                      ))}
                    </Wrap>
                  </Box>
                )}
              </Box>
            );
          })}
      </>
    </VStack>
  );
};
