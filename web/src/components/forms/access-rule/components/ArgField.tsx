import {
  FormControl,
  FormLabel,
  HStack,
  Input,
  Text,
  Select,
} from "@chakra-ui/react";
import React from "react";
import { useFormContext } from "react-hook-form";

import { MultiSelect } from "../components/Select";
import { useListProviderArgOptions } from "../../../../utils/backend-client/admin/admin";
import FilterView from "./FilterView";

interface ArgFieldProps {
  argId: string;
  data: any;
  providerId: string;
}

const ArgField = (props: ArgFieldProps) => {
  const { argId, data: value, providerId } = props;
  const { register } = useFormContext();
  const { data: argOptions } = useListProviderArgOptions(providerId, argId);

  return (
    <>
      <FormControl w="100%">
        <>
          {!!value?.filters ? (
            <>
              {Object.values(value.filters).map((filter: any) => (
                <FilterView
                  filter={filter}
                  providerId={providerId}
                  argId={argId}
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
                fieldName={`target.with.${value.title}`}
                options={argOptions?.options || []}
                shouldAddSelectAllOption={true}
              />
            </HStack>
          ) : (
            <Input
              id="provider-vault"
              bg="white"
              placeholder={""}
              {...register(`target.withText.${value.argId}`)}
            />
          )}
        </div>
        <FormLabel htmlFor="target.providerId.filters.filterId"></FormLabel>
      </FormControl>
    </>
  );
};

export default ArgField;
