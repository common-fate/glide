import React from "react";
import { FormLabel, Text } from "@chakra-ui/react";

import { MultiSelect } from "./Select";
import { useListProviderArgFilters } from "../../../../utils/backend-client/admin/admin";

interface FilterViewProps {
  providerId: string;
  argId: string;
  groupDetail: {
    id: string;
    title:string;
  };
}

const ArgGroupView = (props: FilterViewProps) => {
  const { groupDetail, providerId, argId } = props;

  const { data: argGroupingValues } = useListProviderArgFilters(
    providerId,
    argId,
    groupDetail.id
  );

  if (!argGroupingValues?.options){
    return null
  }

    return (
      <>
        <FormLabel>
          <Text textStyle={"Body/Medium"}>{groupDetail.title}</Text>
        </FormLabel>
        <MultiSelect
          fieldName={`target.withFilter.${argId}.${groupDetail.id}`}
          options={argGroupingValues?.options || []}
        />
      </>
    );
};

export default ArgGroupView;
