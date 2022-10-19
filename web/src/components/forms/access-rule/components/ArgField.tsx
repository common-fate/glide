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

  console.log({
    effectiveOptions,
    touchings: formState.touchedFields,
  });

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
          // border="1px solid"
          // borderColor="gray.300"
          // rounded="md"
          mt={4}
          pl={6}
          // py={4}
          // pt={8}
          // overflow="clip"
          pos="relative"
          w={{ base: "100%", md: "100%" }}
          minW={{ base: "100%", md: "400px", lg: "500px" }}
          // zIndex={1}
        >
          {/* <Tag
            size="md"
            pos="absolute"
            top={"-2px"}
            left={"-2px"}
            // zIndex={0}

            //       fontWeight: "400",
            // color: "#2D2F30",
            // fontSize: "16px",
            // lineHeight: "22.4px",
            fontWeight="400"
            color="#767676"
            fontSize="16px"
            lineHeight="22.4px"
            roundedBottomLeft="0"
            roundedTopRight="0"
            // borderLeft="none"
            // borderTop="white"
            // variant="outline"
            // colorScheme="yellow"
            pl={4}
            // grayscale emoji
            filter="grayscale(1);"
          >
            Dynamic Fields ⚡️
          </Tag> */}
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
                        bg="gray.200"
                        rounded="full"
                        filter="grayscale(1);"
                        transition="all .2s ease"
                        _hover={{
                          filter: "grayscale(0);",
                        }}
                      >
                        {"⚡️"}
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
                  return <CopyableOption label={c.label} value={c.value} />;
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
