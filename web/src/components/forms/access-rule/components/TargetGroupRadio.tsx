import {
  Accordion,
  AccordionButton,
  AccordionIcon,
  AccordionItem,
  AccordionPanel,
  Box,
  Button,
  Spinner,
  VStack,
} from "@chakra-ui/react";
import React, { useState, useEffect } from "react";

import { ProviderIcon, ShortTypes } from "../../../icons/providerIcon";

import { TargetGroup } from "../../../../utils/backend-client/types";
import {
  useFieldArray,
  UseFieldArrayRemove,
  useFormContext,
} from "react-hook-form";
import { useAdminListTargetGroups } from "../../../../utils/backend-client/admin/admin";
import { TargetGroupField } from "./TargetGroupField";
import ReactSelect from "react-select";
import { AccessRuleFormData } from "../CreateForm";

interface TargetGroupDropdownProps {
  item: Record<"id", string>;
  targetGroups: TargetGroup[];
  remove: UseFieldArrayRemove;
  index: number;
}

const TargetGroupDropdown: React.FC<TargetGroupDropdownProps> = (props) => {
  const { remove, item, index } = props;
  const [selectedTargetgroup, setSelectedTargetgroup] = useState<TargetGroup>();

  const methods = useFormContext<AccessRuleFormData>();
  const targetgroups = methods.watch("targetgroups");

  // CreateOption will exclude already selected targetgroups from new targetgroup dropdown.
  const createOptions = () => {
    const excludingItemTargetgroupIds = targetgroups
      ? Object.keys(targetgroups).filter((e) => e)
      : [];

    const excludedList = props.targetGroups.filter(
      (t) => !excludingItemTargetgroupIds.includes(t.id)
    );

    return excludedList.map((t) => ({
      value: t.id,
      label: t.id,
    }));
  };

  return (
    <>
      <Accordion defaultIndex={[0]} allowMultiple>
        <AccordionItem key={item.id}>
          <AccordionButton>
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
                  props.targetGroups.find((t) => t.id === val?.value)
                );
              }}
            />
            {index != 0 && (
              <Button onClick={() => remove(index)}>Delete </Button>
            )}
            <AccordionIcon />
          </AccordionButton>
          <AccordionPanel>
            <VStack>
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
          </AccordionPanel>
        </AccordionItem>
      </Accordion>
    </>
  );
};

interface MultiTargetGroupSelectorProps {
  // field: ControllerRenderProps<AccessRuleFormData, any>;
  // control: Control<AccessRuleFormData, any>;
  field: any;
  control: any;
}

export const MultiTargetGroupSelector: React.FC<
  MultiTargetGroupSelectorProps
> = (props) => {
  const { data } = useAdminListTargetGroups();

  const { fields, append, remove, insert } = useFieldArray({
    control: props.control,
    name: "targetFieldMap",
  });

  useEffect(() => {
    if (!fields.length) {
      insert(0, { targetGroupId: "" });
    }
  }, []);

  if (!data) {
    return <Spinner />;
  }

  return (
    <>
      {fields.map((item, index) => (
        <TargetGroupDropdown
          item={item}
          index={index}
          remove={remove}
          targetGroups={data.targetGroups}
        />
      ))}
      <Button type="button" onClick={() => append({})}>
        + Target
      </Button>
    </>
  );
};
