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
  const { register, formState, watch } = useFormContext<AccessRuleFormData>();

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

  const target = watch("target");
  // filter the options by selected or effective from a group selection

  // assign effective options to a state variable
  // useEffect for handling and state setting
  const [effectiveOptions, setEffectiveOptions] = useState<Option[]>([]);

  useEffect(() => {
    console.log({ target, argument });
    if (target.argumentGroups || target.multiSelects) {
      const selected = target.multiSelects?.[argument.id];

      const selectedGroups = target.argumentGroups?.[argument.id];

      console.log({ selected, selectedGroups });
      let effectiveViaGroups: string[] = [];
      if (selectedGroups) {
        effectiveViaGroups = Object.entries(selectedGroups).flatMap(
          ([groupId, selectedGroupValues]) => {
            // get all the accounts for the selected group value
            const group = argOptions?.groups ? argOptions?.groups[groupId] : [];
            return selectedGroupValues.flatMap((groupValue) => {
              return group.find((g) => g.value === groupValue)?.children || [];
            });
          }
        );
        console.log({ effectiveViaGroups });
      }
      const effectiveAccountIds = [...selected, ...effectiveViaGroups].filter(
        (v, i, a) => a.indexOf(v) === i
      );

      const effectiveOptions =
        argOptions?.options.filter((option) => {
          return effectiveAccountIds.includes(option.value);
        }) || [];
      setEffectiveOptions(effectiveOptions);
    }
    return () => {
      setEffectiveOptions([]);
    };
  }, [target.argumentGroups, target.multiSelects, argOptions, argument]);

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
  // console.log({ target, argument });
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
              rules={{ required: true, minLength: 1 }}
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
                        rules={{ required: true, minLength: 1 }}
                        fieldName={`target.argumentGroups.${argument.id}.${group.id}`}
                        options={argOptions.groups[group.id] || []}
                        shouldAddSelectAllOption={true}
                      />
                    </HStack>
                  </div>
                  <FormLabel htmlFor="target.providerId.filters.filterId"></FormLabel>
                  {/* TODO: msg will eventually be more detailed (one or more options) */}
                  <FormErrorMessage>
                    {group.title} is required{" "}
                  </FormErrorMessage>
                </FormControl>
              );
            })}
          </>
        ) : null}
      </Box>

      <Box>
        <Heading size="md">{"Effective " + ""}</Heading>
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
