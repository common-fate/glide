import { DeleteIcon } from "@chakra-ui/icons";
import {
  Box,
  FormControl,
  FormLabel,
  IconButton,
  Text,
  VStack,
} from "@chakra-ui/react";
import React, { useState } from "react";
import { Controller, useFormContext } from "react-hook-form";
import { useAdminListTargetGroups } from "../../../../utils/backend-client/admin/admin";
import { CreateAccessRuleTargetFieldFilterExpessions } from "../../../../utils/backend-client/types";
import { MultiTargetGroupSelector } from "../components/TargetGroupRadio";

import { AccessRuleFormData } from "../CreateForm";
import { FormStep } from "./FormStep";
import SelectMultiGeneric from "../../../SelectMultiGeneric";
import {
  ProviderIcon,
  ShortTypes,
  shortTypeValues,
} from "../../../icons/providerIcon";

export const TargetStep: React.FC = () => {
  const methods = useFormContext<AccessRuleFormData>();
  // const targets = methods.watch("targetgroups");

  type MockProvider = {
    name: string;
    shortType: ShortTypes;
  };

  const [filteredInput, setFilteredInput] = useState<MockProvider[]>([]);

  const [selectedProviders2, setSelectedProviders2] = useState<MockProvider[]>(
    Object.entries(shortTypeValues).map(([shortType, name]) => ({
      name: name,
      shortType: shortType as ShortTypes,
    }))
  );

  return (
    <FormStep
      heading="Target"
      subHeading="The permissions that the rule gives access to"
      fields={["targetgroups", "targetFieldMap"]}
      // preview={<Preview target={target} provider={provider} />}
      isFieldLoading={false}
    >
      <>
        <FormControl isInvalid={false}>
          <FormLabel htmlFor="target.providerId">
            <Text textStyle={"Body/Medium"}>Target</Text>
          </FormLabel>
          {/* <SelectMultiGeneric
            keyUsedForFilter="name"
            inputArray={selectedProviders2}
            selectedItems={filteredInput}
            setSelectedItems={setFilteredInput}
            boxProps={{ mt: 4 }}
            renderFnTag={(item) => [
              <ProviderIcon shortType={item.shortType} />,
              item.name,
            ]}
            renderFnMenuSelect={(item) => [
              <ProviderIcon shortType={item.shortType} mr={2} />,
              item.name,
            ]}
          /> */}
          <Controller
            control={methods.control}
            name={"targets"}
            render={({ field }) => {
              return (
                <MultiTargetGroupSelector
                  field={field}
                  control={methods.control}
                />
              );
            }}
          />
        </FormControl>
      </>
    </FormStep>
  );
};
