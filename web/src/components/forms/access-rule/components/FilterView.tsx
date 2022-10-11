import React from "react";
import { FormLabel, Text } from "@chakra-ui/react";

import { MultiSelect } from "../components/Select";
import { useListProviderArgFilters } from "../../../../utils/backend-client/admin/admin";

interface FilterViewProps {
  providerId: string;
  argId: string;
  filter: any;
}

const FilterView = (props: FilterViewProps) => {
  const { filter, providerId, argId } = props;

  const { data: argFilterValues } = useListProviderArgFilters(
    providerId,
    argId,
    filter.id
  );

  return (
    <>
      <FormLabel>
        <Text textStyle={"Body/Medium"}>{filter.title}</Text>
      </FormLabel>
      <MultiSelect
        fieldName={`target.withFilter.${argId}.${filter.id}`}
        options={argFilterValues?.options || []}
      />
    </>
  );
};

export default FilterView;
