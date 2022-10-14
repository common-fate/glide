import { FormControl, FormLabel, HStack, Input, Text } from "@chakra-ui/react";
import React from "react";
import { useFormContext } from "react-hook-form";

import ArgGroupView from "./ArgGroupView";
import { MultiSelect } from "../components/Select";
import { useListProviderArgOptions } from "../../../../utils/backend-client/admin/admin";

interface ArgDetails {
  description: string;
  title: string;
  id: string;
  type: string;
  filters?: {
    [key: string]: {
      id: string;
      title: string;
    };
  };
}

interface ArgFieldProps {
  argument: ArgDetails;
  providerId: string;
}

const ArgField = (props: ArgFieldProps) => {
  const { argument, providerId } = props;
  const { register } = useFormContext();
  const { data: argOptions } = useListProviderArgOptions(
    providerId,
    argument.id
  );

  return (
    <>
      <FormControl w="100%">
        <>
          {!!argument?.filters ? (
            <>
              {Object.values(argument.filters).map((group) => (
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
          {argOptions?.hasOptions ? (
            <HStack>
              <MultiSelect
                rules={{ required: true, minLength: 1 }}
                fieldName={`target.with.${argument.id}`}
                options={argOptions?.options || []}
                shouldAddSelectAllOption={true}
              />
            </HStack>
          ) : (
            <Input
              id="provider-vault"
              bg="white"
              placeholder={""}
              {...register(`target.withText.${argument.id}`)}
            />
          )}
        </div>
        <FormLabel htmlFor="target.providerId.filters.filterId"></FormLabel>
      </FormControl>
    </>
  );
};

export default ArgField;
