import {
  FormControl,
  FormErrorMessage,
  FormLabel,
  HStack,
  Input,
  Text,
} from "@chakra-ui/react";
import React from "react";
import { useFormContext } from "react-hook-form";

import ArgGroupView from "./ArgGroupView";
import { MultiSelect } from "../components/Select";
import { AccessRuleFormData } from "../CreateForm";
import { useListProviderArgOptions } from "../../../../utils/backend-client/admin/admin";
import {
  Argument,
  ArgumentFormElement,
} from "../../../../utils/backend-client/types/accesshandler-openapi.yml";

interface ArgFieldProps {
  argument: Argument;
  providerId: string;
}

const ArgField = (props: ArgFieldProps) => {
  const { argument, providerId } = props;
  const { register, formState } = useFormContext<AccessRuleFormData>();

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
          {...register(`target.withText.${argument.id}`)}
        />
      </FormControl>
    );
  }

  const withError = formState.errors.target?.with;

  return (
    <>
      <FormControl
        w="100%"
        isInvalid={withError && withError[argument.id] !== undefined}
      >
        <>
          {argument.groups ? (
            <>
              {Object.values(argument.groups).map((group) => {
                // catch the unexpected case where there are no options for group
                if (
                  argOptions?.groups == undefined ||
                  !argOptions.groups[group.id]
                ) {
                  return null;
                }
                return (
                  <ArgGroupView
                    argId={argument.id}
                    group={group}
                    options={argOptions.groups[group.id]}
                    providerId={providerId}
                  />
                );
              })}
            </>
          ) : null}
        </>
        <div>
          <FormLabel htmlFor="target.providerId">
            <Text textStyle={"Body/Medium"}>{argument.title}</Text>
          </FormLabel>
          <HStack>
            <MultiSelect
              rules={{ required: true, minLength: 1 }}
              fieldName={`target.with.${argument.id}`}
              options={argOptions?.options || []}
              shouldAddSelectAllOption={true}
            />
          </HStack>
        </div>
        <FormLabel htmlFor="target.providerId.filters.filterId"></FormLabel>
        <FormErrorMessage> {argument.title} is required </FormErrorMessage>
      </FormControl>
    </>
  );
};

export default ArgField;
