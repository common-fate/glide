import {
  Accordion,
  AccordionButton,
  AccordionIcon,
  AccordionItem,
  AccordionPanel,
  Box,
  Button,
  Spinner,
  Stack,
  Text,
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
import SelectMultiGeneric from "../../../SelectMultiGeneric";

interface TargetGroupDropdownProps {
  item: Record<"id", string>;
  targetGroups: TargetGroup[];
  remove: UseFieldArrayRemove;
  index: number;
}

const TargetGroupDropdown: React.FC<TargetGroupDropdownProps> = (props) => {
  const { remove, item, index } = props;
  const [selectedTargetgroup, setSelectedTargetgroup] = useState<TargetGroup[]>(
    []
  );

  const methods = useFormContext<AccessRuleFormData>();
  const targetgroups = methods.watch("targetgroups");

  // methods.setValue(`targetgroups[${index}].id`, false);

  // CreateOption will exclude already selected targetgroups from new targetgroup dropdown.
  // const createOptions = () => {
  //   const excludingItemTargetgroupIds = targetgroups
  //     ? Object.keys(targetgroups).filter((e) => e)
  //     : [];

  //   const excludedList = props.targetGroups.filter(
  //     (t) => !excludingItemTargetgroupIds.includes(t.id)
  //   );

  //   return excludedList.map((t) => ({
  //     value: t.id,
  //     label: t.id,
  //   }));
  // };

  return (
    <Box bg="neutrals.100" borderColor="neutrals.300" rounded="lg">
      <Accordion defaultIndex={[0]} allowMultiple>
        <AccordionItem key={item.id} border="none">
          <AccordionButton
            p={2}
            bg="white"
            // bg="neutrals.100"
            roundedTop="md"
            borderColor="neutrals.300"
            borderWidth="1px"
            sx={{
              "&[aria-expanded='false']": {
                roundedBottom: "md",
              },
            }}
            pos="relative"
          >
            <AccordionIcon mr={1} />

            <Box
              as="span"
              flex="1"
              textAlign="left"
              sx={{
                p: { lineHeight: "120%", textStyle: "Body/Extra Small" },
              }}
            >
              <Text color="neutrals.700">Select a target</Text>
            </Box>

            {index != 0 && (
              <Button
                variant="brandSecondary"
                size="xs"
                onClick={() => remove(index)}
                position="absolute"
                right={1}
                top={1}
              >
                Delete{" "}
              </Button>
            )}
          </AccordionButton>
          <AccordionPanel
            borderColor="neutrals.300"
            borderTop="none"
            roundedBottom="md"
            borderWidth="1px"
            bg="white"
            p={4}
          >
            <VStack align="stretch" spacing={2}>
              <SelectMultiGeneric
                keyUsedForFilter="id"
                inputArray={props.targetGroups}
                selectedItems={selectedTargetgroup}
                setSelectedItems={setSelectedTargetgroup}
                boxProps={{ w: "100%" }}
                renderFnTag={(item) => [
                  <ProviderIcon mr={1} shortType={item?.icon} />,
                  item.id,
                ]}
                renderFnMenuSelect={(item) => [
                  <ProviderIcon mr={1} shortType={item.icon} mr={2} />,
                  item.id,
                ]}
              />
              <Box w="100%">
                {selectedTargetgroup[0]?.schema &&
                  Object.values(selectedTargetgroup[0].schema).map((schema) => {
                    return (
                      <TargetGroupField
                        targetGroup={selectedTargetgroup[0]}
                        fieldSchema={schema}
                      />
                    );
                  })}
              </Box>
            </VStack>
          </AccordionPanel>
        </AccordionItem>
      </Accordion>
    </Box>
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
      <Stack>
        {fields.map((item, index) => (
          <TargetGroupDropdown
            item={item}
            index={index}
            remove={remove}
            targetGroups={data.targetGroups}
          />
        ))}
      </Stack>
      <Button
        variant="brandSecondary"
        size="sm"
        mt={2}
        type="button"
        onClick={() => append({})}
      >
        Target +
      </Button>
    </>
  );
};
