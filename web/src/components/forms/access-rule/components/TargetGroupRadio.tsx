import { CheckCircleIcon } from "@chakra-ui/icons";
import {
  Box,
  RadioProps,
  Spinner,
  useRadioGroup,
  UseRadioGroupProps,
  VStack,
} from "@chakra-ui/react";
import React, { useState } from "react";

import { ProviderIcon, ShortTypes } from "../../../icons/providerIcon";

import {
  CreateAccessRuleTargetFieldFilterExpessions,
  TargetGroup,
} from "../../../../utils/backend-client/types";
import { useFormContext } from "react-hook-form";
import { AccessRuleFormData } from "../CreateForm";
import { useAdminListTargetGroups } from "../../../../utils/backend-client/admin/admin";
import { TargetGroupField } from "./TargetGroupField";
import ReactSelect from "react-select";

interface TargetGroupRadioProps extends RadioProps {
  targetGroups: TargetGroup[];
}

const TargetGroupDropdown: React.FC<TargetGroupRadioProps> = (props) => {
  const { targetGroups } = props;

  const [selectedTargetgroup, setSelectedTargetgroup] = useState<TargetGroup>();

  const createOptions = () => {
    return targetGroups.map((t) => ({
      value: t.id,
      label: t.id,
    }));
  };

  return (
    <VStack>
      <Box as="label">
        <ReactSelect
          styles={{
            control: (provided, state) => ({
              ...provided,
              width: 420,
            }),
          }}
          options={createOptions()}
          onChange={(val) => {
            setSelectedTargetgroup(
              targetGroups.find((t) => t.id === val?.value)
            );
          }}
        />
      </Box>
      <Box>
        {!!selectedTargetgroup?.schema &&
          Object.values(selectedTargetgroup.schema).map((schema) => {
            return (
              <TargetGroupField
                targetGroup={selectedTargetgroup}
                fieldSchema={schema}
              />
            );
          })}
      </Box>
    </VStack>
  );
};

export const TargetGroupRadioSelector: React.FC<UseRadioGroupProps> = (
  props
) => {
  const { data } = useAdminListTargetGroups();
  const { getRootProps, getRadioProps } = useRadioGroup(props);
  const group = getRootProps();
  if (!data) {
    return <Spinner />;
  }

  return <TargetGroupDropdown targetGroups={data.targetGroups} />;
};
