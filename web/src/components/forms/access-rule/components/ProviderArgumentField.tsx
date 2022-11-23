import {
  Box,
  Circle,
  FormControl,
  FormErrorMessage,
  FormLabel,
  HStack,
  Input,
  Text,
  Tooltip,
  VStack,
  Wrap,
} from "@chakra-ui/react";
import React from "react";
import { useFormContext } from "react-hook-form";
import { useListProviderArgOptions } from "../../../../utils/backend-client/admin/admin";
import {
  Argument,
  ArgumentRuleFormElement,
  GroupOption,
  Option,
} from "../../../../utils/backend-client/types/accesshandler-openapi.yml";
import { DynamicOption } from "../../../DynamicOption";
import { BoltIcon } from "../../../icons/Icons";
import { AccessRuleFormData } from "../CreateForm";
import { RefreshButton } from "../steps/Provider";
import { MultiSelect } from "./Select";

interface ProviderArgumentFieldProps {
  argument: Argument;
  providerId: string;
}

const ProviderArgumentField: React.FC<ProviderArgumentFieldProps> = ({
  argument,
  providerId,
}) => {
  switch (argument.ruleFormElement) {
    case ArgumentRuleFormElement.MULTISELECT:
      return (
        <ProviderFormElementMultiSelect
          argument={argument}
          providerId={providerId}
        />
      );
    case ArgumentRuleFormElement.INPUT:
      return (
        <ProviderFormElementInput argument={argument} providerId={providerId} />
      );
    default:
      return (
        <ProviderFormElementInput argument={argument} providerId={providerId} />
      );
  }
};

export default ProviderArgumentField;

const ProviderFormElementInput: React.FC<ProviderArgumentFieldProps> = ({
  argument,
  providerId,
}) => {
  const { formState, register, trigger } = useFormContext<AccessRuleFormData>();
  const inputs = formState.errors.target?.inputs;
  const { onBlur, ...rest } = register(`target.inputs.${argument.id}`, {
    minLength: 1,
    required: true,
  });
  return (
    <FormControl
      w="100%"
      isInvalid={inputs && inputs[argument.id] !== undefined}
    >
      <FormLabel htmlFor="target.providerId" display="inline">
        <Text textStyle={"Body/Medium"}>{argument.title}</Text>
      </FormLabel>
      <Input
        id="provider-vault"
        bg="white"
        placeholder={"example"}
        onBlur={(e) => {
          void trigger(`target.inputs.${argument.id}`);
          void onBlur(e);
        }}
        {...rest}
      />
      <FormErrorMessage> {argument.title} is required </FormErrorMessage>
    </FormControl>
  );
};
const ProviderFormElementMultiSelect: React.FC<ProviderArgumentFieldProps> = ({
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
          const group = argOptions?.groups ? argOptions?.groups[groupId] : [];
          return selectedGroupValues.flatMap((groupValue) => {
            return group.find((g) => g.value === groupValue)?.children || [];
          });
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

  // Will be true if no groups have been selected
  const argumentSelectionRequired =
    Object.entries(argumentGroups || {}).find(([k, v]) => v.length > 0) ===
    undefined;

  // true if no argument values have been selected
  const groupSelectionRequired =
    multiSelects === undefined || multiSelects.length === 0;
  console.log({
    id: argument.id,
    argumentSelectionRequired,
    groupSelectionRequired,
    argumentGroups,
    multiSelects,
    watch: watch(),
    argument,
    argOptions,
  });

  return (
    <VStack
      data-testid="argumentField"
      border="1px solid"
      borderColor="gray.300"
      rounded="md"
      p={4}
      // py={6}
      w="100%"
      spacing={4}
      justifyContent="start"
      alignItems="start"
    >
      <FormControl
        w="100%"
        isInvalid={
          multiSelectsError && multiSelectsError[argument.id] !== undefined
        }
      >
        <FormLabel htmlFor="target.providerId">
          <Text textStyle={"Body/Medium"}>{argument.title}s</Text>
          {argument.description && (
            <Text textStyle={"Body/Medium"} color="neutrals.500">
              {argument.description}
            </Text>
          )}
        </FormLabel>
        <HStack w="90%">
          <MultiSelect
            rules={{ required: argumentSelectionRequired, minLength: 1 }}
            fieldName={`target.multiSelects.${argument.id}`}
            options={argOptions?.options || []}
            shouldAddSelectAllOption={true}
            id="providedArgumentField"
          />
          <RefreshButton argId={argument.id} providerId={providerId} mx={20} />
        </HStack>

        {!argument.groups && (
          <FormErrorMessage> {argument.title} is required </FormErrorMessage>
        )}
      </FormControl>

      {/* @TODO: consider adding skeleton group or improving CLS */}
      {argument.groups && (
        <Box
          pos="relative"
          w={{ base: "100%", md: "100%" }}
          minW={{ base: "100%", md: "400px", lg: "500px" }}
        >
          {argument?.groups &&
            Object.values(argument.groups).map((group) => {
              // catch the unexpected case where there are no options for group
              if (
                argOptions?.groups == undefined ||
                !argOptions.groups?.[group.id]
              ) {
                return null;
              }
              return (
                <FormControl
                  w="100%"
                  isInvalid={
                    multiSelectsError &&
                    multiSelectsError[argument.id] !== undefined
                  }
                >
                  <>
                    <FormLabel htmlFor="target.providerId">
                      <Text display="inline" textStyle={"Body/Medium"}>
                        {group.title}{" "}
                      </Text>{" "}
                      <Tooltip label="Dynamic Field" hasArrow={true}>
                        <Circle
                          display="inline-flex"
                          size="24px"
                          px={1}
                          // bg="gray.200"
                          rounded="full"
                        >
                          <BoltIcon boxSize="12px" color="neutrals.400" />
                        </Circle>
                      </Tooltip>
                      {group.description && (
                        <Text textStyle={"Body/Medium"} color="neutrals.500">
                          {group.description}
                        </Text>
                      )}
                    </FormLabel>
                    <HStack w="90%">
                      <MultiSelect
                        rules={{
                          required: groupSelectionRequired,
                          minLength: 1,
                        }}
                        fieldName={`target.argumentGroups.${argument.id}.${group.id}`}
                        options={argOptions.groups[group.id] || []}
                        shouldAddSelectAllOption={true}
                      />
                    </HStack>
                  </>
                </FormControl>
              );
            })}
        </Box>
      )}
      {argOptions?.groups &&
        Object.entries(argOptions?.groups ?? {}).length > 0 && (
          <Box>
            <Wrap>
              {uniqueRes &&
                uniqueRes.map((c) => {
                  return (
                    <DynamicOption
                      label={c.option.label}
                      value={c.option.value}
                      isParentGroup={!!c.parentGroup}
                    />
                  );
                })}
            </Wrap>
          </Box>
        )}
      {effectiveOptions.length == 0 &&
        Object.entries(argOptions?.groups ?? {}).length > 0 &&
        (formState.touchedFields.target?.argumentGroups?.[argument.id] ||
          formState.touchedFields.target?.multiSelects?.[argument.id]) && (
          <Text color="red.500" fontSize="sm">
            {"At least one " + argument.title + " is required"}
          </Text>
        )}
    </VStack>
  );
};
