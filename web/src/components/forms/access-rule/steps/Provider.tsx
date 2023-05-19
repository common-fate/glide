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
