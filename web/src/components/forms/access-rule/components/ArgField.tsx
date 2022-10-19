import {
  Box,
  Circle,
  FormControl,
  FormErrorMessage,
  FormLabel,
  Heading,
  HStack,
  Input,
  Tag,
  Text,
  Tooltip,
  VStack,
  Wrap,
} from "@chakra-ui/react";
import React, { useState, useEffect } from "react";
import { useFormContext } from "react-hook-form";
import ArgGroupView from "./ArgGroupView";
import { MultiSelect } from "../components/Select";
import { AccessRuleFormData } from "../CreateForm";
import { useListProviderArgOptions } from "../../../../utils/backend-client/admin/admin";
import {
  Argument,
  ArgumentFormElement,
  GroupOption,
  Option,
} from "../../../../utils/backend-client/types/accesshandler-openapi.yml";
import { CopyableOption } from "../../../CopyableOption";
import { DynamicOption } from "../../../DynamicOption";
import { BoltIcon } from "../../../icons/Icons";

interface ArgFieldProps {
  argument: Argument;
  providerId: string;
}

const ArgField = (props: ArgFieldProps) => {
  const { argument, providerId } = props;
  const {
    register,
    formState,
    getValues,
    watch,
  } = useFormContext<AccessRuleFormData>();

  const { data: argOptions } = useListProviderArgOptions(
    providerId,
    argument.id,
    {},
    {
      swr: {
        // don't call API if arg doesn't have options
        enabled: argument.formElement !== ArgumentFormElement.INPUT,
      },
    }
  );

  const multiSelectsError = formState.errors.target?.multiSelects;

  const [argumentGroups, multiSelects] = watch([
    `target.argumentGroups.${argument.id}`,
    `target.multiSelects.${argument.id}`,
  ]);

  // TODO: Form input error is not handled for input type.
  if (argument.formElement === ArgumentFormElement.INPUT) {
    return (
      <FormControl w="100%">
        <FormLabel htmlFor="target.providerId">
          <Text textStyle={"Body/Medium"}>{argument.title}</Text>
        </FormLabel>
        <Input
          id="provider-vault"
          bg="white"
          placeholder={`default-${argument.title}`}
          {...register(`target.inputs.${argument.id}`)}
        />
      </FormControl>
    );
  }

  type Obj = {
    option: Option;
    parentGroup?: GroupOption;
  };

  type ResArr = (Obj | Option)[];

  /** get all the group children (for aws these are the accounts for the OUs) */
  const effectiveViaGroups = Object.entries(argumentGroups || {}).flatMap(
    ([groupId, selectedGroupValues]) => {
      // get all the accounts for the selected group value
      const group = argOptions?.groups ? argOptions?.groups[groupId] : [];
      console.log({ groupId, selectedGroupValues, group });
      return selectedGroupValues.flatMap((groupValue) => {
        return group.find((g) => g.value === groupValue)?.children || [];
      });
    }
  );

  // This type allows us to store the options in a flat array (with additional information for group fields)
  // type EffectiveOptions = [string, GroupOption];

  // Desired output: an array of Options (value and label) with ID(s) that can lookup into the richer map object
  // Step 1: get all the options from the groups
  // Step 2: get all the options from the multi-selects
  // Step 3: store the options in an array of type Obj
  // Step 4: remove any duplicate Option.key Option.value paires

  const effectiveGroups = Object.entries(argumentGroups || {}).flatMap(
    ([groupId, selectedGroupValues]) => {
      // get all the accounts for the selected group value
      const group = argOptions?.groups ? argOptions?.groups[groupId] : [];
      console.log({ groupId, selectedGroupValues, group });
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

  console.log({ effectiveGroups, multiSelects });

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

  // console.log({
  //   effectiveOptions,
  //   touchings: formState.touchedFields,
  // });

  return (
    <VStack
      border="1px solid"
      borderColor="gray.300"
      rounded="md"
      p={4}
      py={6}
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
        <div>
          <FormLabel htmlFor="target.providerId">
            <Text textStyle={"Body/Medium"}>
              Individual&nbsp;{argument.title}s
            </Text>
          </FormLabel>
          <HStack>
            <MultiSelect
              rules={{ required: required, minLength: 1 }}
              fieldName={`target.multiSelects.${argument.id}`}
              options={argOptions?.options || []}
              shouldAddSelectAllOption={true}
            />
          </HStack>
        </div>
        <FormLabel htmlFor="target.providerId.filters.filterId"></FormLabel>
        {/* TODO: msg will eventually be more detailed (one or more options) */}
        {!argument.groups && (
          <FormErrorMessage> {argument.title} is required </FormErrorMessage>
        )}
      </FormControl>

      {argument.groups && (
        <Box
          mt={4}
          pos="relative"
          w={{ base: "100%", md: "100%" }}
          minW={{ base: "100%", md: "400px", lg: "500px" }}
        >
          {Object.values(argument.groups).map((group) => {
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
                  <FormLabel
                    htmlFor="target.providerId"
                    display="inline"
                    mb={4}
                  >
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
                        filter="grayscale(1);"
                        transition="all .2s ease"
                        _hover={{
                          filter: "grayscale(0);",
                        }}
                      >
                        <BoltIcon boxSize="12px" color="brandGreen.200" />
                      </Circle>
                    </Tooltip>
                  </FormLabel>
                  <HStack>
                    <MultiSelect
                      rules={{ required: required, minLength: 1 }}
                      fieldName={`target.argumentGroups.${argument.id}.${group.id}`}
                      options={argOptions.groups[group.id] || []}
                      shouldAddSelectAllOption={true}
                    />
                  </HStack>
                </>
                {/* <FormLabel htmlFor="target.providerId.filters.filterId"></FormLabel> */}
                {/* TODO: msg will eventually be more detailed (one or more options) */}
              </FormControl>
            );
          })}
        </Box>
      )}
      {effectiveOptions.length > 0 &&
        argOptions?.groups &&
        Object.entries(argOptions?.groups ?? {}).length > 0 && (
          <Box mt={2}>
            {/* <Text textStyle={"Body/Medium"}>{argument.title + "s"}</Text> */}
            <Wrap>
              {effectiveOptions &&
                effectiveOptions.map((c) => {
                  // console.log(c);
                  return <DynamicOption label={c.label} value={c.value} />;
                })}
            </Wrap>
          </Box>
        )}
      {effectiveOptions.length == 0 &&
        Object.entries(argOptions?.groups ?? {}).length > 0 &&
        (formState.touchedFields.target?.argumentGroups?.[argument.id] ||
          formState.touchedFields.target?.multiSelects?.[argument.id]) && (
          <Text color="red.500" fontSize="sm">
            {"At least one effective " + argument.title + " is required"}
          </Text>
        )}
    </VStack>
  );
};

export default ArgField;
