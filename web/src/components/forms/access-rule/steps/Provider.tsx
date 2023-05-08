import { DeleteIcon } from "@chakra-ui/icons";
import {
  Box,
  FormControl,
  FormLabel,
  IconButton,
  Text,
  VStack,
} from "@chakra-ui/react";
import React from "react";
import { Controller, useFormContext } from "react-hook-form";
import { useAdminListTargetGroups } from "../../../../utils/backend-client/admin/admin";
import { CreateAccessRuleTargetFieldFilterExpessions } from "../../../../utils/backend-client/types";
import { TargetGroupRadioSelector } from "../components/TargetGroupRadio";

import { AccessRuleFormData } from "../CreateForm";
import { FormStep } from "./FormStep";

export const TargetStep: React.FC = () => {
  const methods = useFormContext<AccessRuleFormData>();
  const targets = methods.watch("targets");
  return (
    <FormStep
      heading="Target"
      subHeading="The permissions that the rule gives access to"
      fields={["targets", "target.providerId"]}
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
            render={({ field: { ref, onChange, value, ...rest } }) => {
              return (
                // The implemenbtation here is currently a reduced scope where only one target group can be selected for an access rule, and there are no resource filtering options
                <TargetGroupRadioSelector
                  value={value.length > 0 ? value[0].targetGroupId : undefined}
                  onChange={(targetGroupId: string) => {
                    onChange([{ targetGroupId, fieldFilterExpessions: {} }]);
                  }}
                />
              );
            }}
          />
        </FormControl>
      </>
    </FormStep>
  );
};
