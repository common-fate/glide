import { FormControl, FormLabel, HStack, Input, Text } from "@chakra-ui/react";
import React from "react";
import { useFormContext } from "react-hook-form";

import ArgGroupView from "./ArgGroupView";
import { MultiSelect } from "../components/Select";
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
  const { register } = useFormContext();
  const { data: argOptions } = useListProviderArgOptions(
    providerId,
    argument.id
  );

  if (argument.formElement === ArgumentFormElement.INPUT) {
    return (
      <FormControl w="100%">
        <FormLabel htmlFor="target.providerId">
          <Text textStyle={"Body/Medium"}>{argument.title}</Text>
        </FormLabel>
        <Input
          id="provider-vault"
          bg="white"
          placeholder={""}
          {...register(`target.withText.${argument.id}`)}
        />
      </FormControl>
    );
  }

  return (
    <>
      <FormControl w="100%">
        <>
          {argOptions?.groups ? (
            <>
              {Object.values(argOptions.groups).map((group) => (
                <ArgGroupView
                  argId={argument.id}
                  groupDetail={group}
                  providerId={providerId}
                />
              ))}
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
      </FormControl>
    </>
  );
};

export default ArgField;
