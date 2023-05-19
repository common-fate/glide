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
import React, { useEffect, useState } from "react";

import { ProviderIcon, ShortTypes } from "../../../icons/providerIcon";

import {
  useFieldArray,
  UseFieldArrayRemove,
  useFormContext,
} from "react-hook-form";
import { useAdminListTargetGroups } from "../../../../utils/backend-client/admin/admin";
import { TargetGroup } from "../../../../utils/backend-client/types";
import SelectMultiGeneric from "../../../SelectMultiGeneric";
import { AccessRuleFormData } from "../CreateForm";
import { TargetGroupField } from "./TargetGroupField";

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

  const { setError, clearErrors, trigger } =
    useFormContext<AccessRuleFormData>();

  // we want to set a form error when there is no selectedTargetgroup
  useEffect(() => {
    if (selectedTargetgroup.length === 0) {
      // @ts-ignore; weird circular error
      setError("targetgroups", {
        message: "Please select at least one target",
        type: "required",
      });
    } else {
      clearErrors("targetgroups");
    }
  }, [selectedTargetgroup]);

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
                onlyOne={true}
                keyUsedForFilter="id"
                inputArray={props.targetGroups}
                selectedItems={selectedTargetgroup}
                setSelectedItems={setSelectedTargetgroup}
                boxProps={{ w: "100%" }}
                renderFnTag={(item) => [
                  <ProviderIcon mr={1} shortType={item?.icon as ShortTypes} />,
                  item.id,
                ]}
                renderFnMenuSelect={(item) => [
                  <ProviderIcon shortType={item.icon as ShortTypes} mr={2} />,
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
