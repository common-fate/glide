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
  argId: string;
  value: ArgDetails;
  providerId: string;
}

const ArgField = (props: ArgFieldProps) => {
  const { argId, value, providerId } = props;
  const { register } = useFormContext();
  const { data: argOptions } = useListProviderArgOptions(providerId, argId);

  return (
    <>
      <FormControl w="100%">
        <>
          {!!value?.filters ? (
            <>
              {Object.values(value.filters).map((group) => (
                <ArgGroupView
                  argId={argId}
                  groupDetail={group}
                  providerId={providerId}
                />
              ))}
            </>
          ) : null}
        </>
        <div>
          <FormLabel htmlFor="target.providerId">
            <Text textStyle={"Body/Medium"}>{value.title}</Text>
          </FormLabel>
          {argOptions?.hasOptions ? (
            <HStack>
              <MultiSelect
                rules={{ required: true, minLength: 1 }}
                fieldName={`target.with.${value.id}`}
                options={argOptions?.options || []}
                shouldAddSelectAllOption={true}
              />
            </HStack>
          ) : (
            <Input
              id="provider-vault"
              bg="white"
              placeholder={""}
              {...register(`target.withText.${value.id}`)}
            />
          )}
        </div>
        <FormLabel htmlFor="target.providerId.filters.filterId"></FormLabel>
      </FormControl>
    </>
  );
};

export default ArgField;
