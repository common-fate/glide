import React from "react";
import { FormLabel, Text } from "@chakra-ui/react";

import { MultiSelect } from "./Select";
import { Group } from "../../../../utils/backend-client/types/accesshandler-openapi.yml";

interface FilterViewProps {
  providerId: string;
  argId: string;
  groupDetail: Group;
}

const ArgGroupView = (props: FilterViewProps) => {
  const { groupDetail, providerId, argId } = props;
  return (
    <>
      <FormLabel>
        <Text textStyle={"Body/Medium"}>{groupDetail.title}</Text>
      </FormLabel>
      <MultiSelect
        fieldName={`target.withFilter.${argId}.${groupDetail.id}`}
        options={groupDetail.options || []}
      />
    </>
  );
};

export default ArgGroupView;
