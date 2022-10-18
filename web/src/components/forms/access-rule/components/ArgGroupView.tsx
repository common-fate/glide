import React from "react";
import { FormLabel, Text } from "@chakra-ui/react";

import { MultiSelect } from "./Select";
import {
  Group,
  GroupOption,
} from "../../../../utils/backend-client/types/accesshandler-openapi.yml";

interface FilterViewProps {
  providerId: string;
  argId: string;
  group: Group;
  options: GroupOption[];
}

const ArgGroupView = (props: FilterViewProps) => {
  const { group, options, providerId, argId } = props;
  return (
    <>
      <FormLabel>
        <Text textStyle={"Body/Medium"}>{group.title}</Text>
      </FormLabel>
      <MultiSelect
        fieldName={`target.withFilter.${argId}.${group.id}`}
        options={options}
      />
    </>
  );
};

export default ArgGroupView;
