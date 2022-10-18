import {
  Box,
  FormControl,
  FormErrorMessage,
  FormLabel,
  Heading,
  HStack,
  Input,
  Text,
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
  Option,
} from "../../../../utils/backend-client/types/accesshandler-openapi.yml";
import { CopyableOption } from "../../../CopyableOption";

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

  const effectiveViaGroups = Object.entries(argumentGroups || {}).flatMap(
    ([groupId, selectedGroupValues]) => {
      // get all the accounts for the selected group value
      const group = argOptions?.groups ? argOptions?.groups[groupId] : [];
      return selectedGroupValues.flatMap((groupValue) => {
        return group.find((g) => g.value === groupValue)?.children || [];
      });
    }
  );

  const effectiveAccountIds = [
    ...(multiSelects || []),
    ...effectiveViaGroups,
  ].filter((v, i, a) => a.indexOf(v) === i);

  const effectiveOptions =
    argOptions?.options.filter((option) => {
      return effectiveAccountIds.includes(option.value);
    }) || [];
  const required = effectiveOptions.length === 0;

  return (
    <Box border="1px solid" borderColor="black" p={2}>
      <FormControl
        w="100%"
        isInvalid={
          multiSelectsError && multiSelectsError[argument.id] !== undefined
        }
      >
        <div>
          <FormLabel htmlFor="target.providerId">
            <Text textStyle={"Body/Medium"}>{argument.title}</Text>
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
        <FormErrorMessage> {argument.title} is required </FormErrorMessage>
      </FormControl>

      <Box border="1px solid" borderColor="black" mt={4} p={2}>
        <Heading size="md">Groups</Heading>
        {argument.groups ? (
          <>
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
                  <div>
                    <FormLabel htmlFor="target.providerId">
                      <Text textStyle={"Body/Medium"}>{group.title}</Text>
                    </FormLabel>
                    <HStack>
                      <MultiSelect
                        rules={{ required: required, minLength: 1 }}
                        fieldName={`target.argumentGroups.${argument.id}.${group.id}`}
                        options={argOptions.groups[group.id] || []}
                        shouldAddSelectAllOption={true}
                      />
                    </HStack>
                  </div>
                  <FormLabel htmlFor="target.providerId.filters.filterId"></FormLabel>
                  {/* TODO: msg will eventually be more detailed (one or more options) */}
                </FormControl>
              );
            })}
          </>
        ) : null}
      </Box>

      <Box>
        <Heading size="md">{"Effective " + argument.title + "s"}</Heading>
        <Wrap>
          {effectiveOptions &&
            effectiveOptions.map((c) => {
              return <CopyableOption label={c.label} value={c.value} />;
            })}
        </Wrap>
      </Box>
    </Box>
  );
};

export default ArgField;
